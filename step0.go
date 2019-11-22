package main

import (
	"goexcel/service"
	"log"
	"os"
)

func main(){
	configmap :=service.InitConfig("tempconfig.ini")
	//先创建文件夹
	foldermap :=configmap["folder"].(map[string]string)
	for _,v :=range foldermap{
		_dir := "./"+v
		exist, err := service.PathExists(_dir)
		if err != nil {
			log.Println("get dir error!", err)
			return
		}

		if !exist {
			log.Println("no dir", _dir)
			// 创建文件夹
			err := os.MkdirAll(_dir, os.ModePerm)
			if err != nil {
				log.Println("mkdir failed!", err)
			} else {
				log.Println("mkdir success!")
			}
		}
	}

	_dir := "./log"
	exist, err := service.PathExists(_dir)
	if err != nil {
		log.Println("get dir error!", err)
		return
	}

	if !exist {
		log.Println("no dir", _dir)
		// 创建文件夹
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			log.Println("mkdir failed!", err)
		} else {
			log.Println("mkdir success!")
		}
	}


}


