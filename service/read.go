package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tealeg/xlsx"
	"goexcel/model"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//ReadInExcel is a function
func ReadInExcel(configmap map[string]interface{}) map[string][]model.Tempuser{
	inputfolder :=configmap["folder"].(map[string]string)["input"]
	rst := make(map[string][]model.Tempuser)
	if _,ok :=configmap["changhuaexcelfile"];ok{
		excelfile :=inputfolder+"/"+configmap["changhuaexcelfile"].(string)
		xlsx, err := excelize.OpenFile(excelfile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		sheet :="考勤记录"
		rows ,_:= xlsx.GetRows(sheet)
		len :=len(rows)


		var k int
		for k=5; k<=len; k+=2 {
			useridaxis := "C"+strconv.Itoa(k)
			userid,_:= xlsx.GetCellValue(sheet,useridaxis)
			usernameaxis := "K"+strconv.Itoa(k)
			username,_:= xlsx.GetCellValue(sheet,usernameaxis)
			username = strings.Replace(username, " ", "", -1)
			var col int
			valflag := false
			for col = 1;col<33;col++{
				colline,_:=excelize.ColumnNumberToName(col)
				colaxis := colline+strconv.Itoa(k+1)
				val,_:= xlsx.GetCellValue(sheet,colaxis)
				if val!=""{
					valflag = true
					break
				}
			}
			if valflag {
				tempuser := model.Tempuser{
					Userid:   userid,
					Username: username,
				}
				rst["changhua"] = append(rst["changhua"],tempuser)
			}


		}
	}
	if _,ok :=configmap["pudongexcelfile"];ok{
		excelfile :=inputfolder+"/"+configmap["pudongexcelfile"].(string)
		xlsx, err := excelize.OpenFile(excelfile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		sheet :="考勤记录"
		rows ,_:= xlsx.GetRows(sheet)
		len :=len(rows)


		var k int
		for k=5; k<=len; k+=2 {
			useridaxis := "C"+strconv.Itoa(k)
			userid,_:= xlsx.GetCellValue(sheet,useridaxis)
			usernameaxis := "K"+strconv.Itoa(k)
			username,_:= xlsx.GetCellValue(sheet,usernameaxis)
			username = strings.Replace(username, " ", "", -1)
			var col int
			valflag := false
			for col = 1;col<33;col++{
				colline,_:=excelize.ColumnNumberToName(col)
				colaxis := colline+strconv.Itoa(k+1)
				val,_:= xlsx.GetCellValue(sheet,colaxis)
				if val!=""{
					valflag = true
					break
				}
			}
			if valflag {
				tempuser :=model.Tempuser{
					Userid:   userid,
					Username: username,
				}
				rst["pudong"] = append(rst["pudong"],tempuser)
			}


		}
	}
	if _,ok :=configmap["maomingexcelfile"];ok{
		excelfile :=inputfolder+"/"+configmap["maomingexcelfile"].(string)
		xlsx, err := excelize.OpenFile(excelfile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		sheet :="grid_hr_at_card_rec"
		rows ,_:= xlsx.GetRows(sheet)
		len :=len(rows)


		var k int
		userid := ""
		for k=3; k<=len; k+=1 {
			useridaxis := "A"+strconv.Itoa(k)
			tempuserid,_ := xlsx.GetCellValue(sheet,useridaxis)
			if tempuserid != userid {
				userid = tempuserid
				usernameaxis := "B"+strconv.Itoa(k)
				username,_:= xlsx.GetCellValue(sheet,usernameaxis)
				username = strings.Replace(username, " ", "", -1)
				tempuser :=model.Tempuser{
					Userid:   userid,
					Username: username,
				}
				rst["maoming"] = append(rst["maoming"],tempuser)
			}
		}
	}
	if _,ok :=configmap["dingdingexcelfile"];ok{
		excelfile :=inputfolder+"/"+configmap["dingdingexcelfile"].(string)
		xlsx, err := excelize.OpenFile(excelfile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		sheet :="打卡时间"
		rows ,_:= xlsx.GetRows(sheet)
		len :=len(rows)


		var k int
		userid := ""
		for k=4; k<=len; k+=1 {
			useridaxis := "E"+strconv.Itoa(k)
			tempuserid,_:= xlsx.GetCellValue(sheet,useridaxis)
			if tempuserid != userid {
				userid = tempuserid
				usernameaxis := "A"+strconv.Itoa(k)
				username,_:= xlsx.GetCellValue(sheet,usernameaxis)
				index :=strings.Index(strings.Replace(username, " ", "", -1), "（")
				if index>-1 {
					username = username[:index]
				}
				username = strings.Replace(username, " ", "", -1)

				var col int
				valflag := false
				for col = 6;col<38;col++{
					colline,_:=excelize.ColumnNumberToName(col)
					colaxis := colline+strconv.Itoa(k)
					val,_:= xlsx.GetCellValue(sheet,colaxis)
					if val!=""{
						valflag = true
						break
					}
				}
				if valflag {
					tempuser :=model.Tempuser{
						Userid:   userid,
						Username: username,
					}
					rst["dingding"] = append(rst["dingding"],tempuser)
				}
			}
		}
	}
	log.Println("changhua",rst["changhua"],"maoming",rst["maoming"],"pudong",rst["pudong"],"dingding",rst["dingding"])

	return rst

}

// ReadTealegXlsx 读取excel
//func ReadTealegXlsx(activearr map[int64]string) {
func ReadTealegXlsx(activemap map[int64]int64,excelFileName string,start string,end string) map[int64]*Attendancelist {
	//excelFileName := "test.xlsx"
	// 根据指定的xlsx名称，返回xlsx文件结构
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	//start,end := GetConfigStart("config.ini")
	//start:="2019-09-23"
	//end:="2019-10-22"
	datearr :=PrDates(start,end)
	// xlFile.Sheets :一个xlsx文件可能有多个Sheet
	//pattern := "/^[0-9:]*$/"
	pattern := "[0-9]+:[0-9]+"
	attendancelists := make(map[int64]*Attendancelist)
	//newattendancelists := make(map[int64]*Attendancelist)
	for _, sheet := range xlFile.Sheets {
		// sheet.Rows:得到sheet的每一行
		if strings.Index(strings.Replace(sheet.Name, " ", "", -1), "考勤记录") > -1 {
			//fmt.Println(sheet.Row(2).Cells[2].Value)// 获取某一个单元格的值
			fmt.Println("+++++++++++++")
			var lastuserid int64=0
			for keyrow, row := range sheet.Rows {

				if keyrow<4 {
					continue
				}
				var i int64 =0
				attendancelist := &Attendancelist{}
				signedtime :=make(map[string]*Signedtimelist)
				// row.Cells:每一行的每一列
				for keycol, cell := range row.Cells {
					fmt.Println("------")
					val := cell.String()
					fmt.Println("value is :",val)
					fmt.Println("llllllll",lastuserid)
					//fmt.Print(text, "\t")
					if strings.Index(strings.Replace(val, " ", "", -1), "工号") > -1 {
						fmt.Println(11111111)
						//fmt.Println(sheet.Row(keyrow).Cells[2].Value)
						//这个获取的templastuserid 是直接从方法里获取的
						templastuserid, err := strconv.ParseInt(sheet.Row(keyrow).Cells[2].Value, 10, 64)
						if err!=nil {
							fmt.Println("this has error come from string to int64")
						}

						lastuserid = templastuserid
						fmt.Println("lastuserid is :",templastuserid,"username is:",sheet.Row(keyrow).Cells[10].Value)
						attendancelist.Userid = templastuserid
						attendancelist.Username = sheet.Row(keyrow).Cells[10].Value
						attendancelists[templastuserid] = attendancelist
						fmt.Println("poijkl:",attendancelists)




						//attendancelists.Attendances = append(attendancelists.Attendances,attendancelist)
						//attendancelists.Attendances = append(attendancelists.Attendances,Attendancelist{Userid:lastuserid,Username:sheet.Row(keyrow).Cells[10].Value})

						//fmt.Println(lastuserid)
						break
						//fmt.Println(sheet.Row(2).Cells[keycol].Value)

						//fmt.Println(lastuserid)
					}
					fmt.Println("//////////////////////",lastuserid)
					signedtimelist := &Signedtimelist{}
					if val=="" || GetCompileData(val,pattern) {


						_, exists := datearr[i]
						if exists {//存在
							fmt.Println(val)
							fmt.Println("@@@@@")
							staticendtimestamp :=GetTimstamp(datearr[i]+" 04:00")
							staticstarttimestamp :=GetTimstamp(datearr[i]+" 16:00")
							fmt.Println("static end:",datearr[i]+" 04:00",",static start:",datearr[i]+" 16:00")

							firstsigned := ""
							lastsigned := ""
							if val!="" {
								signedarr := Getsplit(val)
								fmt.Println(signedarr,len(signedarr))
								if len(signedarr)>1 {
									firstsigned = signedarr[0]
									firsttimestamp :=GetTimstamp(datearr[i]+" "+firstsigned)

									if firsttimestamp<staticendtimestamp {
										_,ok := datearr[i-1]
										if ok {
											attendancelists[lastuserid].Signedtime[datearr[i-1]].Lastsigned = "23:59"
										}
										delete(signedarr,0)
										firstsigned = signedarr[0]
									}
									if firsttimestamp>staticstarttimestamp {
										firstsigned = ""
									}
									lastsigned = signedarr[int64(len(signedarr)-1)]
									lasttimestamp :=GetTimstamp(datearr[i]+" "+lastsigned)
									if lasttimestamp<staticstarttimestamp {
										lastsigned = ""
									}

								}
								if len(signedarr)==1 {

									tempsigned := signedarr[0]
									temptimestamp :=GetTimstamp(datearr[i]+" "+tempsigned)

									if temptimestamp<staticendtimestamp {
										_,ok := datearr[i-1]
										if ok {
											attendancelists[lastuserid].Signedtime[datearr[i-1]].Lastsigned = "23:59"
										}
										delete(signedarr,0)
									}
									if len(signedarr)==1 {

										//firstsigned = signedarr[0]
										tempsigned := signedarr[0]
										temptimestamp :=GetTimstamp(datearr[i]+" "+tempsigned)
										fmt.Println("firstsigned:",datearr[i]+" "+tempsigned,",lastsigned:",lastsigned,",date i is:",datearr[i])
										if temptimestamp<staticstarttimestamp {
											fmt.Println("ccccccc")
											firstsigned = tempsigned
										}else{
											fmt.Println("ddddddd")
											lastsigned = tempsigned
										}
										//if lastuserid==2 {
										//	fmt.Println("aaaaaaaaaaaaaaaaaaaa")
										//	fmt.Println(datearr[i],firstsigned,lastsigned)
										//}

									}

								}
							}else{
								fmt.Println("val is empty,row is:",keyrow,",column:",keycol)
							}
							fmt.Println("kkkk",lastuserid,"attendancelists:",attendancelists[lastuserid])
							fmt.Println("lllllll",datearr[i],",type is:",reflect.TypeOf(datearr[i]),"dawdaw:",attendancelists[lastuserid].Signedtime)

							//signedtimelist := &Signedtimelist{
							//	Firstsigned: firstsigned,
							//	Lastsigned:  lastsigned,
							//	Lackperiod:  "",
							//	Period:      "",
							//}
							signedtimelist.Firstsigned =firstsigned
							signedtimelist.Lastsigned = lastsigned
							signedtime[datearr[i]] = signedtimelist
							fmt.Println("signedtime:",signedtime)



							fmt.Println("kkkk",lastuserid,"attendancelists:",attendancelists)
							i++
						}else{
							fmt.Println("datearr 中缺少：",i)
						}
					}else{
						fmt.Println("adawdawdwdwa")
					}
				}
				attendancelists[lastuserid].Signedtime = signedtime

				fmt.Println()

			}
			fmt.Println("++++++++++++++++")
			fmt.Println(attendancelists)



			//map[1:0xc0008bcb00 2:0xc000483fe0 3:0xc000482b40 4:0xc00007c5e0
			//5:0xc000229280 6:0xc00000d600 7:0xc0007780a0 8:0xc000b768e0 9:0xc000941a60
			//10:0xc0002970e0 11:0xc000ff1420 12:0xc000fe6320 13:0xc000f9f160 14:0xc000f61980
			//15:0xc000f603e0 16:0xc00025bea0 17:0xc000f30820 18:0xc000cd7ca0 19:0xc000e85d80 20:0xc000e847e0 21:0xc000e173c0
			//22:0xc000dc3f00 23:0xc000dc2d80 24:0xc000d65a00 25:0xc000900e20 26:0xc000d52b40 27:0xc000d20ec0 28:0xc000d18400
			//29:0xc000ce8500 30:0xc000c93ae0 31:0xc000c92540 32:0xc000c5cf80 33:0xc000c2fb60 34:0xc000c1c8c0 35:0xc000bb33a0
			//36:0xc000bb2000 37:0xc0008bd680 38:0xc0008bc3a0 39:0xc000b26b40 40:0xc000b76d00 41:0xc000aacd40 42:0xc000a7f560
			//43:0xc000778260 44:0xc000a4fe40 55:0xc00000c140 69:0xc0002d8940 76:0xc0001f9120 82:0xc000233c00 83:0xc000940460
			//86:0xc00025a8e0 93:0xc00007d9c0 106:0xc000998f20 116:0xc0009c0380 118:0xc00032e000 122:0xc0009c1d00
			//208:0xc000a49420 231:0xc000e3b400 246:0xc000233020 247:0xc000268540 248:0xc000297260 250:0xc00032e900]

			//for the activearr part
			//for _,val :=range attendancelists{
			//
			//}

			//删除不必要的人。

			for k,_ :=range attendancelists{
				_,ok := activemap[k]
				if !ok {
					delete(attendancelists,k)
				}else{
					attendancelists[k].Userid = activemap[k]
				}
				//else{
				//	newattendancelists[k].Userid = activemap[k]
				//	newattendancelists[k].Userid = v.Userid
				//	newattendancelists[k].Signedtime = v.Signedtime
				//}

			}
			fmt.Println("++++++++++++++++")
			fmt.Println(attendancelists)

		}
	}
	return attendancelists
}
//type Attendancelists struct {
//	Attendances map[int64]*Attendancelist
//}


//ReadDingding is a function
func ReadDingding(activemap map[string]int64,excelFileName string,start string,end string) map[int64]*Attendancelist {
	//map [bns_169:169 bns_223:223 bns_241:241 bns_386:386 bns_391:391 bns_398:398 bns_419:419 bns_448:448 bns_453:453 bns_487:487 bns_494:494 bns_501:501 bns_502:502 bns_504:504 bns_505:505 bns_506:506]
	//excelFileName := "test.xlsx"
	// 根据指定的xlsx名称，返回xlsx文件结构
	//start,end := GetConfigStart("config.ini")
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	//start:="2019-09-23"
	//end:="2019-10-22"
	datearr :=PrDates(start,end)
	// xlFile.Sheets :一个xlsx文件可能有多个Sheet
	//pattern := "/^[0-9:]*$/"
	pattern := "[0-9]+:[0-9]+"
	attendancelists := make(map[int64]*Attendancelist)
	//newattendancelists := make(map[int64]*Attendancelist)
	for _, sheet := range xlFile.Sheets {
		// sheet.Rows:得到sheet的每一行
		if strings.Index(strings.Replace(sheet.Name, " ", "", -1), "打卡时间") > -1 {
			//fmt.Println(sheet.Row(2).Cells[2].Value)// 获取某一个单元格的值
			fmt.Println("+++++++++++++")
			var templastuserid int64 = 0
			tempsignedtime :=make(map[string]*Signedtimelist)
			for keyrow, row := range sheet.Rows {

				if keyrow<3 {
					continue
				}
				var i int64 =0
				attendancelist := &Attendancelist{}
				signedtime :=make(map[string]*Signedtimelist)
				// row.Cells:每一行的每一列
				//这个获取的templastuserid 是直接从方法里获取的
				tempval := sheet.Row(keyrow).Cells[4].Value
				lastuserid := activemap[tempval]

				fmt.Println("templastuserid is :",tempval,"lastuserid is:",lastuserid,"username is:",sheet.Row(keyrow).Cells[0].Value)
				log.Println("templastuserid is :",tempval,"lastuserid is:",lastuserid,"username is:",sheet.Row(keyrow).Cells[0].Value)
				attendancelist.Userid = lastuserid
				username := sheet.Row(keyrow).Cells[0].Value
				index :=strings.Index(strings.Replace(username, " ", "", -1), "（")
				if index>-1 {
					username = username[:index]
				}
				attendancelist.Username = username
				attendancelists[lastuserid] = attendancelist
				fmt.Println("poijkl:",attendancelists)

				for keycol, cell := range row.Cells {
					if keycol<5{
						continue
					}
					fmt.Println("------",keycol)
					val := cell.String()
					val = strings.Replace(val, " ", "", -1)
					// 去除换行符
					val = strings.Replace(val, "\n", "", -1)
					fmt.Println("value is :",val)
					//fmt.Print(text, "\t")







					signedtimelist := &Signedtimelist{}
					if val=="" || GetCompileData(val,pattern) {


						_, exists := datearr[i]
						if exists {//存在
							fmt.Println(val)
							fmt.Println("@@@@@")
							staticendtimestamp :=GetTimstamp(datearr[i]+" 04:00")
							staticstarttimestamp :=GetTimstamp(datearr[i]+" 16:00")
							fmt.Println("static end:",datearr[i]+" 04:00",",static start:",datearr[i]+" 16:00")

							firstsigned := ""
							lastsigned := ""
							if val!="" {
								signedarr := Getsplit(val)
								fmt.Println(signedarr,len(signedarr))
								if len(signedarr)>1 {
									firstsigned = signedarr[0]
									firsttimestamp :=GetTimstamp(datearr[i]+" "+firstsigned)

									if firsttimestamp<staticendtimestamp {
										_,ok := datearr[i-1]
										if ok {
											attendancelists[lastuserid].Signedtime[datearr[i-1]].Lastsigned = "23:59"
										}
										delete(signedarr,0)
										firstsigned = signedarr[0]
									}
									if firsttimestamp>staticstarttimestamp {
										firstsigned = ""
									}
									lastsigned = signedarr[int64(len(signedarr)-1)]
									lasttimestamp :=GetTimstamp(datearr[i]+" "+lastsigned)
									if lasttimestamp<staticstarttimestamp {
										lastsigned = ""
									}

								}
								if len(signedarr)==1 {

									tempsigned := signedarr[0]
									temptimestamp :=GetTimstamp(datearr[i]+" "+tempsigned)

									if temptimestamp<staticendtimestamp {
										_,ok := datearr[i-1]
										if ok {
											attendancelists[lastuserid].Signedtime[datearr[i-1]].Lastsigned = "23:59"
										}
										delete(signedarr,0)
									}
									if len(signedarr)==1 {

										//firstsigned = signedarr[0]
										tempsigned := signedarr[0]
										temptimestamp :=GetTimstamp(datearr[i]+" "+tempsigned)
										fmt.Println("firstsigned:",datearr[i]+" "+tempsigned,",lastsigned:",lastsigned,",date i is:",datearr[i])
										if temptimestamp<staticstarttimestamp {
											fmt.Println("ccccccc")
											firstsigned = tempsigned
										}else{
											fmt.Println("ddddddd")
											lastsigned = tempsigned
										}
										//if lastuserid==2 {
										//	fmt.Println("aaaaaaaaaaaaaaaaaaaa")
										//	fmt.Println(datearr[i],firstsigned,lastsigned)
										//}

									}

								}


							}else{
								fmt.Println("val is empty,row is:",keyrow,",column:",keycol)
							}
							fmt.Println("kkkk",lastuserid,"attendancelists:",attendancelists[lastuserid])
							fmt.Println("lllllll",datearr[i],",type is:",reflect.TypeOf(datearr[i]),"dawdaw:",attendancelists[lastuserid].Signedtime)

							//signedtimelist := &Signedtimelist{
							//	Firstsigned: firstsigned,
							//	Lastsigned:  lastsigned,
							//	Lackperiod:  "",
							//	Period:      "",
							//}
							signedtimelist.Firstsigned =firstsigned
							signedtimelist.Lastsigned = lastsigned
							signedtime[datearr[i]] = signedtimelist
							fmt.Println("signedtime:",signedtime)



							fmt.Println("kkkk",lastuserid,"attendancelists:",attendancelists)
							i++
						}else{
							fmt.Println("datearr 中缺少：",i)
						}
					}else{
						fmt.Println("adawdawdwdwa")
					}
				}
				if templastuserid!=lastuserid {
					templastuserid = lastuserid
					tempsignedtime = signedtime
				}else{
					for atk,atv :=range tempsignedtime {
						if atv.Firstsigned!=""{
							fmt.Println(atk,atv.Firstsigned)
							signedtime[atk].Firstsigned = atv.Firstsigned
						}
						if atv.Lastsigned!=""{
							fmt.Println(atk,atv.Lastsigned)
							signedtime[atk].Lastsigned = atv.Lastsigned
						}
					}

				}
				//if lastuserid==419 {
				//	fmt.Println("-----------------------------------------------------------------------------------------11")
				//
				//	fmt.Println(signedtime)
				//	data, err := json.Marshal(signedtime)
				//
				//	if err != nil {
				//		fmt.Println("json.marshal failed, err:", err)
				//	}
				//	fmt.Println(string(data))
				//	fmt.Println("-----------------------------------------------------------------------------------------11")
				//}


				attendancelists[lastuserid].Signedtime = signedtime

				fmt.Println()

			}
			fmt.Println("++++++++++++++++")
			fmt.Println(attendancelists)



			//map[1:0xc0008bcb00 2:0xc000483fe0 3:0xc000482b40 4:0xc00007c5e0
			//5:0xc000229280 6:0xc00000d600 7:0xc0007780a0 8:0xc000b768e0 9:0xc000941a60
			//10:0xc0002970e0 11:0xc000ff1420 12:0xc000fe6320 13:0xc000f9f160 14:0xc000f61980
			//15:0xc000f603e0 16:0xc00025bea0 17:0xc000f30820 18:0xc000cd7ca0 19:0xc000e85d80 20:0xc000e847e0 21:0xc000e173c0
			//22:0xc000dc3f00 23:0xc000dc2d80 24:0xc000d65a00 25:0xc000900e20 26:0xc000d52b40 27:0xc000d20ec0 28:0xc000d18400
			//29:0xc000ce8500 30:0xc000c93ae0 31:0xc000c92540 32:0xc000c5cf80 33:0xc000c2fb60 34:0xc000c1c8c0 35:0xc000bb33a0
			//36:0xc000bb2000 37:0xc0008bd680 38:0xc0008bc3a0 39:0xc000b26b40 40:0xc000b76d00 41:0xc000aacd40 42:0xc000a7f560
			//43:0xc000778260 44:0xc000a4fe40 55:0xc00000c140 69:0xc0002d8940 76:0xc0001f9120 82:0xc000233c00 83:0xc000940460
			//86:0xc00025a8e0 93:0xc00007d9c0 106:0xc000998f20 116:0xc0009c0380 118:0xc00032e000 122:0xc0009c1d00
			//208:0xc000a49420 231:0xc000e3b400 246:0xc000233020 247:0xc000268540 248:0xc000297260 250:0xc00032e900]

			//for the activearr part
			//for _,val :=range attendancelists{
			//
			//}

			//删除不必要的人。

			//for k,_ :=range attendancelists{
			//	_,ok := activemap[k]
			//	if !ok {
			//		delete(attendancelists,k)
			//	}else{
			//		attendancelists[k].Userid = activemap[k]
			//	}
			//	//else{
			//	//	newattendancelists[k].Userid = activemap[k]
			//	//	newattendancelists[k].Userid = v.Userid
			//	//	newattendancelists[k].Signedtime = v.Signedtime
			//	//}
			//
			//}

			if _,ok :=attendancelists[0];ok{
				delete(attendancelists,0)
			}
			fmt.Println("++++++++++++++++")
			fmt.Println(attendancelists)

		}
	}
	return attendancelists
}


func ReadMaomingv3(activemap map[int64]int64,excelFileName string,start string,end string) map[int64]*Attendancelist {
	//excelFileName := "maoming.xlsx"
	// 根据指定的xlsx名称，返回xlsx文件结构

	xlsx, err := excelize.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}



	t := time.Now()
	firstint := t.Year()
	secondint :=firstint/1000%10*10+firstint/100%10
	secondstr := strconv.Itoa(secondint)

	//str := "9/27/19 09:03"
	//firstIndex := strings.Index(str, "/")
	//lastIndex := strings.LastIndex(str, "/")
	//
	//ystr := secondstr+str[lastIndex+1:lastIndex+3]
	//yint, _ := strconv.Atoi(ystr)
	//if yint>firstint {
	//	ystr = strconv.Itoa(secondint-1)+str[lastIndex+1:lastIndex+3]
	//	yint, _ = strconv.Atoi(ystr)
	//}
	//dstr := str[3:lastIndex]
	//dint, _ := strconv.Atoi(dstr)
	//if dint<10 {
	//	dstr = "0"+dstr
	//}
	//
	//mstr := str[0:firstIndex]
	//mint, _ := strconv.Atoi(mstr)
	//if mint<10 {
	//	mstr = "0"+mstr
	//}
	//leftstr := str[lastIndex+3:]
	//newstr := ystr+"-"+mstr+"-"+dstr+leftstr
	//fmt.Println(newstr)


	//1999 2000

	//str :="4/8/2014 22:05"
	//dateparse.Parsedate(newstr)





	//fmt.Println(datearr)

	//attendancelists := make(map[int64]*Attendancelist)
	//
	//attendancelist := &Attendancelist{}
	collect := make(map[int64][]string)

	user := make(map[int64]string)

	//attendancelists := make(map[int64]*Attendancelist)

	// xlFile.Sheets :一个xlsx文件可能有多个Sheet
	oldsheet :="grid_hr_at_card_rec"


	//testmap := make(map[int64][]string)
	rows ,_:= xlsx.GetRows(oldsheet)

	for keyrow, row := range rows {
		if keyrow<2{
			continue
		}
		//for index 是从0开始计算，但是setcellstyle index 是从1开始，所以总是会大一


		col,_:=excelize.ColumnNumberToName(5)
		style, _ := xlsx.NewStyle(`{"number_format": 22}`)
		xlsx.SetCellStyle(oldsheet, col+strconv.Itoa(keyrow+1), col+strconv.Itoa(keyrow+1), style)
		formatday ,_:=xlsx.GetCellValue(oldsheet, col+strconv.Itoa(keyrow+1))
		str := formatday
		//fmt.Println("----------",str)

		firstIndex := strings.Index(str, "/")
		lastIndex := strings.LastIndex(str, "/")

		ystr := secondstr+str[lastIndex+1:lastIndex+3]
		yint, _ := strconv.Atoi(ystr)
		if yint>firstint {
			ystr = strconv.Itoa(secondint-1)+str[lastIndex+1:lastIndex+3]
			yint, _ = strconv.Atoi(ystr)
		}
		dstr := str[firstIndex+1:lastIndex]
		dint, _ := strconv.Atoi(dstr)
		if dint<10 {
			dstr = "0"+dstr
		}

		mstr := str[0:firstIndex]
		mint, _ := strconv.Atoi(mstr)
		if mint<10 {
			mstr = "0"+mstr
		}
		leftstr := str[lastIndex+3:] //  20:23
		datestr := ystr+"-"+mstr+"-"+dstr //2019-10-24
		newstr := datestr+leftstr
		//filterleftstr := strings.Replace(leftstr, " ", "", -1)

		//staticendtimestamp :=GetTimstamp(datestr+" 04:00")
		//staticstarttimestamp :=GetTimstamp(datestr+" 16:00")
		//fmt.Println("static end:",datestr+" 04:00",staticendtimestamp,",static start:",datestr+" 16:00",staticstarttimestamp)
		//newstrtimestamp :=GetTimstamp(newstr)

		//fmt.Println("-------",keyrow,row[0])
		userid, _ := strconv.ParseInt(row[0], 10, 64)

		_, ok := activemap[userid]

		if ok {
			user[activemap[userid]] = row[1]
			collect[activemap[userid]] = append(collect[activemap[userid]],newstr)

			//if newstrtimestamp>staticstarttimestamp {
			//	if _,exist :=collect[activemap[userid]];exist{
			//		collect[activemap[userid]][datestr].Lastsigned = filterleftstr
			//	}else{
			//		c :=make(map[string]*Signedtimelist)
			//		c[datestr].Lastsigned = filterleftstr
			//		collect[activemap[userid]] = c
			//	}
			//}else if newstrtimestamp<staticendtimestamp {
			//	newdatestr := GetSomeday(datestr,1,"day",false)
			//	if _,exist :=collect[activemap[userid]];exist{
			//		collect[activemap[userid]][datestr].Lastsigned = newdatestr+" 23:59"
			//	}else{
			//		c :=make(map[string]*Signedtimelist)
			//		c[datestr].Lastsigned = newdatestr+" 23:59"
			//		collect[activemap[userid]] = c
			//	}
			//
			//}else {
			//	if _,exist :=collect[activemap[userid]];exist{
			//		collect[activemap[userid]][datestr].Firstsigned = filterleftstr
			//	}else{
			//		c :=make(map[string]*Signedtimelist)
			//		c[datestr].Firstsigned = filterleftstr
			//		collect[activemap[userid]] = c
			//	}
			//}



		}else{
			fmt.Printf("你要找的资料不存在。",userid)
		}

	}

	fmt.Println("collect is:",collect[241])
	//start,end := GetConfigStart("config.ini")
	//start := "2019-09-23" //待转化为时间戳的字符串
	//end := "2019-10-22"   //待转化为时间戳的字符串
	//日期转化为时间戳
	datearr := PrDates(start,end)


	attendancelists := make(map[int64]*Attendancelist)
	for dk,dv :=range collect{
		attendancelist := &Attendancelist{}

		signedtime :=make(map[string]*Signedtimelist)
		for _,v :=range datearr {
			signedtimelist := &Signedtimelist{}// 这个特别重要
			signedtimelist.Firstsigned = ""
			signedtimelist.Lastsigned = ""
			signedtime[v] = signedtimelist
		}

		attendancelist.Userid = dk
		attendancelist.Username = user[dk]
		attendancelist.Signedtime = signedtime



		for _,v :=range dv{
			da := v[:10]
			signed := strings.Replace(v[10:], " ", "", -1)
			if _,ok := attendancelist.Signedtime[da];ok{
				damiddle := GetTimstamp(da+" 00:00")
				yesterdayend := GetTimstamp(da+" 04:00")
				dastart := GetTimstamp(da+" 16:00")
				daend := GetTimstamp(da+" 23:59")
				thistime :=GetTimstamp(v)
				yesterday := GetSomeday(da,1,"day",false)

				if thistime>=damiddle && thistime<=yesterdayend {
					//yesterday
					attendancelist.Signedtime[yesterday].Lastsigned = "23:59"
				}else if thistime>yesterdayend && thistime<=dastart {
					//morning
					if attendancelist.Signedtime[da].Firstsigned !="" {
						firstsigned := GetTimstamp(da+" "+attendancelist.Signedtime[da].Firstsigned)
						if thistime<firstsigned {
							attendancelist.Signedtime[da].Firstsigned = signed
						}
					}else{
						attendancelist.Signedtime[da].Firstsigned = signed
					}
				}else if thistime>dastart && thistime<=daend {
					if attendancelist.Signedtime[da].Lastsigned !="" {
						lastsigned := GetTimstamp(da+" "+attendancelist.Signedtime[da].Lastsigned)
						if thistime>lastsigned {
							attendancelist.Signedtime[da].Lastsigned = signed
						}
					}else{
						attendancelist.Signedtime[da].Lastsigned = signed
					}
				}


			}
		}
		//if dk==398 {
		//	xdata, err := json.Marshal(attendancelist)
		//
		//	if err != nil {
		//		fmt.Println("json.marshal failed, err:", err)
		//		return nil
		//	}
		//	fmt.Println("#################")
		//	fmt.Println("---------------",string(xdata))
		//}


		attendancelists[dk] = attendancelist

	}



	return attendancelists
}

// DealExcel is a function
// DealExcel is a function
func DealExcel(configmap map[string]interface{}) map[int64]*Attendancelist{
	start :=configmap["start"].(string)
	end := configmap["end"].(string)
	//configmap :=InitConfig("config.ini")
	changhuapeoplemap := configmap["newchanghuapeople"]
	pudongflip :=configmap["pudongflip"]
	fmt.Println(changhuapeoplemap)
	perdayhour := configmap["perdayhour"]
	perdayhour, _ = strconv.ParseFloat(perdayhour.(string), 64)

	//changhuapeoplemap is interface
	newchanghuapeoplemap := make(map[int64]int64)
	for k, v := range changhuapeoplemap.(map[int64]int64) {
		newchanghuapeoplemap[k] = v
	}
	pudongflipmap := make(map[int64]int64)
	pudongpeoplemap := make(map[int64]int64)
	for k, v := range pudongflip.(map[int64]int64) {
		pudongpeoplemap[v] = k
		pudongflipmap[k] = v
	}
	dingdingsec :=configmap["dingdingpeople"]
	dingdingpeople := make(map[string]int64)
	for k,v :=range dingdingsec.(map[string]string) {
		gv, _ := strconv.ParseInt(v, 10, 64)
		dingdingpeople[k] = gv
	}

	fmt.Println("new changehua people is:",newchanghuapeoplemap)
	fmt.Println("new pudong people is:",pudongflip)
	maomingpeoplemap := configmap["newmaomingpeople"]

	inputfolder :=configmap["folder"].(map[string]string)["input"]

	oldchanghua := make(map[int64]*Attendancelist)
	oldpudong := make(map[int64]*Attendancelist)
	oldmaoming := make(map[int64]*Attendancelist)
	olddingding := make(map[int64]*Attendancelist)


	if _,ok :=configmap["changhuaexcelfile"];ok{

		changhuaexcelfile :=inputfolder+"/"+configmap["changhuaexcelfile"].(string)
		oldchanghua =ReadTealegXlsx(newchanghuapeoplemap,changhuaexcelfile,start,end)
		if _,ok :=configmap["pudongexcelfile"];ok{
			pudongexcelfile :=inputfolder+"/"+configmap["pudongexcelfile"].(string)
			oldpudong =ReadTealegXlsx(pudongpeoplemap,pudongexcelfile,start,end)
		}
		if _,ok :=configmap["maomingexcelfile"];ok{
			maomingexcelfile :=inputfolder+"/"+configmap["maomingexcelfile"].(string)
			oldmaoming =ReadMaomingv3(maomingpeoplemap.(map[int64]int64),maomingexcelfile,start,end)
		}
		if _,ok :=configmap["dingdingexcelfile"];ok{
			fmt.Println("-----------11111111---------")

			dingdingexcelfile :=inputfolder+"/"+configmap["dingdingexcelfile"].(string)
			olddingding =ReadDingding(dingdingpeople,dingdingexcelfile,start,end)
			log.Println("------ddddd-----",olddingding)
		}
	}else {
		if _,ok :=configmap["pudongexcelfile"];ok{
			pudongexcelfile :=inputfolder+"/"+configmap["pudongexcelfile"].(string)
			oldchanghua =ReadTealegXlsx(pudongpeoplemap,pudongexcelfile,start,end)
			if _,ok :=configmap["maomingexcelfile"];ok{
				maomingexcelfile :=inputfolder+"/"+configmap["maomingexcelfile"].(string)
				oldmaoming =ReadMaomingv3(maomingpeoplemap.(map[int64]int64),maomingexcelfile,start,end)
			}
			if _,ok :=configmap["dingdingexcelfile"];ok{
				dingdingexcelfile :=inputfolder+"/"+configmap["dingdingexcelfile"].(string)
				olddingding =ReadDingding(dingdingpeople,dingdingexcelfile,start,end)
			}
		}else{
			if _,ok :=configmap["maomingexcelfile"];ok{
				maomingexcelfile :=inputfolder+"/"+configmap["maomingexcelfile"].(string)
				oldchanghua =ReadMaomingv3(maomingpeoplemap.(map[int64]int64),maomingexcelfile,start,end)
				if _,ok :=configmap["dingdingexcelfile"];ok{
					dingdingexcelfile :=inputfolder+"/"+configmap["dingdingexcelfile"].(string)
					olddingding =ReadDingding(dingdingpeople,dingdingexcelfile,start,end)
				}
			}else{
				if _,ok :=configmap["dingdingexcelfile"];ok{
					dingdingexcelfile :=inputfolder+"/"+configmap["dingdingexcelfile"].(string)
					oldchanghua =ReadDingding(dingdingpeople,dingdingexcelfile,start,end)
				}else{
					log.Println("something wrong ,no existed file")
					return nil
				}
			}


		}
	}

	changhua := make(map[int64]*Attendancelist)
	pudong := make(map[int64]*Attendancelist)
	maoming := make(map[int64]*Attendancelist)
	dingding := make(map[int64]*Attendancelist)
	for _,v :=range oldchanghua {
		changhua[v.Userid] = v
	}

	for _,v :=range oldpudong {
		pudong[v.Userid] = v
	}
	for _,v :=range oldmaoming {
		maoming[v.Userid] = v
	}
	for _,v :=range olddingding {
		dingding[v.Userid] = v
	}






	fmt.Println("ppppppppppppppppppppp")
	fmt.Println("changhua",changhua,"pudong",pudong,"maoming",maoming,"dingding",dingding)

	fmt.Println("pudong flip map",pudongflipmap)
	fmt.Println("maoming flip map",maomingpeoplemap)


	//xdata, err := json.Marshal(pudong)
	//
	//if err != nil {
	//	fmt.Println("json.marshal failed, err:", err)
	//	return
	//}
	//fmt.Println("#################")
	//fmt.Println(string(xdata))


	//[2:0xc0007a8930 10:0xc000836de0 14:0xc0002104b0 16:0xc000d052f0 28:0xc000e94360 29:0xc000e15500 30:0xc000ddd800 32:0xc000d83980 37:0xc000a6e870 38:0xc000a3b7d0 39:0xc000c395f0 40:0xc000a28750 41:0xc000bf3ec0 42:0xc000bf49f0 44:0xc000bbb470 55:0xc00081b680 76:0xc0007a9170 82:0xc00077c4e0 83:0xc000956f60 93:0xc000211440 116:0xc000b5ae10 118:0xc000210030]

	userinfo :=configmap["userinfo"].(map[int64]string)

	//pudong
	for pkey,pvalue :=range changhua {
		//fmt.Println("sasasasa",pvalue.Userid,pvalue.Username)
		_,ok := pudongflipmap[pvalue.Userid]
		if ok {
			if pudong[pvalue.Userid]!=nil {
				for kp,vp :=range pvalue.Signedtime{


					if (vp.Firstsigned=="" && pudong[pvalue.Userid].Signedtime[kp].Firstsigned!="") || (pudong[pvalue.Userid].Signedtime[kp].Firstsigned!="" && vp.Firstsigned!="" && GetTimstamp(vp.Firstsigned)>GetTimstamp(pudong[pvalue.Userid].Signedtime[kp].Firstsigned)) {

						changhua[pkey].Signedtime[kp].Firstsigned = pudong[pvalue.Userid].Signedtime[kp].Firstsigned
					}
					if (vp.Lastsigned=="" && pudong[pvalue.Userid].Signedtime[kp].Lastsigned!="") || (pudong[pvalue.Userid].Signedtime[kp].Lastsigned!="" && vp.Lastsigned!="" && GetTimstamp(vp.Lastsigned)<GetTimstamp(pudong[pvalue.Userid].Signedtime[kp].Lastsigned)) {
						changhua[pkey].Signedtime[kp].Lastsigned = pudong[pvalue.Userid].Signedtime[kp].Lastsigned
					}
				}
			}

		}
	}
	//pudong left people
	for _,pvalue :=range pudong {

		if _,eok := changhua[pvalue.Userid];!eok{
			_,ok := userinfo[pvalue.Userid]
			if ok {
				if pudong[pvalue.Userid]!=nil {
					newattendancelist := &Attendancelist{
						Userid:     pvalue.Userid,
						Username:   pvalue.Username,
						Signedtime: pvalue.Signedtime,
						All:        pvalue.All,
					}
					changhua[pvalue.Userid] = newattendancelist

				}

			}
		}

	}



	//maoming
	for pkey,pvalue :=range changhua {
		//fmt.Println("sasasasa",pvalue.Userid,pvalue.Username)
		if maoming[pvalue.Userid]!=nil {
			for kp,vp :=range pvalue.Signedtime{


				if (vp.Firstsigned=="" && maoming[pvalue.Userid].Signedtime[kp].Firstsigned!="") || (maoming[pvalue.Userid].Signedtime[kp].Firstsigned!="" && vp.Firstsigned!="" && GetTimstamp(vp.Firstsigned)>GetTimstamp(maoming[pvalue.Userid].Signedtime[kp].Firstsigned)) {

					changhua[pkey].Signedtime[kp].Firstsigned = maoming[pvalue.Userid].Signedtime[kp].Firstsigned
				}
				if (vp.Lastsigned=="" && maoming[pvalue.Userid].Signedtime[kp].Lastsigned!="") || (maoming[pvalue.Userid].Signedtime[kp].Lastsigned!="" && vp.Lastsigned!="" && GetTimstamp(vp.Lastsigned)<GetTimstamp(maoming[pvalue.Userid].Signedtime[kp].Lastsigned)) {
					changhua[pkey].Signedtime[kp].Lastsigned = maoming[pvalue.Userid].Signedtime[kp].Lastsigned
				}
			}
		}
	}

	//maoming left people
	for _,pvalue :=range maoming {
		if _,eok := changhua[pvalue.Userid];!eok{
			_,ok := userinfo[pvalue.Userid]
			if ok {
				if maoming[pvalue.Userid]!=nil {
					newattendancelist := &Attendancelist{
						Userid:     pvalue.Userid,
						Username:   pvalue.Username,
						Signedtime: pvalue.Signedtime,
						All:        pvalue.All,
					}
					changhua[pvalue.Userid] = newattendancelist
				}

			}
		}

	}

	//dingding
	for pkey,pvalue :=range changhua {
		//fmt.Println("sasasasa",pvalue.Userid,pvalue.Username)
		if dingding[pvalue.Userid]!=nil {
			for kp,vp :=range pvalue.Signedtime{


				if (vp.Firstsigned=="" && dingding[pvalue.Userid].Signedtime[kp].Firstsigned!="") || (dingding[pvalue.Userid].Signedtime[kp].Firstsigned!="" && vp.Firstsigned!="" && GetTimstamp(vp.Firstsigned)>GetTimstamp(dingding[pvalue.Userid].Signedtime[kp].Firstsigned)) {

					changhua[pkey].Signedtime[kp].Firstsigned = dingding[pvalue.Userid].Signedtime[kp].Firstsigned
				}
				if (vp.Lastsigned=="" && dingding[pvalue.Userid].Signedtime[kp].Lastsigned!="") || (dingding[pvalue.Userid].Signedtime[kp].Lastsigned!="" && vp.Lastsigned!="" && GetTimstamp(vp.Lastsigned)<GetTimstamp(dingding[pvalue.Userid].Signedtime[kp].Lastsigned)) {
					changhua[pkey].Signedtime[kp].Lastsigned = dingding[pvalue.Userid].Signedtime[kp].Lastsigned
				}
			}
		}
	}
	//dingding left people
	for _,pvalue :=range dingding {
		if _,eok := changhua[pvalue.Userid];!eok{
			_,ok := userinfo[pvalue.Userid]
			if ok {
				if dingding[pvalue.Userid]!=nil {
					newattendancelist := &Attendancelist{
						Userid:     pvalue.Userid,
						Username:   pvalue.Username,
						Signedtime: pvalue.Signedtime,
						All:        pvalue.All,
					}
					changhua[pvalue.Userid] = newattendancelist
				}

			}
		}

	}
	datearr :=PrDates(start,end)
	//check someone if it is existed in changhua,just like pudong 徐小玲
	for k,v :=range userinfo{
		if _,ok :=changhua[k];!ok{
			signedtime :=make(map[string]*Signedtimelist)
			signedtimelist :=&Signedtimelist{
				Firstsigned: "",
				Lastsigned:  "",
				Lackperiod:  "",
				Period:      "",
			}
			for _,datev :=range datearr{
				signedtime[datev] = signedtimelist
			}

			newattendancelist := &Attendancelist{
				Userid:     k,
				Username:   v,
				Signedtime: signedtime,
				All:        0,
			}
			changhua[k] = newattendancelist
		}
	}


	fmt.Println("adawdadawdada",perdayhour)

	for k,v :=range changhua {
		all :=0.00
		for k1,v1 :=range v.Signedtime {
			hour :=0
			minute :=0
			timearr := Timediff(k1+" "+v1.Firstsigned, k1+" "+v1.Lastsigned)
			if v1.Firstsigned==""||v1.Lastsigned=="" {
				timearr["day"] = 0
				timearr["hour"] = 0
				timearr["min"] = 0
				timearr["sec"] = 0
			}
			hour = timearr["hour"]
			minute = timearr["min"]
			perall :=Round(float64(hour)+float64(minute)/float64(60),2)
			all+=perall
			all = Round(all,2)
			diff := Calper(perdayhour.(float64), float64(int64(hour)), float64(int64(minute)))
			changhua[k].Signedtime[k1].Lackperiod = diff["diffhour"]+"h"+diff["diffmin"]+"m"
			changhua[k].Signedtime[k1].Period = strconv.Itoa(hour)+"h"+strconv.Itoa(minute)+"m"

		}
		changhua[k].All = all
	}


	//data, err := json.Marshal(changhua)
	//
	//if err != nil {
	//	fmt.Println("json.marshal failed, err:", err)
	//	return
	//}
	//fmt.Println("#################")
	//fmt.Println(string(data))
	return changhua

}

type Attendancelist struct {
	Userid       int64 `json:"userid"`
	Username string `json:"username"`
	Signedtime 	map[string]*Signedtimelist
	All float64 `all`
}
type Signedtimelist struct {
	//Signedday string `json:signedday`
	//Signedlist struct{
	//	Firstsigned string `firstsigned`
	//	Lastsigned string `lastsigned`
	//	Lackperiod string `lackperiod`
	//	Period string `period`
	//}
	Firstsigned string `firstsigned`
	Lastsigned string `lastsigned`
	Lackperiod string `lackperiod`
	Period string `period`
}