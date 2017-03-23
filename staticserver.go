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
	"io/ioutil"
	"encoding/json"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}
type home struct {
	Title string
}
// type ListFiles struct {
// 	Name string `json:"name"`
// 	Size string `json:"size"`
// }

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
	mux["/list"] = filelist
	mux["/files"] = listindex
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

func listindex(w http.ResponseWriter, r *http.Request) {
	title := home{Title: "文件列表"}
	t, _ := template.ParseFiles(Template_Dir + "list.html")
	t.Execute(w, title)
}

func filelist(w http.ResponseWriter, r *http.Request) {
	// lm := make([]ListFiles, 0)
	// arr := []string{"hello", "apple", "python", "golang", "pear"}
	dir, _ := os.Getwd()
	arr, _ := ListDir(dir+"/upload", ".go")
	data, err := json.Marshal(arr)
	if err == nil && data != nil {
		fmt.Fprintf(w, string(data))
		return
	}
	fmt.Fprintf(w, "%v", "暂无内容")
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


//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
 	files = make([]string, 0, 10)
 	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
	  	return nil, err
	}
 	// PthSep := string(os.PathSeparator)
 	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
	   	if fi.IsDir() { // 忽略目录
	   		continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //忽略匹配文件
		   	continue
		}
  		// files = append(files, dirPth+PthSep+fi.Name())
  		files = append(files, fi.Name())
 	}
 	return files, nil
}