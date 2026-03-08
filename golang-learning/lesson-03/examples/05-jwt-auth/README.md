# 05 - JWT 认证

## 说明

演示如何使用 JWT（JSON Web Token）实现用户认证和授权。

## 知识点

- JWT 基本概念
- 生成 JWT Token
- 解析和验证 JWT Token
- JWT 中间件实现
- Bearer Token 认证
- Claims 自定义
- Token 过期时间设置

## 运行方式

```bash
go run main.go
```

## 测试

### 1. 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

响应示例：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin"
  }
}
```

### 2. 访问公开接口（不需要 Token）

```bash
curl http://localhost:8080/api/public
```

### 3. 访问受保护接口（需要 Token）

```bash
# 使用上一步获取的 token
TOKEN="your-token-here"

curl http://localhost:8080/api/protected \
  -H "Authorization: Bearer $TOKEN"

curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 4. 无效 Token 测试

```bash
# 没有 Authorization header
curl http://localhost:8080/api/protected

# 无效的 token
curl http://localhost:8080/api/protected \
  -H "Authorization: Bearer invalid-token"
```

## JWT 结构

JWT 由三部分组成，用 `.` 分隔：

```
Header.Payload.Signature
```

### Header（头部）
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Payload（载荷）
```json
{
  "user_id": 1,
  "username": "admin",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Signature（签名）
```
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret
)
```

## 注意事项

1. **密钥安全**：生产环境中应使用强密钥，并存储在环境变量中
2. **Token 有效期**：本示例设置为 24 小时
3. **HTTPS**：生产环境中必须使用 HTTPS 传输 Token
4. **Token 刷新**：实际应用中应实现 Refresh Token 机制
5. **Token 存储**：客户端不应将 Token 存储在 localStorage，建议使用 httpOnly cookie

## 测试凭证

- 用户名：`admin`
- 密码：`admin123`

