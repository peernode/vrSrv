package main

import (
	"encoding/json"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"time"
	"runtime"
	"os/exec"
)

type Configuration struct {
	HtmlDir    string
    UploadDir    string
    MediaDir   string
	ConvertDir string
	GearDir	     string
	UploadTitle string
	UploadDesc  string
}


type UploadInfo struct{
	videoType string
	videoName string
	outName  string
}

var goos string
var gUploadFileCh = make(chan UploadInfo, 500)
var logFilename = "srvLog.txt"
var gLogger l4g.Logger
var configuration Configuration

var gMedias = NewMediaInfo("./conf/media.json")

func init(){
	goos = runtime.GOOS
}

func initConfig(){
	file, err := os.Open("./conf/conf.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		configuration = Configuration{HtmlDir:"html", UploadDir:"data", MediaDir:"data_out", GearDir:"gearvr"}
		fmt.Println("error:", err)
	}

	fmt.Println(configuration)
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
	router.GET("/vr/getList2", GetList2)  //获取gearvr的上传列表
	router.GET("/vr/getList3", GetList3)  //获取可用媒体列表
	router.GET("/vr/getUploadList", GetUploadList)  //获取上传列表

	router.ServeFiles("/vr/static/*filepath", http.Dir(configuration.HtmlDir))  //下载相应的静态html文件
	router.ServeFiles("/vr/static2/*filepath", http.Dir(configuration.MediaDir))  //下载相应的媒体文件
	router.ServeFiles("/vr/gearvr/*filepath", http.Dir(configuration.GearDir))

	return router
}

func ffmpegTransfer() {
	for file := range gUploadFileCh {
		gLogger.Info("transfer, file: %s", file.videoName)
		now := time.Now()

		var err error
		if goos == "darwin"{
			fmt.Println("darwin transfer...")
			cmd := exec.Command("/bin/bash", "test.sh", file.videoName)
			_, err = cmd.Output()
		}else{
			_, err = CopyFile(file.outName, file.videoName)
			_, err = CopyFile(file.outName+".jpg", "media/vrtest.jpg")
		}

		cost := time.Since(now)
		if err != nil {
			gLogger.Info("transfer error: %s %s, cost: %d", file.videoName, err.Error(), cost)
		} else {
			gLogger.Info("transfer success: %s, cost: %d", file.videoName, cost)

			gMedias.Add(file.videoType, time.Now().String(), configuration.UploadTitle, configuration.UploadDesc, fmt.Sprintf("%s.jpg", file.outName), file.outName)
		}
	}
}

func main() {
	initConfig()
	initLogger()
	go ffmpegTransfer()
	router := initHttpRouter()

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
	}
	log.Fatal(s.ListenAndServe())
}
