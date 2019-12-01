package service

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io/ioutil"
	"goexcel/model"
	"log"
	"strconv"
	"strings"
	"time"
)

// UpdateExport is a function
func UpdateExport(configmap map[string]interface{}) {
	start :=configmap["start"].(string)
	end := configmap["end"].(string)

	importantfilename := configmap["importantfilename"].(string)
	outfolder :=configmap["folder"].(map[string]string)["output"]
	outfile :=outfolder+"/"+"exportmiddle.xlsx"
	readfirstexportexcelfile := outfolder+"/"+configmap["firstexportexcelfile"].(string)
	importantfile := outfolder+"/"+importantfilename
	normaltmpperdayhour,err := strconv.ParseFloat(configmap["perdayhour"].(string),64)
	special_regulationmap :=configmap["special_regulation"].(map[string]string)
	special_regulation :=make(map[int64]float64)
	for k,v :=range special_regulationmap{
		gk, _ := strconv.ParseInt(k, 10, 64)
		gv,_ := strconv.ParseFloat(v,64)
		special_regulation[gk] = gv
	}



	log.Println("1-come into this readexport function，tmpregulation time is:",normaltmpperdayhour,special_regulation)
	sheet :="考勤"
	xlsx, err := excelize.OpenFile(readfirstexportexcelfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("1-can open "+readfirstexportexcelfile+" file")
	//start := "2019-09-23" //待转化为时间戳的字符串
	//end := "2019-10-22"   //待转化为时间戳的字符串
	////日期转化为时间戳
	log.Println("get start time")
	log.Println(start,end)
	slice := PrDates(start,end)
	log.Println("get slice is :",slice)
	datelen := len(slice)
	log.Println("before获取cmap")
	cmap :=GetWhitelist(importantfile)
	log.Println("after 获取cmap:",cmap)

	//fmt.Println(cmap)
	log.Println("获取cmap")
	len :=cmap["len"].(int)
	maxrowlen :=7*len+4
	headslice := []string{"D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","AA","AB","AC","AD","AE","AF","AG","AH"}
	log.Println("2222222")
	//xbgtype := getCellBgType(xlsx, sheet, "T96")
	//log.Println("sssssssssssss-------",xbgtype)

	defaultunsignbgtype := getCellBgType(xlsx, sheet, "G1")
	defaulttiaoxiubgtype := getCellBgType(xlsx, sheet, "J1")
	defaultsundaybgtype := getCellBgType(xlsx, sheet, "M1")
	style, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#00ffe9"],"pattern":1}}`)
	if err != nil {
		fmt.Println(err)
	}
	if defaultunsignbgtype=="solid" {
		xlsx.SetCellStyle(sheet, "G1", "G1", style)
	}
	if defaulttiaoxiubgtype=="solid" {
		style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#836FFF"],"pattern":1}}`)
		if err != nil {
			fmt.Println(err)
		}
		xlsx.SetCellStyle(sheet, "J1", "J1", style)
	}

	if defaultsundaybgtype=="solid" {
		style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#DAA520"],"pattern":1}}`)
		if err != nil {
			fmt.Println(err)
		}
		xlsx.SetCellStyle(sheet, "M1", "M1", style)
	}








	for k := 4; k <= maxrowlen; k+=7 {

		axis := "C"+strconv.Itoa(k)
		useridstring ,_:= xlsx.GetCellValue(sheet,axis)
		userid, _ := strconv.ParseInt(useridstring, 10, 64)
		if userid==0 {
			continue
		}
		for c := 0; c < datelen; c+=1 {

			datenum :=slice[int64(c)]
			weekend :=false
			num :=ZellerFunction2Week(datenum)
			if num==6 || num ==7 {
				weekend = true
			}
			_, ok := special_regulation[userid]
			if ok && weekend {
				style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#DAA520"],"pattern":1},"alignment":{"horizontal":"center","wrap_text":true}}`)
				if err != nil {
					fmt.Println(err)
				}
				dateaxis := headslice[c]+strconv.Itoa(k-1)
				xlsx.SetCellStyle(sheet, dateaxis, dateaxis, style)
			}

			firstaxis := headslice[c]+strconv.Itoa(k+1)
			firstbgtype := getCellBgType(xlsx, sheet, firstaxis)


			if firstbgtype=="solid" {
				firstval ,_:= xlsx.GetCellValue(sheet,firstaxis)
				if strings.Index(firstval, "：") > -1 {
					firstval = strings.Replace(firstval, "：", ":", -1)
				}
				xlsx.SetCellValue(sheet,firstaxis , firstval)

				style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#00ffe9"],"pattern":1}}`)
				if err != nil {
					fmt.Println(err)
				}
				xlsx.SetCellStyle(sheet, firstaxis, firstaxis, style)

			}
			lastaxis := headslice[c]+strconv.Itoa(k+2)
			lastbgtype := getCellBgType(xlsx, sheet, lastaxis)

			if lastbgtype=="solid" {

				lastval ,_:= xlsx.GetCellValue(sheet,lastaxis)
				if strings.Index(lastval, "：") > -1 {
					lastval = strings.Replace(lastval, "：", ":", -1)
				}
				xlsx.SetCellValue(sheet,lastaxis , lastval)
				style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#00ffe9"],"pattern":1}}`)
				if err != nil {
					fmt.Println(err)
				}
				xlsx.SetCellStyle(sheet, lastaxis, lastaxis, style)

			}
			//if userid ==505 {
			//	log.Println("aaaaaa:----",datenum,firstbgtype,lastbgtype)
			//}


		}
		rowline := strconv.Itoa(k+4)

		style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FF0000"],"pattern":1}}`)
		if err != nil {
			fmt.Println(err)
		}
		xlsx.SetCellStyle(sheet, "D"+rowline, "D"+rowline, style)

	}


	//style, err = xlsx.NewStyle(`{"alignment":{"horizontal":"center","wrap_text":true}}`)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//xlsx.SetCellStyle(sheet, "A1", "AH"+strconv.Itoa(maxrowlen), style)
	//

	for index:=1;index<=maxrowlen;index++{
		xlsx.SetRowHeight(sheet,index,30)
	}
	xlsx.SetColWidth(sheet,"A","AH",9)

	// 根据指定路径保存文件
	err = xlsx.SaveAs(outfile)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("最后一步保存excel成功")

}


//ReadExport is a function
func ReadExport(configmap map[string]interface{}) {
	start :=configmap["start"].(string)
	end := configmap["end"].(string)
	userorderslice := configmap["gouserorderslice"].([]int64)
	importantfilename := configmap["importantfilename"].(string)
	outfolder :=configmap["folder"].(map[string]string)["output"]
	outfilename := configmap["secondexportexcelfile"].(string)
	outfile :=outfolder+"/"+outfilename
	readfirstexportexcelfile := outfolder+"/"+configmap["firstexportexcelfile"].(string)
	importantfile := outfolder+"/"+importantfilename
	normaltmpperdayhour,err := strconv.ParseFloat(configmap["perdayhour"].(string),64)
	special_regulationmap :=configmap["special_regulation"].(map[string]string)
	special_regulation :=make(map[int64]float64)
	for k,v :=range special_regulationmap{
		gk, _ := strconv.ParseInt(k, 10, 64)
		gv,_ := strconv.ParseFloat(v,64)
		special_regulation[gk] = gv
	}


	log.Println("1-come into this readexport function，tmpregulation time is:",normaltmpperdayhour,special_regulation)
	sheet :="考勤"
	xlsx, err := excelize.OpenFile(readfirstexportexcelfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("1-can open "+readfirstexportexcelfile+" file")
	//start := "2019-09-23" //待转化为时间戳的字符串
	//end := "2019-10-22"   //待转化为时间戳的字符串
	////日期转化为时间戳
	log.Println("get start time")
	log.Println(start,end)
	slice := PrDates(start,end)
	log.Println("get slice is :",slice)
	datelen := len(slice)
	log.Println("before获取cmap")
	cmap :=GetWhitelist(importantfile)
	log.Println("after 获取cmap:",cmap)

	attendancelist := cmap["attendance"].(map[int64]*model.Attendancelist)
	//fmt.Println(cmap)
	log.Println("获取cmap")
	len :=cmap["len"].(int)
	maxrowlen :=7*len+4
	headslice := []string{"D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","AA","AB","AC","AD","AE","AF","AG","AH"}
	useridmap :=make(map[int64]map[string]interface{})
	log.Println("2222222")


	for k := 4; k <= maxrowlen; k+=7 {
		category :=make(map[string]interface{})
		axis := "C"+strconv.Itoa(k)
		useridstring ,_:= xlsx.GetCellValue(sheet,axis)
		userid, _ := strconv.ParseInt(useridstring, 10, 64)
		if userid==0 {
			continue
		}
		tmpperdayhour :=normaltmpperdayhour
		if _,ok:=special_regulation[userid];ok{
			tmpperdayhour = special_regulation[userid]
		}
		var unsignedtimes int = 0
		var attendance_days = 0
		var monthtiaoxiu = 0
		var afterearlylate30mintimes = 0
		var afterearlylate15mintimes = 0
		var beforeearlylate15mintimes = 0
		all :=0.00


		for c := 0; c < datelen; c+=1 {
			datenum :=slice[int64(c)]
			firstaxis := headslice[c]+strconv.Itoa(k+1)
			firstbgcolor := getCellBgColor(xlsx, sheet, firstaxis)
			firstval, _ := xlsx.GetCellValue(sheet, firstaxis)
			if firstbgcolor=="FF00FFE9" && firstval!="" {
				unsignedtimes+=1
			}
			firstsigned := datenum+" "+firstval
			lastaxis := headslice[c]+strconv.Itoa(k+2)
			lastbgcolor := getCellBgColor(xlsx, sheet, lastaxis)
			lastval, _ := xlsx.GetCellValue(sheet, lastaxis)
			if lastbgcolor=="FF00FFE9" && lastval!="" {
				unsignedtimes+=1

			}
			lastsigned := datenum+" "+lastval
			if !(firstval=="" && lastval =="") {
				attendance_days+=1
			}


			tiaoxiuaxis := headslice[c]+strconv.Itoa(k)
			tiaoxiubgcolor := getCellBgColor(xlsx, sheet, tiaoxiuaxis)
			tiaoxiuval, _ := xlsx.GetCellValue(sheet, tiaoxiuaxis)
			if tiaoxiuval!="" && tiaoxiubgcolor!="FFA020F0" {
				tiaoxiuint, _ := strconv.ParseInt(tiaoxiuval, 10, 64)
				tiaoxiutimestamp := GetTimstamp(lastsigned)
				t := time.Unix(tiaoxiuint*60+tiaoxiutimestamp, 0)
				dateStr := t.Format("2006-01-02 15:04")
				lastsigned = dateStr[11:]
				attendancelist[userid].Signedtime[datenum].Lastsigned = lastsigned
			}
			if tiaoxiubgcolor=="FF836FFF" {
				monthtiaoxiu+=1
			}

			hour :=0
			minute :=0
			timearr := Timediff(firstsigned,lastsigned)

			if firstval==""||lastval=="" {
				timearr["day"] = 0
				timearr["hour"] = 0
				timearr["min"] = 0
				timearr["sec"] = 0
			}
			hour = timearr["hour"]
			minute = timearr["min"]
			periodaxis := headslice[c]+strconv.Itoa(k+3)
			xlsx.SetCellValue(sheet,periodaxis,strconv.Itoa(hour)+"h"+strconv.Itoa(minute)+"m")



			everyperhour :=Round(float64(hour)+float64(minute)/float64(60),2)

			all+=everyperhour
			all = Round(all,2)





			if firstval!="" && lastval!=""{
				if everyperhour<=tmpperdayhour-0.5{
					afterearlylate30mintimes+=1

				}
				if everyperhour>(tmpperdayhour-0.5) && everyperhour<(tmpperdayhour-0.25) {
					afterearlylate15mintimes+=1

				}
				if everyperhour>=(tmpperdayhour-0.25) && everyperhour<tmpperdayhour {
					beforeearlylate15mintimes+=1

				}

			}
			if (firstval=="" && lastval!="") || (firstval!="" && lastval=="") {
				afterearlylate30mintimes+=1

			}


		}
		rowline := strconv.Itoa(k+4)
		xlsx.SetCellValue(sheet,"B"+rowline,all)
		xlsx.SetCellValue(sheet,"D"+rowline,attendance_days)
		xlsx.SetCellValue(sheet,"F"+rowline,unsignedtimes)
		xlsx.SetCellValue(sheet,"J"+rowline,afterearlylate30mintimes)
		xlsx.SetCellValue(sheet,"M"+rowline,afterearlylate15mintimes)
		xlsx.SetCellValue(sheet,"P"+rowline,beforeearlylate15mintimes)

		category["all"] = all
		category["unsignedtimes"] = unsignedtimes
		category["attendance_days"] = attendance_days
		category["monthtiaoxiu"] = monthtiaoxiu
		category["afterearlylate30mintimes"] = afterearlylate30mintimes
		category["afterearlylate15mintimes"] = afterearlylate15mintimes
		category["beforeearlylate15mintimes"] = beforeearlylate15mintimes
		useridmap[userid] = category
	}
	log.Println("333333")
	if _,ok :=useridmap[0];ok{
		delete(useridmap,0)
	}

	log.Println("already get the right data",useridmap)


	//fmt.Println(useridmap[0])
	//return
	newsheet :="考勤汇总"
	xlsx.NewSheet(newsheet)
	parmas := make(map[string]string)
	parmas["A"] = "序号"
	parmas["B"] = "姓名"
	parmas["C"] = "所属单位"
	parmas["D"] = "职务"
	parmas["E"] = "事假"
	parmas["F"] = "病假"
	parmas["G"] = "年假"
	parmas["H"] = "旷工"
	parmas["I"] = "迟到早退15分钟内"
	parmas["J"] = "迟到早退15-30分钟"
	parmas["K"] = "迟到早退半小时以上"
	parmas["L"] = "未打卡"
	parmas["M"] = "跨月调休"
	parmas["N"] = "扣款"
	parmas["O"] = "全勤奖"
	parmas["P"] = "出勤天数"
	parmas["Q"] = "调整后天数"
	parmas["R"] = "常规月工时"
	parmas["S"] = "实际月工时"
	parmas["T"] = "工资"
	parmas["U"] = "出勤天数"
	log.Println("开始设置新的标签页")
	for k,v :=range parmas {
		axis := k+"1"
		xlsx.SetCellValue(newsheet,axis,v)
	}
	state :=make(map[string]string)
	state["Y4"] = "5"
	state["Y5"] = "5"
	state["Z5"] = "0.1"
	state["Y7"] = "3"
	state["Y8"] = "5"
	state["Z8"] = "0.5"
	state["Z9"] = "0.5"
	state["Y12"] = "0.3"
	log.Println("设置其他页码")

	xlsx.SetCellValue(newsheet,"X1","薪资起始日")
	xlsx.SetCellValue(newsheet,"Y1","薪资末日")
	xlsx.SetCellValue(newsheet,"X2",start)
	xlsx.SetCellValue(newsheet,"Y2",end)
	xlsx.SetCellValue(newsheet,"X4","次数")
	xlsx.SetCellValue(newsheet,"Y4",5)
	xlsx.SetCellValue(newsheet,"Z4","6次及以上")
	xlsx.SetCellValue(newsheet,"X5","未打卡")
	xlsx.SetCellValue(newsheet,"Y5",5)
	xlsx.SetCellValue(newsheet,"Z5",0.1)
	xlsx.SetCellValue(newsheet,"X7","次数")
	xlsx.SetCellValue(newsheet,"Y7",3)
	xlsx.SetCellValue(newsheet,"Z7","3次及以上")
	xlsx.SetCellValue(newsheet,"X8","迟到早退30分钟以内")
	xlsx.SetCellValue(newsheet,"Y8",5)
	xlsx.SetCellValue(newsheet,"Z8",0.5)
	xlsx.SetCellValue(newsheet,"X9","迟到早退超过30分钟")
	xlsx.SetCellValue(newsheet,"Z9",0.5)
	xlsx.SetCellValue(newsheet,"X12","病假扣工资比例")
	xlsx.SetCellValue(newsheet,"Y12",0.3)

	style, err := xlsx.NewStyle(`{"font":{"color":"#CD0000"}}`)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("获取配置文件 userorder slice:",userorderslice)
	log.Println("cj for test:",useridmap[419])
	var krowkey int = 2
	for _,userv :=range  userorderslice {
		newkey := strconv.Itoa(krowkey)

		xlsx.SetCellValue(newsheet,"T"+newkey,0)
		xlsx.SetCellValue(newsheet,"A"+newkey,userv)
		xlsx.SetCellValue(newsheet,"B"+newkey,attendancelist[userv].Username)
		xlsx.SetCellValue(newsheet,"S"+newkey,useridmap[userv]["all"])
		xlsx.SetCellValue(newsheet,"P"+newkey,useridmap[userv]["attendance_days"])

		if useridmap[userv]["unsignedtimes"].(int)>0{
			xlsx.SetCellValue(newsheet,"L"+newkey,useridmap[userv]["unsignedtimes"])
			err = xlsx.SetCellStyle(newsheet, "L"+newkey, "L"+newkey, style)
		}

		if useridmap[userv]["beforeearlylate15mintimes"].(int)>0{
			xlsx.SetCellValue(newsheet,"I"+newkey,useridmap[userv]["beforeearlylate15mintimes"])
			err = xlsx.SetCellStyle(newsheet, "I"+newkey, "I"+newkey, style)
		}

		if useridmap[userv]["afterearlylate15mintimes"].(int)>0{
			xlsx.SetCellValue(newsheet,"J"+newkey,useridmap[userv]["afterearlylate15mintimes"])
			err = xlsx.SetCellStyle(newsheet, "J"+newkey, "J"+newkey, style)
		}

		if useridmap[userv]["afterearlylate30mintimes"].(int)>0{
			xlsx.SetCellValue(newsheet,"K"+newkey,useridmap[userv]["afterearlylate30mintimes"])
			err = xlsx.SetCellStyle(newsheet, "K"+newkey, "K"+newkey, style)
		}

		if useridmap[userv]["monthtiaoxiu"].(int)>0{
			xlsx.SetCellValue(newsheet,"M"+newkey,useridmap[userv]["monthtiaoxiu"])
			err = xlsx.SetCellStyle(newsheet, "M"+newkey, "M"+newkey, style)
		}
		formula :="=Round(IF((I"+newkey+"+J"+newkey+")>"+state["Y7"]+",T"+newkey+"*(I"+newkey+"+J"+newkey+"-"+state["Y7"]+")*"+state["Z8"]+"/26+15,(I"+newkey+"+J"+newkey+")*"+state["Y8"]+")+IF(L"+newkey+">"+state["Y4"]+",T"+newkey+"*(L"+newkey+"-"+state["Y4"]+")*"+state["Z5"]+"/27+25,L"+newkey+"*"+state["Y5"]+")+K"+newkey+"*"+state["Z9"]+"*T"+newkey+"/27+E"+newkey+"*T"+newkey+"/27+F"+newkey+"*"+state["Y12"]+"*T"+newkey+"/27+3*H"+newkey+"*T"+newkey+"/27,2)"
		xlsx.SetCellFormula(newsheet,"Q"+newkey,"=P"+newkey+"+M"+newkey)
		xlsx.SetCellFormula(newsheet,"N"+newkey,formula)

		err = xlsx.SetCellStyle(newsheet, "N"+newkey, "N"+newkey, style)
		krowkey+=1

	}
	log.Println("安全通过:")

	// 根据指定路径保存文件
	err = xlsx.SaveAs(outfile)
	if err != nil {
		fmt.Println(err)
	}

	log.Println("最后一步保存excel成功")

}

//GetWhitelist is a function
func GetWhitelist(importantfilename string) map[string]interface{}{
	b, err := ioutil.ReadFile(importantfilename)
	if err != nil {
		fmt.Print(err)
	}
	c := make(map[string]interface{})
	whitelist := make(map[int64]*model.Attendancelist)
	err = json.Unmarshal([]byte(string(b)), &whitelist)
	if err != nil {
		panic(err)
	}
	len :=len(whitelist)
	c["attendance"] = whitelist
	c["len"] = len
	return c

	//fmt.Println(whitelist)
}


func getCellBgColor(xlsx *excelize.File, sheet, axix string) string {
	styleID,_:= xlsx.GetCellStyle(sheet, axix)
	fillID := xlsx.Styles.CellXfs.Xf[styleID].FillID
	fgColor := xlsx.Styles.Fills.Fill[fillID].PatternFill.FgColor
	if fgColor.Theme != nil {
		srgbClr := xlsx.Theme.ThemeElements.ClrScheme.Children[*fgColor.Theme].SrgbClr.Val
		return excelize.ThemeColor(srgbClr, fgColor.Tint)
	}
	return fgColor.RGB
}

func getCellBgType(xlsx *excelize.File, sheet, axix string) string {
	styleID,_:= xlsx.GetCellStyle(sheet, axix)
	fillID := xlsx.Styles.CellXfs.Xf[styleID].FillID
	patterntype := xlsx.Styles.Fills.Fill[fillID].PatternFill.PatternType
	return patterntype
}


