package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HelloXP struct {
}

//用全局变量保存MMD文件列表
var MMDFileList []MMDFileInfo
var MMDLabelList []string

func (helloxp *HelloXP) Router(engine *gin.Engine) {
	engine.GET("/", helloxp.index)
	engine.GET("/index", helloxp.index)
	engine.GET("/video", helloxp.video)
	engine.POST("/scanpath", helloxp.scanPath)
}

func (helloxp *HelloXP) index(context *gin.Context) {
	//context.JSON(200, &MMDFileList)

	search_label := context.QueryArray("label")
	search_performer := context.QueryArray("performer")
	searchedList := SearchLabel(MMDFileList, search_label)
	searchedList = SearchPerformer(searchedList, search_performer)

	for _, v := range searchedList {
		fmt.Println(v.Name + "封面为" + v.CoverUrl)
	}

	context.HTML(200, "index.html", gin.H{
		"MMDFileList": searchedList,
		"labelList":   MMDLabelList,
	})
}

func (helloxp *HelloXP) video(context *gin.Context) {
	vid, _ := strconv.Atoi(context.DefaultQuery("vid", "0"))
	context.HTML(200, "video.html", gin.H{
		"video_url": MMDFileList[vid].Url,
	})
}

func (helloxp *HelloXP) scanPath(context *gin.Context) {
	GetMMDFileList()
	MMDLabelList = GetLabels(MMDFileList)
	context.Redirect(301, "/index")
}
