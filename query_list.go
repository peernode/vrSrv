package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type VideoInfo struct {
	Title    string
	Desc     string
	ImageUrl string
	VideoUrl string
}

type VideoInfoList []VideoInfo

func GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	infoList := make(VideoInfoList, 0)
	for i := 0; i < 5; i++ {
		infoItem := VideoInfo{Title: "test1", Desc: "test1", ImageUrl: "http://127.0.0.1/", VideoUrl: "http://127.0.0.1/"}
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
