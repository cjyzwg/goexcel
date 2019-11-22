package service

import (
	"github.com/Unknwon/goconfig"
	"log"
	"sort"
	"strconv"
	"strings"
)

func InitConfig(path string) map[string]interface{}{
	myMap := make(map[string]interface{})
	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		log.Fatalf("无法加载配置文件：%s", err)
	}
	startvalue, err := cfg.GetValue("da", "start")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "start", err)
	}
	myMap["start"] = startvalue
	endvalue, err := cfg.GetValue("da", "end")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "end", err)
	}
	myMap["end"] = endvalue

	changhuasec, err := cfg.GetSection("changhua")
	//fmt.Println(changhuasec)
	newchanghuapeople := make(map[int64]int64)
	changhuaflip := make(map[int64]int64)
	for k, v := range changhuasec {
		newv,_ := strconv.ParseInt(v, 10, 64)
		newk,_ := strconv.ParseInt(k, 10, 64)
		changhuaflip[newv] = newk
		newchanghuapeople[newk] = newv
	}
	//fmt.Println(changhuaflip)
	myMap["changhuaflip"] = changhuaflip
	myMap["newchanghuapeople"] = newchanghuapeople

	pudongsec, err := cfg.GetSection("pudong")
	//fmt.Println(pudongsec)
	pudongflip := make(map[int64]int64)
	newpudongpeople := make(map[int64]int64)
	for k, v := range pudongsec {
		newv,_ := strconv.ParseInt(v, 10, 64)
		newk,_ := strconv.ParseInt(k, 10, 64)
		pudongflip[newv] = newk
		newpudongpeople[newk] = newv
	}
	//fmt.Println(pudongflip)
	myMap["pudongflip"] = pudongflip
	myMap["newpudongpeople"] = newpudongpeople

	maomingsec, err := cfg.GetSection("maoming")
	//fmt.Println(maomingsec)
	maomingflip := make(map[int64]int64)
	newmaomingpeople := make(map[int64]int64)
	for k, v := range maomingsec {
		newv,_ := strconv.ParseInt(v, 10, 64)
		newk,_ := strconv.ParseInt(k, 10, 64)
		maomingflip[newv] = newk
		newmaomingpeople[newk] = newv
	}
	//fmt.Println(maomingflip)
	myMap["maomingflip"] = maomingflip
	myMap["newmaomingpeople"] = newmaomingpeople
	perdayhour, err := cfg.GetValue("regulation","perdayhour")
	//fmt.Println(perdayhour)
	myMap["perdayhour"] = perdayhour

	translate, err := cfg.GetSection("translate")
	//fmt.Println(translate)
	myMap["translate"] = translate

	//userorder 存入切片
	var users []string
	var userorderslice []int64
	userorder, err := cfg.GetSection("userorder")
	for k, _ := range userorder {
		//gk, _ := strconv.ParseInt(k, 10, 64)
		users = append(users,k)
	}
	sort.Strings(users)
	for _, v := range users {
		gs, _ := strconv.ParseInt(v, 10, 64)
		userorderslice = append(userorderslice,gs)
	}
	//fmt.Println("userorder is:",userorderslice)
	myMap["userorderslice"] = userorderslice

	var gouserorderslice []int64
	gouserorder, err := cfg.GetValue("gouserorder","order")
	str :=strings.Split(gouserorder, ",")
	for _, v := range str {
		gs, _ := strconv.ParseInt(v, 10, 64)
		gouserorderslice = append(gouserorderslice,gs)
	}


	userinfo :=make(map[int64]string)
	userinfomap, err := cfg.GetSection("userorder")

	for k, v := range userinfomap {
		gk, _ := strconv.ParseInt(k, 10, 64)
		userinfo[gk] = v
	}

	myMap["userinfo"] = userinfo

	//fmt.Println("go user order slice is:",gouserorderslice)
	myMap["gouserorderslice"] = gouserorderslice

	special_regulation, err := cfg.GetSection("special_regulation")
	//fmt.Println(special_regulation)
	myMap["special_regulation"] = special_regulation


	dingdingsec, err := cfg.GetSection("dingding")
	myMap["dingdingpeople"] = dingdingsec
	//fmt.Println("dingding",dingdingsec)

	da, err := cfg.GetSection("da")

	myMap["start"] = da["start"]
	myMap["end"] = da["end"]

	importantfile, err := cfg.GetSection("important")
	myMap["importantfilename"] = importantfile["filename"]

	excelfile, err := cfg.GetSection("excelfile")
	if _,ok :=excelfile["changhua"];ok{
		myMap["changhuaexcelfile"] = excelfile["changhua"]
	}
	if _,ok :=excelfile["maoming"];ok{
		myMap["maomingexcelfile"] = excelfile["maoming"]
	}
	if _,ok :=excelfile["pudong"];ok{
		myMap["pudongexcelfile"] = excelfile["pudong"]
	}
	if _,ok :=excelfile["dingding"];ok{
		myMap["dingdingexcelfile"] = excelfile["dingding"]
	}





	exportexcelfile, err := cfg.GetSection("exportexcelfile")
	myMap["firstexportexcelfile"] = exportexcelfile["firstexportexcelfile"]

	folder, err := cfg.GetSection("folder")
	myMap["folder"] = folder


	log.Println(myMap)
	return  myMap
}

