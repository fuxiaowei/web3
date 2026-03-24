# 快速开始指南

## 📌 前置准备

### 1. 安装Go

访问 [Go官网](https://golang.org/dl/) 下载并安装Go（版本1.19+）

```bash
# 验证安装
go version
```

### 2. 配置环境变量（可选）

```bash
# 设置代理（加速下载依赖）
go env -w GOPROXY=https://goproxy.cn,direct

# 查看配置
go env
```

### 3. 安装IDE（推荐VS Code）

下载并安装 [VS Code](https://code.visualstudio.com/)
安装Go扩展：`golang.go`

## 🚀 开始学习

### 方式一：按顺序学习

1. **阅读课件：** 从 `courseware/basics.md` 开始
2. **运行示例：** 跟着课件运行对应的代码示例
3. **实践练习：** 完成课后作业

### 方式二：直接运行示例

如果你想直接看代码效果：

```bash
# 基础语法示例
cd lesson-01/examples/basics
go run 01-types.go
go run 02-control-flow.go
go run 03-function.go
go run 04-struct.go
go run calculator.go

# 进阶特性示例
cd lesson-01/examples/advanced
go run 01-interface.go
go run 02-goroutine.go
go run 03-concurrency-safe.go
go run 04-producer-consumer.go

# 实战项目
cd lesson-01/examples/practice
go run web-crawler.go
```

## 📖 学习路径

### 第1天：基础语法
- [ ] 阅读课件 `courseware/basics.md`
- [ ] 运行 `examples/basics` 下的所有示例
- [ ] 完成计算器练习
- [ ] 完成课后作业

### 第2天：进阶特性（上）
- [ ] 阅读课件 `courseware/advanced.md` 的接口部分
- [ ] 运行 `examples/advanced/01-interface.go`
- [ ] 理解接口和多态
- [ ] 完成接口作业

### 第3天：进阶特性（下）
- [ ] 学习Goroutine和Channel
- [ ] 运行 `examples/advanced/02-goroutine.go`
- [ ] 运行 `examples/advanced/03-concurrency-safe.go`
- [ ] 运行 `examples/advanced/04-producer-consumer.go`
- [ ] 理解并发编程原理

### 第4天：实战项目
- [ ] 运行 `examples/practice/web-crawler.go`
- [ ] 理解并发编程在项目中的应用
- [ ] 尝试修改代码增强功能

## 💡 学习建议

1. **边学边练：** 每学一个知识点就运行对应的代码
2. **修改代码：** 尝试修改示例代码，看看会发生什么
3. **解决问题：** 遇到错误不要慌，先看错误信息，再看文档
4. **记笔记：** 记录重要概念和自己的理解
5. **多实践：** 完成所有作业和练习

## ❓ 遇到问题？

### 常见错误解决

#### 1. "command not found: go"
```bash
# 检查Go是否在PATH中
echo $PATH | grep go
# 如果没有，添加到PATH
export PATH=$PATH:/usr/local/go/bin
```

#### 2. "cannot find package"
```bash
# 设置代理
go env -w GOPROXY=https://goproxy.cn,direct
# 下载依赖
go mod tidy
```

#### 3. "goroutine没有被执行"
检查代码中是否使用了 `time.Sleep()` 或 `sync.WaitGroup` 等待goroutine完成。

#### 4. 编译错误
仔细阅读错误信息，通常错误会指向具体的行号和原因。

## 🎯 检查清单

完成以下清单表示你已经掌握本课程：

- [ ] 能够编写基本的Go程序
- [ ] 理解数据类型的使用
- [ ] 熟练使用控制流
- [ ] 能够编写函数和方法
- [ ] 理解接口的概念和使用
- [ ] 能够创建和使用goroutine
- [ ] 理解channel的通信机制
- [ ] 能够编写并发安全的代码
- [ ] 完成所有课后作业
- [ ] 能够独立编写简单的Go程序

## 📚 下一步

完成本课程后，你可以继续学习：

- **第2期：** GORM数据库ORM完全指南
- **第3期：** Gin Web框架实战开发
- **第4期：** Go-Eth Client包从入门到精通

## 📞 获取帮助

- 查看官方文档：[golang.org](https://golang.org/doc/)
- 搜索问题：[Stack Overflow](https://stackoverflow.com/questions/tagged/go)
- 提问社区：[Go Forum](https://forum.golangbridge.org/)

---

**开始你的Go语言之旅吧！** 🚀
