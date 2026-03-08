package utils

import "github.com/gin-gonic/gin"

// 错误码定义
const (
	SUCCESS        = 200
	ERROR          = 500
	TokenEmpty     = 401
	TokenInvalid   = 402
	UserExist      = 403
	UserNotExist   = 404
	PasswordError  = 405
	PostNotExist   = 406
	PermissionDeny = 407
)

// ErrResponse 统一错误响应格式
func ErrResponse(code int, msg string) gin.H {
	return gin.H{
		"code": code,
		"msg":  msg,
	}
}

// SuccessResponse 统一成功响应格式
func SuccessResponse(data interface{}) gin.H {
	return gin.H{
		"code": SUCCESS,
		"msg":  "操作成功",
		"data": data,
	}
}
