package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gin-examples/project/config"
	"gin-examples/project/handlers"
	"gin-examples/project/middleware"
	"gin-examples/project/models"
	"gin-examples/project/services"
	"gin-examples/project/utils"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化服务
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, []byte(cfg.JWT.Secret))

	// 创建 Gin 引擎
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status": "ok",
		})
	})

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/users/register", userHandler.Register)
		public.POST("/users/login", userHandler.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth([]byte(cfg.JWT.Secret)))
	{
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PUT("/users/me", userHandler.UpdateProfile)
	}

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

