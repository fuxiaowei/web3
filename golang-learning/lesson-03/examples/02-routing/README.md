# 02 - 路由系统

## 说明

全面演示 Gin 路由系统的各种用法，包括路径参数、查询参数、路由分组、RESTful API 设计等。

## 知识点

- 路径参数（`:id`、通配符）
- 查询参数（`c.Query()`、`c.DefaultQuery()`）
- 表单参数（`c.PostForm()`、`c.DefaultPostForm()`）
- JSON 参数（`c.ShouldBindJSON()`）
- 路由分组（`r.Group()`）
- 嵌套分组
- RESTful API 设计
- 各种 HTTP 方法（GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS、ANY）

## 运行方式

```bash
go run main.go
```

## 测试

```bash
# 基础路由
curl http://localhost:8080/

# 路径参数 - 单个参数
curl http://localhost:8080/users/123

# 路径参数 - 多个参数
curl http://localhost:8080/users/123/posts/456

# 通配符参数
curl http://localhost:8080/files/docs/readme.md

# 查询参数
curl "http://localhost:8080/search?keyword=golang&page=2&size=20"

# 表单参数（POST）
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=123456&remember=true"

# JSON 参数（POST）
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":25}'

# 路由分组（含嵌套示例）
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v1/users/123
# 嵌套分组 demo（避免与 RESTful 产品路由冲突）
curl http://localhost:8080/api/v1/demo-products

# RESTful API
curl http://localhost:8080/api/v1/products
curl http://localhost:8080/api/v1/products/1
curl -X POST http://localhost:8080/api/v1/products
curl -X PUT http://localhost:8080/api/v1/products/1
curl -X DELETE http://localhost:8080/api/v1/products/1

# 任意方法
curl -X GET http://localhost:8080/test
curl -X POST http://localhost:8080/test
```

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | / | 首页 |
| GET | /users/:id | 获取用户 |
| GET | /users/:id/posts/:postId | 获取用户的文章 |
| GET | /files/*filepath | 文件路径 |
| GET | /search | 搜索 |
| POST | /login | 登录（表单参数） |
| POST | /api/users | 创建用户（JSON 参数） |
| GET | /api/v1/users | 用户列表 |
| POST | /api/v1/users | 创建用户 |
| GET | /api/v1/demo-products | 嵌套分组示例 |
| GET | /api/v1/products | 商品列表 |
| GET | /api/v1/products/:id | 商品详情 |
| POST | /api/v1/products | 创建商品 |
| PUT | /api/v1/products/:id | 更新商品 |
| DELETE | /api/v1/products/:id | 删除商品 |

