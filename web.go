package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"goexcel/model"
	"goexcel/service"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
)

// 端口
const (
	HTTP_PORT  string = "8089"
)
// 目录
const (
	CSS_CLIENT_PATH   = "/css/"
	DART_CLIENT_PATH  = "/js/"
	IMAGE_CLIENT_PATH = "/image/"

	CSS_SVR_PATH   = "web"
	DART_SVR_PATH  = "web"
	IMAGE_SVR_PATH = "web"
)


// 目录
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	srcFile := "./tempconfig.ini"
	dstFile := "./config.ini"
	//删除文件
	pathflag,_:= service.PathExists(dstFile)
	if pathflag {
		ferr:=os.Remove(dstFile)
		if ferr!=nil{
			fmt.Println(ferr)
		}
	}

	total, fileerr :=service.CopyFile(dstFile,srcFile)
	if fileerr == nil {
		fmt.Println("配置文件拷贝成功，总量是：",total)
	}else{
		fmt.Printf("配置文件未完成拷贝，错误=%v\n",fileerr)
	}


	/*****这部分关键，开始请求***/
	// 先把css和脚本以及图片服务上去
	http.Handle(CSS_CLIENT_PATH, http.FileServer(http.Dir(CSS_SVR_PATH)))
	http.Handle(DART_CLIENT_PATH, http.FileServer(http.Dir(DART_SVR_PATH)))
	http.Handle(IMAGE_CLIENT_PATH, http.FileServer(http.Dir(IMAGE_SVR_PATH)))
	// 网址与处理逻辑对应起来
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/ajax", OnAjax)
	http.HandleFunc("/postdata", OnPostdata)
	http.HandleFunc("/docpostdata", OnDocdata)
	http.HandleFunc("/test", OnTest)
	http.HandleFunc("/export", OnExport)
	http.HandleFunc("/docuplexport", OnDocuplExport)

	fmt.Println("端口开启端口："+HTTP_PORT)
	//开启json文件
	service.InitWriteFile("info.json")

	// 开始服务
	//err := http.ListenAndServe(":"+HTTP_PORT, nil)
	//if err != nil {
	//	fmt.Println("服务失败 ///O ", err)
	//}

	//go func() {
	//	<-time.After(100 * time.Millisecond)
	//	err := exec.Command("explorer", "http://127.0.0.1:"+HTTP_PORT).Run()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()
	//
	//log.Println("running at port localhost:"+HTTP_PORT)
	//log.Fatal(http.ListenAndServe(":"+HTTP_PORT, nil))
	go open("http://localhost:"+HTTP_PORT)
	panic(http.ListenAndServe(":"+HTTP_PORT, nil))
	//panic(http.ListenAndServe(":"+HTTP_PORT, nil))

}

func WriteTemplateToHttpResponse(res http.ResponseWriter, t *template.Template) error {
	if t == nil || res == nil {
		return errors.New("WriteTemplateToHttpResponse: t must not be nil.")
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, nil)
	if err != nil {
		return err
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = res.Write(buf.Bytes())
	return err
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("web/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = WriteTemplateToHttpResponse(res, t)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func OnAjax(res http.ResponseWriter, req *http.Request) {
	configfile := "config.ini"
	configmap :=service.InitConfig(configfile)
	//fmt.Println(configmap)
	//mardownlist 即可生成json,数据

	jsoncategorylists, _ := json.Marshal(configmap)
	//返回的这个是给json用的，需要去掉
	res.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	io.WriteString(res, string(jsoncategorylists))
}

func OnPostdata(res http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	//fmt.Printf("%s\n", result)


	rstmap := make(map[string]interface{})
	if err := json.Unmarshal(result, &rstmap); err == nil {


		fmt.Println("==============json str 转map=======================")
		tempfile := "xxconfig.ini"
		//打卡并清空文件
		file,err:=os.OpenFile(tempfile,os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err!=nil && os.IsNotExist(err){
			file,err = os.Create(tempfile)
			defer file.Close()
		}
		cfg, _ := goconfig.LoadConfigFile(tempfile)



		cfg.SetValue("da", "start", rstmap["start"].(string))

		cfg.SetValue("da", "end", rstmap["end"].(string))
		cfg.SetSectionComments("da", "#考勤的日期(可以改动=右侧)")

		cfg.SetValue("folder", "input", rstmap["input"].(string))
		log.Println("-----1")
		cfg.SetValue("folder", "output", rstmap["output"].(string))
		log.Println("-----2")
		cfg.SetSectionComments("folder", "#考勤记录表存放位置，需要修改日期，input为考勤机的打卡记录存放位置，output为程序输出的文件夹(可以改动=右侧，不能删除)")
		cfg.SetValue("excelfile", "changhua", rstmap["changhuaexcelfile"].(string))
		log.Println("-----3")
		cfg.SetValue("excelfile", "pudong", rstmap["pudongexcelfile"].(string))
		log.Println("-----4")
		cfg.SetValue("excelfile", "maoming", rstmap["maomingexcelfile"].(string))
		log.Println("-----5")
		cfg.SetValue("excelfile", "dingding", rstmap["dingdinngexcelfile"].(string))
		log.Println("-----6")
		cfg.SetSectionComments("excelfile", "#如果没有前面就加#注释掉或者直接删掉,左边不要改动，以下为考勤的几张表，在input文件夹内(可以删除，或者改动=右侧)，最好和input文件夹下的文件同步")
		cfg.SetValue("exportexcelfile", "firstexportexcelfile", rstmap["firstexportexcelfile"].(string))
		log.Println("-----7")
		cfg.SetValue("exportexcelfile", "secondexportexcelfile", rstmap["secondexportexcelfile"].(string))
		log.Println("-----8")
		cfg.SetSectionComments("exportexcelfile", "#左边不要改动，以下为程序输出的两张表在output文件夹内(不能删除只能改动=右侧)")
		cfg.SetValue("regulation", "perdayhour", rstmap["perdayhour"].(string))
		log.Println("-----9")
		cfg.SetSectionComments("regulation", "#每天在馆十个小时(不能删除，或者改动=右侧)")

		log.Println("1")
		rstmapcheck :=rstmap["special_regulation"].([]interface{})
		for _,value :=range rstmapcheck{
			for k1,v1 :=range value.(map[string]interface{}){
				cfg.SetValue("special_regulation", k1, v1.(string))
			}
		}
		log.Println("2")
		cfg.SetSectionComments("special_regulation", "#老师每天在馆时间，如果有不同，则可以删除或者新添加像: userid=时长(可以删除或增加，或者改动=右侧)")

		cfg.SetValue("important", "filename", "a.txt")
		cfg.SetSectionComments("important", "#切勿改动以及删除")

		rstmapcheck =rstmap["userorder"].([]interface{})
		for _,value :=range rstmapcheck{
			gvalue := value.(map[string]interface{})
			cfg.SetValue("userorder", gvalue["cid"].(string), gvalue["name"].(string))
		}
		log.Println("3")
		cfg.SetSectionComments("userorder", "#用户排序,放在考勤记录顺序保持一致,并且user和现有人员保持一致(可以删除或增加，或者改动=右侧,优先处理这个)")

		cfg.SetValue("gouserorder", "order", rstmap["userorderslice"].(string))
		cfg.SetSectionComments("gouserorder", "#用户排序和上面一个相同，需要把所有的id都放在下面(可以删除或增加，或者改动=右侧)")

		rstmapcheck =rstmap["changhua"].([]interface{})
		for _,value :=range rstmapcheck{
			gvalue := value.(map[string]interface{})
			cfg.SetValue("changhua",  gvalue["id"].(string),gvalue["cid"].(string))
		}
		log.Println("4")

		rstmapcheck =rstmap["pudong"].([]interface{})
		for _,value :=range rstmapcheck{
			gvalue := value.(map[string]interface{})
			cfg.SetValue("pudong",  gvalue["id"].(string),gvalue["cid"].(string))
		}
		log.Println("5")

		rstmapcheck =rstmap["maoming"].([]interface{})
		for _,value :=range rstmapcheck{
			gvalue := value.(map[string]interface{})
			cfg.SetValue("maoming",  gvalue["id"].(string),gvalue["cid"].(string))
		}
		log.Println("6")
		rstmapcheck =rstmap["dingding"].([]interface{})
		for _,value :=range rstmapcheck{
			gvalue := value.(map[string]interface{})
			cfg.SetValue("dingding",  gvalue["id"].(string),gvalue["cid"].(string))
		}
		log.Println("7")
		cfg.SetSectionComments("changhua", "#昌化打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧,可以直接从tempconfig.ini拷贝过来)")
		cfg.SetSectionComments("pudong", "#浦东打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧)")
		cfg.SetSectionComments("maoming", "#茂名新打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧)")
		cfg.SetSectionComments("dingding", "#钉钉打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧)")

		err = goconfig.SaveConfigFile(cfg, tempfile)
		if err != nil {
			log.Fatalf("无法保存配置文件：%s", err)
		}
	}


	//for k,v :=range result{
	//	fmt.Println(k,v)
	//}
	//io.WriteString(res, "你大爷的，我收到了")
	//configmap :=service.InitConfig("testconfig.ini")
	////fmt.Println(configmap)
	////mardownlist 即可生成json,数据
	//
	//jsoncategorylists, _ := json.Marshal(configmap)
	////返回的这个是给json用的，需要去掉
	//res.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	//io.WriteString(res, string(jsoncategorylists))
}


func OnDocdata(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	fmt.Println(r.MultipartForm.Value)

	formvalue := r.MultipartForm.Value
	uploadfilename := formvalue["name"][0]
	filetype := formvalue["type"][0]

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	tempfile :="config.ini"
	cfg, _ := goconfig.LoadConfigFile(tempfile)
	//fmt.Println(cfg)

	cfg.SetValue("excelfile", filetype, uploadfilename)
	goconfig.SaveConfigFile(cfg, tempfile)
	//先生成config file
	foldermap :=service.GetConfig(tempfile,"folder")
	//先创建文件夹
	input :=foldermap["input"]

	f, err := os.OpenFile(input+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	//fmt.Println("----",handler.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	jsonstr := "{'data':{'code':'succ'}}"

	//jsonlists, _ := json.Marshal(configmap)
	//返回的这个是给json用的，需要去掉

	io.WriteString(w, jsonstr)
	//fmt.Fprintln(w, "upload ok!")
}


func OnTest(res http.ResponseWriter, req *http.Request) {

	result, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	//fmt.Printf("%s\n", result)


	rstmap := make(map[string]interface{})
	if err := json.Unmarshal(result, &rstmap); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(rstmap,reflect.TypeOf(rstmap),rstmap["skip"],reflect.TypeOf(rstmap["skip"]))
		skip := rstmap["skip"].(float64)
		fmt.Println("skip is :",skip)
		tempfile := "config.ini"
		if skip==1{
			delstring := rstmap["typename"].(string)
			//del json data
			service.DelFirstFile("tempinfo.json",delstring)
			fmt.Println("delete skip 1")
		}else if skip==2{
			delstring := rstmap["typename"].(string)
			//save config file

			//打卡并清空文件
			file,err:=os.OpenFile(tempfile,os.O_WRONLY|os.O_CREATE, 0666)
			if err!=nil && os.IsNotExist(err){
				file,err = os.Create(tempfile)
				defer file.Close()
			}
			cfg, _ := goconfig.LoadConfigFile(tempfile)
			if delstring=="da"{
				if _,ok :=rstmap["start"];ok{
					cfg.SetValue("da", "start", rstmap["start"].(string))
				}
				if _,ok :=rstmap["end"];ok{
					cfg.SetValue("da", "end", rstmap["end"].(string))
				}
			}
			if delstring=="folder"{
				if _,ok :=rstmap["input"];ok {
					input :=rstmap["input"].(string)
					//create folder
					_dir := "./"+input
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
					cfg.SetValue("folder", "input", input)
				}
				if _,ok :=rstmap["output"];ok {
					output :=rstmap["output"].(string)

					//create folder
					_dir := "./"+output
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
					cfg.SetValue("folder", "output", output)
				}
			}
			if delstring=="excelfile"{
				cfg.DeleteSection("excelfile")
				if _,ok :=rstmap["changhuaexcelfile"];ok{
					cfg.SetValue("excelfile", "changhua", rstmap["changhuaexcelfile"].(string))
				}
				if _,ok :=rstmap["pudongexcelfile"];ok{
					cfg.SetValue("excelfile", "pudong", rstmap["pudongexcelfile"].(string))
				}
				if _,ok :=rstmap["maomingexcelfile"];ok{
					cfg.SetValue("excelfile", "maoming", rstmap["maomingexcelfile"].(string))
				}
				if _,ok :=rstmap["dingdinngexcelfile"];ok{
					cfg.SetValue("excelfile", "dingding", rstmap["dingdinngexcelfile"].(string))
				}
			}
			if delstring=="exportexcelfile"{
				if _,ok :=rstmap["firstexportexcelfile"];ok{
					cfg.SetValue("exportexcelfile", "firstexportexcelfile", rstmap["firstexportexcelfile"].(string))
				}
				if _,ok :=rstmap["secondexportexcelfile"];ok{
					cfg.SetValue("exportexcelfile", "secondexportexcelfile", rstmap["secondexportexcelfile"].(string))
				}
			}
			if delstring=="regulation"{
				if _,ok :=rstmap["perdayhour"];ok{
					cfg.SetValue("regulation", "perdayhour", rstmap["perdayhour"].(string))
				}
			}

			if delstring=="userorder"{
				if _,ok :=rstmap["userorder"];ok{
					fmt.Println(rstmap["userorder"],reflect.TypeOf(rstmap["userorder"]))
				}
				rstmapcheck :=rstmap["userorder"].([]interface{})
				for _,value :=range rstmapcheck{
					gvalue := value.(map[string]interface{})
					cfg.SetValue("userorder", gvalue["cid"].(string), gvalue["name"].(string))
				}
				cfg.SetValue("gouserorder", "order", rstmap["userorderslice"].(string))
			}

			if delstring=="special_regulation"{
				rstmapcheck :=rstmap["special_regulation"].([]interface{})
				for _,value :=range rstmapcheck{
					gvalue := value.(map[string]interface{})
					cfg.SetValue("special_regulation", gvalue["cid"].(string), gvalue["hour"].(string))
				}
			}
			//configmap :=service.InitConfig(tempfile)
			//allusermap :=service.ReadInExcel(configmap)
			//for k,v :=range allusermap {
			//	// 对已有的键进行值重写操作，返回值为 bool 类型，表示是否为插入操作
			//	for _,value :=range v {
			//		// 对已有的键进行值重写操作，返回值为 bool 类型，表示是否为插入操作
			//		cfg.SetKeyComments(k, value.Userid, "#"+value.Username)
			//		if _,ok :=userinfoflip[value.Username];ok{
			//			cfg.SetValue(k, value.Userid, userinfoflip[value.Username])
			//		}else{
			//			cfg.SetValue(k, value.Userid, value.Username)
			//		}
			//
			//	}
			//}
			if delstring=="allpeople"{
				//fmt.Println("------kkk1")
				fmt.Println(rstmap["changhua"],reflect.TypeOf(rstmap["changhua"]))

				cfg.DeleteSection("changhua")
				changhuarstmapcheck :=rstmap["changhua"].([]interface{})
				for _,value :=range changhuarstmapcheck{
					gvalue := value.(map[string]interface{})
					cfg.SetValue("changhua", gvalue["id"].(string), gvalue["cid"].(string))
					cfg.SetSectionComments("changhua", "#昌化打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧,可以直接从tempconfig.ini拷贝过来)")
				}
				//fmt.Println("------kkk2",rstmap["pudong"])
				cfg.DeleteSection("pudong")
				pudongrstmapcheck :=rstmap["pudong"].([]interface{})
				for _,value :=range pudongrstmapcheck{
					gvalue := value.(map[string]interface{})
					cfg.SetValue("pudong", gvalue["id"].(string), gvalue["cid"].(string))
					cfg.SetSectionComments("pudong", "#浦东打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧)")
				}

				//fmt.Println("------kkk3",rstmap["maoming"])
				cfg.DeleteSection("maoming")
				maomingrstmapcheck :=rstmap["maoming"].([]interface{})
				for _,value :=range maomingrstmapcheck{
					gvalue := value.(map[string]interface{})
					cfg.SetValue("maoming", gvalue["id"].(string), gvalue["cid"].(string))
					cfg.SetSectionComments("maoming", "#茂名新打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧)")
				}

				//fmt.Println("------kkk4",rstmap["dingding"])
				cfg.DeleteSection("dingding")
				dingdingrstmapcheck :=rstmap["dingding"].([]interface{})
				for _,value :=range dingdingrstmapcheck{
					gvalue := value.(map[string]interface{})
					cfg.SetValue("dingding", gvalue["id"].(string), gvalue["cid"].(string))
					cfg.SetSectionComments("dingding", "#钉钉打卡机的表格中需要：工号=crm 用户id(可以删除或增加，或者改动=右侧)")
				}
			}














			err = goconfig.SaveConfigFile(cfg, tempfile)
			if err != nil {
				log.Fatalf("无法保存配置文件：%s", err)
			}
			service.DelFirstFile("tempinfo.json",delstring)
			fmt.Println("delete skip 2")



		}else{//just return data
			fmt.Println("no delete skip ")

		}
		firstmap :=make(map[string]string)
		firststring := service.ReadJsonFirst("tempinfo.json")
		if firststring== ""{
			parms := &model.Returnparms{
				Code:     "ok",
				Typename: "finish",
			}
			jsoncategorylists, _ := json.Marshal(parms)
			//返回的这个是给json用的，需要去掉
			io.WriteString(res, string(jsoncategorylists))
			return
		}
		//fmt.Println(firststring,reflect.TypeOf(firststring))
		tempmap :=service.GetConfig(tempfile,firststring)
		//fmt.Println(configmap)
		//mardownlist 即可生成json,数据

		if firststring=="excelfile" {
			foldermap :=service.GetConfig(tempfile,"folder")
			//先创建文件夹
			input :=foldermap["input"]
			for k,v :=range tempmap{
				pathflag,_ :=service.PathExists(input+"/"+v)
				if !pathflag {
					tempmap[k] = ""
				}
			}
		}
		firstmap = tempmap
		if firststring=="userorder" {

			foldermap :=service.GetConfig(tempfile,"gouserorder")
			gouserorder :=foldermap["order"]
			//str :=strings.Split(gouserorder, ",")
			firstmap["gouserorder"] = gouserorder
		}
		if firststring=="special_regulation" {
			foldermap :=service.GetConfig(tempfile,"userorder")
			str :="{"
			for k,_:=range tempmap{
				if _,ok :=foldermap[k];ok{
					mstr :="'"+k+"':'"+foldermap[k]+"',"
					str+=mstr
				}
			}
			str = str[:(len(str)-1)]
			str+="}"
			firstmap["userinfo"] = str
		}

		if firststring=="allpeople" {
			configmap :=service.InitConfig(tempfile)
			allusermap :=service.ReadInExcel(configmap)
			//fmt.Println(allusermap)
			userinfoflip := service.GetUserInfo(tempfile)
			if _,ok :=allusermap["changhua"];ok{
				str :="{"
				for _,v:=range allusermap["changhua"]{
					mstr :="'"+v.Userid+"':'"+userinfoflip[v.Username]+"',"
					str+=mstr
					//if _,ok :=userinfoflip[v.Username];ok{
					//
					//}
				}
				str = str[:(len(str)-1)]
				str+="}"
				firstmap["changhua"] = str
			}

			if _,ok :=allusermap["pudong"];ok{
				str :="{"
				for _,v:=range allusermap["pudong"]{
					mstr :="'"+v.Userid+"':'"+userinfoflip[v.Username]+"',"
					str+=mstr
					//if _,ok :=userinfoflip[v.Username];ok{
					//
					//}
				}
				str = str[:(len(str)-1)]
				str+="}"
				firstmap["pudong"] = str
			}
			if _,ok :=allusermap["maoming"];ok{
				str :="{"
				for _,v:=range allusermap["maoming"]{
					mstr :="'"+v.Userid+"':'"+userinfoflip[v.Username]+"',"
					str+=mstr
					//if _,ok :=userinfoflip[v.Username];ok{
					//
					//}
				}
				str = str[:(len(str)-1)]
				str+="}"
				firstmap["maoming"] = str
			}
			if _,ok :=allusermap["dingding"];ok{
				str :="{"
				for _,v:=range allusermap["dingding"]{
					//if _,ok :=userinfoflip[v.Username];ok{
					//
					//}
					mstr :="'"+v.Userid+"':'"+userinfoflip[v.Username]+"',"
					str+=mstr
				}
				str = str[:(len(str)-1)]
				str+="}"
				firstmap["dingding"] = str
			}
		}


		firstmap["typename"] = firststring
		jsoncategorylists, _ := json.Marshal(firstmap)
		//返回的这个是给json用的，需要去掉
		res.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
		io.WriteString(res, string(jsoncategorylists))
	}

	//queryForm, err := url.ParseQuery(req.URL.RawQuery)
	//if err != nil || len(queryForm["skip"]) == 0 {
	//	fmt.Println("query is wrong", err)
	//	return
	//}
	//skip,_:= strconv.Atoi(queryForm["skip"][0])
	//fmt.Println("skip is :",skip)
	//
	//if skip==1{
	//	//del json data
	//	service.DelFirstFile("tempinfo.json",firststring)
	//}else if skip==2{
	//	//save config file
	//
	//}else{//just return data
	//
	//}


	//firstmap :=service.GetConfig("config.ini",firststring)
	////fmt.Println(configmap)
	////mardownlist 即可生成json,数据
	//firstmap["typename"] = firststring
	//jsoncategorylists, _ := json.Marshal(firstmap)
	////返回的这个是给json用的，需要去掉
	//res.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	//io.WriteString(res, string(jsoncategorylists))
}


func OnExport(res http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	//fmt.Printf("%s\n", result)
	parms := &model.Returnparms{
		Code:     "0",
		Typename: "not ok",
	}

	rstmap := make(map[string]interface{})
	if err := json.Unmarshal(result, &rstmap); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(rstmap,reflect.TypeOf(rstmap),rstmap["time"],reflect.TypeOf(rstmap["time"]))
		time := rstmap["time"].(float64)
		fmt.Println("time is :",time)
		configmap :=service.InitConfig("config.ini")
		//第一次export
		outfolder :=configmap["folder"].(map[string]string)["output"]



		filestr := ""
		if time==2{
			//第二次export
			secondexportexcelfile := configmap["secondexportexcelfile"].(string)
			filestr = secondexportexcelfile
			service.UpdateExport(configmap)
			service.ReadExport(configmap)

		}else{
			//第一次export
			firstexportexcelfile := configmap["firstexportexcelfile"].(string)
			filestr = firstexportexcelfile
			outfile := outfolder+"/"+firstexportexcelfile
			data :=service.DealExcel(configmap)
			service.Export(data,outfile)

			importantfilename := configmap["importantfilename"].(string)
			importantfile := outfolder+"/"+importantfilename
			service.WriteResult(data,importantfile)
		}
		str :="请到当前目录的"+outfolder+"文件夹下找"+filestr+"文件"
		parms = &model.Returnparms{
			Code:     "1",
			Typename: str,
		}

	}


	jsoncategorylists, _ := json.Marshal(parms)
	//返回的这个是给json用的，需要去掉
	io.WriteString(res, string(jsoncategorylists))
}

//have upload doc
func OnDocuplExport(res http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(32 << 20)
	fmt.Println(req.MultipartForm.Value)
	parms := &model.Returnparms{
		Code:     "0",
		Typename: "not ok",
	}

	//form.value : timme,name form.file
	formvalue := req.MultipartForm.Value
	if _,ok :=formvalue["time"];ok{
		time := formvalue["time"][0]

		configmap :=service.InitConfig("config.ini")
		//第一次export
		outfolder :=configmap["folder"].(map[string]string)["output"]
		firstexportexcelfile := configmap["firstexportexcelfile"].(string)
		outfile := outfolder+"/"+firstexportexcelfile
		filestr := ""
		if time=="2"{
			//second export need upload first
			uploadfilename := formvalue["name"][0]
			fmt.Println("hhh ,I get upload file name is:",uploadfilename)

			file, _, err := req.FormFile("file")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()



			//删除文件
			pathflag,_:= service.PathExists(outfile)
			if pathflag {
				ferr:=os.Remove(outfile)
				if ferr!=nil{
					fmt.Println(ferr)
				}
			}

			//newinputfile :=outfolder+"/"+uploadfilename
			newinputfile :=outfile
			f, err := os.OpenFile(newinputfile, os.O_WRONLY|os.O_CREATE, 0666)
			//fmt.Println("----",handler.Filename)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)

			//第二次export
			secondexportexcelfile := configmap["secondexportexcelfile"].(string)
			filestr = secondexportexcelfile
			//这些都需要返回值
			service.UpdateExport(configmap)
			service.ReadExport(configmap)
		}else{
			//第一次导出
			data :=service.DealExcel(configmap)
			service.Export(data,outfile)

			importantfilename := configmap["importantfilename"].(string)
			importantfile := outfolder+"/"+importantfilename
			service.WriteResult(data,importantfile)
			filestr = firstexportexcelfile
		}
		str :="请到当前目录的"+outfolder+"文件夹下找"+filestr+"文件"
		parms = &model.Returnparms{
			Code:     "1",
			Typename: str,
		}
	}
	jsoncategorylists, _ := json.Marshal(parms)
	//返回的这个是给json用的，需要去掉
	io.WriteString(res, string(jsoncategorylists))

}



func init() {
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
	file := "./log/" +"log"+ ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[goexcel]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}