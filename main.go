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

	registerRouter(xpManager)
	xpManager.Run(":" + config.Server.Port)
}

//路由设置
func registerRouter(router *gin.Engine) {
	new(controller.HelloXP).Router(router)
}
