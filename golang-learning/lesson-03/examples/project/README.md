# Gin 项目实战示例

这是一个完整的用户管理系统 API 示例，展示了如何使用 Gin 框架构建 RESTful API。

## 功能特性

- ✅ 用户注册
- ✅ 用户登录（JWT 认证）
- ✅ 获取用户信息
- ✅ 更新用户信息
- ✅ 统一响应格式
- ✅ 错误处理
- ✅ 中间件（日志、CORS、认证）

## 项目结构

```
project/
├── main.go              # 程序入口
├── config/              # 配置管理
│   └── config.go
├── handlers/            # 处理器（Controller）
│   └── user_handler.go
├── middleware/          # 中间件
│   ├── auth.go
│   ├── logger.go
│   └── cors.go
├── models/              # 数据模型
│   └── user.go
├── services/            # 业务逻辑层
│   └── user_service.go
└── utils/               # 工具函数
    ├── response.go
    ├── jwt.go
    └── errors.go
```

## 快速开始

### 1. 安装依赖

```bash
cd lesson-03/examples/project
go mod tidy
```

### 2. 运行项目

```bash
go run main.go
```

服务器将在 `http://0.0.0.0:8080` 启动。

### 3. 测试 API

#### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 获取用户信息（需要 Token）

登录后会返回 token，将 `YOUR_TOKEN` 替换为实际的 token：

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 更新用户信息（需要 Token）

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email": "newemail@example.com"
  }'
```

#### 健康检查

```bash
curl http://localhost:8080/health
```

## API 端点

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/users/register` | 用户注册 | 否 |
| POST | `/api/v1/users/login` | 用户登录 | 否 |
| GET | `/api/v1/users/me` | 获取当前用户信息 | 是 |
| PUT | `/api/v1/users/me` | 更新当前用户信息 | 是 |
| GET | `/health` | 健康检查 | 否 |

## 技术栈

- **Web 框架**: Gin v1.10.0
- **ORM**: GORM v1.25.12
- **数据库**: SQLite
- **认证**: JWT (github.com/golang-jwt/jwt/v5)
- **密码加密**: bcrypt (golang.org/x/crypto)

## 注意事项

1. 数据库文件 `users.db` 会在首次运行时自动创建在项目根目录
2. JWT Secret 在生产环境中应该使用环境变量配置
3. 密码使用 bcrypt 加密存储
4. 所有 API 返回统一的 JSON 格式
5. 项目已配置独立的 `go.mod` 文件，可直接在 `project` 目录下运行

## 扩展建议

- [ ] 添加 Swagger API 文档
- [ ] 添加单元测试和集成测试
- [ ] 实现分页查询
- [ ] 添加日志系统（如 logrus）
- [ ] 使用 Viper 进行配置管理
- [ ] 添加 Redis 缓存
- [ ] 实现文件上传功能

