package main

import (
	"blog-system/handlers"
	"blog-system/middleware"
	"blog-system/models"
	"blog-system/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	// 初始化数据库
	models.InitDB()

	// 初始化Gin路由
	r := gin.Default()

	// 公开接口（无需认证）
	public := r.Group("/api")
	{
		public.POST("/register", handlers.Register) // 用户注册
		public.POST("/login", handlers.Login)       // 用户登录
		public.GET("/posts", handlers.GetPosts)     // 获取所有文章
		public.GET("/posts/:id", handlers.GetPost)  // 获取单篇文章
		// 修复：修改评论接口路由，避免冲突
		public.GET("/comments/post/:post_id", handlers.GetComments) // 获取文章评论
	}

	// 私有接口（需JWT认证）
	private := r.Group("/api")
	private.Use(middleware.JWTAuth())
	{
		private.POST("/posts", handlers.CreatePost)       // 创建文章
		private.PUT("/posts/:id", handlers.UpdatePost)    // 更新文章
		private.DELETE("/posts/:id", handlers.DeletePost) // 删除文章
		private.POST("/comments", handlers.CreateComment) // 创建评论
	}

	// 启动服务
	utils.Logger.Info("服务器启动成功，监听端口: 8080")
	r.Run(":8080")
}
