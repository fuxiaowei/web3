package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 引擎实例
	// gin.Default() 包含日志和恢复中间件
	// gin.New() 创建空白引擎，不包含任何中间件
	r := gin.Default()

	// 定义路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	// 返回字符串响应
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 返回 HTML（需要先加载模板）
	// r.LoadHTMLGlob("templates/*")
	// r.GET("/index", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"title": "Home Page",
	// 	})
	// })

	// 启动服务器，默认监听 8080 端口
	// r.Run() 等同于 r.Run(":8080")
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

