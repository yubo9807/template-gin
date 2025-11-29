package main

import (
	"server/configs"
	"server/src/middleware"
	"server/src/router"
	"strconv"

	"github.com/gin-gonic/gin"
)

func server() *gin.Engine {

	app := gin.Default()
	app.Use(middleware.Core)

	// 代理应用
	power := app.Group("/permissions")
	power.Use(middleware.Authorization)
	power.Use(middleware.RoleVerify("0")) // 指定特定的角色可以调以下接口
	power.Any("/*path", middleware.ProxyPermissions)

	// 自身应用
	self := app.Group(configs.Config.Prefix)
	self.Use(middleware.ContextMiddleware)
	self.Use(middleware.Log)

	router.Basic(self.Group("/basic/api"))
	router.V1(self.Group("/v1/api"))

	// 前端静态文件
	staticPrefix := configs.Config.StaticDir
	middleware.Static(app, staticPrefix)
	app.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, staticPrefix+"/")
	})

	return app
}

func main() {

	port := ":" + strconv.Itoa(configs.Config.Port)
	server().Run(port)

}
