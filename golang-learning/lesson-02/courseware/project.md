# 第2期教案：GORM 项目实战 - 电商订单系统

## 1. 学习目标

通过本课程的学习，你将能够：

- 将 GORM 技能应用到真实的电商业务场景，构建完整的订单系统数据层
- 设计合理的数据模型、索引与约束，确保数据一致性和查询性能
- 实现下单、库存扣减、订单支付等核心业务流程
- 掌握事务处理、并发控制、错误处理等关键技能
- 使用 GORM 进行复杂查询和聚合报表统计

## 2. 项目场景

### 2.1 业务需求

我们将构建一个简化的电商订单系统，包含以下核心功能：

1. **用户管理**：用户可以注册和登录
2. **商品管理**：商品信息包括名称、价格、库存等
3. **订单管理**：
   - 用户可以下单购买多个商品
   - 系统需要检查库存，扣减数量，生成订单
   - 支持订单支付状态更新
4. **数据统计**：提供按天聚合的销售报表

### 2.2 核心实体

系统包含以下核心实体：

- **User（用户）**：用户基本信息
- **Product（商品）**：商品信息和库存
- **Order（订单）**：订单主信息
- **OrderItem（订单项）**：订单中的商品明细

## 3. 数据模型设计

### 3.1 模型定义

让我们看看每个模型的定义和设计考虑：

#### User（用户模型）

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"size:128;uniqueIndex;not null"`
    Name      string    `gorm:"size:64;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**设计要点**：
- `Email` 使用 `uniqueIndex` 确保邮箱唯一性，同时提供快速查询
- 如果后续需要按名称搜索，可以添加 `Name` 字段的普通索引

#### Product（商品模型）

```go
type Product struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:128;not null"`
    SKU       string    `gorm:"size:32;uniqueIndex;not null"`
    Price     int64     `gorm:"not null"`
    Stock     int       `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**设计要点**：
- `SKU`（商品编码）使用 `uniqueIndex` 确保唯一性
- `Price` 使用 `int64` 存储，单位为分（避免浮点数精度问题）
- `Stock` 需要在事务中锁定更新，防止并发扣减导致超卖
- `Name` 可以添加普通索引用于商品名称搜索

#### Order（订单模型）

```go
type Order struct {
    ID          uint        `gorm:"primaryKey"`
    OrderNo     string      `gorm:"size:32;uniqueIndex;not null"`
    UserID      uint        `gorm:"index;not null"`
    TotalAmount int64       `gorm:"not null"`
    Status      string      `gorm:"size:16;index;not null"`
    PaidAt      *time.Time  `gorm:"index"`
    Items       []OrderItem
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**设计要点**：
- `OrderNo` 使用 `uniqueIndex` 确保订单号唯一性，支持幂等设计
- `UserID` 使用普通索引，用于查询用户的订单列表
- `Status` 使用普通索引，用于按状态筛选订单（如查询待支付订单）
- `PaidAt` 使用普通索引，用于查询已支付订单的时间范围
- `Items` 是 Has Many 关联，一个订单包含多个订单项
- 状态枚举：`PENDING`（待支付）、`PAID`（已支付）、`CANCELLED`（已取消）

#### OrderItem（订单项模型）

```go
type OrderItem struct {
    ID        uint      `gorm:"primaryKey"`
    OrderID   uint      `gorm:"index;not null"`
    ProductID uint      `gorm:"index;not null"`
    Product   Product
    Quantity  int       `gorm:"not null"`
    UnitPrice int64     `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**设计要点**：
- `OrderID` 使用普通索引，用于查询订单的所有订单项
- `ProductID` 使用普通索引，用于统计商品的销售情况
- `UnitPrice` 存储下单时的商品单价（快照），避免商品价格变更影响历史订单
- `Product` 是 Belongs To 关联，用于预加载商品信息

### 3.2 索引设计总结

| 表名 | 索引字段 | 索引类型 | 用途 |
|------|----------|----------|------|
| `users` | `email` | 唯一索引 | 确保邮箱唯一，快速查询 |
| `products` | `sku` | 唯一索引 | 确保商品编码唯一 |
| `products` | `name` | 普通索引（建议） | 商品名称搜索 |
| `orders` | `order_no` | 唯一索引 | 确保订单号唯一，幂等设计 |
| `orders` | `user_id` | 普通索引 | 查询用户的订单列表 |
| `orders` | `status` | 普通索引 | 按状态筛选订单 |
| `orders` | `paid_at` | 普通索引 | 查询已支付订单的时间范围 |
| `order_items` | `order_id` | 普通索引 | 查询订单的所有订单项 |
| `order_items` | `product_id` | 普通索引 | 统计商品的销售情况 |
| `order_items` | `(product_id, created_at)` | 复合索引（建议） | 按商品和时间范围查询 |

**索引设计原则**：
- 唯一性约束使用唯一索引
- 频繁查询的字段添加普通索引
- 复合查询场景使用复合索引
- 避免过度索引（影响写入性能）

## 4. 核心功能实现

### 4.1 数据库初始化

#### 4.1.1 模型迁移

使用 `AutoMigrate` 自动创建表结构：

```go
func migrate(db *gorm.DB) error {
    return db.AutoMigrate(&User{}, &Product{}, &Order{}, &OrderItem{})
}
```

**要点**：
- `AutoMigrate` 会根据模型定义创建表结构
- 如果表已存在，会根据模型定义更新表结构（添加新字段、索引等）
- 不会删除已存在的字段（安全设计）
- 所有相关的模型必须一起迁移，确保外键关系正确创建

#### 4.1.2 初始化数据

```go
func seedData(db *gorm.DB) error {
    // 检查用户表是否为空
    var count int64
    if err := db.Model(&User{}).Count(&count).Error; err != nil {
        return err
    }
    if count == 0 {
        users := []User{
            {Name: "Alice", Email: "alice@example.com"},
            {Name: "Bob", Email: "bob@example.com"},
        }
        if err := db.Create(&users).Error; err != nil {
            return err
        }
    }
    
    // 检查商品表是否为空
    if err := db.Model(&Product{}).Count(&count).Error; err != nil {
        return err
    }
    if count == 0 {
        products := []Product{
            {Name: "Go 语言权威指南", SKU: "BOOK-001", Price: 13800, Stock: 50},
            {Name: "机械键盘", SKU: "GEAR-201", Price: 39900, Stock: 20},
            {Name: "人体工学鼠标", SKU: "GEAR-301", Price: 22900, Stock: 15},
        }
        if err := db.Create(&products).Error; err != nil {
            return err
        }
    }
    return nil
}
```

**要点**：
- 使用 `Count` 检查避免重复插入
- 支持多次运行测试，不会重复插入数据

### 4.2 下单流程

#### 4.2.1 业务错误定义

首先定义业务错误类型，便于在业务层进行错误判断：

```go
var (
    errNoItems          = errors.New("order must contain at least one item")
    errOutOfStock       = errors.New("product stock is insufficient")
    errOrderAlreadyPaid = errors.New("order already paid")
)
```

**要点**：
- 使用自定义错误类型便于在业务层进行错误判断和处理
- 使用 `errors.Is` 可以判断是否为特定业务错误

#### 4.2.2 下单函数实现

下单流程是电商系统的核心，需要保证数据一致性和并发安全：

```go
func CreateOrder(ctx context.Context, db *gorm.DB, userID uint, items []OrderItemInput) (*Order, error) {
    // 1. 校验订单项不为空
    if len(items) == 0 {
        return nil, errNoItems
    }

    var order Order
    // 2. 使用事务包装整个下单流程
    err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 2.1 加载用户信息
        var user User
        if err := tx.First(&user, userID).Error; err != nil {
            return fmt.Errorf("load user: %w", err)
        }

        // 2.2 锁定并加载商品信息（使用 FOR UPDATE）
        var products []Product
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("id IN ?", productIDs).
            Find(&products).Error; err != nil {
            return fmt.Errorf("load products: %w", err)
        }

        // 2.3 校验库存并扣减
        for _, item := range items {
            // 校验库存是否充足
            if p.Stock < item.Quantity {
                return fmt.Errorf("%w: %s (需要 %d, 当前 %d)", 
                    errOutOfStock, p.Name, item.Quantity, p.Stock)
            }

            // 扣减库存（使用 UpdateColumn 直接更新）
            if err := tx.Model(&Product{}).
                Where("id = ?", p.ID).
                UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
                return fmt.Errorf("update stock: %w", err)
            }
        }

        // 2.4 生成订单号并创建订单
        order = Order{
            OrderNo:     generateOrderNo(),
            UserID:      user.ID,
            TotalAmount: total,
            Status:      "PENDING",
            Items:       orderItems,
        }
        if err := tx.Create(&order).Error; err != nil {
            return fmt.Errorf("create order: %w", err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }
    return &order, nil
}
```

**关键设计点**：

1. **事务保证原子性**：
   - 使用 `Transaction` 方法包装整个下单流程
   - 如果任何步骤失败，事务自动回滚
   - 库存扣减和订单创建要么全部成功，要么全部失败

2. **FOR UPDATE 锁定**：
   - `clause.Locking{Strength: "UPDATE"}` 相当于 SQL 的 `SELECT ... FOR UPDATE`
   - 这会锁定查询到的商品记录，防止其他事务同时修改库存
   - 锁定会持续到事务结束（提交或回滚）
   - 防止并发下单导致超卖

3. **幂等设计**：
   - 订单号使用唯一索引，确保唯一性
   - 重复下单会因唯一约束失败
   - 订单号格式：`ORD-YYYYMMDD-XXXX`

4. **错误处理**：
   - 使用 `fmt.Errorf` 和 `%w` 包装错误，保留错误链
   - 使用 `errors.Is` 可以判断是否为特定业务错误

#### 4.2.3 订单号生成

```go
func generateOrderNo() string {
    return fmt.Sprintf("ORD-%s-%04d", 
        time.Now().Format("20060102"), 
        rand.Intn(10000))
}
```

**要点**：
- 订单号格式：`ORD-YYYYMMDD-XXXX`
- 虽然随机数可能重复，但结合日期后重复概率极低
- 实际项目中建议使用 UUID 或雪花算法生成唯一ID

### 4.3 支付流程

支付流程需要保证并发安全和幂等性：

```go
func MarkOrderPaid(ctx context.Context, db *gorm.DB, orderNo string) error {
    return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. 锁定并加载订单（使用 FOR UPDATE）
        var order Order
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("order_no = ?", orderNo).
            First(&order).Error; err != nil {
            return fmt.Errorf("load order: %w", err)
        }

        // 2. 幂等性检查：如果订单已经支付，返回错误
        if order.Status == "PAID" {
            return errOrderAlreadyPaid
        }

        // 3. 更新订单状态和支付时间
        now := time.Now()
        return tx.Model(&order).Updates(map[string]any{
            "status":  "PAID",
            "paid_at": &now,
        }).Error
    })
}
```

**关键设计点**：

1. **并发控制**：
   - 使用 `FOR UPDATE` 锁定订单，防止并发支付
   - 如果多个请求同时支付同一订单，只有一个会成功

2. **幂等性**：
   - 检查订单状态，如果已经支付，返回错误
   - 防止重复支付的问题

3. **可扩展性**：
   - 可以在这里添加扣减用户余额、写入支付日志等操作

### 4.4 查询与报表

#### 4.4.1 订单查询（使用 Preload）

使用 `Preload` 预加载关联数据，避免 N+1 查询问题：

```go
func fetchOrders(t *testing.T, db *gorm.DB) []Order {
    var orders []Order
    // Preload("Items.Product"): 预加载订单项及其关联的商品
    // 这是嵌套预加载：Order -> OrderItem -> Product
    if err := db.Preload("Items.Product").
        Order("created_at desc").
        Find(&orders).Error; err != nil {
        t.Fatalf("list orders: %v", err)
    }
    return orders
}
```

**要点**：
- `Preload("Items.Product")` 会：
  1. 先查询所有订单
  2. 再查询这些订单的所有订单项
  3. 最后查询这些订单项关联的商品信息
- 这样只需要 3 次 SQL 查询，而不是 N+1 次（N 为订单数量）

#### 4.4.2 销售报表（聚合查询）

使用聚合查询统计销售数据：

```go
func SalesReport(db *gorm.DB) ([]SalesSummary, error) {
    var rows []SalesSummary
    err := db.Table("orders").
        Select(`
            strftime('%Y-%m-%d', orders.created_at) AS day,
            COUNT(DISTINCT orders.id) AS order_count,
            SUM(order_items.quantity) AS item_count,
            SUM(order_items.quantity * order_items.unit_price) AS total_amount`).
        Joins("JOIN order_items ON order_items.order_id = orders.id").
        Where("orders.status = ?", "PAID").
        Group("day").
        Order("day ASC").
        Scan(&rows).Error
    return rows, err
}
```

**要点**：

1. **聚合函数**：
   - `strftime('%Y-%m-%d', ...)` 是 SQLite 的日期格式化函数
   - MySQL 使用 `DATE(...)`，PostgreSQL 使用 `TO_CHAR(...)`
   - `COUNT(DISTINCT ...)` 统计不重复的订单数量
   - `SUM(...)` 统计商品数量和销售总额

2. **JOIN 关联**：
   - 使用 `Joins` 关联订单项表，用于统计商品数量和计算销售总额

3. **GROUP BY 分组**：
   - 使用 `Group("day")` 按日期分组，将同一天的订单聚合在一起

4. **Scan vs Find**：
   - 聚合查询必须使用 `Scan` 而不是 `Find`
   - 因为结果不直接映射到模型，而是映射到自定义结构体

5. **结果结构体**：

```go
type SalesSummary struct {
    Day         string // 日期（格式：YYYY-MM-DD）
    OrderCount  int64  // 订单数量
    ItemCount   int64  // 商品数量（所有订单项的数量总和）
    TotalAmount int64  // 销售总额（单位：分）
}
```

## 5. 完整测试示例

### 5.1 测试流程

让我们看看完整的测试流程：

```go
func TestEcommerceFlow(t *testing.T) {
    ctx := context.Background()
    db := testutil.NewTestDB(t, "ecommerce.db")

    // 1. 数据模型初始化
    if err := migrate(db); err != nil {
        t.Fatalf("migrate: %v", err)
    }
    
    // 2. 初始化测试数据
    if err := seedData(db); err != nil {
        t.Fatalf("seed data: %v", err)
    }

    // 3. 展示初始库存状态
    products := fetchProducts(t, db)
    for _, p := range products {
        t.Logf(" - %s 库存:%d 价格:%.2f", p.Name, p.Stock, float64(p.Price)/100)
    }

    // 4. 下单流程
    order, err := CreateOrder(ctx, db, 1, []OrderItemInput{
        {ProductID: products[0].ID, Quantity: 1},
        {ProductID: products[1].ID, Quantity: 2},
    })
    if err != nil {
        t.Fatalf("create order: %v", err)
    }
    logOrder(t, db, order.OrderNo)

    // 5. 验证库存扣减
    updated := fetchProducts(t, db)
    for _, p := range updated {
        t.Logf(" - %s 库存:%d", p.Name, p.Stock)
    }

    // 6. 库存不足场景
    _, err = CreateOrder(ctx, db, 1, []OrderItemInput{
        {ProductID: products[1].ID, Quantity: 100},
    })
    if !errors.Is(err, errOutOfStock) {
        t.Fatalf("expected out of stock, got %v", err)
    }

    // 7. 订单支付流程
    if err := MarkOrderPaid(ctx, db, order.OrderNo); err != nil {
        t.Fatalf("mark paid: %v", err)
    }
    logOrder(t, db, order.OrderNo)

    // 8. 订单总览
    orders := fetchOrders(t, db)
    for _, o := range orders {
        t.Logf("- %s [%s] items=%d", o.OrderNo, o.Status, len(o.Items))
    }

    // 9. 销售报表
    report, err := SalesReport(db)
    if err != nil {
        t.Fatalf("sales report: %v", err)
    }
    for _, row := range report {
        t.Logf("%s -> 订单:%d 商品:%d 销售额:%.2f",
            row.Day, row.OrderCount, row.ItemCount, float64(row.TotalAmount)/100)
    }
}
```

### 5.2 运行测试

在 `examples` 目录下运行测试：

```bash
cd examples
go test ./project -run TestEcommerceFlow -v
```

**预期输出**：
1. 初始库存与订单为空
2. 下单成功 → 打印订单信息
3. 再次下单时库存不足 → 打印错误
4. 标记订单支付 → 输出状态变化
5. 查看销售报表 → 输出聚合结果

## 6. 关键知识点总结

### 6.1 事务处理

- **使用场景**：需要保证多个操作原子性的场景（如下单、支付）
- **实现方式**：使用 `db.Transaction(func(tx *gorm.DB) error { ... })`
- **自动回滚**：函数返回 error 时自动回滚，返回 nil 时自动提交

### 6.2 并发控制

- **FOR UPDATE 锁定**：使用 `clause.Locking{Strength: "UPDATE"}` 锁定记录
- **使用场景**：防止并发修改导致的数据不一致（如库存扣减、订单支付）
- **锁定范围**：锁定会持续到事务结束

### 6.3 幂等设计

- **订单号唯一索引**：确保订单号唯一，防止重复下单
- **状态检查**：支付前检查订单状态，防止重复支付
- **错误处理**：使用自定义错误类型，便于业务层判断

### 6.4 预加载优化

- **N+1 问题**：避免在循环中查询关联数据
- **解决方案**：使用 `Preload` 预加载关联数据
- **嵌套预加载**：`Preload("Items.Product")` 可以预加载多层关联

### 6.5 聚合查询

- **使用场景**：统计报表、数据分析
- **实现方式**：使用 `Select`、`Group`、`Scan` 等方法
- **注意事项**：聚合查询必须使用 `Scan` 而不是 `Find`

## 7. 实践练习

### 练习 1：完善订单查询

**任务**：实现一个函数，查询指定用户的所有订单，包含订单项和商品信息。

**提示**：
- 使用 `Preload` 预加载关联数据
- 使用 `Where` 过滤用户ID
- 使用 `Order` 按创建时间降序排序

### 练习 2：实现库存预警

**任务**：实现一个函数，查询库存低于指定阈值的商品。

**提示**：
- 使用 `Where` 过滤库存数量
- 使用 `Find` 查询商品列表

### 练习 3：扩展销售报表

**任务**：扩展 `SalesReport` 函数，支持按商品统计销售数据。

**提示**：
- 修改 `Select` 语句，按商品分组
- 创建新的结果结构体
- 使用 `Joins` 关联商品表

### 练习 4：实现订单取消

**任务**：实现 `CancelOrder` 函数，取消订单并恢复库存。

**提示**：
- 使用事务保证原子性
- 检查订单状态（只能取消待支付订单）
- 恢复库存（增加库存数量）
- 更新订单状态为 `CANCELLED`

## 8. 扩展阅读

### 8.1 性能优化

- **索引优化**：根据查询场景添加合适的索引
- **批量操作**：使用 `CreateInBatches` 批量插入数据
- **连接池配置**：合理配置数据库连接池参数

### 8.2 架构设计

- **分层架构**：将项目拆分为 `repository/service` 层
- **依赖注入**：使用接口实现依赖注入
- **错误处理**：统一错误处理机制

### 8.3 进阶功能

- **优惠券系统**：增加优惠券模型和折扣计算
- **消息队列**：使用消息队列保证异步扣库存
- **数据库迁移**：使用迁移工具管理数据库版本
- **RESTful API**：与 Gin 集成，提供 HTTP API

## 9. 总结

通过本课程的学习，我们完成了以下内容：

1. ✅ 设计了完整的电商订单系统数据模型
2. ✅ 实现了下单、支付等核心业务流程
3. ✅ 掌握了事务处理、并发控制等关键技能
4. ✅ 学会了使用 Preload 优化查询性能
5. ✅ 实现了聚合查询和报表统计

**下一步**：
- 完成实践练习，巩固所学知识
- 阅读扩展阅读，深入了解 GORM 的高级特性
- 尝试将项目拆分为更合理的架构

**本节将理论与实践结合，确保你能在真实业务中落地 GORM。** ✅

