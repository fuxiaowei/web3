package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// ========== 全局中间件 ==========
	r.Use(loggerMiddleware())
	r.Use(recoveryMiddleware())

	// ========== 路由 ==========
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	// ========== 分组中间件 ==========
	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		api.GET("/users", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
			})
		})
	}

	// 不需要认证的路由
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "public endpoint",
		})
	})

	r.Run(":8080")
}

// ========== 日志中间件 ==========
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 前置处理
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 进入下一个处理函数
		c.Next()

		// 后置处理
		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Printf("[%s] %s %d %v\n", method, path, status, latency)
	}
}

// ========== 恢复中间件 ==========
func recoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		c.Abort()
	})
}

// ========== 认证中间件（示例） ==========
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// 模拟验证 Token
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("userID", 1)
		c.Set("username", "admin")

		c.Next()
	}
}

// ========== CORS 中间件 ==========
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ========== 限流中间件（简单示例） ==========
var requestCount = make(map[string]int)
var lastReset = time.Now()

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		// 每分钟重置一次
		if now.Sub(lastReset) > time.Minute {
			requestCount = make(map[string]int)
			lastReset = now
		}

		// 检查请求次数
		if requestCount[ip] >= 10 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		requestCount[ip]++
		c.Next()
	}
}

