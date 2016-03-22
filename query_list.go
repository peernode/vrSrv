package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type VideoInfo struct {
	Title    string
	Desc     string
	ImageUrl string
	VideoUrl string
}

type VideoInfoList []VideoInfo

func GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
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

	fmt.Println("vtype", videotype, " page ", page, " pagesize ", pagesize)

	infoList := make(VideoInfoList, 0)

	if pagesize == 0 || page == 0 {
		w.Header().Set("Content-Type", "application/json")
		js, _ := json.Marshal(infoList)
		w.Write(js)
		return
	}

	sip := "127.0.0.1:8080"
	for i := 0; i < pagesize; i++ {
		imageUrl := fmt.Sprintf("http://%s/vr/static/vrimage/vrtest1.jpg", sip)
		videoUrl := fmt.Sprintf("http://%s/vr/static/video/test1.mp4", sip)
		infoItem := VideoInfo{Title: "test1", Desc: "test1", ImageUrl: imageUrl, VideoUrl: videoUrl}
		infoList = append(infoList, infoItem)
	}

	js, err := json.Marshal(infoList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
