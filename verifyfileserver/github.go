package main

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/kataras/iris"
)

func github(ctx iris.Context) {
	addr := ctx.URLParam("addr")
	if strings.HasPrefix(addr, "https://github.com/") == false {
		logger.Println("reject", addr)
		ctx.StatusCode(404)
		return
	}
	res, err := http.DefaultClient.Get(addr)
	if err != nil {
		logger.Println("connect error", addr)
		ctx.StatusCode(404)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		logger.Println("not found", addr)
		ctx.StatusCode(404)
		return
	}
	logger.Println("start", addr)

	filename := path.Base(addr[8:])
	ctx.Header("Content-Disposition", "attachment;filename="+filename)
	//ctx.StatusCode(200)
	var writer io.Writer = ctx.ResponseWriter()
	// if ctx.ClientSupportsGzip() {
	// 	writer = ctx.GzipResponseWriter()
	// } else {
	// 	writer = ctx.ResponseWriter()
	// }

	var buf [3000]byte
	for {
		n, _ := io.CopyBuffer(writer, res.Body, buf[:])
		if n == 0 {
			break
		}
	}

	logger.Println("end", addr)
}
