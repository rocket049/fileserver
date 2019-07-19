package main

import (
	"io"
	"os/exec"
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
	ctx.StatusCode(200)
	writer := ctx.ResponseWriter()
	for {
		n, _ := io.Copy(writer, stdout)
		if n == 0 {
			break
		}
	}
	logger.Println("end", addr)
}
