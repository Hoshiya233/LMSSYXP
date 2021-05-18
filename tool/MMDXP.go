package tool

import (
	"log"
	"path/filepath"
	"strings"
)

type MMDFileInfo struct {
	Name      string
	Dir       string
	Path      string
	Performer []string
	Bgm       string
	Label     []string
}

func GetMMDFileList() []MMDFileInfo {
	var filePathNames []string
	for _, path := range MMDPaths {
		f, err := filepath.Glob(filepath.Join(path, "*"))
		if err != nil {
			log.Fatal(err)
		}
		filePathNames = append(filePathNames, f...)
	}

	var fileList []MMDFileInfo
	var mmdFileInfo MMDFileInfo
	for i := range filePathNames {
		filePathName := filePathNames[i]
		dir, fileName := filepath.Split(filePathName)
		mmdFileInfo.Dir = dir
		mmdFileInfo.Name = fileName
		mmdFileInfo.Path = filePathName
		mmdFileInfo.Label = readLabel(fileName)
		mmdFileInfo.Performer = readPerformer(fileName)
		mmdFileInfo.Bgm = readBgm(fileName)

		fileList = append(fileList, mmdFileInfo)
	}
	return fileList
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

func GetLabels(fileList []MMDFileInfo) []string {
	var labels []string
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
		labels = append(labels, k)
	}
	return labels
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
