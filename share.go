package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type shareDir struct {
	dir string
}
type dirItem struct {
	N    int
	Name string
	Typ  string
}

func (s *shareDir) GetRealPath(path1 string) string {
	return path.Join(s.dir, path1)
}

//GetType 0-not exist,1-file,2-dir
func (s *shareDir) GetType(realPath string) int {
	if strings.HasPrefix(realPath, "../") || strings.Contains(realPath, "/..") {
		return 0
	}
	stat1, err := os.Stat(realPath)
	if err != nil {
		log.Println(err)
		return 0
	}
	if stat1.IsDir() {
		return 2
	}
	return 1
}

//ListDir 必须先确定为目录
func (s *shareDir) ListDir(realPath string) []dirItem {
	dir1, _ := os.Open(realPath)
	infos, _ := dir1.Readdir(0)
	res := []dirItem{}
	var typ string
	for i, v := range infos {
		if v.IsDir() {
			typ = "D"
		} else {
			typ = "F"
		}
		res = append(res, dirItem{i + 1, v.Name(), typ})
	}
	return res
}

var share *shareDir

const shareBase = "/share/"
const lenShareBase = len(shareBase)

func setShareDir(dir string) {
	http.HandleFunc(shareBase, handleShare)
	share = &shareDir{dir}
}

func handleShare(w http.ResponseWriter, r *http.Request) {
	realPath := share.GetRealPath(r.URL.Path[lenShareBase:])
	switch share.GetType(realPath) {
	case 0:
		w.WriteHeader(404)
	case 1:
		http.ServeFile(w, r, realPath)
	case 2:
		// buf := bytes.NewBufferString(r.URL.Path + "\n")
		list1 := share.ListDir(realPath)
		// for _, v := range list1 {
		// 	buf.WriteString(fmt.Sprintf("%d . %s  %v\n", v.N, v.Typ, v.Name))
		// }
		t := template.New("")
		t.Parse(tmplDir)
		data := make(map[string]interface{})
		data["title"] = r.URL.Path
		data["links"] = list1
		w.WriteHeader(200)
		t.Execute(w, data)
	}

}
