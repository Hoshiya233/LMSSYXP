package controller

import (
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"example.com/m/v2/tool"
)

type MMDFileInfo struct {
	Id        int
	Name      string
	Dir       string
	Path      string
	Url       string
	Performer []string
	Bgm       string
	Label     []string
	CoverUrl  string
}

func GetMMDFileList() {
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
		//mmdFileInfo.CoverUrl = getCoverUrl(filePathName, fileName)

		MMDFileList = append(MMDFileList, mmdFileInfo)
	}
	//初始化一个控制池,设置并发数量
	pool := tool.NewPool(16, len(MMDFileList))
	//计算执行时间
	now := time.Now()
	//并发处理
	for i := range MMDFileList {
		go func(fileinfo *MMDFileInfo) {
			pool.AddOne() // 向并发控制池中添加一个, 一旦池满则此处阻塞
			//任务处理
			fileinfo.CoverUrl = getCoverUrl(fileinfo.Path, fileinfo.Name)
			//fileinfo.CoverUrl = "123"
			pool.DelOne() // 从并发控制池中释放一个, 之后其他被阻塞的可以进入池中
		}(&MMDFileList[i])
	}
	pool.WG.Wait()
	//计算执行时间
	next := time.Now()
	log.Println("提取封面花费时间:", next.Sub(now))
}

func SearchLabel(fileList []MMDFileInfo, labels []string) []MMDFileInfo {
	if len(labels) == 0 {
		return fileList
	}
	var searchedList []MMDFileInfo
	for _, fileinfo := range fileList {
		for _, search_label := range labels {
			for _, file_label := range fileinfo.Label {
				if search_label == file_label {
					searchedList = append(searchedList, fileinfo)
				}
			}
		}
	}
	return searchedList
}

func SearchPerformer(fileList []MMDFileInfo, performers []string) []MMDFileInfo {
	if len(performers) == 0 {
		return fileList
	}
	var searchedList []MMDFileInfo
	for _, fileinfo := range fileList {
		for _, search_performer := range performers {
			for _, file_performer := range fileinfo.Performer {
				if search_performer == file_performer {
					searchedList = append(searchedList, fileinfo)
				}
			}
		}
	}
	return searchedList
}

func GetLabelList(fileList []MMDFileInfo) []string {
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

func GetPerformerList(fileList []MMDFileInfo) []string {
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

	a := strings.Split(filename, "-")
	aa := strings.Split(a[0], "&")
	for i := range aa {
		performer = append(performer, strings.TrimSpace(aa[i]))
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

	bgm = strings.TrimSpace(strings.Split(a[1], "[")[0])

	return bgm
}

func getCoverUrl(filepathname string, filename string) string {
	//只能在windows上执行，如果要在linux上执行，需要改目录分隔符，如果有集成golang的方案就好了
	//获取视频封面的命令
	//ffmpeg.exe -i '.\弱音 - メランコリック.mp4' -y -f image2 -t 0.001 a.jpg
	var coverUrl string

	// 去掉文件名后缀
	x := strings.LastIndexByte(filename, '.')
	if x != -1 {
		filename = filename[:x]
	}
	coverPath := `.\static\tmp\cover\` + filename + `.jpg`
	//cmd := `-i '` + filename + `'-y -f image2 -ss 5 -t 0.001 a.jpg`
	out := exec.Command(`C:\Program Files\ffmpeg\bin\ffmpeg`, "-i", filepathname, "-y", "-f", "image2", "-ss", "5", "-t", "0.001", coverPath)
	out.Output()
	log.Println("已获取" + filename + "封面")
	coverUrl = strings.ReplaceAll(coverPath[1:], "\\", "/")
	return coverUrl
}
