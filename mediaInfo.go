package main

import (
	"encoding/json"
	"os"
	"fmt"
	"io"
	"time"
	"sync"
	"log"
)

func deleteUseAppend(s []MediaInfo, i int) []MediaInfo{
	s = append(s[:i], s[i+1:]...)
	return s
}

func Exist(fileName string) bool{
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}

func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}


type MediaInfo struct{
	Datum string
	Title  string
	Desc	string
	ImgUrl	string
	VideoUrl	string
}

type MediaInfos struct{
	mu    sync.RWMutex
	info map[string][]MediaInfo
}

func NewMediaInfo(filename string) *MediaInfos{
	s := &MediaInfos{info: make(map[string][]MediaInfo)}

	s.load(filename)

	go s.checkFileStatus()  //启动协程检查文件是否丢失

	return s
}

func (s *MediaInfos)load(filename string){
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening mediajson:", err)
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&s.info)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
	fmt.Println(s.info)
}

func (s *MediaInfos)save(filename string){
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Error opening mediajson:", err)
	}

	decoder := json.NewEncoder(f)
	err = decoder.Encode(&s.info)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
}

func (s *MediaInfos) Add(videoType, Datum, Title, Desc, imgUrl, videoUrl string){
	newItem := MediaInfo{Datum: time.Now().String(), Title: configuration.UploadTitle, Desc: configuration.UploadDesc, ImgUrl: imgUrl, VideoUrl: videoUrl}

	s.mu.Lock()
	s.info[videoType]=append(s.info[videoType], newItem)
	s.mu.Unlock()

	s.save("./conf/media.json")
}

// resp for getList
type VideoInfo struct {
	Title    string
	Desc     string
	ImageUrl string
	VideoUrl string
}

type VideoInfoList []*VideoInfo

type VideoInfoResp struct {
	Result   string
	InfoList VideoInfoList
}

func (s *MediaInfos) Get(pagesize, page int, videoType, host string) ([]byte, error){
	result := "ok"
	if pagesize == 0 || page == 0 {
		result = "page info error"
	}

	s.mu.RLock()
	mediaInfos, exist := s.info[videoType]
	s.mu.RUnlock()
	if !exist{
		result = "videoType info error"
	}

	if result !="ok"{
		var r = VideoInfoResp{Result:result, InfoList: make(VideoInfoList, 0)}
		js, _ := json.Marshal(r)
		return js, nil
	}

	infoList := make(VideoInfoList, 0)
	for i :=len(mediaInfos)-1; i>=0; i--{
		if len(infoList) >= pagesize{
			break
		}
		imageUrl := fmt.Sprintf("http://%s/vr/static2/%s", host, mediaInfos[i].ImgUrl)
		videoUrl := fmt.Sprintf("http://%s/vr/static2/%s", host, mediaInfos[i].VideoUrl)
		videoTitle := mediaInfos[i].Title
		videoDesc := mediaInfos[i].Desc
		infoList = append(infoList, &VideoInfo{Title: videoTitle, Desc: videoDesc, ImageUrl: imageUrl, VideoUrl: videoUrl})
	}

	var resp = VideoInfoResp{Result: "ok", InfoList: infoList}
	js, err := json.Marshal(resp)
	return js, err
}

func (s *MediaInfos)checkFileStatus(){
	for{
		lost := false
		s.mu.Lock()
		for k, v:= range s.info{
			for i:=0; i<len(v); i++{
				if !Exist(v[i].ImgUrl) || !Exist(v[i].VideoUrl){
					gLogger.Info("lose file: %s", v[i].VideoUrl)
					v = deleteUseAppend(v, i)
					i--
					lost = true
				}
			}
			if lost{
				s.info[k]=v
			}
		}
		s.mu.Unlock()

		if lost{
			gLogger.Info("lose some media file")
			s.save("./conf/media.json")
		}

		time.Sleep(time.Second*60)
	}
}