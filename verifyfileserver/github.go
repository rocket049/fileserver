package main

import (
	"io"
	"os/exec"
	"path"
	"strings"

	"github.com/kataras/iris"
)

func github(ctx iris.Context) {
	addr := ctx.URLParam("addr")
	if strings.HasPrefix(addr, "https://github.com/") == false {
		logger.Println("reject", addr)
		return
	} else {
		logger.Println("start", addr)
	}
	cmd := exec.Command("wget", "-o", "wget.log", "-O", "-", addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Println(err)
		ctx.StatusCode(404)
		return
	}
	defer stdout.Close()
	err = cmd.Start()
	if err != nil {
		logger.Println(err)
		ctx.StatusCode(404)
		return
	}
	defer cmd.Process.Kill()

	filename := path.Base(addr[8:])
	ctx.Header("Content-Disposition", "attachment;filename="+filename)
	ctx.StatusCode(200)
	var writer io.Writer
	if ctx.ClientSupportsGzip() {
		writer = ctx.GzipResponseWriter()
	} else {
		writer = ctx.ResponseWriter()
	}

	var buf [3000]byte
	for {
		n, _ := io.CopyBuffer(writer, stdout, buf[:])
		if n == 0 {
			break
		}
	}

	logger.Println("end", addr)
}
