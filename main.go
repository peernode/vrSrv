package main

import (
	"encoding/json"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Configuration struct {
	htmlDir    string
    uploadDir    string
    convertDir   string
	gearDir	     string
}

type MediaInfo struct{
	datum string
	title  string
	desc	string
	imgUrl	string
	videoUrl	string
}
type MediaInfos struct{
	yuanchuang []MediaInfo
	meinv	   []MediaInfo
	lvyou	   []MediaInfo
	kongbu	   []MediaInfo
}

var gUploadFileCh = make(chan string, 500)
var logFilename = "srvLog.txt"
var gLogger l4g.Logger
var configuration Configuration
var medias MediaInfos

func initConfig(){
	file, err := os.Open("./conf/conf.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		configuration = Configuration{htmlDir:"html", uploadDir:"data", convertDir:"data_out", gearDir:"gearvr"}
		fmt.Println("error:", err)
	}
}

func initMedias(){
	file, err := os.Open("./conf/media.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&medias)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
	fmt.Println(medias)
}

//init for logger
func initLogger() {
	gLogger = make(l4g.Logger)

	gLogger.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
	flw := l4g.NewFileLogWriter(logFilename, true)
	flw.SetFormat("[%D %T] [%L] %M")
	flw.SetRotateDaily(true)
	gLogger.AddFilter("logfile", l4g.FINEST, flw)

	gLogger.Info("Init logger! The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}

func initHttpRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.POST("/vr/upload", Upload)
	router.GET("/vr/getList", GetList)  //按类型获取媒体列表
	router.GET("/vr/getList2", GetList2)  //获取上传列表
	router.GET("/vr/getUploadList", GetUploadList)  //获取gearvr的上传列表

	router.ServeFiles("/vr/static/*filepath", http.Dir(configuration.htmlDir))  //下载相应的静态html文件
	router.ServeFiles("/vr/static2/*filepath", http.Dir(configuration.convertDir))  //下载上传的媒体文件
	router.ServeFiles("/vr/gearvr/*filepath", http.Dir(configuration.gearDir))

	return router
}

func ffmpegTransfer() {
	for file := range gUploadFileCh {
		gLogger.Info("transfer, file: %s", file)
		now := time.Now()
		cmd := exec.Command("/bin/bash", "test.sh", file)
		bytes, err := cmd.Output()
		cost := time.Since(now)
		if err != nil {
			gLogger.Info("transfer error: %s %s, cost: %d", file, err.Error(), cost)
		} else {
			gLogger.Info("transfer success: %s %s, cost: %d", file, string(bytes), cost)
		}
	}
}

func main() {
	initConfig()
	initMedias()
	initLogger()
	go ffmpegTransfer()
	router := initHttpRouter()

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
	}
	log.Fatal(s.ListenAndServe())
}
