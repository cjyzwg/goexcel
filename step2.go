package main

import (
	"goexcel/service"
	"log"
	"os"
)


func main() {


	configmap :=service.InitConfig("config.ini")
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

	outfolder :=configmap["folder"].(map[string]string)["output"]
	firstexportexcelfile := configmap["firstexportexcelfile"].(string)
	outfile := outfolder+"/"+firstexportexcelfile
	data :=service.DealExcel(configmap)
	service.Export(data,outfile)

	importantfilename := configmap["importantfilename"].(string)
	importantfile := outfolder+"/"+importantfilename
	service.WriteResult(data,importantfile)

}

func init() {
	file := "./log/" +"firstbuildlog"+ ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[goexcel]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}





