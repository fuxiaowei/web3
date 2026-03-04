# 第2期教案：GORM 基础与 CRUD

## 1. 学习目标

通过本课程的学习，你将能够：

- 理解 ORM 与 GORM 的基本概念与优势
- 熟悉 GORM 的安装、初始化、模型定义流程
- 掌握自动迁移和基础 CRUD 的完整用法
- 能够编写规范的条件查询与链式调用

## 2. 核心知识点

### 2.1 ORM 与 GORM 概述

#### 什么是 ORM？

ORM（Object-Relational Mapping，对象关系映射）是一种编程技术，用于在面向对象编程语言和关系型数据库之间建立映射关系。

#### 为什么需要 ORM？

相比直接编写 SQL，ORM 提供了以下优势：

1. **提高生产效率**：减少重复的 SQL 编写，专注于业务逻辑
2. **模型同步**：代码中的模型定义与数据库结构保持一致
3. **组合能力**：通过链式调用灵活组合查询条件
4. **跨数据库支持**：同一套代码可以适配不同的数据库
5. **类型安全**：编译时检查，减少运行时错误

#### GORM 简介

GORM 是 Go 语言中最流行的 ORM 框架，具有以下特点：

- Go 社区使用最广泛的 ORM
- 支持多种数据库驱动（SQLite、MySQL、PostgreSQL 等）
- 提供自动迁移、钩子、事务、关联等强大功能
- 活跃的社区和详尽的文档

**官方文档**：https://gorm.io/

### 2.2 安装与初始化

#### 项目结构

```
lesson-02/examples/
├── basics/          # 基础示例
│   ├── setup_test.go
│   ├── crud_test.go
│   └── query_builder_test.go
├── db/              # 数据库文件存储目录
├── go.mod           # 依赖管理
└── .env             # 环境变量配置
```

#### 安装依赖

```bash
cd lesson-02/examples
go mod tidy
```

#### 数据库驱动

GORM 支持多种数据库，本课程默认使用 SQLite（无需额外安装），生产环境可替换为 MySQL 或 PostgreSQL。

**支持的数据库驱动**：
- SQLite: `gorm.io/driver/sqlite`
- MySQL: `gorm.io/driver/mysql`
- PostgreSQL: `gorm.io/driver/postgres`

#### 数据库连接配置

参考示例：`examples/testutil/helpers.go`

**核心配置项**：

1. **Logger**：控制 SQL 日志输出
   - `logger.Silent`：无日志
   - `logger.Error`：仅错误
   - `logger.Warn`：错误和警告
   - `logger.Info`：所有 SQL 查询（开发推荐）

2. **NamingStrategy**：自定义命名策略
   - `TablePrefix`：表名前缀
   - `SingularTable`：使用单数表名
   - `NoLowerCase`：禁用自动小写

3. **连接池配置**（在 `*sql.DB` 上配置）：
   - `SetMaxIdleConns`：最大空闲连接数
   - `SetMaxOpenConns`：最大打开连接数
   - `SetConnMaxLifetime`：连接最大生存时间

#### 环境变量配置

通过 `.env` 文件配置数据库类型和连接信息：

```env
# 数据库类型：sqlite, mysql, postgres
TEST_DB_TYPE=sqlite

# MySQL 连接字符串
TEST_MYSQL_DSN=root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local

# PostgreSQL 连接字符串
TEST_POSTGRES_DSN=host=localhost user=postgres password=password dbname=testdb port=5432 sslmode=disable
```

#### 模型定义

参考示例：`examples/basics/setup_test.go`

**常用 GORM 标签**：

- `primaryKey`：主键
- `size:64`：字段大小
- `not null`：非空约束
- `uniqueIndex`：唯一索引
- `index`：普通索引
- `default:value`：默认值
- `autoCreateTime`：自动创建时间
- `autoUpdateTime`：自动更新时间

**示例**：

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:64;not null"`
    Email     string    `gorm:"size:128;uniqueIndex;not null"`
    Age       uint8     `gorm:"not null"`
    Status    string    `gorm:"size:16;default:active;index"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

#### 自动迁移（AutoMigrate）

**作用**：
- 自动创建表（如果不存在）
- 添加新列（如果结构体有新字段）
- 创建索引（根据标签）

**注意事项**：
- ⚠️ **不会删除已存在的列**
- ⚠️ **不会修改现有数据**
- ⚠️ **不会删除索引**

**使用示例**：

```go
if err := db.AutoMigrate(&User{}); err != nil {
    log.Fatalf("auto migrate: %v", err)
}
```

### 2.3 基础 CRUD 操作

参考示例：`examples/basics/crud_test.go`

#### Create（创建）

**单条插入**：

```go
user := User{Name: "Alice", Email: "alice@example.com", Age: 28}
if err := db.Create(&user).Error; err != nil {
    // 处理错误
}
// Create 后，user.ID 会自动填充
```

**批量插入**：

```go
users := []User{
    {Name: "Alice", Email: "alice@example.com"},
    {Name: "Bob", Email: "bob@example.com"},
}
if err := db.Create(&users).Error; err != nil {
    // 处理错误
}
```

**控制字段插入**：

- `Select`：只插入指定字段
- `Omit`：排除指定字段

#### Read（查询）

**First**：获取第一条记录（找不到返回错误）

```go
var user User
// 使用条件
if err := db.Where("email = ?", "alice@example.com").First(&user).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // 记录不存在
    }
}
// 使用主键
db.First(&user, 1)
```

**Take**：获取一条记录（不要求条件，找不到不报错）

```go
var user User
db.Take(&user)  // 获取任意一条记录
```

**Find**：获取所有匹配记录（找不到返回空切片）

```go
var users []User
db.Where("status = ?", "active").Find(&users)
```

**Scan**：扫描到自定义结构体或 map

```go
type UserSummary struct {
    Name   string
    Email  string
}
var summaries []UserSummary
db.Model(&User{}).Select("name", "email").Scan(&summaries)
```

#### Update（更新）

**Save**：更新所有字段（包括零值）

```go
user.Name = "New Name"
db.Save(&user)
```

**Updates**：更新指定字段（忽略零值）

```go
// 使用结构体
db.Model(&user).Updates(User{Age: 31, Status: "vip"})

// 使用 map（推荐，避免零值问题）
db.Model(&user).Updates(map[string]any{"age": 31, "status": "vip"})
```

**Select + Updates**：只更新指定字段

```go
db.Model(&user).Select("Age", "Status").Updates(User{Age: 31, Status: "vip"})
```

**批量更新**：

```go
db.Model(&User{}).Where("status = ?", "inactive").Updates(map[string]any{"status": "pending"})
```

#### Delete（删除）

**删除单条记录**：

```go
// 方式1：使用实例
db.Delete(&user)

// 方式2：使用主键
db.Delete(&User{}, user.ID)
```

**批量删除**：

```go
db.Where("status = ?", "inactive").Delete(&User{})
```

**验证删除**：

```go
err := db.First(&User{}, user.ID).Error
if !errors.Is(err, gorm.ErrRecordNotFound) {
    // 记录仍然存在
}
```

### 2.4 条件查询，链式调用，原生sql的使用

参考示例：`examples/basics/query_builder_test.go`

#### Where 条件查询

**基本条件**：

```go
db.Where("status = ?", "active").Find(&users)
```

**多条件**：

```go
db.Where("status = ? AND age > ?", "active", 25).Find(&users)
```

**LIKE 查询**：

```go
db.Where("email LIKE ?", "a%").Find(&users)  // 以 'a' 开头
```

**IN 查询**：

```go
db.Where("status IN ?", []string{"active", "pending"}).Find(&users)
```

**BETWEEN 查询**：

```go
db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)
```

#### Select 指定字段

```go
// 只查询指定字段（提高性能）
db.Select("id", "name", "email").Find(&users)
```

#### Order 排序

```go
db.Order("created_at desc").Find(&users)  // 降序
db.Order("age asc").Find(&users)          // 升序
```

#### Limit 和 Offset（分页）

```go
// Limit: 限制返回记录数
db.Limit(10).Find(&users)

// Offset: 跳过记录数
db.Offset(10).Limit(10).Find(&users)  // 第2页，每页10条
```

#### Scopes（作用域）

Scopes 允许将通用查询条件提取为可复用函数。

**定义 Scope**：

```go
func activeUsers() func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", "active")
    }
}
```

**使用 Scope**：

```go
db.Scopes(activeUsers()).Find(&users)
```

**分页 Scope 示例**：

```go
func paginate(page, size int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if page <= 0 {
            page = 1
        }
        if size <= 0 {
            size = 10
        }
        offset := (page - 1) * size
        return db.Offset(offset).Limit(size)
    }
}

// 使用
db.Scopes(paginate(1, 10)).Find(&users)
```

**组合多个 Scope**：

```go
db.Scopes(activeUsers(), paginate(1, 10)).Find(&users)
```

#### 聚合查询

**Count**：统计数量

```go
var count int64
db.Model(&User{}).Where("status = ?", "active").Count(&count)
```

**Group By**：分组统计

```go
type StatusCount struct {
    Status string
    Total  int64
}
var counts []StatusCount
db.Model(&User{}).
    Select("status, COUNT(*) as total").
    Group("status").
    Scan(&counts)
```

#### 链式调用的特点

1. **顺序无关**：查询条件的构建顺序不影响最终 SQL
2. **可组合**：每个方法返回 `*gorm.DB`，可以继续链式调用
3. **延迟执行**：只有在调用 `Find`、`First`、`Scan` 等方法时才执行 SQL

**示例**：

```go
// 以下两种写法等价
db.Where("status = ?", "active").Order("age desc").Limit(10).Find(&users)
db.Limit(10).Order("age desc").Where("status = ?", "active").Find(&users)
```

#### 原生 SQL 使用方式与示例

当链式查询无法满足复杂需求时，可以直接执行原生 SQL。常用方式如下：

- `db.Raw(sql, args...).Scan(dest)`：查询语句，将结果映射到结构体或切片
- `db.Exec(sql, args...)`：执行更新、删除、插入等非查询语句
- `db.Raw(sql).Rows()` / `db.Raw(sql).Row()`：需要手动遍历 `*sql.Rows` 或获取单行记录时使用

**查询示例**：

```go
type StatusSummary struct {
    Status string
    Total  int64
    AvgAge float64
}

var stats []StatusSummary
start := time.Now().AddDate(0, -1, 0)
end := time.Now()

err := db.Raw(`
    SELECT status, COUNT(*) AS total, AVG(age) AS avg_age
    FROM users
    WHERE created_at BETWEEN ? AND ?
    GROUP BY status
`, start, end).Scan(&stats).Error

if err != nil {
    log.Fatalf("query failed: %v", err)
}
```

**执行语句示例**：

```go
threshold := time.Now().AddDate(0, 0, -30)
result := db.Exec(
    "UPDATE users SET status = ? WHERE last_login_at < ?",
    "inactive",
    threshold,
)

if result.Error != nil {
    log.Fatalf("exec failed: %v", result.Error)
}

fmt.Printf("affected rows: %d\n", result.RowsAffected)
```

## 3. 关键概念总结

### 3.1 ORM 三大优势

1. **生产效率**：减少重复代码，提高开发速度
2. **模型同步**：代码模型与数据库结构保持一致
3. **组合能力**：通过链式调用灵活组合查询

### 3.2 重要区别

**Create / Save / Updates 的区别**：

- `Create`：插入新记录
- `Save`：保存记录（插入或更新所有字段）
- `Updates`：更新指定字段（忽略零值）

**Find vs Scan 的区别**：

- `Find`：查询到相同模型的结构体
- `Scan`：查询到自定义结构体、map 或原始值（用于聚合查询）

### 3.3 最佳实践

1. **错误处理**：始终检查错误，特别是使用 `First` 时要检查 `gorm.ErrRecordNotFound`
2. **连接池配置**：根据应用负载合理配置连接池参数
3. **使用 Select**：只查询需要的字段，提高性能
4. **使用 Scopes**：提取通用查询逻辑，提高代码复用性
5. **环境变量**：使用 `.env` 文件管理配置，便于不同环境切换

## 4. 实践练习

### 练习 1：扩展用户模型

根据示例扩展用户模型，增加以下字段：
- `Phone`：电话号码（字符串，唯一索引）
- `LastLoginAt`：最后登录时间（时间类型）

### 练习 2：实现用户操作函数

完成以下操作函数：

1. **新增用户**：创建用户并默认开启激活状态
   ```go
   func CreateUser(db *gorm.DB, name, email string) (*User, error) {
       // 你的实现
   }
   ```

2. **模糊查询**：根据邮箱模糊查询用户列表（支持分页）
   ```go
   func SearchUsersByEmail(db *gorm.DB, emailPattern string, page, size int) ([]User, error) {
       // 你的实现
   }
   ```

3. **批量更新状态**：批量更新用户状态
   ```go
   func UpdateUserStatus(db *gorm.DB, ids []uint, status string) error {
       // 你的实现
   }
   ```

4. **删除过期用户**：删除超过 30 天未登录的用户
   ```go
   func DeleteInactiveUsers(db *gorm.DB) error {
       // 你的实现（注意：软删除将在进阶模块讲解）
   }
   ```

### 练习 3：使用 Scopes

创建一个 `youngUsers` scope，筛选年龄在 18-30 岁之间的用户，并实现分页查询。

## 5. 运行示例代码

### 运行基础示例

```bash
cd lesson-02/examples

# 初始化数据库与自动迁移
go test ./basics -run TestSetupDemo -v

# 完整 CRUD 流程
go test ./basics -run TestCRUDDemo -v

# 条件查询、分页、排序
go test ./basics -run TestQueryBuilderDemo -v
```

### 查看 SQL 日志

在 `examples/testutil/helpers.go` 中的 `NewTestDB`/`newSQLiteDB` 设置 Logger 为 `logger.Info`，即可看到所有 SQL 查询：

```go
Logger: logger.Default.LogMode(logger.Info),
```

### 查看数据库文件

SQLite 数据库文件存储在 `examples/db/` 目录下，可以使用 SQLite 查看工具打开查看。

## 6. 常见问题

### Q: 运行时提示找不到驱动？

A: 确认已经执行 `go mod tidy`，并在代码中正确导入了驱动库。切换到 MySQL/PostgreSQL 时，需要在 `go.mod` 中添加对应驱动。

### Q: AutoMigrate 会删除字段吗？

A: 不会。AutoMigrate 只会添加新字段，不会删除已存在的字段。如需删除字段，需要手动执行 SQL 或使用迁移工具。

### Q: 如何查看生成的 SQL？

A: 将 Logger 设置为 `logger.Info`，运行测试时会输出所有 SQL 语句。

### Q: Find 和 Scan 什么时候用哪个？

A: 
- 使用 `Find` 查询完整模型到结构体
- 使用 `Scan` 查询部分字段到自定义结构体或进行聚合查询

## 7. 下一步学习

- 深入阅读示例代码中的注释
- 完成实践练习
- 学习进阶内容：关联关系、事务、钩子、软删除等
- 将 GORM 与 Web 框架（如 Gin）结合，构建完整的应用

---

**完成以上内容即可掌握 GORM 的基础 CRUD 操作！** 🚀
