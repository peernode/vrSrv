package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"time"
	"strings"
)

/*
<form enctype="multipart/form-data" action="http://127.0.0.1:9090/upload" method="post">
  <input type="file" name="uploadfile" />
  <input type="hidden" name="token" value="{{.}}"/>
  <input type="submit" value="upload" />
</form>
*/

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	gLogger.Info("httpreq|Upload|method: %s", r.Method)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		fmt.Println("curtime: ", curtime)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			gLogger.Info("upload fail: %s", err.Error())
			fmt.Fprintf(w, "上传失败， %s", err.Error())
			return
		}
		defer file.Close()

		usrName := ""
		id := r.MultipartForm.Value["id"]
		if len(id) > 0 {
			usrName = id[0]
		}
		datum := "20150101"
		datumValue := r.MultipartForm.Value["datum"]
		if len(datumValue) > 0 {
			datum = datumValue[0]
		}

		uploadFile := fmt.Sprintf("%s_%s_%s", datum, usrName, handler.Filename)
		fileName := fmt.Sprintf("%s/%s", configuration.UploadDir, uploadFile)  //相对路径
		outfileName := fmt.Sprintf("%s/%s_%s_%s_fin.mp4", configuration.ConvertDir, datum, usrName, handler.Filename)
		f, err := os.Create(fileName)
		if err != nil {
			gLogger.Info("upload fail, id: %s, name: %s, err: %s", usrName, handler.Filename, err.Error())
			fmt.Fprintf(w, "%s 上传失败， %s", handler.Filename, err.Error())
			return
		}
		//defer f.Close()

		io.Copy(f, file)
		fmt.Fprintf(w, "%s 上传成功！", handler.Filename)
		f.Close()


		if strings.HasSuffix(fileName, "mp4") || strings.HasSuffix(fileName, "MP4") || strings.HasSuffix(fileName, "mov") || strings.HasSuffix(fileName, "MOV"){
			gUploadFileCh <- UploadInfo{videoType:"YuanChuang", videoName:uploadFile, outName:outfileName}
		}

		gLogger.Info("upload success, id: %s, name: %s", usrName, handler.Filename)
	}
}
