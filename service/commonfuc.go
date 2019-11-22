package service

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//GetConfigStart is a functionn
func GetConfigStart(filename string) (s string,e string){
	configmap :=InitConfig(filename)
	start :=configmap["start"].(string)
	end := configmap["end"].(string)
	return start,end
}

// PrDates 列出所有日期
func PrDates(start string, end string) map[int64]string {
	log.Println("1")
	//日期转化为时间戳
	timeLayout := "2006-01-02"
	var loc = time.FixedZone("CST", 8*3600)       // 东八
	if runtime.GOOS !="windows" {
		loc, _ = time.LoadLocation("Asia/Shanghai") //获取时区
	}

	log.Println("2",loc,runtime.GOOS)
	starttmp, _ := time.ParseInLocation(timeLayout, start, loc)
	starttimestamp := starttmp.Unix()
	endtmp, _ := time.ParseInLocation(timeLayout, end, loc)
	log.Println("3",starttimestamp)
	endtimestamp := endtmp.Unix()
	log.Println("4",endtimestamp)
	//var slice []string

	slicemap := make(map[int64]string)
	var i int64 = 0
	for {
		startdatetime := time.Unix(starttimestamp, 0).Format(timeLayout)
		starttimestamp += 24 * 60 * 60
		//slice = append(slice, startdatetime)
		slicemap[i] = startdatetime
		i++
		if starttimestamp > endtimestamp {
			break
		}
	}
	log.Println("5")
	//fmt.Println(slice)
	return slicemap
}
//GetCompileData is a compare function
func GetCompileData(searchIn string,pat string) bool{
	fmt.Println("before in ge compile data",searchIn)
	var b bool
	b = false
	//searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	//pat := "[0-9]+.[0-9]+"          //正则
	//f := func(s string) string{
	//	v, _ := strconv.ParseFloat(s, 32)
	//	return strconv.FormatFloat(v * 2, 'f', 2, 32)
	//}
	if ok, _ := regexp.Match(pat, []byte(searchIn)); ok {
		fmt.Println("Match Found!")
		b = true
	}
	fmt.Println("after in ge compile data",searchIn)
	return b
	//re, _ := regexp.Compile(pat)
	////将匹配到的部分替换为 "##.#"
	//str := re.ReplaceAllString(searchIn, "##.#")
	//fmt.Println(str)
	////参数为函数时
	//str2 := re.ReplaceAllStringFunc(searchIn, f)
	//fmt.Println(str2)
}

// Getsplit is a split time function string:"10:2312:0922:10"
func Getsplit(s string) map[int64]string{
	fmt.Println("get split from common function ",s)
	mapper :=make(map[int64]string)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)

	var i int64=0
	for {
		s1 :=s[:5]
		//fmt.Println("s1字符串类型是：",reflect.TypeOf(s1))
		mapper[i] = s1
		s = s[5:]
		i++
		if len(s)<=0 {
			break
		}
	}
	return mapper
}

// GetTimstamp is a function
func GetTimstamp(start string)  int64 {
	timeLayout := "2006-01-02 15:04"
	var loc = time.FixedZone("CST", 8*3600)       // 东八
	if runtime.GOOS !="windows" {
		loc, _ = time.LoadLocation("Asia/Shanghai") //获取时区
	}

	if len(start)==10 {
		start = start+" 00:00"
	}
	starttmp, _ := time.ParseInLocation(timeLayout, start, loc)
	starttimestamp := starttmp.Unix()
	return starttimestamp
}

//GetPreday is a function
func GetSomeday(start string,duration int,perstr string,addflag bool) string {
	starttimestamp :=GetTimstamp(start)
	durationtimestamp :=0
	switch perstr {
		case "year": durationtimestamp = 86400*365*duration
		case "month": durationtimestamp = 86400*30*duration
		case "day" : durationtimestamp = 86400*duration
		case "hour" : durationtimestamp = 3600*duration
		case "min" : durationtimestamp = 60*duration
		default: durationtimestamp = duration
	}
	if addflag {
		starttimestamp += int64(durationtimestamp)
	}else{
		starttimestamp -= int64(durationtimestamp)
	}
	t := time.Unix(starttimestamp, 0)
	dateStr := t.Format("2006-01-02 15:04:05")
	return  dateStr
}


//var weekday = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
var weekday = [7]int{7,1,2,3,4,5,6}
//ZellerFunction2Week 计算日期是周几
func ZellerFunction2Week(start string) int {

	s :=strings.Split(start,"-")
	yearint, _ := strconv.Atoi(s[0])
	monthint, _ := strconv.Atoi(s[1])
	dayint, _ := strconv.Atoi(s[2])
	year := uint16(yearint)
	month := uint16(monthint)
	day := uint16(dayint)
	//fmt.Println(year,month,day)



	var y, m, c uint16
	if month >= 3 {
		m = month

		y = year % 100

		c = year / 100

	} else {

		m = month + 12

		y = (year - 1) % 100

		c = (year - 1) / 100

	}

	week := y + (y / 4) + (c / 4) - 2*c + ((26 * (m + 1)) / 10) + day - 1

	if week < 0 {

		week = 7 - (-week)%7

	} else {

		week = week % 7

	}

	which_week := int(week)

	return weekday[which_week]

}


// Calper is a function
func Calper(perdayhour float64,hour float64,min float64) map[string]string{
	arrmap :=make(map[string]string)
	diffhour := hour-perdayhour
	diffmin := min
	if diffhour<0 && min>0 {
		diffhour = diffhour+1
		diffmin = 60-min
	}

	//arrmap["diffhour"] = fmt.Sprintf("%0.2f", diffhour)
	//arrmap["diffmin"] = fmt.Sprintf("%0.2f", diffmin)
	arrmap["diffhour"] = strconv.Itoa(int(diffhour))
	arrmap["diffmin"] = strconv.Itoa(int(diffmin))
	return arrmap
}

// Timediff is a function
func Timediff(start string,end string) map[string]int{
	tempstarttimestamp :=GetTimstamp(start)
	tempendtimestamp :=GetTimstamp(end)
	starttimestamp := tempstarttimestamp
	endtimestamp := tempendtimestamp
	if(starttimestamp>endtimestamp){
		starttimestamp = tempendtimestamp
		endtimestamp = tempstarttimestamp
	}
	//计算天数
	timediff :=endtimestamp-starttimestamp

	days := int(math.Ceil(float64(timediff/86400)))
	//计算小时数
	remain := timediff%86400
	hours :=int(math.Ceil(float64(remain/3600)))
	//计算分钟数
	remain = remain%3600
	mins :=int(math.Ceil(float64(remain/60)))
	//计算秒数
	secs := remain%60

	res :=make(map[string]int)
	res["day"] = days
	res["hour"] = hours
	res["min"] = mins
	res["sec"] = int(secs)
	return res
}

func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}
//ConvertToFormatDay is a function
func ConvertToFormatDay(excelDaysString string)string{
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b,_ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond + realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}

//PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


