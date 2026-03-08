# 第3期教案：完整项目实战

## 1. 学习目标

通过本课程的学习，你将能够：

- 设计符合 RESTful 规范的 API
- 实现统一的错误处理和响应格式
- 构建完整的用户管理系统 API
- 编写单元测试和 API 测试
- 配置和部署 Web 服务

## 2. 核心知识点

### 2.1 RESTful API 设计规范

参考示例：`examples/project/` 目录

> 💡 **演示提示**：
>
> ```bash
> cd lesson-03/examples/project
> go run main.go
> ```
>
> 这是一个完整的用户管理系统，包含了前面所有示例的知识点

#### RESTful 原则

REST（Representational State Transfer）是一种架构风格，遵循以下原则：

1. **资源导向**：URL 表示资源，而不是动作
2. **HTTP 方法**：使用标准 HTTP 方法（GET、POST、PUT、DELETE）
3. **无状态**：每个请求包含所有必要信息
4. **统一接口**：使用标准 HTTP 状态码和响应格式

#### URL 设计规范

**好的设计**：
```
GET    /api/v1/users          # 获取用户列表
GET    /api/v1/users/:id      # 获取单个用户
POST   /api/v1/users          # 创建用户
PUT    /api/v1/users/:id      # 更新用户（完整更新）
PATCH  /api/v1/users/:id      # 更新用户（部分更新）
DELETE /api/v1/users/:id      # 删除用户
```

**不好的设计**：
```
GET    /api/v1/getUsers
POST   /api/v1/createUser
POST   /api/v1/updateUser
POST   /api/v1/deleteUser
```

#### HTTP 状态码

**常用状态码**：

- `200 OK`：请求成功
- `201 Created`：资源创建成功
- `204 No Content`：请求成功，无返回内容
- `400 Bad Request`：请求参数错误
- `401 Unauthorized`：未认证
- `403 Forbidden`：无权限
- `404 Not Found`：资源不存在
- `409 Conflict`：资源冲突
- `422 Unprocessable Entity`：验证失败
- `500 Internal Server Error`：服务器错误

#### 统一响应格式

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{
        Code:    200,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Code:    code,
        Message: message,
        Error:   message,
    })
}

func ValidationError(c *gin.Context, errors map[string]string) {
    c.JSON(422, Response{
        Code:    422,
        Message: "validation failed",
        Error:   errors,
    })
}
```

#### 错误处理规范

```go
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func HandleError(c *gin.Context, err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        Error(c, appErr.Code, appErr.Message)
        return
    }
    
    // 未知错误，记录日志但不暴露给客户端
    log.Printf("Internal error: %v", err)
    Error(c, 500, "Internal server error")
}
```

### 2.2 用户系统 API 实现

#### 项目结构

```
project/
├── main.go
├── config/
│   ├── config.go
│   └── viper.go
├── models/
│   └── user.go
├── handlers/
│   ├── user_handler.go
│   └── auth_handler.go
├── middleware/
│   ├── auth.go
│   ├── logger.go
│   └── cors.go
├── services/
│   └── user_service.go
├── utils/
│   ├── response.go
│   └── jwt.go
└── tests/
    └── api_test.go
```

#### 用户模型

```go
type User struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"uniqueIndex;not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
    Email string `json:"email" binding:"omitempty,email"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type UserResponse struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 用户服务层

```go
type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
    // 检查用户名是否已存在
    var existingUser User
    if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
        return nil, &AppError{Code: 409, Message: "Username already exists"}
    }
    
    // 检查邮箱是否已存在
    if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        return nil, &AppError{Code: 409, Message: "Email already exists"}
    }
    
    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // 创建用户
    user := User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
    }
    
    if err := s.db.Create(&user).Error; err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (s *UserService) GetUserByID(id int) (*User, error) {
    var user User
    if err := s.db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &AppError{Code: 404, Message: "User not found"}
        }
        return nil, err
    }
    return &user, nil
}

func (s *UserService) Authenticate(username, password string) (*User, error) {
    var user User
    if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &AppError{Code: 401, Message: "Invalid credentials"}
        }
        return nil, err
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, &AppError{Code: 401, Message: "Invalid credentials"}
    }
    
    return &user, nil
}
```

#### 用户处理器

```go
type UserHandler struct {
    userService *UserService
    jwtSecret   []byte
}

func NewUserHandler(userService *UserService, jwtSecret []byte) *UserHandler {
    return &UserHandler{
        userService: userService,
        jwtSecret:   jwtSecret,
    }
}

func (h *UserHandler) Register(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ValidationError(c, parseValidationErrors(err))
        return
    }
    
    user, err := h.userService.CreateUser(req)
    if err != nil {
        HandleError(c, err)
        return
    }
    
    Success(c, UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    })
}

func (h *UserHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ValidationError(c, parseValidationErrors(err))
        return
    }
    
    user, err := h.userService.Authenticate(req.Username, req.Password)
    if err != nil {
        HandleError(c, err)
        return
    }
    
    token, err := generateToken(h.jwtSecret, user.ID, user.Username)
    if err != nil {
        HandleError(c, err)
        return
    }
    
    Success(c, gin.H{
        "token": token,
        "user": UserResponse{
            ID:        user.ID,
            Username:  user.Username,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
    })
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    user, err := h.userService.GetUserByID(userID.(int))
    if err != nil {
        HandleError(c, err)
        return
    }
    
    Success(c, UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    })
}
```

#### 路由注册

```go
func setupRoutes(r *gin.Engine, db *gorm.DB, jwtSecret []byte) {
    // 初始化服务
    userService := NewUserService(db)
    userHandler := NewUserHandler(userService, jwtSecret)
    
    // 公开路由
    public := r.Group("/api/v1")
    {
        public.POST("/users/register", userHandler.Register)
        public.POST("/users/login", userHandler.Login)
    }
    
    // 需要认证的路由
    protected := r.Group("/api/v1")
    protected.Use(authMiddleware(jwtSecret))
    {
        protected.GET("/users/me", userHandler.GetProfile)
        protected.PUT("/users/me", userHandler.UpdateProfile)
    }
}
```

### 2.3 项目部署与测试

#### 单元测试

```go
func TestUserService_CreateUser(t *testing.T) {
    // 创建测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    service := NewUserService(db)
    
    req := CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    user, err := service.CreateUser(req)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Username, user.Username)
    assert.Equal(t, req.Email, user.Email)
}
```

#### API 测试

```go
func TestRegisterAPI(t *testing.T) {
    router := setupTestRouter()
    
    req := CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    body, _ := json.Marshal(req)
    w := httptest.NewRecorder()
    httpReq, _ := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(body))
    httpReq.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, httpReq)
    
    assert.Equal(t, 200, w.Code)
    
    var response Response
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, 200, response.Code)
}
```

#### 部署配置

**Dockerfile**：

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./main"]
```

**docker-compose.yml**：

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=3306
    depends_on:
      - db
  
  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=mydb
    volumes:
      - db_data:/var/lib/mysql

volumes:
  db_data:
```

#### 环境变量配置

```bash
# .env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_NAME=mydb
JWT_SECRET=your-secret-key
```

#### 健康检查

```go
func healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
        "timestamp": time.Now().Unix(),
    })
}

r.GET("/health", healthCheck)
```

## 3. 总结

### 核心要点

1. **RESTful 设计**：遵循 REST 原则，使用标准 HTTP 方法和状态码
2. **统一响应**：定义统一的响应格式，便于前端处理
3. **错误处理**：统一错误处理逻辑，记录日志但不暴露敏感信息
4. **分层架构**：Handler → Service → Model，职责清晰
5. **测试覆盖**：编写单元测试和 API 测试，保证代码质量

### 项目结构

```
project/
├── main.go              # 程序入口
├── config/              # 配置管理
│   └── config.go
├── handlers/            # 处理器层（Controller）
│   └── user_handler.go
├── services/            # 业务逻辑层
│   └── user_service.go
├── models/              # 数据模型层
│   └── user.go
├── middleware/          # 中间件
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
└── utils/               # 工具函数
    ├── jwt.go
    ├── response.go
    └── errors.go
```

### 教学提示

1. **项目演示顺序**：
   - 先展示整体目录结构，说明分层架构
   - 从 `main.go` 开始，展示如何组织路由
   - 依次讲解 models → services → handlers
   - 最后讲解 middleware 和 utils

2. **代码讲解方式**：
   - 从 Model 层开始：数据结构定义
   - 到 Service 层：业务逻辑实现
   - 到 Handler 层：HTTP 请求处理
   - 最后是 Middleware：横切关注点

3. **实践演示**：
   ```bash
   # 启动项目
   cd lesson-03/examples/project
   go run main.go
   
   # 测试完整流程
   # 1. 注册用户
   curl -X POST http://localhost:8080/api/v1/register \
     -H "Content-Type: application/json" \
     -d '{"username":"test","email":"test@example.com","password":"123456"}'
   
   # 2. 登录获取 token
   curl -X POST http://localhost:8080/api/v1/login \
     -H "Content-Type: application/json" \
     -d '{"username":"test","password":"123456"}'
   
   # 3. 使用 token 访问受保护的接口
   curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/v1/profile
   ```

4. **扩展练习**：
   - 添加用户角色管理（普通用户、管理员）
   - 添加用户资料更新接口
   - 添加用户列表接口（需要管理员权限）
   - 添加用户头像上传功能
   - 集成数据库（MySQL/PostgreSQL）

### 最佳实践

1. **API 版本控制**：使用 `/api/v1/` 前缀
2. **参数验证**：使用 binding 标签进行验证
3. **密码安全**：使用 bcrypt 加密，不存储明文密码
4. **Token 安全**：设置合理的过期时间，使用 HTTPS
5. **日志记录**：记录关键操作和错误信息
6. **配置管理**：使用环境变量覆盖配置文件

### 常见问题

**Q: 如何实现 API 版本控制？**
A: 使用路由分组，如 `/api/v1/` 和 `/api/v2/`，可以同时支持多个版本。

**Q: 如何处理数据库连接池？**
A: 在应用启动时创建数据库连接，使用连接池配置，避免频繁创建连接。

**Q: 如何实现分页？**
A: 使用 `LIMIT` 和 `OFFSET`，或者使用游标分页（cursor-based pagination）。

**Q: 如何实现限流？**
A: 使用中间件实现限流，可以使用 `golang.org/x/time/rate` 或 Redis 实现分布式限流。

---

**恭喜完成 Gin Web 框架学习！** 🎉

