package main

import (
	"example.com/m/v2/controller"
	"example.com/m/v2/tool"
	"github.com/gin-gonic/gin"
)

func main() {
	tool.AllConfig = tool.ParseConfig(`config.yml`)
	engine := gin.Default()

	xpManager := engine
	xpManager.LoadHTMLGlob("./templates/*")
	xpManager.Static("/static", "./static")

	//挂载媒体目录

	// var path_head_list []string
	// for _, path := range config.MMDpaths {
	// 	a := strings.Split(path, ":")[0]
	// 	path_head_list = append(path_head_list, a)
	// }
	// path_head_list = tool.DeleteRepeatList(path_head_list)
	// for _, a := range path_head_list {
	// 	xpManager.Static("/"+a, a+":\\")
	// }
	// 以上适用windows系统，把盘符挂载进来

	// bug: 在linux中，如果挂载的目录有空格，会报404
	xpManager.Static("/mnt", "/mnt")

	registerRouter(xpManager)

	xpManager.Run(":" + tool.AllConfig.Server.Port)
}

//路由设置
func registerRouter(router *gin.Engine) {
	new(controller.HelloXP).Router(router)
}
