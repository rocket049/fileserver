package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
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
	for _, if1 := range ifs {
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
				fmt.Printf("Access URL: http://%s:6868/index\n", vs[0])
				dir1, err := os.UserCacheDir()
				if err != nil {
					panic(err)
				}
				png := filepath.Join(dir1, fmt.Sprintf("fileserver-%d.png", i))
				qrcode.WriteFile(fmt.Sprintf("http://%s:6868/index", vs[0]), qrcode.Highest, 400, png)
				open.Start(png)
			}
		}
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
