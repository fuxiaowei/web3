# 第3期教案：Gin 基础与路由

## 1. 学习目标

通过本课程的学习，你将能够：

- 理解 Web 框架与 Gin 的基本概念与优势
- 熟悉 Gin 的安装、初始化、项目结构设计
- 掌握路由系统的完整用法（路由定义、分组、参数）
- 能够处理各种类型的 HTTP 请求和响应

## 2. 核心知识点

### 2.1 Web 框架与 Gin 概述

#### 什么是 Web 框架？

Web 框架是一个软件框架，用于简化 Web 应用程序的开发。它提供了一系列工具和库，帮助开发者处理路由、请求解析、响应生成等常见任务。

#### 为什么需要 Web 框架？

相比直接使用标准库 `net/http`，Web 框架提供了以下优势：

1. **提高生产效率**：减少重复代码，专注于业务逻辑
2. **统一规范**：提供统一的代码组织和开发规范
3. **丰富功能**：内置路由、中间件、参数绑定等功能
4. **社区支持**：活跃的社区和丰富的插件生态
5. **性能优化**：框架级别的性能优化

#### Gin 简介

Gin 是 Go 语言中最流行的 Web 框架，具有以下特点：

- **高性能**：基于 `httprouter`，性能接近 `net/http`
- **简洁易用**：API 设计简洁，学习曲线平缓
- **功能丰富**：路由、中间件、参数绑定、JSON 验证等
- **活跃社区**：GitHub 60k+ stars，大量中间件和插件
- **生产就绪**：被众多公司用于生产环境

**官方文档**：https://gin-gonic.com/docs/

**GitHub**：https://github.com/gin-gonic/gin

### 2.2 安装与项目初始化

#### 项目结构

> **重要**：每个示例现在都是独立的工程，都有自己的 `go.mod` 文件！

```
lesson-03/examples/
├── README.md              # 示例总览
├── 01-hello-world/        # 独立工程
│   ├── main.go
│   ├── go.mod
│   └── README.md
├── 02-routing/            # 独立工程
│   ├── main.go
│   ├── go.mod
│   └── README.md
├── 03-request-handling/   # 独立工程
│   ├── main.go
│   ├── go.mod
│   └── README.md
└── project/               # 综合项目
    ├── main.go
    └── ...
```

#### 运行示例

每个示例都是独立的，可以直接运行：

```bash
# 进入示例目录
cd lesson-03/examples/01-hello-world

# 第一次运行会自动下载依赖
go run main.go

# 或者先下载依赖
go mod download
go run main.go
```

#### 第一个 Hello World

参考示例：`examples/01-hello-world/` 目录

**最简单的 Gin 应用**：

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    // 创建 Gin 引擎实例
    r := gin.Default()
    
    // 定义路由
    r.GET("/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello World!",
        })
    })
    
    // 启动服务器，默认监听 8080 端口
    r.Run() // 等同于 r.Run(":8080")
}
```

**运行**：
```bash
cd lesson-03/examples/01-hello-world
go run main.go
```

**访问**：http://localhost:8080/hello

**演示步骤**：
1. 打开 `01-hello-world` 目录
2. 查看 `main.go` 文件内容
3. 运行 `go run main.go`
4. 使用 `curl http://localhost:8080/hello` 测试
5. 查看 `README.md` 了解更多细节

**关键点**：
- `gin.Default()`：创建默认引擎（包含日志和恢复中间件）
- `gin.New()`：创建空白引擎（不包含任何中间件）
- `r.Run()`：启动服务器，默认端口 8080
- `c.JSON()`：返回 JSON 响应

#### 项目结构设计

**推荐的目录结构**：

```
project/
├── main.go              # 程序入口
├── config/              # 配置管理
│   └── config.go
├── handlers/            # 处理器（Controller）
│   ├── user_handler.go
│   └── product_handler.go
├── middleware/          # 中间件
│   ├── auth.go
│   └── logger.go
├── models/              # 数据模型
│   └── user.go
├── services/            # 业务逻辑层
│   └── user_service.go
├── utils/               # 工具函数
│   └── response.go
└── go.mod
```

### 2.3 路由系统

参考示例：`examples/02-routing/` 目录

运行示例：

```bash
cd lesson-03/examples/02-routing
go run main.go
```

#### 路由定义

**HTTP 方法**：

```go
r := gin.Default()

// GET 请求
r.GET("/users", getUsers)

// POST 请求
r.POST("/users", createUser)

// PUT 请求
r.PUT("/users/:id", updateUser)

// DELETE 请求
r.DELETE("/users/:id", deleteUser)

// PATCH 请求
r.PATCH("/users/:id", patchUser)

// 任意方法：匹配所有 HTTP 方法（GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS 等）
// 无论客户端使用什么 HTTP 方法请求该路径，都会调用同一个处理函数
r.Any("/users", handleUsers)

// 指定方法：只匹配特定的 HTTP 方法（这里是 "GET"）
// r.Handle 是通用方法，可以指定任意 HTTP 方法字符串
// r.GET、r.POST 等实际上是 r.Handle 的快捷方式
r.Handle("GET", "/users", getUsers)
```

#### 路径参数

**单个参数**：

```go
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})
```

**多个参数**：

```go
r.GET("/users/:id/posts/:postId", func(c *gin.Context) {
    userId := c.Param("id")
    postId := c.Param("postId")
    c.JSON(200, gin.H{
        "userId": userId,
        "postId": postId,
    })
})
```

**通配符参数**：

通配符参数使用 `/*` 符号，可以匹配路径中剩余的所有部分（包括多个斜杠）。

**与普通路径参数的区别**：
- `:id` 只能匹配单个路径段（如 `/users/:id` 匹配 `/users/123`，但不匹配 `/users/123/posts`）
- `/*filepath` 可以匹配多个路径段（如 `/files/*filepath` 可以匹配 `/files/docs/readme.md`）

**使用场景**：
- 文件路径：`/files/*filepath` 可以匹配 `/files/docs/readme.md`、`/files/images/photo.jpg` 等
- 静态资源：`/static/*filepath` 可以匹配所有静态资源路径
- 代理转发：将剩余路径转发到其他服务

**注意事项**：
- 通配符参数必须放在路由路径的最后
- 获取参数时使用 `c.Param("filepath")`，参数名不包含 `*` 号
- 获取到的值会包含前导斜杠（如 `"/docs/readme.md"`）

```go
r.GET("/files/*filepath", func(c *gin.Context) {
    filepath := c.Param("filepath")
    c.JSON(200, gin.H{"filepath": filepath})
})
```

**访问示例**：
- `GET /files/docs/readme.md` → `filepath = "/docs/readme.md"`
- `GET /files/images/photo.jpg` → `filepath = "/images/photo.jpg"`
- `GET /files/a/b/c/d.txt` → `filepath = "/a/b/c/d.txt"`

#### 查询参数

```go
r.GET("/search", func(c *gin.Context) {
    // 获取查询参数
    keyword := c.Query("keyword")
    page := c.DefaultQuery("page", "1")  // 带默认值
    size := c.Query("size")
    
    c.JSON(200, gin.H{
        "keyword": keyword,
        "page": page,
        "size": size,
    })
})
```

**访问示例**：`/search?keyword=golang&page=1&size=10`

#### 表单参数（POST）

表单参数通常用于 HTML 表单提交或 `application/x-www-form-urlencoded` 格式的 POST 请求。

**基础用法**：

```go
r.POST("/login", func(c *gin.Context) {
    // 获取表单参数
    username := c.PostForm("username")
    password := c.PostForm("password")
    // 带默认值
    remember := c.DefaultPostForm("remember", "false")
    
    c.JSON(200, gin.H{
        "username": username,
        "password": password,
        "remember": remember,
    })
})
```

**使用 ShouldBind 绑定到结构体**：

```go
type LoginRequest struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required"`
    Remember bool   `form:"remember"`
}

r.POST("/login", func(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理登录逻辑...
    c.JSON(200, gin.H{"message": "Login success"})
})
```

**访问示例**：
```bash
# 使用 curl 发送表单数据
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=123456&remember=true"
```

**注意事项**：
- `c.PostForm()` 获取单个表单字段
- `c.DefaultPostForm()` 获取表单字段，如果不存在则返回默认值
- `c.ShouldBind()` 会自动识别 Content-Type，支持表单、JSON、XML 等多种格式
- 表单参数适用于 `application/x-www-form-urlencoded` 和 `multipart/form-data` 格式

#### JSON 参数（POST）

JSON 参数用于接收 JSON 格式的请求体，这是现代 API 开发中最常用的方式。

**基础用法**：

```go
r.POST("/users", func(c *gin.Context) {
    // 获取原始 JSON 数据
    jsonData, err := c.GetRawData()
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 手动解析 JSON（不推荐）
    // var user map[string]interface{}
    // json.Unmarshal(jsonData, &user)
    
    c.JSON(200, gin.H{"data": string(jsonData)})
})
```

**使用 ShouldBindJSON 绑定到结构体（推荐）**：

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=120"`
}

r.POST("/users", func(c *gin.Context) {
    var req CreateUserRequest
    // ShouldBindJSON 会自动解析 JSON 并验证
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理业务逻辑...
    c.JSON(201, gin.H{
        "id":    1,
        "name":  req.Name,
        "email": req.Email,
        "age":   req.Age,
    })
})
```

**访问示例**：
```bash
# 使用 curl 发送 JSON 数据
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":25}'
```

**注意事项**：
- `c.ShouldBindJSON()` 会自动解析 JSON 并验证字段
- 结构体标签 `json:"name"` 指定 JSON 字段名
- `binding:"required"` 等标签用于验证（详见参数绑定部分）
- Content-Type 必须是 `application/json`
- 如果 JSON 格式错误或验证失败，会返回错误信息

**表单参数 vs JSON 参数**：

| 特性 | 表单参数 | JSON 参数 |
|------|---------|----------|
| Content-Type | `application/x-www-form-urlencoded` | `application/json` |
| 获取方法 | `c.PostForm()` / `c.ShouldBind()` | `c.ShouldBindJSON()` |
| 适用场景 | HTML 表单、文件上传 | RESTful API、前后端分离 |
| 数据格式 | 键值对 | JSON 对象 |
| 嵌套数据 | 支持（有限） | 完全支持 |

#### 路由分组

**基础分组**：

```go
v1 := r.Group("/api/v1")
{
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
    v1.GET("/products", getProducts)
}
```

**带中间件的分组**：

```go
v1 := r.Group("/api/v1")
v1.Use(authMiddleware())  // 应用到整个分组
{
    v1.GET("/users", getUsers)
    v1.POST("/users", createUser)
}

public := r.Group("/api/public")
{
    public.GET("/products", getProducts)  // 不需要认证
}
```

**嵌套分组**：

```go
api := r.Group("/api")
{
    v1 := api.Group("/v1")
    {
        v1.GET("/users", getUsers)
    }
    
    v2 := api.Group("/v2")
    {
        v2.GET("/users", getUsersV2)
    }
}
```

#### RESTful API 设计

**商品管理 API 示例**：

```go
products := r.Group("/api/v1/products")
{
    products.GET("", listProducts)           // GET /api/v1/products
    products.GET("/:id", getProduct)         // GET /api/v1/products/:id
    products.POST("", createProduct)         // POST /api/v1/products
    products.PUT("/:id", updateProduct)      // PUT /api/v1/products/:id
    products.DELETE("/:id", deleteProduct)   // DELETE /api/v1/products/:id
}
```

### 2.4 请求处理

参考示例：`examples/03-request-handling/` 目录

运行示例：

```bash
cd lesson-03/examples/03-request-handling
go run main.go
```

#### Handler 函数

**函数签名**：

```go
func handler(c *gin.Context) {
    // 处理逻辑
}
```

**gin.Context 核心方法**：

- `c.Request`：原始 `*http.Request`
- `c.Writer`：原始 `http.ResponseWriter`
- `c.Param(key)`：获取路径参数
- `c.Query(key)`：获取查询参数
- `c.PostForm(key)`：获取表单参数
- `c.GetRawData()`：获取原始请求体
- `c.ShouldBindJSON(obj)`：绑定 JSON 到结构体
- `c.ShouldBindQuery(obj)`：绑定查询参数到结构体
- `c.ShouldBindUri(obj)`：绑定路径参数到结构体
- `c.JSON(code, obj)`：返回 JSON 响应
- `c.XML(code, obj)`：返回 XML 响应
- `c.String(code, format, values...)`：返回字符串响应
- `c.HTML(code, name, obj)`：返回 HTML 响应
- `c.Redirect(code, location)`：重定向
- `c.File(filepath)`：返回文件
- `c.Data(code, contentType, data)`：返回原始数据

#### 参数绑定

**JSON 绑定**：

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=120"`
}

func createUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理请求...
    c.JSON(200, gin.H{"message": "User created"})
}
```

**查询参数绑定**：

```go
type ListProductsRequest struct {
    Page  int    `form:"page" binding:"gte=1"`
    Size  int    `form:"size" binding:"gte=1,lte=100"`
    Keyword string `form:"keyword"`
}

func listProducts(c *gin.Context) {
    var req ListProductsRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理请求...
}
```

**路径参数绑定**：

```go
type GetUserRequest struct {
    ID int `uri:"id" binding:"required"`
}

func getUser(c *gin.Context) {
    var req GetUserRequest
    if err := c.ShouldBindUri(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理请求...
}
```

**表单绑定**：

```go
type LoginRequest struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required"`
}

func login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理登录...
}
```

#### 验证标签

**常用验证标签**：

- `required`：必填
- `email`：邮箱格式
- `url`：URL 格式
- `min=10`：最小值
- `max=100`：最大值
- `gte=0`：大于等于
- `lte=120`：小于等于
- `len=10`：长度等于
- `oneof=red green blue`：枚举值

**自定义验证**：

```go
import "github.com/go-playground/validator/v10"

// 注册自定义验证器
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("customtag", customValidator)
}
```

#### 响应处理

**JSON 响应**：

```go
c.JSON(200, gin.H{
    "code": 200,
    "message": "success",
    "data": user,
})
```

**自定义响应结构**：

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{
        Code:    200,
        Message: "success",
        Data:    data,
    })
}

func error(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Code:    code,
        Message: message,
    })
}
```

**XML 响应**：

```go
c.XML(200, gin.H{"message": "Hello"})
```

**字符串响应**：

```go
c.String(200, "Hello %s", name)
```

**HTML 响应**：

```go
r.LoadHTMLGlob("templates/*")
c.HTML(200, "index.html", gin.H{"title": "Home"})
```

### 2.5 实践：商品 API 设计

**需求**：
- 商品列表（分页、搜索）
- 商品详情
- 创建商品
- 更新商品
- 删除商品

**实现步骤**：

1. 定义商品模型
2. 定义请求/响应结构
3. 实现各个 Handler
4. 注册路由

**参考代码**：`examples/02-routing/` 目录

**实践步骤**：
1. 打开 `02-routing/main.go`
2. 找到商品管理相关代码
3. 运行并测试商品 API

**测试命令**：

```bash
# 测试商品列表
curl "http://localhost:8080/api/v1/products?page=1&size=10"

# 测试商品详情
curl http://localhost:8080/api/v1/products/1

# 测试创建商品
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"商品名称","price":99.99}'
```

## 3. 总结

### 核心要点

1. **Gin 安装**：每个示例都是独立工程，使用 `go run main.go` 即可
2. **创建引擎**：`gin.Default()` 或 `gin.New()`
3. **路由定义**：`r.GET()`, `r.POST()` 等
4. **路由分组**：`r.Group()` 组织路由
5. **参数获取**：`c.Param()`, `c.Query()`, `c.ShouldBindJSON()`
6. **响应返回**：`c.JSON()`, `c.String()`, `c.HTML()`

### 示例索引

| 示例 | 目录 | 主要内容 |
|------|------|----------|
| Hello World | `01-hello-world/` | Gin 基础、第一个应用 |
| 路由系统 | `02-routing/` | 路由定义、分组、参数 |
| 请求处理 | `03-request-handling/` | 参数绑定、验证、响应 |

### 实践建议

1. **学习顺序**：按 01 → 02 → 03 的顺序学习
2. **实践方式**：每个示例都先运行起来，再理解代码
3. **测试工具**：使用 curl 或 Postman 测试 API
4. **动手练习**：基于示例实现自己的 API

### 最佳实践

1. **统一响应格式**：定义统一的响应结构
2. **参数验证**：使用 binding 标签进行验证
3. **错误处理**：统一错误处理逻辑
4. **路由分组**：按功能模块分组路由
5. **代码组织**：按 MVC 模式组织代码结构

### 常见问题

**Q: Gin 和标准库 net/http 有什么区别？**
A: Gin 基于 net/http，提供了路由、中间件、参数绑定等高级特性，开发效率更高。

**Q: 如何设置自定义端口？**
A: `r.Run(":3000")` 或使用环境变量 `PORT`。

**Q: 如何处理 404？**
A: 使用 `r.NoRoute()` 定义 404 处理函数。

**Q: 如何获取原始请求体？**
A: 使用 `c.GetRawData()` 或 `c.Request.Body`。

---

**下一节：中间件与高级特性** 🚀

