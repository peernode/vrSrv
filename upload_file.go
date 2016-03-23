package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"time"
)

/*
<form enctype="multipart/form-data" action="http://127.0.0.1:9090/upload" method="post">
  <input type="file" name="uploadfile" />
  <input type="hidden" name="token" value="{{.}}"/>
  <input type="submit" value="upload" />
</form>
*/

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("methond: ", r.Method)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		fmt.Println("curtime: ", curtime)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
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

		fmt.Println("upload request, id: ", usrName, " datum: ", datum)

		fmt.Fprintf(w, "%v", handler.Header)
		fileName := fmt.Sprintf("%s_%s_%s", datum, usrName, handler.Filename)
		f, err := os.OpenFile("./resource/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
	}
}
