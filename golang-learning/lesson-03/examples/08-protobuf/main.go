package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	"protobuf/pb"
)

func main() {
	r := gin.Default()

	// Protobuf API 路由
	api := r.Group("/api/proto")
	{
		// 获取单个用户（返回 Protobuf 格式）
		api.GET("/user/:id", getUserProto)

		// 获取用户列表（返回 Protobuf 格式）
		api.GET("/users", getUserListProto)

		// 创建用户（接收和返回 Protobuf 格式）
		api.POST("/user", createUserProto)

		// 对比：返回 JSON 格式的用户（用于对比）
		api.GET("/user/:id/json", getUserJSON)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Protobuf example server is running",
		})
	})

	log.Println("Server starting on :8080")
	log.Println("Try these endpoints:")
	log.Println("  GET  http://localhost:8080/api/proto/user/1")
	log.Println("  GET  http://localhost:8080/api/proto/users")
	log.Println("  POST http://localhost:8080/api/proto/user (with protobuf body)")
	log.Println("  GET  http://localhost:8080/api/proto/user/1/json (JSON format for comparison)")
	r.Run(":8080")
}

// getUserProto 返回 Protobuf 格式的用户信息
func getUserProto(c *gin.Context) {
	id := c.Param("id")

	// 模拟从数据库获取用户
	user := &pb.User{
		Id:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Age:      25,
		Active:   true,
		Tags:     []string{"admin", "developer"},
		Metadata: map[string]string{
			"department": "engineering",
			"location":    "Beijing",
		},
	}

	// 如果提供了 ID，可以设置不同的 ID
	if id == "2" {
		user.Id = 2
		user.Username = "bob"
		user.Email = "bob@example.com"
		user.Age = 30
		user.Tags = []string{"user", "tester"}
	}

	// 序列化 Protobuf
	data, err := proto.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to marshal protobuf: %v", err),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/x-protobuf")
	c.Data(http.StatusOK, "application/x-protobuf", data)
}

// getUserListProto 返回 Protobuf 格式的用户列表
func getUserListProto(c *gin.Context) {
	// 模拟用户列表数据
	users := []*pb.User{
		{
			Id:       1,
			Username: "alice",
			Email:    "alice@example.com",
			Age:      25,
			Active:   true,
			Tags:     []string{"admin", "developer"},
		},
		{
			Id:       2,
			Username: "bob",
			Email:    "bob@example.com",
			Age:      30,
			Active:   true,
			Tags:     []string{"user", "tester"},
		},
		{
			Id:       3,
			Username: "charlie",
			Email:    "charlie@example.com",
			Age:      28,
			Active:   false,
			Tags:     []string{"user"},
		},
	}

	userList := &pb.UserList{
		Users: users,
		Total: int32(len(users)),
	}

	// 序列化 Protobuf
	data, err := proto.Marshal(userList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to marshal protobuf: %v", err),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/x-protobuf")
	c.Data(http.StatusOK, "application/x-protobuf", data)
}

// createUserProto 接收 Protobuf 格式的请求，创建用户并返回 Protobuf 响应
func createUserProto(c *gin.Context) {
	// 读取原始请求数据
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to read request body: %v", err),
		})
		return
	}

	// 反序列化 Protobuf 请求
	var req pb.CreateUserRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to unmarshal protobuf: %v", err),
		})
		return
	}

	// 验证请求数据
	if req.Username == "" || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username and email are required",
		})
		return
	}

	// 创建用户（模拟）
	user := &pb.User{
		Id:       100, // 模拟生成的 ID
		Username: req.Username,
		Email:    req.Email,
		Age:      req.Age,
		Active:   true,
		Tags:     []string{"user"},
		Metadata: map[string]string{
			"created_at": "2024-01-01",
		},
	}

	// 构建响应
	resp := &pb.CreateUserResponse{
		User:    user,
		Success: true,
		Message: "User created successfully",
	}

	// 序列化响应
	respData, err := proto.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to marshal response: %v", err),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/x-protobuf")
	c.Data(http.StatusOK, "application/x-protobuf", respData)
}

// getUserJSON 返回 JSON 格式的用户信息（用于对比）
func getUserJSON(c *gin.Context) {
	user := gin.H{
		"id":       1,
		"username": "alice",
		"email":    "alice@example.com",
		"age":      25,
		"active":   true,
		"tags":     []string{"admin", "developer"},
		"metadata": map[string]string{
			"department": "engineering",
			"location":    "Beijing",
		},
	}

	c.JSON(http.StatusOK, user)
}

