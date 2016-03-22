package main

import (
	"fmt"
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

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./public/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
	}
}
