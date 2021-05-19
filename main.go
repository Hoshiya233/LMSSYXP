package main

import (
	"example.com/m/v2/controller"
	"example.com/m/v2/tool"
	"github.com/gin-gonic/gin"
)

func main() {
	config := tool.ParseConfig(`config.yml`)
	engine := gin.Default()

	xpManager := engine
	xpManager.LoadHTMLGlob("./templates/*")
	xpManager.Static("/static", "./static")

	//挂载媒体目录
	//需要修复重复问题
	// for _, path := range config.MMDpaths {
	// 	b := strings.SplitAfter(path, ":")[0]
	// 	a := strings.Split(b, ":")[0]
	// 	xpManager.Static("/"+a, b+"\\")
	// }

	xpManager.Static("/E", "E:\\")
	registerRouter(xpManager)
	xpManager.Run(":" + config.Server.Port)
}

//路由设置
func registerRouter(router *gin.Engine) {
	new(controller.HelloXP).Router(router)
}
