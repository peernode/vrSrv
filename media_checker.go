package main

import (
	"encoding/json"
	"os"
	"fmt"
	"io"
	"time"
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

func unserializeMediaInfo(){
	file, err := os.Open("./conf/media.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}
	defer file.Close()

	medias.mu.RLock()
	defer medias.mu.RUnlock()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&medias.info)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
	fmt.Println(medias.info)
}

func serializeMediaInfo(){
	file, err := os.Create("./conf/media.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}
	defer file.Close()

	medias.mu.Lock()
	defer medias.mu.Unlock()
	decoder := json.NewEncoder(file)
	err = decoder.Encode(&medias.info)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
}

func checkFileStatus(){
	for{
		lost := false
		medias.mu.Lock()
		for k, v:= range medias.info{
			for i:=0; i<len(v); i++{
				if !Exist(v[i].ImgUrl) || !Exist(v[i].VideoUrl){
					v = deleteUseAppend(v, i)
					i--
					lost = true
				}
			}
			if lost{
				medias.info[k]=v
			}
		}
		medias.mu.Unlock()

		if lost{
			gLogger.Info("lose some media file")
			serializeMediaInfo()
		}

		time.Sleep(time.Second*60)
	}
}