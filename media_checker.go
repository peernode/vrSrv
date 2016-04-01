package main

import (
	"encoding/json"
	"os"
	"fmt"
	"time"
)

func deleteUseAppend(s []MediaInfo, i int) {
	s = append(s[:i], s[i+1:]...)
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
	file, err := os.Open("./conf/media.json")
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
	fmt.Println(medias)
}

func checkFileStatus(){
	for{
		lost := false
		for _, v:= range medias{
			for i, info := range v{
				if !Exist(info.ImgUrl) || !Exist(info.VideoUrl){
					deleteUseAppend(v, i)
					lost = true
				}
			}
		}

		if lost{
			fmt.Println("lost sth")
			serializeMediaInfo()
		}

		time.Sleep(time.Second*60)
	}
}