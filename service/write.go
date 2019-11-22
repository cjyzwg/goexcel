package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strconv"
)

//WriteInExcel is a function
func WriteInExcel() {
	xlsx := excelize.NewFile()

	index := xlsx.NewSheet("Sheet1")
	xlsx.SetCellValue("Sheet1", "A1", "姓名")
	xlsx.SetCellValue("Sheet1", "B1", "年龄")
	xlsx.SetCellValue("Sheet1", "A2", "狗子")
	xlsx.SetCellValue("Sheet1", "B2", "18")
	xlsx.SetCellValue("Sheet1", "A3", "胆子")
	xlsx.SetCellValue("Sheet1", "B3", "20")
	xlsx.SetCellValue("Sheet1", "A4", "wda")
	xlsx.SetCellValue("Sheet1", "B4", "33")
	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	err := xlsx.SaveAs("test_write.xls")
	if err != nil {
		fmt.Println(err)
	}
}

// CreateTealegXlsx 创建excel

func CreateTealegXlsx() {
	// 创建一个xlsx文件
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row, row1, row2 *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.SetHeightCM(1)
	cell = row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "年龄"

	row1 = sheet.AddRow()
	row1.SetHeightCM(1)
	cell = row1.AddCell()
	cell.Value = "狗子"
	cell = row1.AddCell()
	cell.Value = "18"

	row2 = sheet.AddRow()
	row2.SetHeightCM(1)
	cell = row2.AddCell()
	cell.Value = "蛋子"
	cell = row2.AddCell()
	cell.Value = "28"

	err = file.Save("test_write.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

var (
	outFile = "./out_student.xlsx"
)



func Export(attendancelists map[int64]*Attendancelist,outfile string){
	//attendancelists map[int64]*Attendancelist
	//headarr :=[]string{"D","E","F","G",'H','I','J','K','L','M', 'N','O','P','Q','R','S','T','U','V','W','X','Y','Z','AA','AB','AC','AD','AE','AF','AG','AH'}
	configmap :=InitConfig("config.ini")
	gouserorderslice :=configmap["gouserorderslice"]
	special_regulation :=configmap["special_regulation"]
	specialarr :=make(map[int64]float64)
	for k,v :=range special_regulation.(map[string]string) {
		gk, _ := strconv.ParseInt(k, 10, 64)
		gv,_ := strconv.ParseFloat(v,64)
		specialarr[gk] = gv
	}
	start :=configmap["start"].(string)
	end := configmap["end"].(string)


	xlsx := excelize.NewFile()

	oldsheet :="Sheet1"
	newsheet :="考勤"

	index :=xlsx.GetSheetIndex(oldsheet)
	attendancelistslen :=len(attendancelists)*7


	//index := xlsx.NewSheet(newsheet)
	//设置字体自动换行
	style, err := xlsx.NewStyle(`{"alignment":{"horizontal":"center","wrap_text":true}}`)
	if err != nil {
		fmt.Println(err)
	}
	xlsx.SetCellStyle(oldsheet, "A1", "AH"+strconv.Itoa(attendancelistslen), style)



	xlsx.SetRowHeight(oldsheet,1,30)
	xlsx.SetRowHeight(oldsheet,2,30)
	xlsx.SetColWidth(oldsheet,"A","AH",7)
	xlsx.SetCellValue(oldsheet, "A1", "考勤日期")
	xlsx.SetCellValue(oldsheet, "B1", start)
	//style, err := xlsx.NewStyle(`{"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"vertical":"","wrap_text":true}}`)

	xlsx.SetCellValue(oldsheet, "C1", "2019-10-22")
	xlsx.SetCellValue(oldsheet, "F1", "未打卡")
	style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#00ffe9"],"pattern":1}}`)
	if err != nil {
		fmt.Println(err)
	}
	xlsx.SetCellStyle(oldsheet, "G1", "G1", style)

	xlsx.SetCellValue(oldsheet, "I1", "跨月调休")
	style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#836FFF"],"pattern":1}}`)
	if err != nil {
		fmt.Println(err)
	}
	xlsx.SetCellStyle(oldsheet, "J1", "J1", style)
	xlsx.SetCellValue(oldsheet, "L1", "周六日")

	style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#DAA520"],"pattern":1}}`)
	if err != nil {
		fmt.Println(err)
	}
	xlsx.SetCellStyle(oldsheet, "M1", "M1", style)



	headslice := []string{"D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","AA","AB","AC","AD","AE","AF","AG","AH"}
	datearr :=PrDates(start,end)


	i:=2
	userorderslice := gouserorderslice.([]int64)
	for _,userv :=range  userorderslice {
		//fmt.Println("这个user 是存在的,",userv)
		_, exists := attendancelists[userv]
		if exists {
			fmt.Println("这个user 是存在的,",userv)
			//for _,attendancelist :=range attendancelists {
			//
			//}
			attendancelist := attendancelists[userv]
			//标注颜色
			style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#0000FF"],"pattern":1}}`)
			if err != nil {
				fmt.Println(err)
			}
			xlsx.SetCellStyle(oldsheet, "A"+strconv.Itoa(i), "AH"+strconv.Itoa(i), style)

			attendancedays :=0
			for k,v :=range datearr {

				weekend :=false
				num :=ZellerFunction2Week(v)
				if num==6 || num ==7 {
					weekend = true
				}
				xlsx.SetRowHeight(oldsheet,i+1,30)
				axis := headslice[k]+strconv.Itoa(i+1)
				_, ok := specialarr[userv]
				if ok && len(specialarr)>0 && weekend {
					style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#DAA520"],"pattern":1},"alignment":{"horizontal":"center","wrap_text":true}}`)
					if err != nil {
						fmt.Println(err)
					}
					xlsx.SetCellStyle(oldsheet, axis, axis, style)
				}

				xlsx.SetCellValue(oldsheet,axis , v)

				signed :=attendancelist.Signedtime[v]


				firstaxis := headslice[k]+strconv.Itoa(i+3)
				xlsx.SetCellValue(oldsheet,firstaxis , signed.Firstsigned)
				lastaxis := headslice[k]+strconv.Itoa(i+4)
				xlsx.SetCellValue(oldsheet,lastaxis , signed.Lastsigned)
				periodaxis := headslice[k]+strconv.Itoa(i+5)
				xlsx.SetCellValue(oldsheet,periodaxis , signed.Period)
				//00ffe9

				if signed.Firstsigned=="" {
					style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#00ffe9"],"pattern":1}}`)
					if err != nil {
						fmt.Println(err)
					}
					xlsx.SetCellStyle(oldsheet, firstaxis, firstaxis, style)
				}
				if signed.Lastsigned=="" {
					style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#00ffe9"],"pattern":1}}`)
					if err != nil {
						fmt.Println(err)
					}
					xlsx.SetCellStyle(oldsheet, lastaxis, lastaxis, style)
				}
				if !(signed.Firstsigned=="" && signed.Lastsigned=="") {
					attendancedays++
				}
			}

			xlsx.SetCellValue(oldsheet, "A"+strconv.Itoa(i+2), "姓名")
			xlsx.SetCellValue(oldsheet, "B"+strconv.Itoa(i+2), attendancelist.Username)
			xlsx.SetCellValue(oldsheet, "C"+strconv.Itoa(i+2), attendancelist.Userid)
			xlsx.SetCellValue(oldsheet, "A"+strconv.Itoa(i+3), "上班打卡")
			xlsx.SetCellValue(oldsheet, "A"+strconv.Itoa(i+4), "下班打卡")
			xlsx.SetCellValue(oldsheet, "A"+strconv.Itoa(i+5), "时长")
			xlsx.SetCellValue(oldsheet, "A"+strconv.Itoa(i+6), "总工时")
			xlsx.SetCellValue(oldsheet, "B"+strconv.Itoa(i+6),  attendancelist.All)
			xlsx.SetCellValue(oldsheet, "C"+strconv.Itoa(i+6), "出勤天数")
			xlsx.SetCellValue(oldsheet, "D"+strconv.Itoa(i+6), attendancedays)
			style, err = xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FF0000"],"pattern":1}}`)
			if err != nil {
				fmt.Println(err)
			}
			xlsx.SetCellStyle(oldsheet, "D"+strconv.Itoa(i+6), "D"+strconv.Itoa(i+6), style)
			xlsx.SetCellValue(oldsheet, "E"+strconv.Itoa(i+6), "未打卡次数")
			xlsx.SetCellValue(oldsheet, "F"+strconv.Itoa(i+6), 0)
			xlsx.SetCellValue(oldsheet, "I"+strconv.Itoa(i+6), "30m上")
			xlsx.SetCellValue(oldsheet, "J"+strconv.Itoa(i+6), 0)
			xlsx.SetCellValue(oldsheet, "L"+strconv.Itoa(i+6), "15上")
			xlsx.SetCellValue(oldsheet, "M"+strconv.Itoa(i+6), 0)
			xlsx.SetCellValue(oldsheet, "O"+strconv.Itoa(i+6), "15m下")
			xlsx.SetCellValue(oldsheet, "P"+strconv.Itoa(i+6), 0)
			i = i+7


		}
	}

















	// Set active sheet of the workbook.
	xlsx.SetSheetName(oldsheet,newsheet)
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	log.Println("outfile is:",outfile)
	err = xlsx.SaveAs(outfile)
	if err != nil {
		fmt.Println(err)
	}
}



func WriteResult(attendancelists map[int64]*Attendancelist, outfile string) error {
	data, err := json.Marshal(attendancelists)

	if err != nil {
		fmt.Println("json.marshal failed, err:", err)
		return nil
	}
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println("writer",err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(string(data))
	writer.WriteString("\n")
	writer.Flush()


	return err
}
