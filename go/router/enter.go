package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Static("/static", "./static")
	router.Static("/assets", "./static/assets")
	router.StaticFile("/vite.svg", "./static/vite.svg")
	router.LoadHTMLGlob("static/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	// 其他路由
	routerGroup := router.Group("/api")
	routerGroup2 := router.Group("/")
	InitPollHttpRouter(routerGroup, db)
	InitWebSocketRouter(routerGroup2)

	return router
}
