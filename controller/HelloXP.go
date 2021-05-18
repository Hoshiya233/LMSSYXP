package controller

import (
	"strconv"

	"example.com/m/v2/tool"
	"github.com/gin-gonic/gin"
)

type HelloXP struct {
}

//用全局变量保存MMD文件列表
var MMDFileList []tool.MMDFileInfo

func (helloxp *HelloXP) Router(engine *gin.Engine) {
	engine.GET("/", helloxp.index)
	engine.GET("/index", helloxp.index)
	engine.GET("/video", helloxp.video)
}

func (helloxp *HelloXP) index(context *gin.Context) {
	MMDFileList = tool.GetMMDFileList()
	//context.JSON(200, &MMDFileList)

	labelList := tool.GetLabels(MMDFileList)

	search_label := context.QueryArray("label")
	search_performer := context.QueryArray("performer")
	searchedList := tool.SearchLabel(MMDFileList, search_label)
	searchedList = tool.SearchPerformer(searchedList, search_performer)

	context.HTML(200, "index.html", gin.H{
		"MMDFileList": searchedList,
		"labelList":   labelList,
	})
}

func (helloxp *HelloXP) video(context *gin.Context) {
	vid, _ := strconv.Atoi(context.DefaultQuery("vid", "0"))
	context.HTML(200, "video.html", gin.H{
		"video_url": MMDFileList[vid].Url,
	})
}
