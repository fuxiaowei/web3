package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ========== 基础路由 ==========
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome"})
	})

	// ========== 路径参数 ==========
	// 单个参数
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
		})
	})

	// 多个参数
	r.GET("/users/:id/posts/:postId", func(c *gin.Context) {
		userId := c.Param("id")
		postId := c.Param("postId")
		c.JSON(http.StatusOK, gin.H{
			"user_id": userId,
			"post_id": postId,
		})
	})

	// 通配符参数
	// 通配符参数使用 /* 符号，可以匹配路径中剩余的所有部分（包括多个斜杠）
	// 与普通路径参数 :id 的区别：
	//   - :id 只能匹配单个路径段（如 /users/:id 匹配 /users/123，但不匹配 /users/123/posts）
	//   - /*filepath 可以匹配多个路径段（如 /files/*filepath 可以匹配 /files/docs/readme.md）
	//
	// 使用场景：
	//   - 文件路径：/files/*filepath 可以匹配 /files/docs/readme.md、/files/images/photo.jpg 等
	//   - 静态资源：/static/*filepath 可以匹配所有静态资源路径
	//   - 代理转发：将剩余路径转发到其他服务
	//
	// 注意：
	//   - 通配符参数必须放在路由路径的最后
	//   - 获取参数时使用 c.Param("filepath")，参数名不包含 * 号
	//   - 获取到的值会包含前导斜杠（如 "/docs/readme.md"）
	//
	// 示例请求：
	//   GET /files/docs/readme.md        -> filepath = "/docs/readme.md"
	//   GET /files/images/photo.jpg      -> filepath = "/images/photo.jpg"
	//   GET /files/a/b/c/d.txt          -> filepath = "/a/b/c/d.txt"
	r.GET("/files/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.JSON(http.StatusOK, gin.H{
			"filepath": filepath,
		})
	})

	// ========== 查询参数 ==========
	r.GET("/search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		page := c.DefaultQuery("page", "1") // 带默认值
		size := c.Query("size")

		c.JSON(http.StatusOK, gin.H{
			"keyword": keyword,
			"page":    page,
			"size":    size,
		})
	})

	// ========== 表单参数（POST） ==========
	// 表单参数用于接收 application/x-www-form-urlencoded 格式的数据
	// 常用于 HTML 表单提交
	r.POST("/login", func(c *gin.Context) {
		// 获取表单参数
		username := c.PostForm("username")
		password := c.PostForm("password")
		remember := c.DefaultPostForm("remember", "false") // 带默认值

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
			"remember": remember,
		})
	})

	// ========== JSON 参数（POST） ==========
	// JSON 参数用于接收 application/json 格式的请求体
	// 这是现代 API 开发中最常用的方式
	r.POST("/api/users", func(c *gin.Context) {
		// 方式1：使用 ShouldBindJSON 绑定到结构体（推荐）
		type CreateUserRequest struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
			Age   int    `json:"age" binding:"gte=0,lte=120"`
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":    1,
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		})
	})

	// ========== 路由分组 ==========
	// 基础分组
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUsers)
		v1.GET("/users/:id", getUser)
		v1.POST("/users", createUser)
		v1.PUT("/users/:id", updateUser)
		v1.DELETE("/users/:id", deleteUser)
	}

	// 嵌套分组
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Nested group demo; use a different path to avoid clashing with the RESTful products group below.
			v1.GET("/demo-products", getProducts)
		}

		v2 := api.Group("/v2")
		{
			v2.GET("/demo-products", getProductsV2)
		}
	}

	// ========== RESTful API 示例：商品管理 ==========
	products := r.Group("/api/v1/products")
	{
		products.GET("", listProducts)         // GET /api/v1/products
		products.GET("/:id", getProduct)       // GET /api/v1/products/:id
		products.POST("", createProduct)       // POST /api/v1/products
		products.PUT("/:id", updateProduct)    // PUT /api/v1/products/:id
		products.DELETE("/:id", deleteProduct) // DELETE /api/v1/products/:id
	}

	// ========== HTTP 方法 ==========
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.PATCH("/users/:id", patchUser)
	r.DELETE("/users/:id", deleteUser)
	r.HEAD("/users", headUsers)
	r.OPTIONS("/users", optionsUsers)

	// 任意方法
	r.Any("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method": c.Request.Method,
		})
	})

	// 启动服务器
	r.Run(":8080")
}

// ========== Handler 函数 ==========

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []gin.H{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		},
	})
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "User " + id,
	})
}

func createUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
	})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User " + id + " updated",
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User " + id + " deleted",
	})
}

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"products": []gin.H{
			{"id": 1, "name": "Product 1"},
			{"id": 2, "name": "Product 2"},
		},
	})
}

func getProductsV2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "v2",
		"products": []gin.H{
			{"id": 1, "name": "Product 1", "price": 99.99},
		},
	})
}

func listProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")

	c.JSON(http.StatusOK, gin.H{
		"page":     page,
		"size":     size,
		"products": []gin.H{},
	})
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "Product " + id,
	})
}

func createProduct(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created",
	})
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Product " + id + " updated",
	})
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Product " + id + " deleted",
	})
}

func patchUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User " + id + " patched",
	})
}

func headUsers(c *gin.Context) {
	c.Status(http.StatusOK)
}

func optionsUsers(c *gin.Context) {
	c.Header("Allow", "GET, POST, PUT, DELETE, OPTIONS")
	c.Status(http.StatusOK)
}
