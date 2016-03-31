package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type VideoInfo struct {
	Title    string
	Desc     string
	ImageUrl string
	VideoUrl string
}

type VideoInfoList []VideoInfo

type VideoInfoResp struct {
	Result   string
	InfoList VideoInfoList
}

func GetList2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, _ := os.Open("program2.json")
	js, _ := ioutil.ReadAll(file)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}

func GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	result := "ok"
	videotype := ""
	page := 0
	pagesize := 0
	if len(r.Form["type"]) > 0 {
		videotype = r.Form["type"][0]
	}
	if len(r.Form["page"]) > 0 {
		i, err := strconv.Atoi(r.Form["page"][0])
		if err == nil {
			page = i
		}
	}
	if len(r.Form["pagesize"]) > 0 {
		i, err := strconv.Atoi(r.Form["pagesize"][0])
		if err == nil {
			pagesize = i
		}
	}

	gLogger.Info("httpreq|GetList|vtype: %s, page: %d, pagesize %d", videotype, page, pagesize)
	fmt.Println("Host: ", r.Host)

	infoList := make(VideoInfoList, 0)
	var resp VideoInfoResp

	if pagesize == 0 || page == 0 {
		result = "page info error"
	}

	if result != "ok" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp.Result = result
		resp.InfoList = infoList
		js, _ := json.Marshal(resp)
		w.Write(js)
		return
	}

	infos, _ := getFileInfo("./data_out/", true)
	sort.Sort(infos)

	for i := 0; i < pagesize; i++ {
		imageUrl := fmt.Sprintf("http://%s/vr/static/vrimage/vrtest1.jpg", r.Host)
		videoUrl := fmt.Sprintf("http://%s/vr/static/video/test1.mp4", r.Host)

		if len(infos) > 0 {
			imageUrl = fmt.Sprintf("http://%s/vr/static2/%s.jpg", r.Host, infos[i%len(infos)].Name)
			videoUrl = fmt.Sprintf("http://%s/vr/static2/%s", r.Host, infos[i%len(infos)].Name)
		}

		videoTitle := fmt.Sprintf("test中文%d", i)
		videsDesc := fmt.Sprintf("test中文%d", i)

		infoItem := VideoInfo{Title: videoTitle, Desc: videsDesc, ImageUrl: imageUrl, VideoUrl: videoUrl}
		infoList = append(infoList, infoItem)
	}

	fmt.Println(infoList)

	resp.Result = result
	resp.InfoList = infoList
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}

type fileInfo struct {
	Name    string
	Size    int64
	ModTime string
	Utc     int64
}

type FileInfos []fileInfo

func (s FileInfos) Len() int {
	return len(s)
}

func (s FileInfos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FileInfos) Less(i, j int) bool {
	return s[i].Utc > s[j].Utc
}

func getFileInfo(root string, mp4Only bool) (FileInfos, error) {
	infos := make(FileInfos, 0)
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if mp4Only && !strings.HasSuffix(f.Name(), ".mp4") {
			return nil
		}

		if !strings.HasPrefix(f.Name(), ".") {
			var item fileInfo
			item.Name = f.Name()
			item.Size = f.Size()
			item.ModTime = f.ModTime().String()
			item.Utc = f.ModTime().Unix()
			infos = append(infos, item)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return infos, nil
}

func GetUploadList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	infos, _ := getFileInfo("./data/", false)
	sort.Sort(infos)
	fmt.Println(infos)
	gLogger.Info("httpreq|GetUploadList| %d items", len(infos))

	/*
	   	const tpl = `
	   <!DOCTYPE html>
	   <html>
	   	<head>
	   		<meta charset="UTF-8">
	   		<title>上传列表</title>
	   	</head>
	   	<body>
	   		<table border="1">
	   		  <tr>
	   		    <th>文件名</th>
	   		    <th>大小</th>
	   		    <th>创建时间</th>
	   		  </tr>

	   		  {{range .}}
	   		  <tr>
	   		    <td>{{.Name}}</td>
	   		    <td>{{.Size}}</td>
	   		    <td>{{.ModTime}}</td>
	   		  </tr>
	   		  {{end}}
	   		</table>
	   	</body>
	   </html>`

	   	check := func(err error) {
	   		if err != nil {
	   			fmt.Println("err3: ", err)
	   		}
	   	}
	   	t, err := template.New("webpage").Parse(tpl)
	   	check(err)

	   	err = t.Execute(w, infos)
	   	check(err)

	   	return
	*/

	t, err := template.ParseFiles("uploadList.html")
	if err != nil {
		fmt.Println("template parse error: ", err)
		gLogger.Info("get upload list error: %s", err.Error())
		fmt.Fprintf(w, "get upload list error: %s", err.Error())
		return
	}
	err = t.Execute(w, infos)
	if err != nil {
		gLogger.Info("get upload list error2: %s", err.Error())
		fmt.Fprintf(w, "get upload list error2: %s", err.Error())
		return
	}

	//fmt.Fprintf(w, "get infos ok")
}
