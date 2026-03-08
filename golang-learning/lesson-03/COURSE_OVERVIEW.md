# 第3期课程总览

## 🎯 课程目标

通过本课程的学习，全面掌握 Gin Web 框架的核心能力，能够构建高性能的 RESTful API，实现认证授权、文件处理、配置管理等企业级功能，为后续的区块链 Web 服务开发打下坚实基础。

## 📊 课程概览

| 部分 | 内容 | 时长 | 难度 |
|------|------|------|------|
| 基础与路由 | Gin 安装、路由系统、请求处理、参数绑定 | 35 分钟 | ⭐ |
| 中间件与高级特性 | 中间件机制、JWT 认证、文件上传、配置管理、Protobuf | 50 分钟 | ⭐⭐ |
| 项目实战 | RESTful API 设计、用户系统、测试部署 | 40 分钟 | ⭐⭐⭐ |

## 📚 课件资源

### 基础与路由课件 (`courseware/basics.md`)

**第一部分：Gin 安装与项目初始化 (10分钟)**
- Gin 安装与配置
- 项目结构设计
- 第一个 Hello World

**第二部分：路由系统 (15分钟)**
- 路由定义
- 路由分组
- 路径参数
- 查询参数
- 实践：商品 API 设计

**第三部分：请求处理 (10分钟)**
- Handler 函数
- 参数绑定
- 请求体解析

### 中间件与高级特性课件 (`courseware/advanced.md`)

**第一部分：中间件机制 (20分钟)**
- 中间件概念
- 自定义中间件
- 中间件链
- 常用中间件：日志、CORS、限流
- 实践：JWT 认证中间件

**第二部分：文件上传下载 (10分钟)**
- 文件上传处理
- 文件下载
- 静态文件服务

**第三部分：配置管理 (10分钟)**
- Viper 配置管理
- 环境变量配置
- 配置热加载

**第四部分：Protocol Buffers (15分钟)**
- Protobuf 基础概念
- 定义 .proto 文件
- 在 Gin 中使用 Protobuf
- Protobuf vs JSON 性能对比

### 项目实战课件 (`courseware/project.md`)

**第一部分：RESTful API 设计规范 (10分钟)**
- API 设计原则
- 错误处理规范
- 统一响应格式

**第二部分：用户系统 API 实现 (20分钟)**
- 用户注册登录
- Token 生成和验证
- 用户信息管理
- 权限控制

**第三部分：项目部署与测试 (10分钟)**
- 单元测试
- API 测试
- 部署配置

## 💻 代码示例

### 基础示例 (`examples/basics`)

1. **01-hello-world.go** - Gin 入门示例
   - 最简单的 Gin 应用
   - 路由定义和响应

2. **02-routing.go** - 路由系统
   - 路由分组
   - 路径参数和查询参数
   - RESTful 路由设计

3. **03-request-handling.go** - 请求处理
   - 参数绑定
   - JSON/XML/表单数据解析
   - 文件上传基础

### 进阶示例 (`examples/advanced`)

1. **01-middleware.go** - 中间件机制
   - 自定义中间件
   - 中间件链
   - 日志中间件
   - CORS 中间件

2. **02-jwt-auth.go** - JWT 认证
   - Token 生成
   - Token 验证
   - 认证中间件
   - 权限控制

3. **03-file-upload.go** - 文件处理
   - 单文件上传
   - 多文件上传
   - 文件下载
   - 静态文件服务

4. **04-config.go** - 配置管理
   - Viper 基础使用
   - 环境变量配置
   - 配置文件热加载

5. **05-protobuf.go** - Protocol Buffers
   - Protobuf 序列化/反序列化
   - 在 Gin 中使用 Protobuf
   - Protobuf 与 JSON 对比示例

### 项目示例 (`examples/project`)

1. **main.go** - 主程序入口
   - 服务启动
   - 路由注册
   - 中间件配置

2. **handlers/** - 处理器目录
   - user_handler.go：用户相关 API
   - auth_handler.go：认证相关 API

3. **middleware/** - 中间件目录
   - auth.go：JWT 认证中间件
   - logger.go：日志中间件
   - cors.go：跨域中间件

4. **models/** - 模型目录
   - user.go：用户模型

5. **config/** - 配置目录
   - config.go：配置结构
   - viper.go：Viper 初始化

## 🎓 学习路径

### 初学者

1. ✅ 阅读 `courseware/basics.md`
2. ✅ 运行所有基础示例
3. ✅ 完成商品 API 练习
4. ✅ 阅读 `courseware/advanced.md`
5. ✅ 运行所有进阶示例
6. ✅ 完成 JWT 认证练习
7. ✅ 运行并理解项目示例
8. ✅ 完成用户管理系统 API

### 有经验者

1. ⚡ 快速浏览课件重点
2. ⚡ 聚焦中间件和认证机制
3. ⚡ 结合项目示例思考架构设计
4. ⚡ 尝试集成 GORM 和数据库
5. ⚡ 实现微服务架构

## ✅ 学习成果检验

完成本课程后，你应该可以：

- [x] 使用 Gin 构建 Web 服务
- [x] 设计 RESTful API 路由
- [x] 实现参数绑定和请求解析
- [x] 编写自定义中间件
- [x] 实现 JWT 认证和权限控制
- [x] 处理文件上传下载
- [x] 使用 Viper 管理配置
- [x] 理解和使用 Protocol Buffers
- [x] 编写单元测试和 API 测试
- [x] 部署 Web 服务到生产环境

## 📈 进度跟踪

### 基础与路由
- [ ] 完成 Gin 安装和 Hello World
- [ ] 完成路由系统学习
- [ ] 完成请求处理学习
- [ ] 完成商品 API 练习

### 中间件与高级特性
- [ ] 完成中间件机制学习
- [ ] 完成 JWT 认证实现
- [ ] 完成文件上传下载
- [ ] 完成配置管理学习
- [ ] 完成 Protobuf 学习

### 项目实战
- [ ] 运行项目示例
- [ ] 完成用户系统 API
- [ ] 完成单元测试
- [ ] 完成部署配置

## 🎯 下一步计划

完成本课程后，建议继续学习：

1. **第4期：Go-Eth Client 包从入门到精通**
   - 区块链数据交互
   - 事件监听与交易处理
   - 结合 Gin 构建区块链 API 服务

2. **第5期：区块链数据同步核心技术**
   - 数据同步架构
   - Reorg 处理机制
   - 结合 Gin 提供数据查询 API

## 📚 推荐资源

- **Gin 官方文档：** https://gin-gonic.com/docs/
- **Gin GitHub：** https://github.com/gin-gonic/gin
- **JWT 官方文档：** https://jwt.io/
- **Viper 文档：** https://github.com/spf13/viper
- **RESTful API 设计：** https://restfulapi.net/

## 💬 获取帮助

- 📧 课程技术支持邮箱
- 💬 课程讨论群（含 Web 开发心得）
- 🌐 课程官网常见问题

---

**期待你在 Web 开发上更进一步！** 🚀

