# Gin 框架示例代码

本目录包含 Gin 框架的各种示例，每个示例都是独立的工程，可以单独运行。

## 📚 示例列表

### 基础示例

| 序号 | 名称 | 说明 | 目录 |
|------|------|------|------|
| 01 | Hello World | 第一个 Gin 应用 | [01-hello-world](./01-hello-world) |
| 02 | 路由系统 | 路由定义、参数、分组 | [02-routing](./02-routing) |
| 03 | 请求处理 | 参数绑定、验证、响应 | [03-request-handling](./03-request-handling) |

### 高级示例

| 序号 | 名称 | 说明 | 目录 |
|------|------|------|------|
| 04 | 中间件 | 自定义中间件、日志、认证 | [04-middleware](./04-middleware) |
| 05 | JWT 认证 | JWT Token 生成和验证 | [05-jwt-auth](./05-jwt-auth) |
| 06 | 文件上传 | 单文件、多文件上传和下载 | [06-file-upload](./06-file-upload) |
| 07 | 配置管理 | Viper 配置管理 | [07-config](./07-config) |

### 综合项目

| 名称 | 说明 | 目录 |
|------|------|------|
| RESTful API 项目 | 完整的用户管理系统 | [project](./project) |

## 🚀 快速开始

每个示例都是独立的，可以单独运行。以 Hello World 为例：

```bash
# 进入示例目录
cd 01-hello-world

# 下载依赖
go mod download

# 运行示例
go run main.go
```

## 📋 学习路线

建议按照以下顺序学习：

1. **01-hello-world** - 了解 Gin 基本概念
2. **02-routing** - 掌握路由系统
3. **03-request-handling** - 学习请求和响应处理
4. **04-middleware** - 理解中间件机制
5. **05-jwt-auth** - 实现 JWT 认证
6. **06-file-upload** - 处理文件上传
7. **07-config** - 配置管理
8. **project** - 综合实践

## 💡 运行方式

每个示例支持两种运行方式：

### 方式一：直接运行

```bash
cd <example-directory>
go run main.go
```

### 方式二：编译后运行

```bash
cd <example-directory>
go build
./<binary-name>
```

## 📖 示例说明

每个示例目录包含：

- `main.go` - 主程序文件
- `go.mod` - Go 模块文件
- `README.md` - 详细说明和测试方法

## 🔧 依赖安装

每个示例的依赖已在 `go.mod` 中定义。首次运行时会自动下载依赖，也可以手动执行：

```bash
cd <example-directory>
go mod download
```

## 📝 测试工具

推荐使用以下工具测试 API：

- **curl** - 命令行工具（示例中已提供 curl 命令）
- **Postman** - 图形化工具
- **httpie** - 更友好的命令行工具

## 🌟 主要技术栈

- **Gin** - Web 框架
- **jwt-go** - JWT 认证
- **Viper** - 配置管理
- **GORM** - ORM（在 project 中使用）

## 🎯 学习目标

通过学习这些示例，你将掌握：

- ✅ Gin 框架的基本用法
- ✅ RESTful API 设计
- ✅ 中间件的编写和使用
- ✅ JWT 认证实现
- ✅ 文件上传处理
- ✅ 配置管理最佳实践
- ✅ 完整项目的组织结构

## 📚 参考资料

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Go 官方文档](https://golang.org/doc/)

## 🤝 贡献

如果你发现问题或有改进建议，欢迎提交 Issue 或 Pull Request。

## 📄 许可证

MIT License

