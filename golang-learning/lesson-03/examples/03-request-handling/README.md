# 03 - 请求处理

## 说明

演示如何处理各种类型的 HTTP 请求和响应，包括参数绑定、验证、以及不同格式的响应。

## 知识点

- JSON 请求绑定（`c.ShouldBindJSON()`）
- 查询参数绑定（`c.ShouldBindQuery()`）
- 路径参数绑定（`c.ShouldBindUri()`）
- 表单绑定（`c.ShouldBind()`）
- 参数验证标签（`binding:"required,email,gte,lte"`）
- 统一响应格式
- 多种响应类型（JSON、XML、String、Redirect、Data）

## 运行方式

```bash
go run main.go
```

## 测试

```bash
# JSON 绑定 - 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":25}'

# 验证失败示例
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","email":"invalid-email","age":150}'

# 查询参数绑定
curl "http://localhost:8080/api/products?page=1&size=20&keyword=golang"

# 路径参数绑定
curl http://localhost:8080/api/users/123

# 表单绑定（登录）
curl -X POST http://localhost:8080/api/login \
  -d "username=admin&password=admin123"

# 原始请求体
curl -X POST http://localhost:8080/api/raw \
  -H "Content-Type: text/plain" \
  -d "This is raw data"

# JSON 响应
curl http://localhost:8080/json

# XML 响应
curl http://localhost:8080/xml

# 字符串响应
curl http://localhost:8080/string

# 重定向
curl -L http://localhost:8080/redirect

# 二进制数据响应
curl http://localhost:8080/data
```

## 参数验证标签

| 标签 | 说明 | 示例 |
|------|------|------|
| required | 必填 | `binding:"required"` |
| email | 邮箱格式 | `binding:"email"` |
| min | 最小值 | `binding:"min=10"` |
| max | 最大值 | `binding:"max=100"` |
| gte | 大于等于 | `binding:"gte=0"` |
| lte | 小于等于 | `binding:"lte=120"` |
| len | 长度 | `binding:"len=10"` |
| omitempty | 可选 | `binding:"omitempty,email"` |

## 响应格式

统一响应结构：
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

