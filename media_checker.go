package main

import (
	"encoding/json"
	"os"
	"fmt"
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

func unserializeMediaInfo(){
	file, err := os.Open("./conf/media.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&medias)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
	fmt.Println(medias)
}

func serializeMediaInfo(){
	file, err := os.Create("./conf/media.json")
	if err != nil{
		fmt.Println("err: ", err.Error())
		return
	}
	defer file.Close()

	decoder := json.NewEncoder(file)
	err = decoder.Encode(&medias)
	if err != nil{
		fmt.Println("error: ", err.Error())
	}
//	fmt.Println(medias)
}

func checkFileStatus(){
	for{
		lost := false
		for k, v:= range medias{
			for i:=0; i<len(v); i++{
				if !Exist(v[i].ImgUrl) || !Exist(v[i].VideoUrl){
					v = deleteUseAppend(v, i)
					i--
					lost = true
				}
			}
			if lost{
				medias[k]=v
			}
		}

		if lost{
			fmt.Println("lost sth")
			serializeMediaInfo()
		}

		time.Sleep(time.Second*60)
	}
}