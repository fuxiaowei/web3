# 04 - 中间件

## 说明

演示如何编写和使用 Gin 中间件，包括全局中间件、分组中间件、以及常见的中间件实现。

## 知识点

- 中间件基本概念
- 全局中间件（`r.Use()`）
- 分组中间件
- 日志中间件
- 恢复中间件（panic recovery）
- 认证中间件
- CORS 中间件
- 限流中间件
- `c.Next()` 和 `c.Abort()` 的使用

## 运行方式

```bash
go run main.go
```

## 测试

```bash
# 访问公开路由（不需要认证）
curl http://localhost:8080/public

# 访问测试路由
curl http://localhost:8080/test

# 访问需要认证的路由（无 token）
curl http://localhost:8080/api/users

# 访问需要认证的路由（有效 token）
curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/users

# 访问需要认证的路由（无效 token）
curl -H "Authorization: Bearer invalid-token" http://localhost:8080/api/users
```

## 中间件执行顺序

```
请求 → 日志中间件（前） → 恢复中间件（前） → 认证中间件（前） → 
Handler → 
认证中间件（后） → 恢复中间件（后） → 日志中间件（后） → 响应
```

## 中间件函数签名

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        // ...
        
        c.Next() // 调用下一个中间件/处理函数
        
        // 后置处理
        // ...
    }
}
```

## 关键方法

- `c.Next()`：继续执行后续的中间件和处理函数
- `c.Abort()`：终止请求，不再执行后续的处理函数
- `c.Set(key, value)`：在上下文中设置值
- `c.Get(key)`：从上下文中获取值
- `c.GetHeader(key)`：获取请求头
- `c.Header(key, value)`：设置响应头

