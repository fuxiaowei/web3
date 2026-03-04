# 第2期教案：GORM 关联关系与事务

## 1. 学习目标

通过本课程的学习，你将能够：

- 理解 GORM 支持的三类关联关系及其实现方式
- 掌握预加载、懒加载、自定义 Preload 的用法
- 熟练运用事务控制（含 SavePoint、嵌套事务）
- 了解钩子生命周期、软删除、乐观锁的使用场景

## 2. 核心知识点

### 2.1 关联关系概述

#### 示例代码文件结构

关联关系的示例代码已拆分为多个文件，便于分模块学习和讲解：

- **`associations_models.go`**：包含所有模型定义（user, profile, order, orderItem, product, role）
- **`associations_basics_test.go`**：基础关联关系示例（Has One, Has Many, Belongs To）
- **`associations_many2many_test.go`**：Many-to-Many 关系示例
- **`associations_api_test.go`**：Association API 操作示例（Append, Replace, Delete, Clear, Count, Find）
- **`associations_preload_test.go`**：Preload 操作示例（条件预加载、嵌套预加载、clause.Associations）

#### 什么是关联关系？

关联关系用于描述数据库表之间的逻辑关系，GORM 支持三种主要的关联类型：

1. **Has One / Belongs To**：一对一关系
2. **Has Many**：一对多关系
3. **Many to Many**：多对多关系（通过中间表）

#### 关联关系的方向

- **Has One**：一个模型拥有另一个模型（如：User Has One Profile）
- **Belongs To**：一个模型属于另一个模型（如：Profile Belongs To User）
- **Has Many**：一个模型拥有多个另一个模型（如：User Has Many Orders）
- **Many to Many**：两个模型互相拥有多个（如：User Many to Many Roles）

#### 外键（Foreign Key）

外键是建立关联关系的关键，GORM 默认使用 `ModelName + ID` 的命名规则：

- `User` 模型 → 外键通常是 `UserID`
- `Order` 模型 → 外键通常是 `OrderID`

**自定义外键**：使用 `gorm:"foreignKey:CustomKey"` 标签

#### 模型标签

在定义关联关系时，可以使用以下标签：

- `foreignKey`：指定外键字段名
- `references`：指定引用的字段（默认是主键）
- `constraint`：指定外键约束（如：`OnUpdate:CASCADE,OnDelete:SET NULL`）

参考示例：
- `examples/advanced/associations_models.go` - 模型定义
- `examples/advanced/associations_basics_test.go` - 基础关联关系（Has One, Has Many, Belongs To）
- `examples/advanced/associations_many2many_test.go` - Many-to-Many 关系
- `examples/advanced/associations_api_test.go` - Association API 操作
- `examples/advanced/associations_preload_test.go` - Preload 操作

### 2.2 预加载（Preload）

#### 什么是预加载？

预加载（Eager Loading）是在查询主记录时，同时加载关联记录的技术。相比懒加载（Lazy Loading），预加载可以避免 N+1 查询问题。

#### N+1 查询问题

**什么是 N+1 查询问题？**

N+1 查询问题是 ORM 中常见的性能问题，指的是：
- **1 次查询**：获取主记录列表（如查询所有用户）
- **N 次查询**：在循环中为每条主记录查询关联数据（如为每个用户查询订单）
- **总计**：1 + N 次数据库查询

当主记录数量（N）很大时，会产生大量数据库查询，严重影响性能。

**懒加载导致 N+1 问题的示例**：

```go
// 场景：查询所有用户及其订单
var users []User
db.Find(&users)  // 第 1 次查询：SELECT * FROM users

// 假设查询到 100 个用户
for _, user := range users {
    // 为每个用户单独查询订单（100 次查询）
    db.Model(&user).Association("Orders").Find(&user.Orders)
    // 第 2 次查询：SELECT * FROM orders WHERE user_id = 1
    // 第 3 次查询：SELECT * FROM orders WHERE user_id = 2
    // ...
    // 第 101 次查询：SELECT * FROM orders WHERE user_id = 100
}
// 总共：1 + 100 = 101 次查询！
```

**性能影响**：

假设每次数据库查询耗时 10ms：
- **懒加载**：101 次查询 × 10ms = **1010ms**（超过 1 秒）
- **预加载**：2 次查询 × 10ms = **20ms**（仅 20 毫秒）

性能提升：**50 倍**！

**为什么会出现 N+1 问题？**

1. **懒加载机制**：ORM 默认使用懒加载，只有在访问关联字段时才执行查询
2. **循环访问**：在循环中访问关联字段，每次循环都会触发一次查询
3. **缺乏批量优化**：ORM 无法自动识别批量查询场景，无法合并查询

**预加载如何解决 N+1 问题？**

预加载（Eager Loading）的核心思想是：**在查询主记录时，一次性批量加载所有关联数据**。

**预加载的解决方案**：

```go
// 使用 Preload 预加载关联数据
var users []User
db.Preload("Orders").Find(&users)
// 第 1 次查询：SELECT * FROM users
// 第 2 次查询：SELECT * FROM orders WHERE user_id IN (1, 2, 3, ..., 100)
// 总共：2 次查询！

// 现在可以直接访问关联数据，无需额外查询
for _, user := range users {
    fmt.Println(user.Name, len(user.Orders))  // 不会触发额外查询
}
```

**预加载的工作原理**：

1. **第一步**：查询主记录
   ```sql
   SELECT * FROM users;
   ```

2. **第二步**：收集所有主记录的 ID，批量查询关联数据
   ```sql
   SELECT * FROM orders WHERE user_id IN (1, 2, 3, ..., 100);
   ```

3. **第三步**：在内存中将关联数据匹配到对应的主记录

**性能对比总结**：

| 方式 | 查询次数 | 耗时（100 用户） | 适用场景 |
|------|---------|-----------------|---------|
| 懒加载 | 1 + N | ~1010ms | 只访问少量记录的关联数据 |
| 预加载 | 2 | ~20ms | 需要访问多条记录的关联数据（推荐） |

**最佳实践**：

- ✅ **总是使用 Preload**：当需要访问关联数据时，优先使用 Preload
- ✅ **条件预加载**：使用条件预加载减少数据量（如只加载已支付的订单）
- ❌ **避免循环查询**：不要在循环中使用 `Association().Find()` 查询关联数据

#### 基础 Preload

**加载单个关联**：
```go
var user User
db.Preload("Profile").First(&user, 1)
// 这会执行两个查询：
// 1. SELECT * FROM users WHERE id = 1
// 2. SELECT * FROM profiles WHERE user_id = 1
```

**加载多个关联**：
```go
db.Preload("Profile").Preload("Orders").First(&user, 1)
```

#### 条件预加载

只加载满足条件的关联记录：

```go
// 只加载状态为 PAID 的订单
db.Preload("Orders", "status = ?", "PAID").First(&user, 1)
```

#### 嵌套预加载

加载多层嵌套的关联：

```go
// 加载订单，以及订单中的商品项，以及商品项对应的商品
db.Preload("Orders.Items.Product").First(&user, 1)
```

#### 预加载所有关联

使用 `clause.Associations` 预加载所有直接关联：

```go
var order Order
db.Preload(clause.Associations).First(&order, 1)
// 这会预加载 Order 的所有直接关联（如 Items, User 等）
// 注意：只预加载直接关联，不包括嵌套关联
```

#### Preload vs Joins

**Preload**：
- 使用多个独立的 SQL 查询
- 适合复杂的关联关系
- 不会导致数据重复

**Joins**：
- 使用单个 SQL 查询（JOIN）
- 适合简单的关联关系
- 可能导致数据重复（需要 DISTINCT）

**选择建议**：
- 简单的一对一、一对多：可以使用 Joins
- 复杂的嵌套关联：使用 Preload
- 需要聚合查询：使用 Joins

### 2.3 创建带关联的记录

#### 使用 Create 创建关联

当创建主记录时，可以同时创建关联记录：

```go
user := User{
    Name:  "Alice",
    Email: "alice@example.com",
    Profile: profile{
        Nickname: "阿狸",
        Phone:    "13800000000",
    },
    Orders: []order{
        {OrderNo: "ORDER-001", Status: "PAID"},
        {OrderNo: "ORDER-002", Status: "PENDING"},
    },
}

// 使用 FullSaveAssociations 确保所有关联都被保存
db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&user)
```

#### FullSaveAssociations

`FullSaveAssociations: true` 的作用：
- 保存所有关联记录，即使它们是零值
- 确保嵌套的关联结构被正确创建
- 自动设置外键关系

#### 使用 Select 控制关联保存

```go
// 只保存指定的关联
db.Select("Profile", "Orders").Create(&user)

// 排除指定的关联
db.Omit("Orders").Create(&user)
```

### 2.4 事务（Transaction）

#### 为什么需要事务？

事务用于确保一组数据库操作要么全部成功，要么全部失败，保证数据一致性。

**典型场景**：
- 转账操作：扣款和加款必须同时成功或失败
- 订单创建：创建订单和扣减库存必须原子性
- 批量操作：多个相关操作必须作为一个整体

#### 自动事务（推荐）

使用 `db.Transaction` 自动管理事务：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // 所有操作都在事务中
    if err := tx.Create(&order).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    return nil  // 返回 nil 会自动提交
})
```

**特点**：
- 自动提交：函数返回 `nil` 时自动提交
- 自动回滚：函数返回错误时自动回滚
- 自动处理 panic：panic 也会触发回滚

#### 手动事务

手动控制事务的开始、提交和回滚：

```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&order).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Commit().Error; err != nil {
    tx.Rollback()
    return err
}
```

#### SavePoint（保存点）

SavePoint 允许在事务中创建检查点，可以回滚到特定点而不回滚整个事务：

```go
db.Transaction(func(tx *gorm.DB) error {
    // 操作1
    tx.Create(&order)
    
    // 创建保存点
    tx.SavePoint("sp1")
    
    // 操作2
    if err := tx.Create(&item).Error; err != nil {
        // 回滚到保存点（不影响操作1）
        tx.RollbackTo("sp1")
        return err
    }
    
    return nil
})
```

#### 嵌套事务

GORM 支持嵌套事务，内层事务失败不会影响外层事务：

```go
db.Transaction(func(tx1 *gorm.DB) error {
    // 外层事务
    tx1.Create(&order)
    
    return tx1.Transaction(func(tx2 *gorm.DB) error {
        // 内层事务
        tx2.Create(&item)
        return nil
    })
})
```

#### 幂等性设计

在事务中实现幂等性，防止重复操作：

```go
func transfer(db *gorm.DB, input transferInput) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 检查是否已存在相同的转账记录（幂等性）
        var existing transferRecord
        if err := tx.Where("reference = ?", input.Reference).First(&existing).Error; err == nil {
            return errDuplicateTransfer  // 已存在，返回错误
        }
        
        // 执行转账逻辑
        // ...
        return nil
    })
}
```

参考示例：`examples/advanced/transactions_test.go`

### 2.5 钩子（Hooks）

#### 示例代码文件结构

钩子、软删除和乐观锁的示例代码已拆分为多个测试函数，便于分模块学习和讲解：

- **`hooks_softdelete_test.go`**：包含所有钩子、软删除和乐观锁的示例
  - `TestHookBeforeCreate`：BeforeCreate 钩子示例
  - `TestHookBeforeUpdate`：BeforeUpdate 钩子示例
  - `TestHookBeforeDelete`：BeforeDelete 钩子示例
  - `TestSoftDeleteBasic`：软删除基本行为示例
  - `TestSoftDeleteUnscoped`：Unscoped 查询已删除记录示例
  - `TestSoftDeleteHardDelete`：永久删除示例
  - `TestOptimisticLock`：乐观锁示例

#### 什么是钩子？

钩子是 GORM 提供的生命周期回调函数，允许在特定时机执行自定义逻辑。

#### 钩子类型

**创建钩子**：
- `BeforeCreate`：创建前执行
- `AfterCreate`：创建后执行

**更新钩子**：
- `BeforeUpdate`：更新前执行
- `AfterUpdate`：更新后执行

**查询钩子**：
- `AfterFind`：查询后执行

**删除钩子**：
- `BeforeDelete`：删除前执行
- `AfterDelete`：删除后执行

#### 钩子示例

**完整的模型定义**（包含审计字段和软删除）：

```go
type AuditFields struct {
    CreatedBy string
    UpdatedBy string
    DeletedBy string
}

type Article struct {
    ID        uint           `gorm:"primaryKey"`
    Title     string         `gorm:"size:128;not null"`
    Content   string         `gorm:"size:512"`
    Version   int            `gorm:"version"`           // 乐观锁字段
    Audit     AuditFields    `gorm:"embedded"`           // 嵌入审计字段
    DeletedAt gorm.DeletedAt `gorm:"index"`            // 软删除字段
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**BeforeCreate 钩子**：

```go
func (a *Article) BeforeCreate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.CreatedBy = user
    a.Audit.UpdatedBy = user
    // 对于 embedded 字段，使用扁平化的字段名（snake_case）
    tx.Statement.SetColumn("created_by", user)
    tx.Statement.SetColumn("updated_by", user)
    return nil
}
```

**BeforeUpdate 钩子**：

```go
func (a *Article) BeforeUpdate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.UpdatedBy = user
    // 对于 embedded 字段，使用扁平化的字段名（snake_case）
    tx.Statement.SetColumn("updated_by", user)
    return nil
}
```

**BeforeDelete 钩子**：

```go
func (a *Article) BeforeDelete(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.DeletedBy = user
    // 对于 embedded 字段，使用扁平化的字段名（snake_case）
    tx.Statement.SetColumn("deleted_by", user)
    return nil
}
```

**注意**：对于软删除，`BeforeDelete` 钩子会被触发，但 `SetColumn` 在软删除的 UPDATE 语句中可能不起作用。实际使用时，建议在删除前先更新 `deleted_by` 字段，然后再执行删除操作。

#### 使用 Context 传递额外信息

钩子中可以通过 Context 获取额外信息（如当前用户）：

```go
type ctxKey string

const ctxKeyOperator ctxKey = "operator"

func withOperator(name string) context.Context {
    return context.WithValue(context.Background(), ctxKeyOperator, name)
}

func getCurrentUser(tx *gorm.DB) string {
    if tx != nil && tx.Statement != nil && tx.Statement.Context != nil {
        if v, ok := tx.Statement.Context.Value(ctxKeyOperator).(string); ok && v != "" {
            return v
        }
    }
    return "system"
}

// 使用示例
ctx := withOperator("alice")
db.WithContext(ctx).Create(&article)  // BeforeCreate 钩子会获取到 "alice"
```

#### 使用 SetColumn 设置字段

在钩子中使用 `SetColumn` 设置字段值（特别是 embedded 字段）：

```go
func (a *Article) BeforeCreate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    // 对于 embedded 字段，使用扁平化的字段名（snake_case）
    tx.Statement.SetColumn("created_by", user)
    tx.Statement.SetColumn("updated_by", user)
    return nil
}
```

#### 钩子执行顺序

```
BeforeCreate → SQL INSERT → AfterCreate
BeforeUpdate → SQL UPDATE → AfterUpdate
BeforeDelete → SQL DELETE → AfterDelete
```

### 2.6 软删除（Soft Delete）

#### 什么是软删除？

软删除是指不真正删除数据库记录，而是通过标记字段（如 `deleted_at`）来标记记录已被删除。

#### 启用软删除

使用 `gorm.DeletedAt` 类型启用软删除：

```go
type Article struct {
    ID        uint
    Title     string
    DeletedAt gorm.DeletedAt `gorm:"index"`  // 软删除字段
}
```

#### 软删除的行为

**删除操作**：
```go
db.Delete(&article)  // 不会真正删除，而是设置 deleted_at
// SQL: UPDATE articles SET deleted_at = NOW() WHERE id = 1
```

**查询操作**：
```go
db.First(&article, 1)  // 自动过滤已删除的记录
// SQL: SELECT * FROM articles WHERE id = 1 AND deleted_at IS NULL
// 如果记录已被软删除，会返回 ErrRecordNotFound
```

**查询已删除的记录**：
```go
db.Unscoped().First(&article, 1)  // 包含已删除的记录
// SQL: SELECT * FROM articles WHERE id = 1
```

**永久删除**：
```go
db.Unscoped().Delete(&article)  // 真正删除记录
// SQL: DELETE FROM articles WHERE id = 1
```

#### 软删除与 BeforeDelete 钩子

软删除时会触发 `BeforeDelete` 钩子，但需要注意：

1. **SetColumn 的限制**：在软删除的 UPDATE 语句中，`SetColumn` 可能不起作用
2. **推荐做法**：在删除前先更新 `deleted_by` 字段，然后再执行删除操作

```go
// 推荐做法：先更新 deleted_by，再删除
ctx := withOperator("charlie")
db.WithContext(ctx).
    Model(&article{}).
    Where("id = ?", articleID).
    Update("deleted_by", "charlie")
db.WithContext(ctx).Delete(&article{}, articleID)
```

#### GORM 如何区分软删除和硬删除？

GORM 通过以下机制自动区分软删除和硬删除：

**1. 模型定义检测**

当模型包含 `gorm.DeletedAt` 类型的字段时，GORM 会自动为该字段注册 `DeleteClauses`：

```go
type Article struct {
    ID        uint
    Title     string
    DeletedAt gorm.DeletedAt `gorm:"index"`  // 关键：gorm.DeletedAt 类型
}
```

**2. DeleteClauses 机制**

`gorm.DeletedAt` 类型实现了 `DeleteClauses` 方法，返回 `SoftDeleteDeleteClause`。当调用 `Delete()` 方法时：

- GORM 会检查模型的 Schema 中是否有 `DeleteClauses`
- 如果有，会调用这些 clauses 的 `ModifyStatement` 方法

**3. Unscoped 标志判断**

在 `SoftDeleteDeleteClause.ModifyStatement` 中，关键判断逻辑如下：

```go
// 伪代码展示核心逻辑
if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
    // 执行软删除：将 DELETE 语句转换为 UPDATE 语句
    // UPDATE articles SET deleted_at = NOW() WHERE id = 1
    stmt.AddClause(clause.Set{{Column: "deleted_at", Value: curTime}})
    stmt.AddClause(clause.Update{})
} else {
    // 执行硬删除：使用 DELETE 语句
    // DELETE FROM articles WHERE id = 1
    stmt.AddClause(clause.Delete{})
}
```

**4. Unscoped() 方法的作用**

`Unscoped()` 方法会设置 `Statement.Unscoped = true`：

```go
func (db *DB) Unscoped() (tx *DB) {
    tx = db.getInstance()
    tx.Statement.Unscoped = true  // 关键：设置标志为 true
    return
}
```

**5. 执行流程对比**

**软删除流程**（默认，`Unscoped = false`）：
```
db.Delete(&article, 1)
  ↓
检查模型是否有 gorm.DeletedAt 字段 → 有
  ↓
调用 SoftDeleteDeleteClause.ModifyStatement
  ↓
检查 Unscoped 标志 → false
  ↓
构建 UPDATE 语句：UPDATE articles SET deleted_at = NOW() WHERE id = 1
```

**硬删除流程**（使用 `Unscoped()`）：
```
db.Unscoped().Delete(&article, 1)
  ↓
设置 Statement.Unscoped = true
  ↓
检查模型是否有 gorm.DeletedAt 字段 → 有
  ↓
调用 SoftDeleteDeleteClause.ModifyStatement
  ↓
检查 Unscoped 标志 → true（跳过软删除逻辑）
  ↓
构建 DELETE 语句：DELETE FROM articles WHERE id = 1
```

**6. 查询时的软删除过滤**

在查询时，`gorm.DeletedAt` 也会自动添加过滤条件：

```go
// 普通查询
db.First(&article, 1)
// SQL: SELECT * FROM articles WHERE id = 1 AND deleted_at IS NULL

// Unscoped 查询
db.Unscoped().First(&article, 1)
// SQL: SELECT * FROM articles WHERE id = 1  （不包含 deleted_at 过滤）
```

**总结**：

- **软删除**：模型有 `gorm.DeletedAt` 字段 + `Unscoped = false`（默认） → 执行 UPDATE
- **硬删除**：使用 `Unscoped()` 设置 `Unscoped = true` → 执行 DELETE

#### 软删除的优势

1. **数据恢复**：可以恢复"删除"的数据
2. **审计追踪**：保留历史数据用于审计
3. **关联完整性**：不会破坏外键约束

参考示例：
- `examples/advanced/hooks_softdelete_test.go` - `TestSoftDeleteBasic`：软删除基本行为
- `examples/advanced/hooks_softdelete_test.go` - `TestSoftDeleteUnscoped`：查询已删除记录
- `examples/advanced/hooks_softdelete_test.go` - `TestSoftDeleteHardDelete`：永久删除

### 2.7 乐观锁（Optimistic Locking）

#### 什么是乐观锁？

乐观锁是一种并发控制机制，通过版本号来检测数据是否被其他事务修改。

#### 实现乐观锁

使用 `gorm:"version"` 标签启用乐观锁：

```go
type Article struct {
    ID      uint
    Title   string
    Version int `gorm:"version"`  // 乐观锁字段
}
```

#### 乐观锁的工作原理

1. **读取时**：获取当前版本号
2. **更新时**：检查版本号是否匹配
3. **版本不匹配**：更新失败，影响 0 行

#### 使用 Updates 方法更新（推荐）

使用 `Updates` 方法时，需要手动递增版本号，并通过检查 `RowsAffected` 来判断是否更新成功：

```go
// 第一次更新，版本号从 0 变为 1
if err := db.First(&article, articleID).Error; err != nil {
    return err
}

result := db.Model(&article).
    Select("content", "updated_by", "version").
    Updates(map[string]any{
        "content": "新内容",
        "version": gorm.Expr("version + 1"), // 手动递增版本号
    })

if result.Error != nil {
    return result.Error
}
// 版本号已更新为 1

// 使用旧版本号尝试更新（会失败）
stale := Article{ID: articleID, Version: 0}
result = db.Model(&stale).
    Where("version = ?", 0). // 使用旧版本号
    Updates(map[string]any{"content": "尝试使用旧版本更新"})

if result.Error != nil {
    return result.Error
}
// 检查更新的行数，如果版本号不匹配，应该更新 0 行
if result.RowsAffected == 0 {
    return errors.New("optimistic lock conflict: version mismatch")
}
```

**关键点**：
- `Updates` 方法不会自动递增版本号，需要手动使用 `gorm.Expr("version + 1")`
- 版本号不匹配时，更新会影响 0 行，但不会返回错误
- 需要通过检查 `RowsAffected` 来判断是否更新成功

#### 乐观锁 vs 悲观锁

**乐观锁**：
- 假设冲突很少发生
- 通过版本号检测冲突
- 性能较好，适合读多写少场景
- 推荐使用

**悲观锁**：
- 假设冲突经常发生
- 通过数据库锁机制（SELECT FOR UPDATE）
- 性能较差，适合写多场景
- ⚠️ **不推荐使用**，原因如下：
  1. **死锁风险**：多个事务以不同顺序锁定资源时容易产生死锁
  2. **性能问题**：
     - 阻塞其他事务：被锁定的记录会阻塞其他需要修改该记录的事务
     - 锁持有时间长：锁会持续到事务结束，如果事务中有慢操作，锁持有时间会更长
     - 并发性能差：高并发场景下，大量事务会排队等待锁释放
  3. **资源浪费**：即使事务最终可能失败，锁也会一直持有直到事务结束
  4. **扩展性差**：随着并发量增加，性能会急剧下降

**注意**：
- `SELECT FOR UPDATE` 是**行锁**，不是表锁，但如果锁定的行很多，影响范围也会很大
- 如果必须使用悲观锁，确保：
  - 锁定顺序一致（避免死锁）
  - 事务尽可能短（减少锁持有时间）
  - 只锁定必要的记录（避免锁范围过大）
  - 考虑使用超时机制（避免长时间等待）

### 2.8 审计字段（Audit Fields）

#### 什么是审计字段？

审计字段用于记录数据的创建人、更新人、删除人等信息，便于追踪数据变更历史。

#### 实现审计字段

```go
type AuditFields struct {
    CreatedBy string
    UpdatedBy string
    DeletedBy string
}

type Article struct {
    ID        uint
    Title     string
    Audit     AuditFields `gorm:"embedded"`  // 嵌入审计字段
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### 在钩子中设置审计字段

```go
func (a *Article) BeforeCreate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.CreatedBy = user
    a.Audit.UpdatedBy = user
    // 对于 embedded 字段，使用扁平化的字段名（snake_case）
    tx.Statement.SetColumn("created_by", user)
    tx.Statement.SetColumn("updated_by", user)
    return nil
}

func (a *Article) BeforeUpdate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.UpdatedBy = user
    tx.Statement.SetColumn("updated_by", user)
    return nil
}

func (a *Article) BeforeDelete(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.DeletedBy = user
    tx.Statement.SetColumn("deleted_by", user)
    return nil
}
```

**注意**：对于软删除，`BeforeDelete` 钩子会被触发，但 `SetColumn` 在软删除的 UPDATE 语句中可能不起作用。实际使用时，建议在删除前先更新 `deleted_by` 字段。

参考示例：
- `examples/advanced/hooks_softdelete_test.go` - `TestHookBeforeCreate`：BeforeCreate 钩子设置审计字段
- `examples/advanced/hooks_softdelete_test.go` - `TestHookBeforeUpdate`：BeforeUpdate 钩子设置审计字段
- `examples/advanced/hooks_softdelete_test.go` - `TestHookBeforeDelete`：BeforeDelete 钩子设置审计字段

## 3. 关键概念总结

### 3.1 关联关系要点

1. **外键命名**：默认使用 `ModelName + ID`，可通过标签自定义
2. **关联字段**：必须导出（首字母大写）
3. **迁移顺序**：所有相关模型必须一起迁移
4. **N+1 问题**：使用 Preload 避免懒加载导致的性能问题

### 3.2 预加载要点

1. **Preload vs Joins**：Preload 适合复杂关联，Joins 适合简单关联
2. **条件预加载**：可以过滤关联记录
3. **嵌套预加载**：使用点号分隔（如 `Orders.Items`）
4. **clause.Associations**：预加载所有直接关联

### 3.3 事务要点

1. **自动事务**：使用 `db.Transaction` 简化事务管理
2. **错误处理**：返回错误会自动回滚
3. **Panic 处理**：panic 也会触发自动回滚
4. **嵌套事务**：支持多层嵌套
5. **SavePoint**：支持部分回滚

### 3.4 钩子要点

1. **执行顺序**：Before → SQL → After
2. **SetColumn**：用于设置字段值（特别是 embedded 字段）
3. **Context**：可以通过 Context 传递额外信息（如当前用户）

### 3.5 软删除要点

1. **自动过滤**：查询时自动过滤已删除记录
2. **Unscoped**：使用 `Unscoped()` 查询已删除记录
3. **永久删除**：使用 `Unscoped().Delete()` 真正删除
4. **BeforeDelete 钩子**：软删除时会触发，但 `SetColumn` 可能不起作用，建议在删除前先更新 `deleted_by` 字段

### 3.6 乐观锁要点

1. **版本字段**：使用 `gorm:"version"` 标签
2. **版本递增**：使用 `Updates` 时需要手动递增版本号（`gorm.Expr("version + 1")`）
3. **冲突检测**：版本不匹配时更新会影响 0 行，需要通过 `RowsAffected` 判断是否成功
4. **最佳实践**：使用 `Select` 指定要更新的字段，包括版本号字段

## 4. 实践练习

### 练习 1：设计博客系统关联关系

设计以下模型的关联关系：
- `User`：用户
- `Post`：文章
- `Comment`：评论
- `Tag`：标签

**要求**：
- User Has Many Posts
- Post Belongs To User
- Post Has Many Comments
- Comment Belongs To Post
- Post Many to Many Tags

### 练习 2：实现查询功能

1. **查询用户最新文章**：查询用户发表的最新 10 篇文章（含标签）
   ```go
   func GetUserLatestPosts(db *gorm.DB, userID uint) ([]Post, error) {
       // 你的实现
   }
   ```

2. **统计评论数量**：使用 Preload + Count 统计每篇文章的评论数量
   ```go
   func GetPostsWithCommentCount(db *gorm.DB) ([]PostWithCount, error) {
       // 你的实现
   }
   ```

### 练习 3：实现事务操作

在事务中实现文章发布 + 标签绑定：

```go
func PublishPostWithTags(db *gorm.DB, post *Post, tagIDs []uint) error {
    // 在事务中：
    // 1. 创建文章
    // 2. 绑定标签（Many to Many）
    // 3. 更新用户文章数量
    // 你的实现
}
```

### 练习 4：实现软删除

为评论新增软删除，提供"彻底清除"功能：

```go
// 软删除评论
func SoftDeleteComment(db *gorm.DB, commentID uint) error {
    // 你的实现
}

// 彻底删除评论
func HardDeleteComment(db *gorm.DB, commentID uint) error {
    // 你的实现
}
```

## 5. 运行示例代码

### 运行关联关系示例

关联关系的示例代码已拆分为多个测试文件，便于学习和讲解：

```bash
cd lesson-02/examples

# 1. 基础关联关系（Has One, Has Many, Belongs To）
go test ./advanced -run TestAssociationsBasics -v

# 2. Many-to-Many 关系
go test ./advanced -run TestAssociationsManyToMany -v

# 3. Association API 操作（Append, Replace, Delete, Clear, Count, Find）
go test ./advanced -run TestAssociationsAPI -v

# 4. Preload 操作（条件预加载、嵌套预加载、clause.Associations）
go test ./advanced -run TestAssociationsPreload -v

# 运行所有关联关系测试
go test ./advanced -run TestAssociations -v
```

### 运行事务示例

```bash
# 事务控制
go test ./advanced -run TestTransactionsDemo -v
```

### 运行钩子与软删除示例

```bash
# 运行所有钩子、软删除、乐观锁测试
go test ./advanced -run "TestHook|TestSoftDelete|TestOptimisticLock" -v

# 运行单个测试
go test ./advanced -run TestHookBeforeCreate -v      # BeforeCreate 钩子
go test ./advanced -run TestHookBeforeUpdate -v      # BeforeUpdate 钩子
go test ./advanced -run TestHookBeforeDelete -v      # BeforeDelete 钩子
go test ./advanced -run TestSoftDeleteBasic -v       # 软删除基本行为
go test ./advanced -run TestSoftDeleteUnscoped -v     # Unscoped 查询已删除记录
go test ./advanced -run TestSoftDeleteHardDelete -v   # 永久删除
go test ./advanced -run TestOptimisticLock -v         # 乐观锁
```

### 查看 SQL 日志

在 `testutil/helpers.go` 中将 Logger 设置为 `logger.Info`，即可看到所有 SQL 查询：

```go
Logger: logger.Default.LogMode(logger.Info),
```

## 6. 常见问题

### Q: 如何避免 N+1 查询问题？

A: 使用 Preload 预加载关联记录，而不是在循环中查询。

### Q: Preload 和 Joins 什么时候用哪个？

A: 
- 简单的一对一、一对多：可以使用 Joins
- 复杂的嵌套关联：使用 Preload
- 需要聚合查询：使用 Joins

### Q: 事务中可以使用 goroutine 吗？

A: 不建议。事务绑定到特定的数据库连接，goroutine 可能导致连接混乱。

### Q: 软删除后如何恢复数据？

A: 使用 `Unscoped().Update()` 将 `deleted_at` 设置为 NULL。

### Q: 乐观锁更新失败怎么办？

A: 重新读取最新数据，获取新的版本号，然后重试更新。

### Q: 钩子中可以访问数据库吗？

A: 可以，但要注意避免无限递归。建议使用传入的 `tx *gorm.DB` 参数。

## 7. 最佳实践

1. **关联关系**：
   - 明确定义关联方向（Has One / Belongs To / Has Many）
   - 使用 Preload 避免 N+1 问题
   - 合理使用条件预加载减少数据量

2. **事务**：
   - 优先使用 `db.Transaction` 自动事务
   - 保持事务尽可能短
   - 避免在事务中执行长时间操作

3. **钩子**：
   - 钩子中不要执行复杂逻辑
   - 使用 Context 传递额外信息
   - 注意 embedded 字段的 SetColumn 用法

4. **软删除**：
   - 重要数据使用软删除
   - 定期清理过期的软删除数据
   - 使用索引优化查询性能

5. **乐观锁**：
   - 读多写少场景使用乐观锁
   - 实现重试机制处理版本冲突
   - 记录版本冲突日志用于分析

## 8. 下一步学习

- 深入阅读示例代码中的注释
- 完成实践练习
- 学习项目实战：电商系统完整实现
- 将 GORM 与 Web 框架（如 Gin）结合，构建完整的应用

---

**完成以上内容即可掌握 GORM 的关联关系、事务、钩子等高级功能！** 🚀
