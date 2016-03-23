package main

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var logFilename = "srvLog.txt"
var GLogger l4g.Logger

//init for logger
func initLogger() {
	GLogger = make(l4g.Logger)

	GLogger.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())

	flw := l4g.NewFileLogWriter(logFilename, true)
	flw.SetFormat("[%D %T] [%L] %M")
	GLogger.AddFilter("logfile", l4g.FINEST, flw)
	flw.SetRotateDaily(true)

	//GLogger.AddFilter("logfile", l4g.FINEST, l4g.NewFileLogWriter(logFilename, true))
	GLogger.Info("Init logger! The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}

var gUploadFileCh = make(chan string, 500)

func ffmpegTransfer() {
	for file := range gUploadFileCh {
		GLogger.Info("transfer, file: %s", file)

		curtime := time.Now().Unix()
		cmd := exec.Command("/bin/bash", "test.sh", file)
		//cmd := exec.Command("/bin/bash", "t.sh")

		bytes, err := cmd.Output()
		cost := time.Now().Unix() - curtime
		if err != nil {
			GLogger.Info("transfer error: %s %s, cost: %d", file, err.Error(), cost)
		} else {
			GLogger.Info("transfer success: %s %s, cost: %d", file, string(bytes), cost)
		}
	}

}

func main() {
	initLogger()

	go ffmpegTransfer()

	router := httprouter.New()
	initPathRouter(router)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1000 * time.Second,
		WriteTimeout: 500 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
