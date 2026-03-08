package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ========== 请求结构体 ==========

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=0,lte=120"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   int    `json:"age" binding:"omitempty,gte=0,lte=120"`
}

type ListProductsRequest struct {
	Page    int    `form:"page" binding:"gte=1"`
	Size    int    `form:"size" binding:"gte=1,lte=100"`
	Keyword string `form:"keyword"`
}

type GetUserRequest struct {
	ID int `uri:"id" binding:"required"`
}

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// ========== 响应结构体 ==========

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

func main() {
	r := gin.Default()

	// ========== JSON 绑定 ==========
	r.POST("/api/users", func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		success(c, gin.H{
			"id":    1,
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		})
	})

	// ========== 查询参数绑定 ==========
	r.GET("/api/products", func(c *gin.Context) {
		var req ListProductsRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		success(c, gin.H{
			"page":    req.Page,
			"size":    req.Size,
			"keyword": req.Keyword,
			"products": []gin.H{
				{"id": 1, "name": "Product 1"},
				{"id": 2, "name": "Product 2"},
			},
		})
	})

	// ========== 路径参数绑定 ==========
	r.GET("/api/users/:id", func(c *gin.Context) {
		var req GetUserRequest
		if err := c.ShouldBindUri(&req); err != nil {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		success(c, gin.H{
			"id":   req.ID,
			"name": "User " + string(rune(req.ID)),
		})
	})

	// ========== 表单绑定 ==========
	r.POST("/api/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBind(&req); err != nil {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// 模拟登录验证
		if req.Username == "admin" && req.Password == "admin123" {
			success(c, gin.H{
				"token": "fake-jwt-token",
				"user": gin.H{
					"username": req.Username,
				},
			})
		} else {
			errorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		}
	})

	// ========== 原始请求体 ==========
	r.POST("/api/raw", func(c *gin.Context) {
		data, err := c.GetRawData()
		if err != nil {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		success(c, gin.H{
			"raw_data": string(data),
		})
	})

	// ========== 多种绑定方式 ==========
	r.POST("/api/mixed", func(c *gin.Context) {
		// 路径参数
		id := c.Param("id")

		// 查询参数
		page := c.DefaultQuery("page", "1")

		// JSON 体
		var body CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			errorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		success(c, gin.H{
			"id":   id,
			"page": page,
			"body": body,
		})
	})

	// ========== 响应类型 ==========

	// JSON 响应
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "JSON response"})
	})

	// XML 响应
	r.GET("/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "XML response"})
	})

	// 字符串响应
	r.GET("/string", func(c *gin.Context) {
		c.String(http.StatusOK, "String response")
	})

	// 重定向
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/json")
	})

	// 文件响应
	// r.GET("/file", func(c *gin.Context) {
	// 	c.File("./file.txt")
	// })

	// 数据流响应
	r.GET("/data", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/octet-stream", []byte("binary data"))
	})

	r.Run(":8080")
}

