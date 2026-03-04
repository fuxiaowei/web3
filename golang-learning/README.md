# Golang入门到精通系列课程

## 课程概述

本课程分为多个阶段，涵盖从Golang基础入门到框架实战，再到区块链开发的完整学习路径。

## 课程阶段

### 第一阶段：Golang基础语法与进阶特性

**课程数量：** 第1期课程  
**学习目标：** 掌握Golang核心语法和基础编程能力。  
**涵盖内容：** 包括Go语言环境搭建、基础数据类型、控制流语句、函数与结构体、接口和多态、并发编程（Goroutine和Channel）、并发安全（Mutex和RWMutex）、标准库使用等基础知识，并完成Web爬虫实战项目。

**课程详情：** [第1期课程](./lesson-01/README.md)

---

### 第二阶段：GORM数据库ORM完全指南

**课程数量：** 第2期课程  
**学习目标：** 深入学习数据库ORM框架和持久化层开发。  
**涵盖内容：** 包括GORM安装配置、模型定义、基础CRUD操作、条件查询、关联关系（一对一、一对多、多对多）、预加载与懒加载、事务处理、钩子函数、软删除、乐观锁等进阶主题，并完成电商订单系统实战项目。

**课程详情：** [第2期课程](./lesson-02/README.md)

---

### 第三阶段：Gin Web框架实战开发

**课程数量：** 第3期课程  
**学习目标：** 掌握专业Web开发框架和实战项目开发。  
**涵盖内容：** 学习Gin Web框架、路由系统、请求处理、中间件机制、JWT认证、文件上传下载、配置管理（Viper）、Protocol Buffers等开发技能，通过用户管理系统、RESTful API设计等完整项目提升实战能力。

**课程详情：** [第3期课程](./lesson-03/README.md)

---

### 第四阶段：区块链基础与 Go-Eth Client 实战

**课程数量：** 第4期课程  
**学习目标：** 掌握以太坊基础架构与 Go-Eth Client 客户端开发能力，具备从 Go 代码中“读写以太坊”的工程实践能力。  
**涵盖内容：** 包括以太坊账户模型、交易与区块结构、Gas 机制、智能合约与网络类型（主网/测试网/本地链）等基础概念；Go-Ethereum (`go-ethereum`) 客户端安装配置与节点连接；区块与交易查询、交易回执解析与账户余额查询；订阅新区块、pending 交易与合约日志事件及断线重连；使用 ABI 与 Go 代码进行智能合约方法调用、发送合约交易与事件解析，并完成迷你区块浏览器 / ERC-20 转账监听实战项目。

**课程详情：** [第4期课程](./lesson-04/README.md)

---

## 快速开始

### 第1期：Golang语言

```bash
cd lesson-01/examples/basics
go run 01-types.go
```

[查看第1期快速开始指南](./lesson-01/QUICK_START.md)

### 第2期：GORM框架

```bash
cd lesson-02/examples
go mod tidy
go test ./basics -run TestSetupDemo -v
```

[查看第2期快速开始指南](./lesson-02/QUICK_START.md)

### 第3期：Gin框架

```bash
cd lesson-03/examples/01-hello-world
go run main.go
```

[查看第3期快速开始指南](./lesson-03/QUICK_START.md)

### 第4期：Go-Eth Client 与区块链基础

```bash
cd lesson-04
# 按 QUICK_START.md 中说明配置 ETH_RPC_URL 等环境变量
cat QUICK_START.md
```

[查看第4期快速开始指南](./lesson-04/QUICK_START.md)

## 学习路径建议

1. **初学者：** 按照阶段顺序学习，每个阶段完成所有示例和实践项目
2. **有经验者：** 可以快速浏览课件，重点关注实战项目和最佳实践
3. **进阶学习：** 完成前三阶段后，可以继续学习区块链相关课程（第4期 Go-Eth Client、后续数据同步与多链支持等）

## 课程特色

- ✅ **系统化学习路径**：从基础到进阶，循序渐进
- ✅ **丰富的代码示例**：每个知识点都配有可运行的代码示例
- ✅ **实战项目驱动**：通过真实项目巩固所学知识
- ✅ **最佳实践指导**：涵盖企业级开发规范和性能优化

## 目录结构

```
new-golang/
├── README.md              # 本文档
├── lesson-01/             # 第1期：Golang语言
│   ├── README.md
│   ├── COURSE_OVERVIEW.md
│   ├── QUICK_START.md
│   ├── courseware/        # 课件
│   └── examples/          # 代码示例
├── lesson-02/             # 第2期：GORM框架
│   ├── README.md
│   ├── COURSE_OVERVIEW.md
│   ├── QUICK_START.md
│   ├── courseware/        # 课件
│   └── examples/          # 代码示例
├── lesson-03/             # 第3期：Gin框架
│   ├── README.md
│   ├── COURSE_OVERVIEW.md
│   ├── QUICK_START.md
│   ├── courseware/        # 课件
│   └── examples/          # 代码示例
└── lesson-04/             # 第4期：区块链基础与 Go-Eth Client
    ├── README.md
    ├── COURSE_OVERVIEW.md
    ├── QUICK_START.md
    ├── courseware/        # 课件
    └── examples/          # 代码示例
```

## 学习成果

完成本系列课程后，你将能够：

1. ✅ 熟练掌握Golang核心语法和特性
2. ✅ 理解并应用并发编程模式
3. ✅ 使用GORM构建健壮的数据库访问层
4. ✅ 使用Gin框架开发高性能Web服务
5. ✅ 设计并实现RESTful API
6. ✅ 完成企业级项目开发

## 相关资源

- [Go官方文档](https://golang.org/doc/)
- [GORM官方文档](https://gorm.io/)
- [Gin官方文档](https://gin-gonic.com/docs/)
- [Effective Go](https://golang.org/doc/effective_go.html)

---

**祝学习愉快，成功掌握Golang开发！** 🚀

