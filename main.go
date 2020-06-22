package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/rocket049/discover-go"
	"github.com/skratchdot/open-golang/open"

	qrcode "github.com/skip2/go-qrcode"
)

func runDiscover() *discover.DiscoverClient {
	client := discover.NewClient()
	ok := client.Append("http", 6868, "index", "FileServer", "share files")
	if ok == false {
		server := discover.NewServer()
		go server.Serve(false)
		time.Sleep(time.Millisecond * 100)
		client.Append("http", 6868, "index", "FileServer", "share files")
	}
	res := client.Query()
	fmt.Println("Servers:")
	for _, v := range res {
		if v.Name == "FileServer" {
			fmt.Println(v.Href, v.Name, v.Title)
		}
	}
	return client
}

func listServers() {
	client := discover.NewClient()
	res := client.Query()
	fmt.Println("Servers:")
	for _, v := range res {
		if v.Name == "FileServer" {
			fmt.Println(v.Href, v.Name, v.Title)
		}
	}
}

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
	var browse = flag.Bool("browse", false, "List server in this lan.")
	flag.Parse()
	if *browse {
		listServers()
		return
	}
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(mainPage))
	})
	setShareDir(*share)
	setUploadDir(*upload)
	showAddr()

	c := runDiscover()
	defer c.Remove("http", 6868, "index")

	go http.ListenAndServe(":6868", nil)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Exit")
}
