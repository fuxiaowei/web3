# 01 - Hello World

## 说明

第一个 Gin 应用示例，演示如何创建最简单的 Web 服务器。

## 知识点

- 创建 Gin 引擎实例（`gin.Default()` vs `gin.New()`）
- 定义基础路由（`r.GET()`）
- 返回 JSON 响应（`c.JSON()`）
- 返回字符串响应（`c.String()`）
- 启动服务器（`r.Run()`）

## 运行方式

```bash
# 方式一：直接运行
go run main.go

# 方式二：编译后运行
go build
./hello-world
```

## 测试

```bash
# 测试 JSON 响应
curl http://localhost:8080/hello

# 测试字符串响应
curl http://localhost:8080/ping
```

## 预期输出

访问 `/hello`：
```json
{
  "message": "Hello World!"
}
```

访问 `/ping`：
```
pong
```

