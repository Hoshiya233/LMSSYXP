package controller

import (
	"log"
	"net/http"
	"strconv"

	"example.com/m/v2/tool"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type HelloXP struct {
}

// websocket需要
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 消息队列
var Message chan string

//用全局变量保存MMD文件列表
var MMDFileList []MMDFileInfo
var MMDLabelList []string

var MMDPerformerList []string

func (helloxp *HelloXP) Router(engine *gin.Engine) {
	engine.GET("/", helloxp.index)
	engine.GET("/index", helloxp.index)
	engine.GET("/video", helloxp.video)
	engine.GET("/iteachyou", helloxp.iteachyou)
	engine.POST("/scanpath", helloxp.scanPath)
	engine.POST("/extractcover", helloxp.extractcover)
	engine.GET("/msg", helloxp.msg)

	//初始化操作
	tool.ReadStructFromJson("MMDFileList.json", &MMDFileList)
	MMDLabelList = getLabelList(MMDFileList)
	MMDPerformerList = getPerformerList(MMDFileList)

}

func (helloxp *HelloXP) index(context *gin.Context) {
	//context.JSON(200, &MMDFileList)

	search_label := context.QueryArray("label")
	search_performer := context.QueryArray("performer")
	searchedList := searchLabel(MMDFileList, search_label)
	searchedList = searchPerformer(searchedList, search_performer)

	context.HTML(200, "index.html", gin.H{
		"MMDFileList":   searchedList,
		"labelList":     MMDLabelList,
		"performerList": MMDPerformerList,
	})
}

func (helloxp *HelloXP) video(context *gin.Context) {
	vid, _ := strconv.Atoi(context.DefaultQuery("vid", "0"))
	context.HTML(200, "video.html", gin.H{
		"video_url": MMDFileList[vid].Url,
	})
}

func (helloxp *HelloXP) scanPath(context *gin.Context) {
	getMMDFileList()

	//将数据写入到缓存中
	tool.WriteStructToJson("MMDFileList.json", MMDFileList)

	MMDLabelList = getLabelList(MMDFileList)
	MMDPerformerList = getPerformerList(MMDFileList)

	context.Redirect(301, "/index")
}

func (helloxp *HelloXP) extractcover(context *gin.Context) {
	context.Redirect(301, "/index")
	extractCover()
}

func (helloxp *HelloXP) iteachyou(context *gin.Context) {
	context.HTML(200, "iteachyou.html", gin.H{})
}

func (helloxp *HelloXP) msg(context *gin.Context) {
	log.Println("正在运行msg")
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//读取ws数据
		// mt, _, err := ws.ReadMessage()
		// if err != nil {
		// 	break
		// }
		log.Println("正在监听chan")
		a := <-Message
		//写入ws数据
		err = ws.WriteMessage(200, []byte(a))
		if err != nil {
			break
		}
	}
}
