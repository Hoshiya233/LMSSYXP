package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

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

//用全局变量保存MMD文件列表
var MMDFileList []MMDFileInfo

var MMDLabelList []string
var MMDPerformerList []string

func (h *HelloXP) Router(engine *gin.Engine) {
	engine.GET("/", h.index)
	engine.GET("/index", h.index)
	engine.GET("/video-list", h.videoList)
	engine.GET("/video", h.video)
	engine.GET("/iteachyou", h.iteachyou)
	engine.GET("/scanpath", h.scanPath)
	//engine.POST("/extractcover", h.extractcover)
	engine.GET("/ws-jobplan", h.wsJobPlan)

	//初始化操作
	tool.ReadStructFromJson("MMDFileList.json", &MMDFileList)
	MMDLabelList = getLabelList(MMDFileList)
	MMDPerformerList = getPerformerList(MMDFileList)

}

func (h *HelloXP) index(context *gin.Context) {
	//context.JSON(200, &MMDFileList)

	context.HTML(200, "index.html", gin.H{})
}

func (h *HelloXP) videoList(context *gin.Context) {

	search_label := context.QueryArray("label")
	search_performer := context.QueryArray("performer")
	searchedList := searchLabel(MMDFileList, search_label)
	searchedList = searchPerformer(searchedList, search_performer)

	data := map[string]interface{}{
		"labelList":     MMDLabelList,
		"performerList": MMDPerformerList,
		"videoList":     &searchedList,
	}
	context.JSON(200, data)

	// context.HTML(200, "video_list.html", gin.H{
	// 	"MMDFileList":   searchedList,
	// 	"labelList":     MMDLabelList,
	// 	"performerList": MMDPerformerList,
	// })
}

func (h *HelloXP) video(context *gin.Context) {
	vid, _ := strconv.Atoi(context.DefaultQuery("vid", "0"))
	context.HTML(200, "video.html", gin.H{
		"video_url": MMDFileList[vid].Url,
	})
}

func (h *HelloXP) scanPath(context *gin.Context) {
	getMMDFileList()

	//将数据写入到缓存中
	tool.WriteStructToJson("MMDFileList.json", MMDFileList)

	MMDLabelList = getLabelList(MMDFileList)
	MMDPerformerList = getPerformerList(MMDFileList)

	context.Redirect(301, "/index")
}

// func (h *HelloXP) extractcover(context *gin.Context) {
// 	context.Redirect(301, "/index")
// 	//extractCover()
// }

func (h *HelloXP) iteachyou(context *gin.Context) {
	context.HTML(200, "iteachyou.html", gin.H{})
}

func (h *HelloXP) wsJobPlan(context *gin.Context) {
	log.Println("正在运行wss")
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	//检查是否需要关闭连接
	do_close := false

	msgList := tool.NewMessageList()

	go func() {
		for {
			//读取ws数据
			t, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("读取ws数据出错：", err)
			}
			if t == -1 {
				do_close = true
				return
			}

			// 提取视频封面
			if string(message) == "extractcover" {
				extractCover(msgList)
			}

			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for {
			//log.Println("正在监听chan")
			msg := msgList.Read()
			//log.Println("从chan中取出：", a)
			//写入ws数据
			err = ws.WriteJSON(msg)
			if err != nil {
				log.Println("写入ws数据出错：", err)
			}
		}
	}()

	for {
		if do_close {
			return
		}
		time.Sleep(time.Second * 10)
	}
}
