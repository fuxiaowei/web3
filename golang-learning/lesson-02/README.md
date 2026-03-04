# 第2期：GORM数据库ORM完全指南

## 📚 课程介绍

本课程是《Golang区块链开发》系列的第二期，将系统讲解 GORM 的核心能力，从基础 CRUD 到复杂关联关系与事务处理，再到完整的业务项目实战，帮助你构建稳健的持久化层。

## 📖 课程内容

## 第一阶段: 基础入门

### 2.1 GORM安装与配置

**状态:** ✅  
**学时:** 60分钟  
**教学内容:**
- 数据库连接配置
- 连接池设置
- 模型定义
- 了解GORM支持的数据库驱动

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6942571de4b0694ca15a78c6?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [GORM基础 - 安装与配置(课件).md](./courseware/basics.md#第一部分安装与配置)

**代码示例:**
- [setup_test.go](./examples/basics/setup_test.go)

---

### 2.2 CRUD操作详解

**状态:** ✅  
**学时:** 90分钟  
**教学内容:**
- 自动迁移与数据初始化
- 基础 CRUD：Create / First & Find / Updates / Delete
- 条件查询与链式调用
- 实战：构建用户管理服务

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_69410f11e4b0694ca159ac0d?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [GORM基础 - CRUD操作(课件).md](./courseware/basics.md#第二部分CRUD操作)

**代码示例:**
- [crud_test.go](./examples/basics/crud_test.go)

---

### 2.3 高级查询

**状态:** ✅  
**学时:** 90分钟  
**教学内容:**
- 链式查询
- 条件构建
- 原生SQL
- 聚合查询与分页

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_693a688ee4b0694c5b60cc88?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [GORM基础 - 高级查询(课件).md](./courseware/basics.md#第三部分高级查询)

**代码示例:**
- [query_builder_test.go](./examples/basics/query_builder_test.go)

---

## 第二阶段: 进阶特性

### 2.4 关联关系

**状态:** ✅  
**学时:** 120分钟  
**教学内容:**
- 一对一关系建模
- 一对多关系建模
- 多对多关系建模
- 预加载Preload
- 懒加载与自定义预加载条件

**课程地址:**
- 关联关系: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_693a68a1e4b0694c5b60cc90?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)
- 预加载: [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6942572de4b0694ca15a78d8?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [GORM进阶 - 关联关系(课件).md](./courseware/advanced.md#第一部分关联关系)

**代码示例:**
- [associations_basics_test.go](./examples/advanced/associations_basics_test.go)
- [associations_many2many_test.go](./examples/advanced/associations_many2many_test.go)
- [associations_preload_test.go](./examples/advanced/associations_preload_test.go)
- [associations_api_test.go](./examples/advanced/associations_api_test.go)

---

### 2.5 事务处理

**状态:** ✅  
**学时:** 90分钟  
**教学内容:**
- 事务开启、提交、回滚
- 嵌套事务
- 乐观锁 / 悲观锁方案
- 实战：账户转账事务完整实现

**课程地址:** [观看课程](https://appibxs98ig9955.h5.xet.citv.cn/p/course/video/v_6943e229e4b0694ca15b8d9d?product_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&course_id=course_2xArq45EVWERVb4oHPYEmCLzHG5&sub_course_id=subcourse_35y60lohGZFZ6cO3dA09xkRlkfj)

**文档资料:**
- 课件: [GORM进阶 - 事务处理(课件).md](./courseware/advanced.md#第二部分事务处理)

**代码示例:**
- [transactions_test.go](./examples/advanced/transactions_test.go)

---

### 2.6 软删除与钩子

**状态:** ⏳  
**学时:** 60分钟  
**教学内容:**
- 软删除机制
- 模型钩子函数
- BeforeCreate / AfterCreate
- BeforeUpdate / AfterUpdate
- BeforeDelete / AfterDelete

**课程地址:** 暂无

**文档资料:**
- 课件: [GORM进阶 - 软删除与钩子(课件).md](./courseware/advanced.md#第三部分软删除与钩子)

**代码示例:**
- [hooks_softdelete_test.go](./examples/advanced/hooks_softdelete_test.go)

---

## 第三阶段: 项目实战

### 2.7 电商项目实战

**状态:** ✅  
**学时:** 120分钟  
**教学内容:**
- 电商领域模型设计（用户、商品、订单、订单项）
- 聚合查询、分页、软删除与审计字段
- 日志与错误处理规范
- 性能优化：索引、批量操作、连接池调优
- 实战：下单流程与库存扣减

**课程地址:** 暂无

**文档资料:**
- 课件: [GORM项目实战 - 电商系统(课件).md](./courseware/project.md)

**代码示例:**
- [ecommerce_test.go](./examples/project/ecommerce_test.go)

## 📁 目录结构

```
lesson-02/
├── README.md                  # 本文档
├── COURSE_OVERVIEW.md         # 课程总览
├── QUICK_START.md             # 快速上手指南
├── DIRECTORY_STRUCTURE.md     # 目录结构说明
├── courseware/                # 课件目录
│   ├── basics.md              # 基础与 CRUD 课件
│   ├── advanced.md            # 关联关系与事务课件
│   └── project.md             # 项目实战课件
└── examples/                  # 代码示例目录
    ├── basics/                # 基础示例
    ├── advanced/              # 进阶示例
    └── project/               # 项目级示例
```

## 🚀 快速开始

### 运行基础示例

```bash
# 进入示例目录
cd examples

# 安装依赖
go mod tidy

# 运行安装与配置示例
go test ./basics -run TestSetupDemo -v

# 运行CRUD示例
go test ./basics -run TestCRUDDemo -v

# 运行高级查询示例
go test ./basics -run TestQueryBuilderDemo -v
```

### 运行进阶特性示例

```bash
# 运行关联关系基础示例
go test ./advanced -run TestAssociationsBasics -v

# 运行多对多关系示例
go test ./advanced -run TestAssociationsMany2Many -v

# 运行预加载示例
go test ./advanced -run TestAssociationsPreload -v

# 运行事务处理示例
go test ./advanced -run TestTransactionsDemo -v

# 运行软删除与钩子示例
go test ./advanced -run TestHooksSoftDelete -v
```

### 运行项目实战示例

```bash
# 运行电商项目示例
go test ./project -run TestEcommerceFlow -v
```

> 所有测试用例都会在临时目录中创建 SQLite 数据库，无需额外安装数据库服务。

## 🎯 学习目标

完成本课程后，你应该能够：

1. ✅ 熟练使用 GORM 完成基础 CRUD
2. ✅ 设计并实现复杂关联关系
3. ✅ 编写可靠的事务逻辑并处理异常
4. ✅ 构建完整的电商持久化层
5. ✅ 优化数据库访问性能并规避常见坑

## 📚 示例代码与课件对应关系

| 示例文件 | 对应课件章节 | 主要知识点 | 运行命令 |
|---------|------------|-----------|---------|
| **基础示例** | | | |
| `setup_test.go` | 第一部分：安装与配置 | 数据库连接、连接池、模型定义、自动迁移 | `go test ./basics -run TestSetupDemo -v` |
| `crud_test.go` | 第二部分：CRUD操作 | Create、Read、Update、Delete | `go test ./basics -run TestCRUDDemo -v` |
| `query_builder_test.go` | 第三部分：高级查询 | 链式查询、条件构建、原生SQL | `go test ./basics -run TestQueryBuilderDemo -v` |
| **进阶示例** | | | |
| `associations_basics_test.go` | 第一部分：关联关系 | 一对一、一对多关系 | `go test ./advanced -run TestAssociationsBasics -v` |
| `associations_many2many_test.go` | 第一部分：关联关系 | 多对多关系 | `go test ./advanced -run TestAssociationsMany2Many -v` |
| `associations_preload_test.go` | 第一部分：关联关系 | 预加载Preload | `go test ./advanced -run TestAssociationsPreload -v` |
| `associations_api_test.go` | 第一部分：关联关系 | 关联关系API操作 | `go test ./advanced -run TestAssociationsAPI -v` |
| `transactions_test.go` | 第二部分：事务处理 | 事务开启、提交、回滚、嵌套事务 | `go test ./advanced -run TestTransactionsDemo -v` |
| `hooks_softdelete_test.go` | 第三部分：软删除与钩子 | 软删除机制、模型钩子函数 | `go test ./advanced -run TestHooksSoftDelete -v` |
| **项目实战** | | | |
| `ecommerce_test.go` | 项目实战 | 电商系统完整实现 | `go test ./project -run TestEcommerceFlow -v` |

## 💡 实践建议

### 课后练习

1. **基础练习：**
   - 实现一个博客系统（用户-文章-标签）
   - 实现一个学生管理系统（学生-课程-成绩）

2. **进阶练习：**
   - 实现账户转账系统（使用事务保证数据一致性）
   - 实现多对多关系的标签系统

3. **综合练习：**
   - 尝试将 SQLite 替换为 MySQL 或 PostgreSQL
   - 结合 Gin 构建一个提供 RESTful API 的数据服务
   - 实现完整的电商系统（用户、商品、订单、支付）

### 扩展阅读

- [GORM官方文档](https://gorm.io/docs/)
- [GORM最佳实践](https://gorm.io/docs/best_practices.html)
- [数据库设计原则](https://en.wikipedia.org/wiki/Database_normalization)

## 📞 常见问题

### Q: AutoMigrate 会删除现有数据吗？
A: 不会。AutoMigrate 只会创建新表、添加新字段和索引，不会删除现有数据或字段。

### Q: 什么时候使用预加载 Preload，什么时候使用懒加载？
A: 如果确定需要关联数据，使用 Preload 可以避免 N+1 查询问题。如果关联数据可能不需要，可以使用懒加载。

### Q: 事务中发生错误会自动回滚吗？
A: 是的。如果在事务中发生 panic 或返回错误，GORM 会自动回滚事务。也可以手动调用 `Rollback()`。

### Q: 软删除和硬删除有什么区别？
A: 软删除只是标记记录为已删除（更新 DeletedAt 字段），数据仍然存在于数据库中。硬删除是真正从数据库中删除记录。

## 🔗 相关资源

- [课件文档](./courseware/)
- [代码示例](./examples/)
- [快速开始指南](./QUICK_START.md)
- [目录结构说明](./DIRECTORY_STRUCTURE.md)
- [课程总览](./COURSE_OVERVIEW.md)

---

**祝学习愉快，掌握 Go + GORM 的数据库开发核心技能！** 🚀

