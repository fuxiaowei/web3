# 第3期教案：中间件与高级特性

## 1. 学习目标

通过本课程的学习，你将能够：

- 理解中间件的概念和工作原理
- 编写自定义中间件并理解中间件链
- 实现 JWT 认证和权限控制
- 处理文件上传下载
- 使用 Viper 进行配置管理
- 理解和使用 Protocol Buffers（Protobuf）
- 理解和使用 gRPC 构建高性能 RPC 服务

## 2. 核心知识点

### 2.1 中间件机制

参考示例：`examples/04-middleware/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/04-middleware
> go run main.go
> ```
>
> 边运行边观察控制台的日志输出，帮助理解中间件执行顺序

#### 什么是中间件？

中间件（Middleware）是在请求处理前后执行的函数，用于处理横切关注点（Cross-cutting Concerns），如日志记录、认证、CORS、限流等。

#### 中间件的工作原理

```
请求 → 中间件1 → 中间件2 → Handler → 中间件2 → 中间件1 → 响应
```

**执行顺序**：
1. 请求到达时，按注册顺序执行中间件的前置代码
2. 调用 `c.Next()` 进入下一个中间件或 Handler
3. Handler 处理完成后，按相反顺序执行中间件的后置代码

#### 自定义中间件

**基础中间件**：

```go
func loggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        start := time.Now()
        path := c.Request.URL.Path
        
        // 进入下一个处理函数
        c.Next()
        
        // 后置处理
        latency := time.Since(start)
        status := c.Writer.Status()
        fmt.Printf("[%s] %s %d %v\n", c.Request.Method, path, status, latency)
    }
}
```

**使用中间件**：

```go
r := gin.Default()
r.Use(loggerMiddleware())  // 全局中间件

r.GET("/test", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "test"})
})
```

#### 中间件链

**多个中间件**：

```go
r := gin.Default()
r.Use(middleware1())
r.Use(middleware2())
r.Use(middleware3())

r.GET("/test", handler)
```

**执行顺序**：middleware1 → middleware2 → middleware3 → handler → middleware3 → middleware2 → middleware1

**分组中间件**：

```go
r := gin.Default()

// 全局中间件
r.Use(loggerMiddleware())

// 分组中间件
api := r.Group("/api")
api.Use(authMiddleware())
{
    api.GET("/users", getUsers)
}

// 不需要认证的路由
r.GET("/public", publicHandler)
```

#### 常用中间件

**1. 日志中间件**

```go
func loggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}
```

**2. 恢复中间件**

恢复中间件（Recovery Middleware）是 Web 应用中**必不可少**的中间件，用于捕获和处理 panic，防止程序崩溃。

**作用**：

1. **防止程序崩溃**：当处理函数（Handler）或中间件发生 panic 时，恢复中间件会捕获它，避免整个服务崩溃
2. **优雅的错误处理**：捕获 panic 后，返回统一的错误响应（如 500 错误），而不是直接崩溃或返回空响应
3. **提升服务稳定性**：即使某个请求处理出错，也不会影响其他请求，服务可以继续正常运行
4. **生产环境必需**：在生产环境中，恢复中间件是必需的，可以避免单个请求错误导致整个服务不可用

**使用场景**：

```go
// 没有恢复中间件时
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    // 如果发生数组越界、空指针等错误，会导致 panic
    userID := id[10] // panic！整个服务崩溃 💥
})

// 有恢复中间件时
r.Use(recoveryMiddleware())
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    userID := id[10] // panic 发生
    // 恢复中间件捕获 panic，返回 500 错误
    // 服务继续运行 ✅
})
```

**实现方式**：

```go
func recoveryMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        // recovered 是 panic 的值，可以用于日志记录
        // 例如：log.Printf("Panic recovered: %v", recovered)
        
        c.JSON(500, gin.H{
            "error": "Internal Server Error",
        })
        c.Abort()
    })
}
```

**Gin 的默认行为**：

- `gin.Default()`：自动包含日志和恢复中间件
- `gin.New()`：不包含任何中间件，需要手动添加恢复中间件

**最佳实践**：

1. 在生产环境中，**必须**使用恢复中间件
2. 可以在恢复中间件中记录 panic 信息，便于排查问题
3. 返回统一的错误格式，不要暴露内部错误详情给客户端

**3. CORS 中间件**

```go
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
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

**或使用官方 CORS 中间件**：

```bash
go get github.com/gin-contrib/cors
```

```go
import "github.com/gin-contrib/cors"

config := cors.DefaultConfig()
config.AllowOrigins = []string{"http://localhost:3000"}
config.AllowCredentials = true
r.Use(cors.New(config))
```

**4. 限流中间件**

```go
import "golang.org/x/time/rate"

func rateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(10, 20) // 每秒10个请求，突发20个
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 2.2 JWT 认证

参考示例：`examples/05-jwt-auth/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/05-jwt-auth
> go run main.go
> ```
>
> 准备好测试命令，演示完整的认证流程

#### JWT 简介

JWT（JSON Web Token）是一种开放标准（RFC 7519），用于在各方之间安全地传输信息。JWT 由三部分组成：
- Header（头部）
- Payload（载荷）
- Signature（签名）

#### 安装依赖

```bash
go get github.com/golang-jwt/jwt/v5
```

#### Token 生成

```go
import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key")

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
```

#### Token 验证

```go
func parseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
```

#### 认证中间件

```go
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Header 获取 Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // 提取 Token（Bearer <token>）
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // 验证 Token
        claims, err := parseToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // 将用户信息存储到 Context
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        
        c.Next()
    }
}
```

#### 登录接口

```go
func login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 验证用户名密码（示例，实际应从数据库查询）
    if req.Username == "admin" && req.Password == "admin123" {
        token, err := generateToken(1, req.Username)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to generate token"})
            return
        }
        
        c.JSON(200, gin.H{
            "token": token,
            "user": gin.H{
                "id":       1,
                "username": req.Username,
            },
        })
        return
    }
    
    c.JSON(401, gin.H{"error": "Invalid credentials"})
}
```

#### 使用认证中间件

```go
r := gin.Default()

// 公开路由
r.POST("/api/login", login)

// 需要认证的路由
api := r.Group("/api")
api.Use(authMiddleware())
{
    api.GET("/profile", getProfile)
    api.PUT("/profile", updateProfile)
}
```

#### 权限控制

```go
func requireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists || userRole != role {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用
admin := api.Group("/admin")
admin.Use(requireRole("admin"))
{
    admin.GET("/users", listUsers)
    admin.DELETE("/users/:id", deleteUser)
}
```

### 2.3 文件上传下载

参考示例：`examples/06-file-upload/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/06-file-upload
> go run main.go
> ```
>
> 准备好测试文件，实时演示上传和下载功能

#### 单文件上传

```go
func uploadFile(c *gin.Context) {
    // 获取上传的文件
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 保存文件
    filename := filepath.Base(file.Filename)
    dst := filepath.Join("uploads", filename)
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "message": "File uploaded successfully",
        "filename": filename,
    })
}

r.POST("/upload", uploadFile)
```

#### 多文件上传

```go
func uploadFiles(c *gin.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    files := form.File["files"]
    var filenames []string
    
    for _, file := range files {
        filename := filepath.Join("uploads", filepath.Base(file.Filename))
        if err := c.SaveUploadedFile(file, filename); err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        filenames = append(filenames, filename)
    }
    
    c.JSON(200, gin.H{
        "message": "Files uploaded successfully",
        "files": filenames,
    })
}

r.POST("/upload-multiple", uploadFiles)
```

#### 文件下载

```go
func downloadFile(c *gin.Context) {
    filename := c.Param("filename")
    filepath := filepath.Join("uploads", filename)
    
    // 检查文件是否存在
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"error": "File not found"})
        return
    }
    
    // 设置响应头
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Header("Content-Type", "application/octet-stream")
    
    // 返回文件
    c.File(filepath)
}

r.GET("/download/:filename", downloadFile)
```

#### 静态文件服务

```go
// 提供静态文件服务
r.Static("/static", "./static")

// 提供单个文件
r.StaticFile("/favicon.ico", "./static/favicon.ico")

// 提供文件系统
r.StaticFS("/files", http.Dir("./uploads"))
```

### 2.4 配置管理

参考示例：`examples/07-config/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/07-config
> # 查看配置文件
> cat config.yaml
> # 运行示例
> go run main.go
> ```
>
> 可以现场修改 `config.yaml` 演示配置的作用

#### Viper 简介

Viper 是 Go 语言的配置管理库，支持：
- JSON、TOML、YAML、HCL、envfile、Java properties
- 环境变量
- 命令行参数
- 远程配置系统（etcd、Consul）
- 配置热加载

#### 安装依赖

```bash
go get github.com/spf13/viper
```

#### 基础使用

```go
import "github.com/spf13/viper"

func init() {
    // 设置配置文件名称（不含扩展名）
    viper.SetConfigName("config")
    // 设置配置文件类型
    viper.SetConfigType("yaml")
    // 添加配置文件搜索路径
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME/.app")
    
    // 读取配置文件
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }
}

// 使用配置
port := viper.GetString("server.port")
dbHost := viper.GetString("database.host")
```

#### 配置文件示例（config.yaml）

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  mode: "debug"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  dbname: "mydb"

jwt:
  secret: "your-secret-key"
  expire: 24h
```

#### 环境变量支持

Viper 支持从环境变量读取配置，有两种方式：

**方式一：使用 SetEnvPrefix（自动映射）**

```go
// 自动读取环境变量
viper.AutomaticEnv()

// 设置环境变量前缀
viper.SetEnvPrefix("APP")
// 设置前缀后，Viper 会自动将配置键转换为环境变量名：
// server.port -> APP_SERVER_PORT
// database.host -> APP_DATABASE_HOST
// jwt.secret -> APP_JWT_SECRET
```

**工作原理：**
- Viper 会将配置键中的 `.` 替换为 `_`
- 自动添加前缀 `APP_`
- 转换为大写字母

**使用示例：**

```bash
# 设置环境变量
export APP_SERVER_PORT=9000
export APP_DATABASE_HOST=192.168.1.100
export APP_JWT_SECRET=my-secret-key
```

```go
// 获取配置值（与从配置文件获取方式相同）
port := viper.GetString("server.port")        // 从 APP_SERVER_PORT 读取，值为 "9000"
host := viper.GetString("database.host")      // 从 APP_DATABASE_HOST 读取
secret := viper.GetString("jwt.secret")       // 从 APP_JWT_SECRET 读取
```

**方式二：使用 BindEnv（手动绑定）**

```go
// 绑定环境变量：将配置键 server.port 绑定到环境变量 PORT（不含前缀）
viper.BindEnv("server.port", "PORT")
// 这样 server.port 会直接从环境变量 PORT 读取，而不是 APP_SERVER_PORT
```

**使用示例：**

```bash
# 设置环境变量（注意：这里不需要 APP_ 前缀）
export PORT=9000
```

```go
// 获取配置值
port := viper.GetString("server.port")  // 从 PORT 环境变量读取，值为 "9000"
```

**两种方式的区别：**

| 方式 | 环境变量名 | 是否需要前缀 | 使用场景 |
|------|-----------|-------------|---------|
| `SetEnvPrefix` | `APP_SERVER_PORT` | 是 | 统一管理，避免环境变量冲突 |
| `BindEnv` | `PORT` | 否 | 使用标准环境变量名（如 PORT、HOST） |

**完整示例：**

```go
func init() {
    // 自动读取环境变量
    viper.AutomaticEnv()
    
    // 方式一：设置前缀（推荐）
    viper.SetEnvPrefix("APP")
    // 此时 APP_SERVER_PORT 对应 server.port
    
    // 方式二：手动绑定（可选，会覆盖方式一）
    // viper.BindEnv("server.port", "PORT")
    // 此时 PORT 对应 server.port
    
    // 设置默认值
    viper.SetDefault("server.port", "8080")
}

// 获取配置的多种方式
func getConfig() {
    // 方式1：直接获取字符串
    port := viper.GetString("server.port")
    
    // 方式2：获取并指定类型
    portInt := viper.GetInt("server.port")
    
    // 方式3：获取原始值
    portRaw := viper.Get("server.port")
    
    // 方式4：检查是否存在
    if viper.IsSet("server.port") {
        port := viper.GetString("server.port")
    }
    
    // 方式5：获取并处理默认值（如果不存在）
    port := viper.GetString("server.port")
    if port == "" {
        port = "8080"  // 设置默认值
    }
    // 或者使用 SetDefault 预先设置默认值
}
```

#### 配置结构体

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
    Port string `mapstructure:"port"`
    Host string `mapstructure:"host"`
    Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
    Secret string `mapstructure:"secret"`
    Expire string `mapstructure:"expire"`
}

func LoadConfig() (*Config, error) {
    var config Config
    
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

#### 配置热加载

```go
// 监听配置文件变化
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    log.Println("Config file changed:", e.Name)
    // 重新加载配置
    config, _ := LoadConfig()
    // 更新全局配置
})
```

#### 使用配置

```go
func main() {
    config, err := LoadConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    r := gin.Default()
    
    // 使用配置
    r.Run(config.Server.Host + ":" + config.Server.Port)
}
```

### 2.5 Protocol Buffers（Protobuf）

参考示例：`examples/08-protobuf/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/08-protobuf
> # 如果修改了 user.proto，需要重新生成代码
> protoc --go_out=. --go_opt=paths=source_relative user.proto
> # 运行示例
> go run main.go
> ```
>
> 准备好测试命令，演示 Protobuf 的序列化和反序列化，并与 JSON 格式进行对比

#### Protobuf 简介

Protocol Buffers（简称 Protobuf）是 Google 开发的一种数据序列化格式，具有以下特点：

- **高效**：比 JSON/XML 更小、更快
- **跨语言**：支持多种编程语言
- **类型安全**：强类型定义，编译时检查
- **向后兼容**：支持字段版本演进
- **广泛应用**：gRPC、微服务通信、数据存储

#### 安装工具

**1. 安装 protoc 编译器**

```bash
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler

# 或从源码编译
# https://github.com/protocolbuffers/protobuf/releases
```

**2. 安装 Go 插件**

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**3. 安装依赖**

```bash
go get google.golang.org/protobuf/proto
```

#### 定义 .proto 文件

**示例：user.proto**

```protobuf
syntax = "proto3";

package user;

option go_package = "./pb";

// 用户消息定义
message User {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
  bool active = 5;
  repeated string tags = 6;  // 数组
  map<string, string> metadata = 7;  // 映射
}

// 用户列表
message UserList {
  repeated User users = 1;
  int32 total = 2;
}

// 创建用户请求
message CreateUserRequest {
  string username = 1;
  string email = 2;
  int32 age = 3;
}

// 创建用户响应
message CreateUserResponse {
  User user = 1;
  bool success = 2;
  string message = 3;
}
```

**字段编号规则**：
- 1-15：常用字段（占用 1 字节）
- 16-2047：不常用字段（占用 2 字节）
- 不要使用已删除的字段编号

#### 编译 .proto 文件

```bash
# 生成 Go 代码
protoc --go_out=. --go_opt=paths=source_relative user.proto

# 如果使用 gRPC
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto
```

**生成的文件**：`user.pb.go`

#### 在 Gin 中使用 Protobuf

**1. 返回 Protobuf 响应**

```go
import (
    "github.com/gin-gonic/gin"
    "google.golang.org/protobuf/proto"
    "gin-examples/pb"
)

func getUserProto(c *gin.Context) {
    user := &pb.User{
        Id:       1,
        Username:  "alice",
        Email:     "alice@example.com",
        Age:       25,
        Active:    true,
        Tags:      []string{"admin", "developer"},
        Metadata:  map[string]string{"department": "engineering"},
    }
    
    // 设置响应头
    c.Header("Content-Type", "application/x-protobuf")
    
    // 序列化并返回
    data, err := proto.Marshal(user)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.Data(200, "application/x-protobuf", data)
}
```

**2. 解析 Protobuf 请求**

```go
func createUserProto(c *gin.Context) {
    // 读取原始数据
    data, err := c.GetRawData()
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 反序列化
    var req pb.CreateUserRequest
    if err := proto.Unmarshal(data, &req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 处理请求
    user := &pb.User{
        Id:       1,
        Username:  req.Username,
        Email:     req.Email,
        Age:       req.Age,
        Active:    true,
    }
    
    resp := &pb.CreateUserResponse{
        User:    user,
        Success: true,
        Message: "User created successfully",
    }
    
    // 返回 Protobuf 响应
    data, _ = proto.Marshal(resp)
    c.Data(200, "application/x-protobuf", data)
}
```

**3. 注册路由**

```go
r := gin.Default()

// Protobuf API
r.GET("/api/proto/user", getUserProto)
r.POST("/api/proto/user", createUserProto)
```

#### Protobuf vs JSON

**性能对比**：

| 特性 | Protobuf | JSON |
|------|----------|------|
| 大小 | 小（二进制） | 大（文本） |
| 速度 | 快 | 慢 |
| 可读性 | 差 | 好 |
| 类型安全 | 是 | 否 |
| 浏览器支持 | 需要解码 | 原生支持 |

**使用场景**：

- **Protobuf**：微服务通信、高性能 API、数据存储
- **JSON**：Web API、前端交互、配置文件

#### 版本演进

**向后兼容规则**：

1. ✅ 可以添加新字段（使用新编号）
2. ✅ 可以删除可选字段
3. ✅ 可以重命名字段（编号不变）
4. ❌ 不能修改字段类型
5. ❌ 不能重用已删除的字段编号

**示例**：

```protobuf
// v1
message User {
  int64 id = 1;
  string name = 2;
}

// v2（兼容 v1）
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;  // 新增字段
  // int32 age = 4;   // 已删除，但编号不能重用
}
```

#### 最佳实践

1. **字段编号**：常用字段使用 1-15
2. **命名规范**：使用驼峰命名（CamelCase）
3. **包管理**：使用 `option go_package` 指定包路径
4. **版本控制**：将 .proto 文件纳入版本控制
5. **文档注释**：使用注释说明字段用途

```protobuf
// 用户信息
message User {
  int64 id = 1;  // 用户ID
  string username = 2;  // 用户名
  string email = 3;  // 邮箱地址
}
```

### 2.6 gRPC

参考示例：`examples/09-grpc/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/09-grpc
> # 如果修改了 user.proto，需要重新生成代码
> protoc --go_out=. --go_opt=paths=source_relative \
>        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
>        user.proto
> # 启动服务端（终端1）
> go run . -mode=server
> # 运行客户端（终端2）
> go run . -mode=client
> ```
>
> 准备好两个终端，分别运行服务端和客户端，演示完整的 gRPC 调用流程

#### gRPC 简介

gRPC（gRPC Remote Procedure Calls）是 Google 开发的高性能、开源的 RPC 框架，具有以下特点：

- **高性能**：基于 HTTP/2，支持多路复用、流式传输
- **跨语言**：支持多种编程语言（Go、Java、Python、C++ 等）
- **类型安全**：使用 Protobuf 定义接口，编译时检查
- **流式传输**：支持服务端流、客户端流、双向流
- **自动代码生成**：从 `.proto` 文件自动生成客户端和服务端代码

**应用场景**：
- 微服务通信
- 高性能 API
- 实时数据推送
- 跨语言服务调用

#### 安装依赖

**1. 安装 protoc 编译器**（如果还没有安装）

```bash
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler
```

**2. 安装 Go 插件**

```bash
# 安装 protobuf 编译器插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 gRPC 插件
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 确保 PATH 包含 Go bin 目录
export PATH="$PATH:$(go env GOPATH)/bin"
```

**3. 安装 gRPC 依赖**

```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf/proto
```

#### 定义 gRPC 服务

在 `.proto` 文件中定义服务：

**示例：user.proto**

```protobuf
syntax = "proto3";

package user;

option go_package = "./pb";

// 用户服务定义
service UserService {
  // 获取单个用户（一元 RPC）
  rpc GetUser(GetUserRequest) returns (User);
  
  // 创建用户（一元 RPC）
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  
  // 获取用户列表（一元 RPC）
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  
  // 流式获取用户（服务端流）
  rpc StreamUsers(StreamUsersRequest) returns (stream User);
  
  // 批量创建用户（客户端流）
  rpc BatchCreateUsers(stream CreateUserRequest) returns (BatchCreateUsersResponse);
  
  // 双向流（聊天式交互）
  rpc ChatUsers(stream ChatMessage) returns (stream ChatMessage);
}

// 用户消息定义
message User {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
  bool active = 5;
}

// 获取用户请求
message GetUserRequest {
  int64 id = 1;
}

// 创建用户请求
message CreateUserRequest {
  string username = 1;
  string email = 2;
  int32 age = 3;
}

// 创建用户响应
message CreateUserResponse {
  User user = 1;
  bool success = 2;
  string message = 3;
}
```

#### 生成 gRPC 代码

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto
```

**生成的文件**：
- `pb/user.pb.go` - Protobuf 消息定义
- `pb/user_grpc.pb.go` - gRPC 服务定义

#### 实现 gRPC 服务端

**1. 实现服务接口**

```go
package main

import (
    "context"
    "log"
    "net"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    "grpc/pb"
)

// UserServiceServer 实现 UserService 接口
type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
    users map[int64]*pb.User
    nextID int64
}

// NewUserServiceServer 创建新的用户服务实例
func NewUserServiceServer() *UserServiceServer {
    return &UserServiceServer{
        users:  make(map[int64]*pb.User),
        nextID: 1,
    }
}

// GetUser 获取单个用户
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, exists := s.users[req.Id]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "User with id %d not found", req.Id)
    }
    return user, nil
}

// CreateUser 创建用户
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    if req.Username == "" || req.Email == "" {
        return nil, status.Errorf(codes.InvalidArgument, "Username and email are required")
    }
    
    user := &pb.User{
        Id:       s.nextID,
        Username: req.Username,
        Email:    req.Email,
        Age:      req.Age,
        Active:   true,
    }
    
    s.users[s.nextID] = user
    s.nextID++
    
    return &pb.CreateUserResponse{
        User:    user,
        Success: true,
        Message: "User created successfully",
    }, nil
}
```

**2. 启动 gRPC 服务器**

```go
func main() {
    // 监听端口
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    // 创建 gRPC 服务器
    grpcServer := grpc.NewServer()
    
    // 注册服务
    userService := NewUserServiceServer()
    pb.RegisterUserServiceServer(grpcServer, userService)
    
    // 启动服务器
    log.Println("gRPC server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

#### 实现 gRPC 客户端

```go
package main

import (
    "context"
    "log"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    "grpc/pb"
)

func main() {
    // 连接到服务器
    conn, err := grpc.NewClient(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // 创建客户端
    client := pb.NewUserServiceClient(conn)
    
    // 创建上下文（带超时）
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 调用 GetUser
    req := &pb.GetUserRequest{Id: 1}
    user, err := client.GetUser(ctx, req)
    if err != nil {
        log.Fatalf("GetUser failed: %v", err)
    }
    log.Printf("User: %+v", user)
    
    // 调用 CreateUser
    createReq := &pb.CreateUserRequest{
        Username: "alice",
        Email:    "alice@example.com",
        Age:      25,
    }
    resp, err := client.CreateUser(ctx, createReq)
    if err != nil {
        log.Fatalf("CreateUser failed: %v", err)
    }
    log.Printf("Created user: %+v", resp)
}
```

#### 四种 RPC 模式

**1. 一元 RPC（Unary RPC）**

客户端发送一个请求，服务端返回一个响应：

```go
// 服务端
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // 处理请求并返回响应
    return user, nil
}

// 客户端
user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
```

**2. 服务端流（Server Streaming）**

客户端发送一个请求，服务端返回一个流：

```go
// 服务端
func (s *UserServiceServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
    for i := 0; i < int(req.Limit); i++ {
        user := &pb.User{Id: int64(i + 1), Username: fmt.Sprintf("user%d", i+1)}
        if err := stream.Send(user); err != nil {
            return err
        }
        time.Sleep(time.Duration(req.IntervalMs) * time.Millisecond)
    }
    return nil
}

// 客户端
stream, err := client.StreamUsers(ctx, &pb.StreamUsersRequest{Limit: 10})
for {
    user, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("StreamUsers failed: %v", err)
    }
    log.Printf("Received user: %+v", user)
}
```

**使用场景**：
- 实时数据推送
- 大文件分块传输
- 服务器端推送通知

**3. 客户端流（Client Streaming）**

客户端发送一个流，服务端返回一个响应：

```go
// 服务端
func (s *UserServiceServer) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
    var users []*pb.User
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        // 处理请求
        user := &pb.User{Username: req.Username, Email: req.Email}
        users = append(users, user)
    }
    return stream.SendAndClose(&pb.BatchCreateUsersResponse{
        Users: users,
        SuccessCount: int32(len(users)),
    })
}

// 客户端
stream, err := client.BatchCreateUsers(ctx)
for _, userReq := range userRequests {
    if err := stream.Send(userReq); err != nil {
        log.Fatalf("Send failed: %v", err)
    }
}
resp, err := stream.CloseAndRecv()
```

**使用场景**：
- 批量上传
- 数据采集
- 日志收集

**4. 双向流（Bidirectional Streaming）**

客户端和服务端都可以随时发送消息：

```go
// 服务端
func (s *UserServiceServer) ChatUsers(stream pb.UserService_ChatUsersServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        // 处理消息并回复
        reply := &pb.ChatMessage{
            UserId:    "server",
            Message:   "Echo: " + msg.Message,
            Timestamp: time.Now().Unix(),
        }
        if err := stream.Send(reply); err != nil {
            return err
        }
    }
}

// 客户端
stream, err := client.ChatUsers(ctx)
// 发送消息
stream.Send(&pb.ChatMessage{Message: "Hello"})
// 接收消息
msg, err := stream.Recv()
```

**使用场景**：
- 实时聊天
- 游戏服务器
- 协作编辑

#### 错误处理

gRPC 使用状态码来表示错误：

```go
import (
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
)

// 服务端返回错误
return nil, status.Errorf(codes.NotFound, "User not found")

// 客户端检查错误
st, ok := status.FromError(err)
if ok {
    switch st.Code() {
    case codes.NotFound:
        log.Println("User not found")
    case codes.InvalidArgument:
        log.Println("Invalid argument")
    }
}
```

**常用状态码**：
- `OK` - 成功
- `InvalidArgument` - 参数无效
- `NotFound` - 资源未找到
- `AlreadyExists` - 资源已存在
- `PermissionDenied` - 权限不足
- `Internal` - 服务器内部错误
- `Unavailable` - 服务不可用

#### gRPC vs REST API

| 特性 | gRPC | REST API |
|------|------|----------|
| 协议 | HTTP/2 | HTTP/1.1 |
| 数据格式 | Protobuf（二进制） | JSON（文本） |
| 性能 | 高 | 中 |
| 流式传输 | 支持 | 有限支持（SSE/WebSocket） |
| 浏览器支持 | 需要 gRPC-Web | 原生支持 |
| 代码生成 | 自动生成 | 手动编写 |
| 类型安全 | 强类型 | 弱类型 |

**使用场景**：
- **gRPC**：微服务通信、高性能 API、实时数据推送
- **REST API**：Web API、前端交互、公开 API

#### 测试工具

**1. grpcurl（命令行工具）**

```bash
# 安装
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 列出服务
grpcurl -plaintext localhost:50051 list

# 调用方法
grpcurl -plaintext -d '{"id": 1}' \
  localhost:50051 user.UserService/GetUser
```

**2. grpcui（Web UI）**

```bash
# 安装
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest

# 启动 Web UI
grpcui -plaintext localhost:50051
```

然后在浏览器中打开 `http://localhost:8080`，可以看到 Web 界面。

#### 安全

**TLS/SSL 加密**（生产环境推荐）：

```go
// 服务端
creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
grpcServer := grpc.NewServer(grpc.Creds(creds))

// 客户端
creds, err := credentials.NewClientTLSFromFile("ca.crt", "")
conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
```

#### 最佳实践

1. **版本控制**：将 `.proto` 文件纳入版本控制
2. **向后兼容**：遵循 Protobuf 版本演进规则
3. **错误处理**：使用标准状态码
4. **超时设置**：为所有 RPC 调用设置超时
5. **连接复用**：客户端应该复用连接
6. **流式传输**：大数据使用流式传输
7. **监控**：监控 RPC 调用指标

## 3. 总结

### 核心要点

1. **中间件**：在请求处理前后执行的函数，用于处理横切关注点
2. **中间件链**：按注册顺序执行，`c.Next()` 进入下一个处理函数
3. **JWT 认证**：使用 JWT 实现无状态认证
4. **文件处理**：使用 `c.FormFile()` 和 `c.SaveUploadedFile()` 处理文件上传
5. **配置管理**：使用 Viper 管理配置文件和环境变量
6. **Protobuf**：高效的数据序列化格式，适用于高性能 API 和微服务通信
7. **gRPC**：高性能 RPC 框架，支持流式传输，适用于微服务通信

### 示例索引

| 示例 | 目录 | 主要内容 |
|------|------|----------|
| 中间件 | `04-middleware/` | 自定义中间件、日志、认证、CORS |
| JWT 认证 | `05-jwt-auth/` | Token 生成、验证、认证中间件 |
| 文件上传 | `06-file-upload/` | 单文件、多文件上传、下载 |
| 配置管理 | `07-config/` | Viper 配置、环境变量 |
| Protobuf | `08-protobuf/` | Protobuf 序列化、Gin 集成 |
| gRPC | `09-grpc/` | gRPC 服务端、客户端、流式传输 |

### 教学提示

1. **演示顺序**：按 04 → 05 → 06 → 07 → 08 → 09 的顺序讲解
2. **实践方式**：
   - 04：观察中间件执行顺序的日志输出
   - 05：完整演示登录-获取token-访问受保护API的流程
   - 06：准备测试文件，实时演示上传下载
   - 07：修改配置文件，演示配置的作用
   - 08：演示 Protobuf 序列化，与 JSON 格式进行对比
   - 09：准备两个终端，分别运行服务端和客户端，演示 gRPC 调用和流式传输
3. **互动环节**：
   - 让学员自己实现一个限流中间件
   - 让学员尝试修改 JWT 过期时间
   - 让学员上传不同类型的文件
4. **作业布置**：
   - 实现一个带权限控制的用户管理系统
   - 实现文件上传的类型和大小验证
   - 使用环境变量配置数据库连接

### 最佳实践

1. **中间件设计**：保持中间件职责单一
2. **认证安全**：使用 HTTPS，Token 设置合理的过期时间
3. **文件安全**：验证文件类型和大小，防止路径遍历攻击
4. **配置管理**：使用结构体管理配置，支持环境变量覆盖

### 常见问题

**Q: 中间件的执行顺序是什么？**
A: 按注册顺序执行前置代码，Handler 执行后按相反顺序执行后置代码。

**Q: JWT Token 应该存储在哪里？**
A: 通常存储在客户端的 localStorage 或 cookie 中。使用 cookie 更安全（HttpOnly），但需要考虑 CSRF 防护。

**Q: 如何实现配置热加载？**
A: 使用 `viper.WatchConfig()` 监听配置文件变化，在回调中重新加载配置。

**Q: 文件上传大小限制？**
A: 使用 `r.MaxMultipartMemory = 8 << 20` 设置最大上传大小（8MB）。

**Q: 什么时候使用 Protobuf，什么时候使用 JSON？**
A: Protobuf 适用于微服务通信、高性能 API、数据存储；JSON 适用于 Web API、前端交互、配置文件。

**Q: Protobuf 如何实现版本兼容？**
A: 可以添加新字段（使用新编号），可以删除可选字段，但不能修改字段类型或重用已删除的字段编号。

**Q: gRPC 和 REST API 有什么区别？**
A: gRPC 基于 HTTP/2 和 Protobuf，性能更高，支持流式传输，但浏览器支持需要 gRPC-Web；REST API 基于 HTTP/1.1 和 JSON，浏览器原生支持，但性能较低。

**Q: 什么时候使用 gRPC，什么时候使用 REST API？**
A: gRPC 适用于微服务通信、高性能 API、实时数据推送；REST API 适用于 Web API、前端交互、公开 API。

**Q: gRPC 支持哪些流式传输模式？**
A: gRPC 支持四种模式：一元 RPC（请求-响应）、服务端流（请求-流式响应）、客户端流（流式请求-响应）、双向流（流式请求-流式响应）。

**Q: 如何测试 gRPC 服务？**
A: 可以使用 `grpcurl` 命令行工具或 `grpcui` Web 界面来测试 gRPC 服务。

---

**下一节：项目实战** 🚀

