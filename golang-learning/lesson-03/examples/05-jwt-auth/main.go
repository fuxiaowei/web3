package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// ========== JWT Claims ==========
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	r := gin.Default()

	// ========== 公开路由 ==========
	r.POST("/api/login", login)
	r.GET("/api/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public endpoint"})
	})

	// ========== 需要认证的路由 ==========
	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		api.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"message":  "protected endpoint",
				"user_id":  userID,
				"username": username,
			})
		})

		api.GET("/profile", getProfile)
	}

	r.Run(":8080")
}

// ========== 生成 Token ==========
func generateToken(userID int, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ========== 解析 Token ==========
func parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ========== 认证中间件 ==========
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// 提取 Token（Bearer <token>）
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// ========== 登录接口 ==========
func login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户名密码（示例，实际应从数据库查询）
	if req.Username == "admin" && req.Password == "admin123" {
		token, err := generateToken(1, req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":       1,
				"username": req.Username,
			},
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Invalid credentials",
	})
}

// ========== 获取用户信息 ==========
func getProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
	})
}

