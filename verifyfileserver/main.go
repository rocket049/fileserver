package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/kataras/iris"

	"database/sql"

	"github.com/kataras/iris/sessions"
	_ "github.com/mattn/go-sqlite3"
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
<h2><a href="/get/win32">Windows 32位</a></h2>
<h2><a href="/get/win64">Windows 64位</a></h2>
<h2><a href="/get/linux">Linux amd64</a></h2>
</body>
</html>
`
)

func realatePath(items ...string) string {
	exe1, _ := os.Executable()
	base := filepath.Dir(exe1)
	paths := append([]string{base}, items...)
	return filepath.Join(paths...)
}

type myServer struct {
	db   *sql.DB
	sess *sessions.Sessions
}

func (s *myServer) Init() error {
	var err error
	s.sess = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
	s.db, err = sql.Open("sqlite3", realatePath("keys.db"))
	if err != nil {
		return err
	}
	_, err = s.db.Exec("create table if not exists keys (keyword text unique,time int);")
	return err
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
	fname := realatePath("files", filename)
	//ctx.SendFile(fname, filename)
	typ := mime.TypeByExtension(".zip")
	fp, err := os.Open(fname)
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	defer fp.Close()
	info, _ := fp.Stat()
	ctx.Header("Content-Type", typ)
	ctx.Header("Content-Length", fmt.Sprint(info.Size()))
	ctx.Header("Content-disposition", "attachment; filename="+filename)
	ctx.StatusCode(200)
	io.Copy(ctx.ResponseWriter(), fp)
}

func main() {
	app := iris.New()

	server := new(myServer)
	server.Init()
	defer server.Close()

	app.Get("/get/{pkg}", func(ctx iris.Context) {
		pkg := ctx.Params().Get("pkg")
		if server.IsVerified(ctx) {
			switch pkg {
			case "win32":
				sendFile(ctx, "secret-diary-win32.zip")

			case "win64":
				sendFile(ctx, "secret-diary-win64.zip")

			case "linux":
				sendFile(ctx, "secret-diary-ubuntu_amd64.zip")

			}
		} else {
			ctx.StatusCode(404)
		}
	})

	app.Get("/{key}", func(ctx iris.Context) {
		key := ctx.Params().Get("key")
		kw := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
		if server.CheckKey(kw) {
			server.SetVerified(ctx, true)
			ctx.WriteString(homepage)
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

	app.Run(iris.Addr(":8080"))
}
