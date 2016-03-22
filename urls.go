package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func initPathRouter(router *httprouter.Router) {
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.PUT("/vr/upload", Upload)
	router.GET("/vr/getList", GetList)

	rootPath := "resource"
	router.ServeFiles("/vr/static/*filepath", http.Dir(rootPath))
	//router.ServeFiles("/log2/video/*filepath", http.Dir(rootPath))
	//router.ServeFiles("/src/*filepath", http.Dir("public"))
}
