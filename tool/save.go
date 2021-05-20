package tool

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func WriteStructToJson(filename string, data interface{}) {
	CreateDir("./static/tmp/save/")
	stu, err := json.Marshal(data)
	if err != nil {
		log.Println("将struct转为json错误：", err)
	}
	err = ioutil.WriteFile("./static/tmp/save/"+filename, stu, 0777)
	if err != nil {
		log.Println("写入文件错误：", err)
	}
}

func ReadStructFromJson(filename string, data interface{}) {
	stu, err := ioutil.ReadFile("./static/tmp/save/" + filename)
	if err != nil {
		log.Println("读取文件错误：", err)
	}
	err = json.Unmarshal(stu, &data)
	if err != nil {
		log.Println("将json转为struct错误：", err)
	}
}

func IsExistDir(dirName string) bool {
	s, err := os.Stat(dirName)
	if err != nil {
		log.Println("判断目录是否存在时发生错误：", err)
		return false
	}
	return s.IsDir()
}

func CreateDir(dirName string) {
	if !IsExistDir(dirName) {
		err := os.MkdirAll(dirName, 0777)
		if err != nil {
			log.Println("创建目录时发生错误：", err)
		}
	}
}
