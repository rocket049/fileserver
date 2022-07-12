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

func runDiscover(ips []string) *discover.DiscoverClient {
	server := discover.NewServer()
	for _, ip := range ips {
		server.Append("http", ip, 6868, "index", "FileServer", "Share Files")
	}

	go server.Serve(true)
	time.Sleep(time.Millisecond * 100)
	client := discover.NewClient()
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

func showAddr() []string {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	res := []string{}
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
				//png := filepath.Join(dir1, fmt.Sprintf("fileserver-%d-%d.png", n, i))
				png := filepath.Join(dir1, fmt.Sprintf("<%d-%d>%s.png", n, i, vs[0]))
				fmt.Println(png)
				var addr string
				if strings.Contains(vs[0], ":") {
					addr = fmt.Sprintf("http://[%s]:6868/index", vs[0])
					res = append(res, fmt.Sprintf("[%s]", vs[0]))
				} else {
					addr = fmt.Sprintf("http://%s:6868/index", vs[0])
					res = append(res, vs[0])
				}
				fmt.Printf("Access URL: %s\n", addr)
				qrcode.WriteFile(addr, qrcode.Highest, 400, png)
				showPng(png, addr)
			}
		}
	}
	return res
}

func showPng(fn, title string) {
	path1, err := exec.LookPath("showImg")
	if err != nil {
		open.Start(fn)
	} else {
		sync1 := exec.Command("sync", fn)
		sync1.Run()
		sync1.Wait()
		//fyne版showImg太快，二维码文件来不及同步到硬盘，因此打开时间延长100毫秒
		time.Sleep(time.Millisecond * 100)
		cmd := exec.Command(path1, "-t", title, "-f", fn)
		cmd.Start()
		go cmd.Wait()
	}
}

func main() {
	var share = flag.String("share", ".", "Share files in this DIR")
	var upload = flag.String("upload", ".", "Upload files to this DIR")
	var browse = flag.Bool("list", false, "List server in this lan.")
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
	ips := showAddr()

	runDiscover(ips)

	go http.ListenAndServe(":6868", nil)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Exit")
}
