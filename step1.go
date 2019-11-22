package main

import (
	"github.com/Unknwon/goconfig"
	"goexcel/service"
	"log"
	"os"
	"strconv"
	"strings"
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
	allusermap :=service.ReadInExcel(configmap)
	tempfile := "config.ini"
	//打卡并清空文件
	file,err:=os.OpenFile(tempfile,os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err!=nil && os.IsNotExist(err){
		file,err = os.Create(tempfile)
		defer file.Close()
	}
	cfg, _ := goconfig.LoadConfigFile(tempfile)


	cfg.SetValue("da", "start", configmap["start"].(string))
	cfg.SetValue("da", "end", configmap["end"].(string))
	cfg.SetValue("folder", "input", configmap["folder"].(map[string]string)["input"])
	cfg.SetValue("folder", "output", configmap["folder"].(map[string]string)["output"])
	if _,ok :=configmap["changhuaexcelfile"];ok{
		cfg.SetValue("excelfile", "changhua", configmap["changhuaexcelfile"].(string))
	}
	if _,ok :=configmap["pudongexcelfile"];ok{
		cfg.SetValue("excelfile", "pudong", configmap["pudongexcelfile"].(string))
	}
	if _,ok :=configmap["maomingexcelfile"];ok{
		cfg.SetValue("excelfile", "maoming", configmap["maomingexcelfile"].(string))
	}
	if _,ok :=configmap["dingdingexcelfile"];ok{
		cfg.SetValue("excelfile", "dingding", configmap["dingdingexcelfile"].(string))
	}



	cfg.SetValue("exportexcelfile", "firstexportexcelfile", configmap["firstexportexcelfile"].(string))



	cfg.SetValue("regulation", "perdayhour", configmap["perdayhour"].(string))



	for k,v :=range configmap["special_regulation"].(map[string]string) {
		cfg.SetValue("special_regulation", k, v)
	}


	cfg.SetValue("important", "filename", configmap["importantfilename"].(string))



	userinfo :=configmap["userinfo"].(map[int64]string)
	userinfoflip :=make(map[string]string)
	gouserorderslice :=configmap["gouserorderslice"].([]int64)
	userstr := ""
	for _,v :=range gouserorderslice{
		gk := strconv.FormatInt(v,10)
		userstr+=gk+","
		cfg.SetValue("userorder", gk, userinfo[v])
		gv :=strings.Replace(userinfo[v], " ", "", -1)
		userinfoflip[gv] = gk
	}


	userstr = userstr[:len(userstr)-1]
	cfg.SetValue("gouserorder", "order", userstr)



	for k,v :=range allusermap {
		// 对已有的键进行值重写操作，返回值为 bool 类型，表示是否为插入操作
		for _,value :=range v {
			// 对已有的键进行值重写操作，返回值为 bool 类型，表示是否为插入操作

			if _,ok :=userinfoflip[value.Username];ok{
				cfg.SetValue(k, value.Userid, userinfoflip[value.Username])
			}else{
				cfg.SetValue(k, value.Userid, value.Username)
			}

		}
	}

	err = goconfig.SaveConfigFile(cfg, tempfile)
	if err != nil {
		log.Fatalf("无法保存配置文件：%s", err)
	}

}
func init() {
	file := "./log/" +"step0log"+ ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[goexcel]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

