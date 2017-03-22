package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	// "path/filepath"
	"regexp"
	// "strconv"
	"time"
	"strings"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}
type home struct {
	Title string
}

const (
	Template_Dir = "./view/"
	Upload_Dir   = "./upload/"
)

func main() {
	server := http.Server{
		Addr:        ":8000",
		Handler:     &Myhandler{},
		ReadTimeout: 10 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index
	mux["/upload"] = upload
	mux["/file"] = StaticServer
	fmt.Println("服务已启动：127.0.0.1:8000")
	server.ListenAndServe()
}

func (*Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	if ok, _ := regexp.MatchString("/css/", r.URL.String()); ok {
		http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))).ServeHTTP(w, r)
	} else {
		http.StripPrefix("/", http.FileServer(http.Dir("./upload/"))).ServeHTTP(w, r)
	}

}

func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t, _ := template.ParseFiles(Template_Dir + "upfile.html")
		t.Execute(w, "上传文件")
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Fprintf(w, "%v", "上传错误  " + err.Error())
			return
		}
		// fileext := filepath.Ext(handler.Filename)
		// if check(fileext) == false {
		// 	fmt.Fprintf(w, "%v", "不允许的上传类型")
		// 	return
		// }
		// filename := strconv.FormatInt(time.Now().Unix(), 10) + fileext
		filename := handler.Filename

		//check upload dir exists or not
		_, err = PathExists(strings.Trim(Upload_Dir,"."))
		if err != nil {
			fmt.Fprintf(w, "%v", "创建目录失败  "+ err.Error())
			return
		}

		//check file exists or not
		fi, err := os.Open(Upload_Dir+filename)
		if err != nil && os.IsNotExist(err) {
			f, _ := os.OpenFile(Upload_Dir+filename, os.O_CREATE|os.O_WRONLY, 0660)
			_, err = io.Copy(f, file)
			if err != nil {
				fmt.Fprintf(w, "%v", "上传失败  "+ err.Error())
				return
			}
			http.Redirect(w, r, "/file", http.StatusFound)
		 }
		defer fi.Close()
		fmt.Fprintf(w, "%v", "文件已经存在，请重新上传")
		
		// filedir, _ := filepath.Abs(Upload_Dir + filename)
		// fmt.Fprintf(w, "%v", filename+"上传完成,服务器地址:"+filedir)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	title := home{Title: "首页"}
	t, _ := template.ParseFiles(Template_Dir + "index.html")
	t.Execute(w, title)
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/file", http.FileServer(http.Dir("./upload/"))).ServeHTTP(w, r)
}

func check(name string) bool {
	ext := []string{".exe", ".js", ".png"}

	for _, v := range ext {
		if v == name {
			return false
		}
	}
	return true
}

func PathExists(path string) (bool, error) {
	dir, _ := os.Getwd()  //当前的目录
	_, err := os.Stat(dir+path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {

		err := os.Mkdir(dir+path, os.ModePerm)
		if err != nil {
		  // fmt.Println(err)
		  return false, err
		 }
		// fmt.Println("创建目录" + dir + path + "成功")
		return false, nil
	}
	return false, err
}
