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
<title>{{.title}}</title>
<head>
<body>
{{.body}}
</body>
</html>`

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
		sendFile(ctx, filename[0])
		log.Printf("%s Down %s\n", ctx.RemoteAddr(), filename[0])
	} else {
		log.Printf("%s ERR %s\n", ctx.RemoteAddr(), pkg)
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

func sendFile(ctx iris.Context, filename string) {
	fname := relatePath("files", filename)
	info, err := os.Stat(fname)
	if err != nil {
		ctx.StatusCode(404)
		log.Printf("File Error:%s", filename)
		return
	}

	ctx.Header("Content-Length", fmt.Sprint(info.Size()))
	ctx.SendFile(fname, filename)
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
	t.Execute(ctx.ResponseWriter(), data)
}

func main() {
	var addr = flag.String("addr", ":8080", "format [IP:Port]")
	flag.Parse()

	app := iris.New()

	server := new(myServer)
	log.Println(server.Init())
	defer server.Close()

	initMdTmpl()

	app.Get("/get/{pkg}", func(ctx iris.Context) {
		pkg := ctx.Params().Get("pkg")
		if server.IsVerified(ctx) {
			server.ServePkg(ctx, pkg)
		} else {
			ctx.StatusCode(404)
		}
	})

	app.Get("/down/{key}", func(ctx iris.Context) {
		key := ctx.Params().Get("key")
		kw := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
		if server.CheckKey(kw) {
			server.SetVerified(ctx, true)
			server.ShowPkgs(ctx)
			log.Printf("%s IP: %s\n", key, ctx.RemoteAddr())
		} else {
			server.SetVerified(ctx, false)
			ctx.StatusCode(404)
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
		if strings.HasSuffix(strings.ToLower(fn), ".md") {
			sendMarkdown(ctx, fn)
		} else {
			ctx.ServeFile(relatePath("files", fn), ctx.ClientSupportsGzip())
		}

		log.Printf("%s Get /%s\n", ctx.RemoteAddr(), fn)
	})

	app.Get("/", func(ctx iris.Context) {
		sendMarkdown(ctx, "index.md")
		log.Printf("%s Get /\n", ctx.RemoteAddr())
	})

	app.Run(iris.Addr(*addr))
}
