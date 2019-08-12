package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kataras/iris"

	"database/sql"

	"github.com/kataras/iris/sessions"
	_ "github.com/mattn/go-sqlite3"
	md "github.com/russross/blackfriday"
)

const (
	cookieNameForSessionID = "clientstatusid"
	homepage               = `<html>
<head>
	<meta http-equiv="content-type" content="text/html;charset=utf-8"/>
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<title>选择对应你的电脑的版本</title>
</head>
<body>
<h1>选择对应你的电脑的版本</h1>
{{range .}}<h2><a href="/get/{{.Pkg}}">{{.Name}}</a></h2>
{{end}}
</body>
</html>
`
)

var mdTmpl = `<html>
<head>
<meta http-equiv="content-type" content="text/html;charset=utf-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<link type="text/css" rel="stylesheet" href="/style.css"/>
<script src="../jquery-3.3.1.min.js"></script>
<title>{{.title}}</title>
<head>
<body>
<div id="box">
    <div class="ad">
        <ul class="slider">
            <li><img src="images/1.png" /></li>
            <li><img src="images/2.png" /></li>
            <li><img src="images/3.png" /></li>
            <li><img src="images/4.png" /></li>
            <li><img src="images/5.png" /></li>
        </ul>
        <ul class="num">
            <li>1</li>
            <li>2</li>
            <li>3</li>
            <li>4</li>
            <li>5</li>
        </ul>
    </div>
</div>
{{.body}}
<div style="background:#E8E7E3;padding:20px 0;">
  <div style="font-size:x-small;text-align:center">备案号：苏ICP备19034936号 &nbsp;&nbsp; 站长邮箱：fuhuizn@163.com</div>
  <div style="width:300px;margin:0 auto;font-size:x-small;">
		 		<a target="_blank" href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=32050702010684" style="display:inline-block;text-decoration:none;height:20px;line-height:20px;"><img src="/bei-an.png" style="float:left;"/><p style="float:left;height:20px;line-height:20px;margin: 0px 0px 0px 5px; color:black;">苏公网安备 32050702010684号</p></a>
		 	</div>
</div>
<script type="text/javascript">
	$(function(){
		var index=0;
		var len=$(".num>li").length;
		var timer;
		$(".num>li").mouseover(function(){
			index=$(this).index()
			toggleImg(index)
		}).eq(0).mouseover()
		$(".ad").hover(function(){
			clearInterval(timer)
		},function(){
			timer=setInterval(function(){
				toggleImg(index);
				index++;
				if(index==len){
					index=0;
				}
			},2000)
		}).trigger("mouseleave")
		
	})
	function toggleImg(index){
		var adHeight=$(".ad").height();
		$(".slider").stop().animate({"top":-adHeight*index},1000)
		$(".num li").removeClass("on").eq(index).addClass("on")
	}
</script>
</body>
</html>`

var logger *log.Logger

func relatePath(items ...string) string {
	exe1, _ := os.Executable()
	base := filepath.Dir(exe1)
	paths := append([]string{base}, items...)
	res := filepath.Join(paths...)
	if strings.Contains(res, "..") {
		return "404"
	}
	return res

}

type myServer struct {
	db    *sql.DB
	sess  *sessions.Sessions
	files map[string][]string
}

func initMdTmpl() {
	tpl, err := ioutil.ReadFile(relatePath("md.tpl"))
	if err != nil {
		return
	} else {
		mdTmpl = string(tpl)
	}
}

func (s *myServer) Init() error {
	var err error
	s.sess = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
	s.db, err = sql.Open("sqlite3", relatePath("keys.db"))
	if err != nil {
		return err
	}
	_, err = s.db.Exec("create table if not exists keys (keyword text unique,time int);")
	if err != nil {
		return err
	}

	s.files = make(map[string][]string)
	jsonData, err := ioutil.ReadFile(relatePath("files.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &s.files)
	if err != nil {
		return err
	}
	return nil
}

func (s *myServer) ServePkg(ctx iris.Context, pkg string) {
	filename, ok := s.files[pkg]
	if ok {
		err := sendFile(ctx, filename[0])
		if err != nil {
			logger.Printf("ERROR %s\n", err.Error())
		} else {
			logger.Printf("%s get %s\n", ctx.RemoteAddr(), filename[0])
		}

	} else {
		logger.Printf("ERROR %s get %s\n", ctx.RemoteAddr(), pkg)
	}
}

type PageItem struct {
	Pkg  string
	Name string
}

func (s *myServer) ShowPkgs(ctx iris.Context) error {
	var items = make([]PageItem, len(s.files))
	i := 0
	for k, v := range s.files {
		items[i].Pkg = k
		items[i].Name = v[1]
		i++
	}
	t := template.New("")
	t, err := t.Parse(homepage)
	if err != nil {
		return err
	}
	return t.Execute(ctx.ResponseWriter(), items)
}

func (s *myServer) Close() {
	s.db.Close()
	s.sess.DestroyAll()
}

func (s *myServer) CheckKey(key string) bool {
	res, err := s.db.Query("select time from keys where keyword=? limit 1;", key)
	if err != nil {
		return false
	}
	defer res.Close()
	return res.Next()
}

func (s *myServer) AddKey(key string) {
	kw := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
	s.db.Exec("insert into keys(keyword,time) values (?,7)", kw)
}

func (s *myServer) SetVerified(ctx iris.Context, b bool) {
	session := s.sess.Start(ctx)
	session.Set("authenticated", b)
}

func (s *myServer) IsVerified(ctx iris.Context) bool {
	b, err := s.sess.Start(ctx).GetBoolean("authenticated")
	if err != nil {
		return false
	}
	return b
}

func getSize(filename string) int64 {
	fh, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return fh.Size()
}

func sendFile(ctx iris.Context, filename string) error {
	fname := relatePath("files", filename)
	info, err := os.Stat(fname)
	if err != nil {
		ctx.StatusCode(404)
		return err
	}

	ctx.Header("Content-Length", fmt.Sprint(info.Size()))
	ctx.SendFile(fname, filename)
	return nil
}

func getTitle(p []byte) string {
	n := bytes.IndexByte(p, '\n')
	if n <= 0 {
		return ""
	}
	line1 := string(p[0:n])
	return strings.Trim(line1, "# \n\r\t")
}

func sendMarkdown(ctx iris.Context, filename string) {
	fname := relatePath("files", filename)
	fstat, err := os.Lstat(fname)
	if err != nil {
		ctx.StatusCode(404)
		return
	}

	cacheName := relatePath("files", filename+".htm")
	cacheStat, err := os.Stat(cacheName)
	if err == nil {
		if cacheStat.ModTime().Unix() > fstat.ModTime().Unix() {
			ctx.ServeFile(cacheName, ctx.ClientSupportsGzip())
			return
		}
	}

	fp, err := os.Create(cacheName)
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	defer fp.Close()

	writer := io.MultiWriter(fp, ctx.ResponseWriter())

	data := make(map[string]string)

	file1, err := os.Open(relatePath("files", filename))
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	stat1, _ := file1.Stat()
	var buf []byte = make([]byte, stat1.Size())
	n, _ := io.ReadFull(file1, buf)
	file1.Close()
	if int64(n) == stat1.Size() {
		data["title"] = getTitle(buf)

		body := md.Run(buf, md.WithExtensions(md.CommonExtensions))
		data["body"] = string(body)

	} else {
		ctx.StatusCode(500)
		return
	}
	t := template.New("")
	t.Parse(mdTmpl)
	t.Execute(writer, data)
}

func adRedirect(ctx iris.Context) {
	ctx.Redirect("https://www.baidu.com/", 302)
}

func main() {
	var addr = flag.String("addr", ":8080", "format [IP:Port]")
	flag.Parse()

	logger = log.New(os.Stdout, "", log.LstdFlags)

	app := iris.New()

	server := new(myServer)
	logger.Println(server.Init())
	defer server.Close()

	initMdTmpl()

	app.Get("/get/{pkg}", func(ctx iris.Context) {
		pkg := ctx.Params().Get("pkg")
		if server.IsVerified(ctx) {
			server.ServePkg(ctx, pkg)
		} else {
			ctx.StatusCode(404)
			logger.Printf("ERROR 404 get %s\n", pkg)
		}
	})

	app.Get("/down/{key}", func(ctx iris.Context) {
		key := ctx.Params().Get("key")
		kw := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
		if server.CheckKey(kw) {
			server.SetVerified(ctx, true)
			server.ShowPkgs(ctx)
			logger.Printf("down %s IP: %s\n", key, ctx.RemoteAddr())
		} else {
			server.SetVerified(ctx, false)
			ctx.StatusCode(404)
			logger.Printf("ERROR down %s IP: %s\n", key, ctx.RemoteAddr())
		}

	})

	app.Get("/add/{key}", func(ctx iris.Context) {
		key := ctx.Params().Get("key")
		switch ctx.RemoteAddr() {
		case "127.0.0.1":
			server.AddKey(key)
		case "::1":
			server.AddKey(key)
		}
	})

	app.Get("/{key:path}", func(ctx iris.Context) {
		fn := strings.TrimSpace(ctx.Params().Get("key"))
		if strings.Contains(fn, "..") {
			logger.Printf("ERROR %s Get /%s\n", ctx.RemoteAddr(), fn)
			ctx.StatusCode(404)
			return
		}
		if strings.HasSuffix(strings.ToLower(fn), ".md") {
			sendMarkdown(ctx, fn)
			logger.Printf("%s Get /%s\n", ctx.RemoteAddr(), fn)
		} else {
			err := ctx.ServeFile(relatePath("files", fn), ctx.ClientSupportsGzip())
			if err != nil {
				logger.Printf("ERROR %s Get /%s\n", ctx.RemoteAddr(), fn)
			} else {
				logger.Printf("%s Get /%s\n", ctx.RemoteAddr(), fn)
			}
		}

	})

	app.Get("/", func(ctx iris.Context) {
		sendMarkdown(ctx, "index.md")
		logger.Printf("%s Get /\n", ctx.RemoteAddr())
	})

	// /github?addr=https://github.com/...
	// app.Get("/github", github)

	app.Get("/redirect", adRedirect)

	app.Run(iris.Addr(*addr))
}
