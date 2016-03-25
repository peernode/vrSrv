package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func initPathRouter(router *httprouter.Router) {
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.POST("/vr/upload", Upload)
	router.GET("/vr/getList", GetList)
	router.GET("/vr/getUploadList", GetUploadList)

	rootPath := "resource"
	router.ServeFiles("/vr/static/*filepath", http.Dir(rootPath), gLogger)
	router.ServeFiles("/vr/static2/*filepath", http.Dir("data_out"), gLogger)
	//router.ServeFiles("/log2/video/*filepath", http.Dir(rootPath))
	//router.ServeFiles("/src/*filepath", http.Dir("public"))
}
