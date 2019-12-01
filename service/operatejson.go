package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Info struct {
	Hobbies []string
}

func InitWriteFile(jsonfile string) {

	//filePtr, err := os.Open("info.json")
	filePtr, err := os.Open(jsonfile)
	if err != nil {
		fmt.Println("Open file failed [Err:%s]", err.Error())
		return
	}
	defer filePtr.Close()

	var person []Info

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&person)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())
	}
	// 创建文件
	ptr, err := os.Create("tempinfo.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(ptr)

	err = encoder.Encode(person)
	if err != nil {
		fmt.Println("init Encoder failed", err.Error())
	}

}

func DelFirstFile(tempjsonfile string,del string) {
	a := []string{"da", "folder", "excelfile","exportexcelfile","regulation","userorder","special_regulation","allpeople"}
	index :=0
	for k,v :=range a{
		if v==del{
			index = k
		}
	}
	//删除第i之前的元素
	a = append([]string{},a[index+1:]...)
	info := []Info{{a}}
	//fmt.Println(a)
	//info := []Info{{[]string{"跑步", "读书", "看电影"}}}

	// 创建文件
	//filePtr, err := os.Create("tempinfo.json")
	filePtr, err := os.Create(tempjsonfile)
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)

	err = encoder.Encode(info)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}

}

func ReadJsonFirst(tempjsonfile string) string{

	//filePtr, err := os.Open("tempinfo.json")
	filePtr, err := os.Open(tempjsonfile)
	if err != nil {
		fmt.Println("Open file failed [Err:%s]", err.Error())
		return ""
	}
	defer filePtr.Close()

	var person []Info
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&person)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())
	}
	if len(person[0].Hobbies)==0 {
		return ""
	}
	return person[0].Hobbies[0]
}

func CopyFile(destName,srcName string)(int,error){
	//1.打开源文件，读取数据
	//2.打开目标文件，写出数据
	srcFile, err:=os.Open(srcName)
	if err!=nil {
		return 0 ,err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(destName,os.O_WRONLY|os.O_CREATE,0777)
	if err != nil{
		return 0, err
	}
	defer destFile.Close()
	//复制数据
	bs := make([] byte, 1024)
	count := 0//每次实际读入数据量
	total :=0// 用于统计读取的数据总量
	for{
		count,err = srcFile.Read(bs)
		if err!=nil && err !=io.EOF{ //EOF
			return total, err
		}else if err == io.EOF{
			fmt.Println("已经到达文件末尾。。")
			break
		}
		destFile.Write(bs[:count])
		total+=count
	}
	return total, nil

}
