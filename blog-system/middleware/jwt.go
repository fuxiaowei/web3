package middleware

import (
	"blog-system/config"
	"blog-system/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.ErrResponse(utils.TokenEmpty, "Token为空"))
			c.Abort()
			return
		}

		// 解析Token格式（Bearer xxxxx）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, utils.ErrResponse(utils.TokenInvalid, "Token格式错误"))
			c.Abort()
			return
		}

		// 验证Token
		token := parts[1]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrResponse(utils.TokenInvalid, "Token无效或过期"))
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", uint(claims["id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Next()
	}
}
