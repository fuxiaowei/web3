# 第2期目录结构说明

```
lesson-02/
├── README.md                  # 课程简介与学习指南
├── COURSE_OVERVIEW.md         # 课程总览（目标、内容、收获）
├── QUICK_START.md             # 快速上手指南
├── DIRECTORY_STRUCTURE.md     # 本文档
├── courseware/                # 课件资料
│   ├── basics.md              # 基础与 CRUD 课件大纲
│   ├── advanced.md            # 关联关系与事务课件大纲
│   └── project.md             # 项目实战课件大纲
├── examples/                  # 代码示例
│   ├── go.mod                 # 示例模块依赖声明
│   ├── go.sum                 # （运行 go mod tidy 后生成）
│   ├── basics/                # 基础示例
│   │   ├── helpers_test.go    # SQLite 测试数据库辅助函数
│   │   ├── setup_test.go      # 数据库连接与自动迁移
│   │   ├── crud_test.go       # 基础 CRUD 示例
│   │   └── query_builder_test.go # 条件查询与分页示例
│   ├── advanced/              # 进阶示例
│   │   ├── helpers_test.go    # 通用测试辅助工具
│   │   ├── associations_test.go   # 多种关联关系示例
│   │   ├── transactions_test.go   # 事务与 SavePoint 示例
│   │   └── hooks_softdelete_test.go # 钩子与软删除示例
│   └── project/               # 项目级示例
│       ├── helpers_test.go    # 测试专用工具
│       └── ecommerce_test.go  # 电商订单流程示例
```

## 目录使用说明

- `courseware/` 目录提供录制课程所需的详细大纲和讲稿提示，可直接转为课件或演示脚本。
- `examples/` 目录中的示例现在通过 `go test` 触发，每个测试都会在临时目录中创建 SQLite 数据库，运行完毕后自动清理。
- 若需要演示其他数据库，只需在测试代码中调整驱动与连接配置，结构无需改动。

---

**按照本结构组织课程材料，能够保证教学内容条理清晰、示例易于查找。** 💡

