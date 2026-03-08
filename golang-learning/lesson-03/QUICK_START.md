# 第3期快速上手指南

## 1. 环境准备

### ✅ 必备工具
- Go 1.22+
- Git
- curl 或 Postman（用于测试 API）

### ✅ 可选工具
- VS Code + Go 插件
- Postman 或 Insomnia（API 测试工具）

## 2. 获取代码

> **重要变更**：每个示例现在都是独立的工程！

每个示例目录都包含自己的 `go.mod` 文件，无需手动初始化。直接进入示例目录运行即可：

```bash
cd /Users/zhouxing/Documents/rcc/new-golang/lesson-03/examples
ls -la
```

你会看到如下独立示例目录：
- `01-hello-world/`
- `02-routing/`
- `03-request-handling/`
- `04-middleware/`
- `05-jwt-auth/`
- `06-file-upload/`
- `07-config/`
- `project/`

## 3. 运行基础示例

### 01 - Hello World

```bash
cd 01-hello-world
go run main.go
```

访问：http://localhost:8080/hello

**测试命令**：
```bash
curl http://localhost:8080/hello
curl http://localhost:8080/ping
```

**查看说明**：
```bash
cat README.md
```

### 02 - 路由系统

```bash
cd ../02-routing
go run main.go
```

**测试路由**：
```bash
# 路径参数
curl http://localhost:8080/users/123

# 查询参数
curl "http://localhost:8080/search?keyword=golang&page=2"

# RESTful API
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v1/products/1
```

### 03 - 请求处理

```bash
cd ../03-request-handling
go run main.go
```

**测试请求**：
```bash
# JSON 绑定
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com","age":25}'

# 验证失败测试
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","email":"invalid-email","age":150}'
```

## 4. 运行进阶示例

### 04 - 中间件

```bash
cd ../04-middleware
go run main.go
```

**测试中间件**：
```bash
# 公开路由
curl http://localhost:8080/public

# 需要认证的路由（无 token）
curl http://localhost:8080/api/users

# 需要认证的路由（有效 token）
curl -H "Authorization: Bearer valid-token" \
  http://localhost:8080/api/users
```

观察控制台的日志输出，理解中间件执行顺序。

### 05 - JWT 认证

```bash
cd ../05-jwt-auth
go run main.go
```

**完整测试流程**：

1. 登录获取 Token：
   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'
   ```

2. 复制返回的 token，然后访问受保护资源：
   ```bash
   TOKEN="your-token-here"
   curl http://localhost:8080/api/protected \
     -H "Authorization: Bearer $TOKEN"
   
   curl http://localhost:8080/api/profile \
     -H "Authorization: Bearer $TOKEN"
   ```

### 06 - 文件上传

```bash
cd ../06-file-upload
go run main.go
```

**测试文件上传**：
```bash
# 创建测试文件
echo "Hello, World!" > test.txt

# 单文件上传
curl -X POST http://localhost:8080/api/upload \
  -F "file=@test.txt"

# 多文件上传
curl -X POST http://localhost:8080/api/upload-multiple \
  -F "files=@test.txt" \
  -F "files=@test2.txt"

# 下载文件
curl http://localhost:8080/api/download/test.txt

# 访问文件列表
curl http://localhost:8080/files/
```

### 07 - 配置管理

```bash
cd ../07-config
go run main.go
```

**测试配置**：
```bash
# 查看配置
curl http://localhost:8080/config

# 健康检查
curl http://localhost:8080/health
```

**测试环境变量**：
```bash
# 使用环境变量覆盖配置
export APP_SERVER_PORT=9000
go run main.go
# 服务器将在 9000 端口启动
```

## 5. 项目实战示例

```bash
cd ../project
go run main.go
```

### API 端点

- **健康检查**：GET http://localhost:8080/health
- **用户注册**：POST http://localhost:8080/api/v1/users/register
  ```json
  {
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }
  ```
- **用户登录**：POST http://localhost:8080/api/v1/users/login
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```
- **获取用户信息**：GET http://localhost:8080/api/v1/users/me
  （需要 Authorization Header）

## 6. 常见问题

### Q: 运行时提示找不到包？
A: 确认已经执行 `go mod tidy`，并且所有依赖都已正确安装。

### Q: 端口被占用怎么办？
A: 修改代码中的端口号，或者停止占用 8080 端口的进程：
```bash
# macOS/Linux
lsof -ti:8080 | xargs kill -9
```

### Q: JWT Token 如何测试？
A: 使用 Postman 或 curl，在请求头中添加：
```
Authorization: Bearer YOUR_TOKEN
```

### Q: 文件上传失败？
A: 检查文件路径是否正确，确保有读取权限。上传的文件会保存在 `uploads/` 目录。

## 7. 下一步

- 深入阅读 `courseware/` 下的课件文档
- 尝试独立实现商品管理 API
- 结合 GORM 构建完整的用户管理系统
- 添加 Swagger API 文档
- 实现单元测试和集成测试

---

**完成以上步骤即可顺利开始第3期 Gin Web 框架实战学习！** 🚀

