# 目录结构说明

```
lesson-01/
│
├── README.md                          # 课程总览和介绍
├── QUICK_START.md                      # 快速开始指南
├── DIRECTORY_STRUCTURE.md             # 本文件，目录结构说明
│
├── courseware/                        # 📚 课件目录
│   ├── basics.md                      # 基础语法课件
│   └── advanced.md                    # 进阶特性课件
│
└── examples/                          # 💻 代码示例目录
    │
    ├── basics/                        # 基础语法示例
    │   ├── 01-types.go               # 数据类型：基本类型、数组、切片、映射、指针
    │   ├── 02-control-flow.go        # 控制流：if/switch/for/defer/panic/recover
    │   ├── 03-function.go             # 函数：基本函数、闭包、多返回值
    │   ├── 04-struct.go              # 结构体：定义、方法、嵌入
    │   └── calculator.go             # 实践：简单计算器项目
    │
    ├── advanced/                      # 进阶特性示例
    │   ├── 01-interface.go           # 接口和多态：接口定义、空接口、类型断言
    │   ├── 02-goroutine.go          # 并发编程：goroutine、channel、select
    │   ├── 03-concurrency-safe.go    # 并发安全：mutex、rwmutex、context
    │   └── 04-producer-consumer.go   # 模式：生产者消费者、worker pool
    │
    └── practice/                      # 实战项目
        └── web-crawler.go            # 综合项目：并发Web爬虫
```

## 📋 文件说明

### 课件文件

#### `courseware/basics.md`
- **内容：** 基础语法课件
- **包含：** 环境搭建、数据类型、控制流、函数和结构体
- **时长：** 约40分钟
- **适合：** 初学者入门

#### `courseware/advanced.md`
- **内容：** 进阶特性课件
- **包含：** 接口、并发编程、并发安全、标准库
- **时长：** 约80分钟
- **适合：** 有基础的开发者

### 基础语法示例

#### `01-types.go` - 数据类型
演示Go语言的基本类型和复合类型的使用

**运行：** `go run examples/basics/01-types.go`

**包含内容：**
- 基本类型：bool, int, float, string
- 数组和切片的区别
- 映射的使用
- 指针的操作

#### `02-control-flow.go` - 控制流
演示Go语言的控制流语句

**运行：** `go run examples/basics/02-control-flow.go`

**包含内容：**
- if/else条件语句
- switch选择语句
- for循环（三种形式）
- defer延迟执行
- panic和recover异常处理

#### `03-function.go` - 函数
演示函数的使用和闭包

**运行：** `go run examples/basics/03-function.go`

**包含内容：**
- 基本函数定义
- 多返回值
- 可变参数
- 闭包的使用
- 高阶函数

#### `04-struct.go` - 结构体
演示结构体和方法的定义

**运行：** `go run examples/basics/04-struct.go`

**包含内容：**
- 结构体定义和初始化
- 值接收者和指针接收者
- 方法定义
- 结构体嵌入
- 接口实现

#### `calculator.go` - 计算器实践
一个完整的计算器程序示例

**运行：** `go run examples/basics/calculator.go`

**功能：**
- 加法、减法、乘法、除法
- 多个数字求和
- 求平均值
- 历史记录

### 进阶特性示例

#### `01-interface.go` - 接口和多态
演示接口的使用和多态编程

**运行：** `go run examples/advanced/01-interface.go`

**包含内容：**
- 接口定义和实现
- 多态性的应用
- 空接口和类型断言
- 接口组合

#### `02-goroutine.go` - 并发编程
演示goroutine和channel的使用

**运行：** `go run examples/advanced/02-goroutine.go`

**包含内容：**
- goroutine的创建
- WaitGroup同步
- 无缓冲和缓冲channel
- select语句的使用
- 超时和非阻塞操作

#### `03-concurrency-safe.go` - 并发安全
演示并发安全的技术

**运行：** `go run examples/advanced/03-concurrency-safe.go`

**包含内容：**
- Mutex互斥锁
- RWMutex读写锁
- Context上下文控制
- 并发安全的日志系统

#### `04-producer-consumer.go` - 生产者消费者
演示并发编程的经典模式

**运行：** `go run examples/advanced/04-producer-consumer.go`

**包含内容：**
- 生产者消费者模式
- Worker Pool模式
- 任务分发和结果收集

### 实战项目

#### `web-crawler.go` - Web爬虫
一个使用并发编程的Web爬虫项目

**运行：** `go run examples/practice/web-crawler.go`

**功能：**
- 并发抓取多个URL
- 超时控制
- 错误处理
- 结果收集

**技术点：**
- goroutine并发
- channel通信
- Worker Pool模式
- HTTP客户端

## 🎯 学习顺序建议

### 初学者路径
1. 阅读 `courseware/basics.md`
2. 运行 `examples/basics` 下的所有示例
3. 阅读 `courseware/advanced.md`
4. 运行 `examples/advanced` 下的所有示例
5. 完成 `examples/practice/web-crawler.go`

### 有基础的学习路径
1. 快速浏览课件
2. 运行感兴趣的示例
3. 阅读源代码理解原理
4. 尝试修改和扩展代码

### 实战导向路径
1. 直接运行实战项目
2. 遇到问题查看相关示例
3. 阅读课件理解原理
4. 优化和改进代码

## 📝 笔记建议

建议在以下目录创建笔记文件：

```
lesson-01/
├── notes/
│   ├── basics-notes.md         # 基础语法笔记
│   ├── advanced-notes.md       # 进阶特性笔记
│   └── practice-notes.md          # 实战项目笔记
```

## 🔍 如何查找示例

### 按主题查找

**数据类型：**
- `examples/basics/01-types.go`

**控制流：**
- `examples/basics/02-control-flow.go`

**函数：**
- `examples/basics/03-function.go`

**结构体：**
- `examples/basics/04-struct.go`

**接口：**
- `examples/advanced/01-interface.go`

**并发编程：**
- `examples/advanced/02-goroutine.go`

**并发安全：**
- `examples/advanced/03-concurrency-safe.go`

**设计模式：**
- `examples/advanced/04-producer-consumer.go`

**实战项目：**
- `examples/practice/web-crawler.go`

## 📖 使用示例代码的最佳实践

1. **先运行：** 查看代码的实际效果
2. **再阅读：** 理解代码的逻辑
3. **后修改：** 尝试修改参数和逻辑
4. **多思考：** 思考为什么这样设计
5. **做笔记：** 记录重要的知识点
