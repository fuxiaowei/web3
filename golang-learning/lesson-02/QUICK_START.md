# 第2期快速上手指南

## 1. 环境准备

### ✅ 必备工具
- Go 1.22+
- Git
- make（可选，用于执行辅助脚本）

### ✅ 数据库
- 默认使用 SQLite，无需额外安装
- 若需使用 MySQL / PostgreSQL，请准备对应数据库并更新连接配置

## 2. 获取代码

```bash
cd /Users/zhouxing/Documents/rcc/new-golang/lesson-02/examples
go mod tidy
```

`go mod tidy` 会自动拉取 `gorm.io/gorm` 和对应驱动 `gorm.io/driver/sqlite`。

## 3. 运行基础示例

```bash
# 初始化数据库与自动迁移
go test ./basics -run TestSetupDemo -v

# 完整 CRUD 流程
go test ./basics -run TestCRUDDemo -v

# 条件查询、分页、排序
go test ./basics -run TestQueryBuilderDemo -v
```

所有测试都会在临时目录下创建 SQLite 数据库文件，执行完成后会自动清理，无需手动删除。

## 4. 运行进阶示例

```bash
# 关联关系与预加载
go test ./advanced -run TestAssociationsDemo -v

# 事务与 SavePoint
go test ./advanced -run TestTransactionsDemo -v

# 钩子、软删除与审计字段
go test ./advanced -run TestHooksAndSoftDeleteDemo -v
```

## 5. 项目实战示例

```bash
go test ./project -run TestEcommerceFlow -v
```

输出将展示从下单、库存扣减、订单支付状态更新到订单报表查询的完整流程。

## 6. 常见问题

### Q: 运行时提示找不到驱动？
A: 确认已经执行 `go mod tidy`，并在示例中正确导入了驱动库，例如 `gorm.io/driver/sqlite`。当切换到 MySQL / PostgreSQL 时，请在 `go.mod` 中添加对应驱动。

### Q: 数据库文件存放在哪里？
A: 当前测试会使用 `testing.T.TempDir()` 创建临时目录，测试结束后会自动清理。如果希望保留数据，可复制示例代码并使用自定义路径。

## 7. 下一步

- 深入阅读 `courseware/` 下的课件文档
- 尝试独立实现博客系统的 CRUD 与关联
- 将项目实战示例与 Gin API 结合，形成完整的 Web 服务

---

**完成以上步骤即可顺利开始第2期 GORM 实战学习！** 🚀

