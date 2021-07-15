package controller

import (
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.com/m/v2/tool"
)

type MMDFileInfo struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Dir       string   `json:"dir"`
	Path      string   `json:"path"`
	Url       string   `json:"url"`
	Performer []string `json:"performer"`
	Bgm       string   `json:"bgm"`
	Label     []string `json:"label"`
	CoverUrl  string   `json:"coverurl"`
}

func getMMDFileList() {
	var filePathNames []string
	MMDFileList = nil
	for _, path := range tool.MMDPaths {
		f, err := filepath.Glob(filepath.Join(path, "*"))
		if err != nil {
			log.Fatal(err)
		}
		filePathNames = append(filePathNames, f...)
	}

	var mmdFileInfo MMDFileInfo
	for i := range filePathNames {
		filePathName := filePathNames[i]
		dir, fileName := filepath.Split(filePathName)
		mmdFileInfo.Id = i
		mmdFileInfo.Name = fileName
		mmdFileInfo.Dir = dir
		mmdFileInfo.Path = filePathName
		mmdFileInfo.Label = readLabel(fileName)
		mmdFileInfo.Performer = readPerformer(fileName)
		mmdFileInfo.Bgm = readBgm(fileName)
		mmdFileInfo.Url = strings.Replace(strings.Replace(filePathName, ":", "", 1), `\`, "/", -1)
		mmdFileInfo.CoverUrl = `/static/tmp/cover/` + strconv.Itoa(i) + `.jpg`

		MMDFileList = append(MMDFileList, mmdFileInfo)
	}

}

func filterLabel(fileList []MMDFileInfo, labels []string) []MMDFileInfo {
	if len(labels) == 0 {
		return fileList
	}
	var filteredList []MMDFileInfo
	for _, fileinfo := range fileList {
		for _, filter_label := range labels {
			for _, file_label := range fileinfo.Label {
				if filter_label == file_label {
					filteredList = append(filteredList, fileinfo)
				}
			}
		}
	}
	return filteredList
}

func filterPerformer(fileList []MMDFileInfo, performers []string) []MMDFileInfo {
	if len(performers) == 0 {
		return fileList
	}
	var filteredList []MMDFileInfo
	for _, fileinfo := range fileList {
		for _, filter_performer := range performers {
			for _, file_performer := range fileinfo.Performer {
				if filter_performer == file_performer {
					filteredList = append(filteredList, fileinfo)
				}
			}
		}
	}
	return filteredList
}

func getLabelList(fileList []MMDFileInfo) []string {
	var labelList []string
	m := make(map[string]int)
	for _, fileinfo := range fileList {
		for _, label := range fileinfo.Label {
			if _, ok := m[label]; ok {
				continue
			} else {
				m[label] = 1
			}
		}
	}
	for k := range m {
		labelList = append(labelList, k)
	}
	sort.Strings(labelList)
	return labelList
}

func getPerformerList(fileList []MMDFileInfo) []string {
	var performerList []string
	m := make(map[string]int)
	for _, fileinfo := range fileList {
		for _, performer := range fileinfo.Performer {
			if _, ok := m[performer]; ok {
				continue
			} else {
				m[performer] = 1
			}
		}
	}
	for k := range m {
		performerList = append(performerList, k)
	}
	sort.Strings(performerList)
	return performerList
}

func readLabel(filename string) []string {
	// 从文件名中读取标签

	var label []string

	// 去掉文件名后缀
	x := strings.LastIndexByte(filename, '.')
	if x != -1 {
		filename = filename[:x]
	}

	a := strings.Split(filename, "[")
	if len(a) == 1 || len(a) == 0 {
		return label
	}

	for i := 1; i < len(a); i++ {
		label = append(label, strings.TrimRight(a[i], "]"))
	}

	return label
}

func readPerformer(filename string) []string {
	//从文件名中读取演出者

	var performer []string

	// 去掉文件名后缀
	x := strings.LastIndexByte(filename, '.')
	if x != -1 {
		filename = filename[:x]
	}

	x = strings.IndexByte(filename, '-')
	if x == -1 {
		return performer
	} else {
		a := filename[:x]
		aa := strings.Split(a, "&")
		for i := range aa {
			performer = append(performer, strings.TrimSpace(aa[i]))
		}
	}

	return performer
}

func readBgm(filename string) string {
	//从文件名中读取BGM

	var bgm string

	// 去掉文件名后缀
	x := strings.LastIndexByte(filename, '.')
	if x != -1 {
		filename = filename[:x]
	}

	a := strings.Split(filename, "-")
	if len(a) == 1 || len(a) == 0 {
		return bgm
	}

	bgm = strings.TrimSpace(strings.TrimRight(strings.TrimRight(a[1], "["), "-"))

	return bgm
}

func extractCover(msgList *tool.MessageList) {
	/*
		提取视频封面
		只能在windows上执行，如果要在linux上执行，需要改目录分隔符，如果有集成golang的方案就好了
		ffmpeg参数说明  -threads n 限制使用CPU核心数
						-ss n 放在-i之前，表示直接从第n秒开始读取；放在-i之后，表示在n之前的正常解码但会被丢弃，浪费CPU
	*/

	tool.CreateDir("./static/tmp/cover/")
	//初始化一个控制池,设置并发数量
	pool := tool.NewPool(4, len(MMDFileList))
	//计算执行时间
	begin := time.Now()
	//并发处理
	for i := range MMDFileList {
		go func(item *MMDFileInfo) {
			pool.AddOne() // 向并发控制池中添加一个, 一旦池满则此处阻塞
			//任务处理
			coverPath := `.\static\tmp\cover\` + strconv.Itoa(item.Id) + `.jpg`
			out := exec.Command(`C:\Program Files\ffmpeg\bin\ffmpeg`, "-threads", "1", "-ss", "5", "-i", item.Path, "-y", "-f", "image2", "-t", "0.001", coverPath)
			out.Output()
			pool.DelOne() // 从并发控制池中释放一个, 之后其他被阻塞的可以进入池中

			log.Println("已获取" + item.Name + "封面")
			log.Println("任务进度：", pool.GetProgressRate())

			//写入ws数据，使用json格式
			msgList.Write(pool.GetProgressRate(), "已获取"+item.Name+"封面")

		}(&MMDFileList[i])
	}
	pool.Wait()
	//计算执行时间
	end := time.Now()
	log.Println("提取封面花费时间:", end.Sub(begin))
	//写入ws数据
	msgList.Write(1, "提取视频封面任务已完成")
}
