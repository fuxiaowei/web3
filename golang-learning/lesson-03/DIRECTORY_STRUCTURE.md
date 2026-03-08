# 第3期目录结构说明

```
lesson-03/
├── README.md                  # 课程简介与学习指南
├── COURSE_OVERVIEW.md         # 课程总览（目标、内容、收获）
├── QUICK_START.md             # 快速上手指南
├── DIRECTORY_STRUCTURE.md     # 本文档
├── courseware/                # 课件资料
│   ├── basics.md              # 基础与路由课件大纲
│   ├── advanced.md            # 中间件与高级特性课件大纲
│   └── project.md             # 项目实战课件大纲
└── examples/                  # 代码示例（每个示例都是独立工程）
    ├── README.md              # 示例总览和学习指南
    ├── 01-hello-world/        # 基础：第一个 Gin 应用
    │   ├── main.go
    │   ├── go.mod
    │   └── README.md
    ├── 02-routing/            # 基础：路由系统
    │   ├── main.go
    │   ├── go.mod
    │   └── README.md
    ├── 03-request-handling/   # 基础：请求处理
    │   ├── main.go
    │   ├── go.mod
    │   └── README.md
    ├── 04-middleware/         # 高级：中间件
    │   ├── main.go
    │   ├── go.mod
    │   └── README.md
    ├── 05-jwt-auth/           # 高级：JWT 认证
    │   ├── main.go
    │   ├── go.mod
    │   └── README.md
    ├── 06-file-upload/        # 高级：文件上传下载
    │   ├── main.go
    │   ├── go.mod
    │   └── README.md
    ├── 07-config/             # 高级：配置管理
    │   ├── main.go
    │   ├── go.mod
    │   ├── config.yaml
    │   └── README.md
    └── project/               # 综合项目：用户管理系统
        ├── main.go            # 主程序入口
        ├── handlers/          # 处理器目录
        │   └── user_handler.go
        ├── middleware/        # 中间件目录
        │   ├── auth.go
        │   ├── logger.go
        │   └── cors.go
        ├── models/            # 模型目录
        │   └── user.go
        ├── services/          # 业务逻辑层
        │   └── user_service.go
        ├── config/            # 配置目录
        │   └── config.go
        └── utils/             # 工具目录
            ├── jwt.go
            ├── response.go
            └── errors.go
```

## 🎯 目录设计理念

### 独立示例设计

与之前的共享 `go.mod` 不同，新的目录结构采用**每个示例独立工程**的设计：

- ✅ **每个示例都有自己的 `go.mod`**：避免依赖冲突
- ✅ **每个示例都可以单独运行**：`go run main.go` 即可
- ✅ **每个示例都有详细的 README**：包含说明、测试方法、知识点
- ✅ **便于讲解**：可以针对单个示例进行深入讲解

### 为什么这样设计？

#### 之前的问题
```
examples/
├── go.mod           # 共享的 go.mod
├── basics/
│   ├── 01-hello-world.go    # package main
│   ├── 02-routing.go        # package main（冲突！）
│   └── 03-request-handling.go  # package main（冲突！）
```

**问题**：
- 多个文件都是 `package main`，都有 `main()` 函数
- 无法同时编译，无法单独运行
- 讲解时需要复制代码到临时文件

#### 现在的解决方案
```
examples/
├── 01-hello-world/     # 独立工程
│   ├── main.go
│   └── go.mod
├── 02-routing/         # 独立工程
│   ├── main.go
│   └── go.mod
```

**优势**：
- ✅ 每个示例都是完整的、可运行的工程
- ✅ 可以直接 `cd` 进去运行
- ✅ 便于教学演示和学习

## 📚 目录使用说明

### 课件目录（courseware/）
提供录制课程所需的详细大纲和讲稿提示，可直接转为课件或演示脚本。

### 示例目录（examples/）
按学习顺序编号，包含三类示例：

#### 1. 基础示例（01-03）
- **01-hello-world**：第一个 Gin 应用，了解基本概念
- **02-routing**：路由系统，掌握路由定义和分组
- **03-request-handling**：请求处理，学习参数绑定和验证

#### 2. 高级示例（04-07）
- **04-middleware**：中间件机制，理解中间件的编写和使用
- **05-jwt-auth**：JWT 认证，实现完整的认证流程
- **06-file-upload**：文件处理，掌握文件上传和下载
- **07-config**：配置管理，使用 Viper 管理配置

#### 3. 综合项目（project/）
完整的 RESTful API 项目，展示企业级项目结构和最佳实践。

## 🚀 示例运行方式

### 独立示例运行

每个示例都是独立的工程，可以单独运行：

```bash
# 方式一：直接运行
cd examples/01-hello-world
go run main.go

# 方式二：编译后运行
cd examples/01-hello-world
go build
./hello-world
```

### 项目示例运行

```bash
cd examples/project
go mod download
go run main.go
```

## 📖 学习建议

1. **按顺序学习**：从 01 到 07，循序渐进
2. **动手实践**：每个示例都要运行起来，亲自测试
3. **阅读 README**：每个示例的 README 都有详细说明
4. **修改代码**：尝试修改代码，理解各个参数的作用
5. **综合应用**：最后学习 project，理解如何组织完整项目

## 💡 教学建议

### 讲解顺序
1. 先讲基础示例（01-03），建立基础概念
2. 再讲高级示例（04-07），深入理解机制
3. 最后讲综合项目，展示最佳实践

### 讲解方式
1. **展示代码**：打开示例目录，展示代码结构
2. **运行演示**：`go run main.go` 启动服务
3. **测试验证**：使用 curl 或 Postman 测试
4. **讲解原理**：结合代码讲解实现原理
5. **扩展练习**：引导学员思考如何改进

## 🔧 技术栈

### 核心框架
- **Gin** - Web 框架
- **Viper** - 配置管理
- **JWT** - 认证授权
- **GORM** - ORM（在 project 中）

### 开发工具
- Go 1.22+
- curl / Postman（API 测试）
- Git（版本控制）

---

**这种独立示例的组织方式，更适合教学演示和学员学习！** 🚀
