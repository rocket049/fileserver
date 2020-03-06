package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/skratchdot/open-golang/open"

	qrcode "github.com/skip2/go-qrcode"
)

func showAddr() {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for n, if1 := range ifs {
		addrs, err := if1.Addrs()
		if err != nil {
			panic(err)
		}
		for i, addr := range addrs {
			if strings.HasPrefix(addr.String(), "127.") {
				continue
			} else if strings.Contains(addr.String(), ":") {
				continue
			} else {
				vs := strings.Split(addr.String(), "/")

				dir1, err := os.UserCacheDir()
				if err != nil {
					panic(err)
				}
				png := filepath.Join(dir1, fmt.Sprintf("fileserver-%d-%d.png", n, i))
				fmt.Println(png)
				var addr string
				if strings.Contains(vs[0], ":") {
					addr = fmt.Sprintf("http://[%s]:6868/index", vs[0])
				} else {
					addr = fmt.Sprintf("http://%s:6868/index", vs[0])
				}
				fmt.Printf("Access URL: %s\n", addr)
				qrcode.WriteFile(addr, qrcode.Highest, 400, png)
				showPng(png, addr)
			}
		}
	}
}

func showPng(fn, title string) {
	path1, err := exec.LookPath("showImg")
	if err != nil {
		open.Start(fn)
	} else {
		cmd := exec.Command(path1, "-t", title, "-f", fn)
		cmd.Start()
	}
}

func main() {
	var share = flag.String("share", ".", "Share files in this DIR")
	var upload = flag.String("upload", ".", "Upload files to this DIR")
	flag.Parse()
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(mainPage))
	})
	setShareDir(*share)
	setUploadDir(*upload)
	showAddr()
	http.ListenAndServe(":6868", nil)
}
