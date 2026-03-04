package advanced

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"lesson02examples/testutil"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 错误定义
var (
	errInsufficientBalance = errors.New("insufficient balance")
	errDuplicateTransfer   = errors.New("duplicate transfer reference")
)

// account 账户模型
// 用于演示转账操作中的账户信息
type account struct {
	ID        uint
	Name      string
	Balance   int64 // 账户余额
	CreatedAt time.Time
	UpdatedAt time.Time
}

// transferRecord 转账记录模型
// 用于记录每次转账的详细信息，Reference 字段用于实现幂等性
type transferRecord struct {
	ID            uint
	Reference     string `gorm:"uniqueIndex"` // 唯一索引，用于幂等性检查
	FromAccountID uint
	ToAccountID   uint
	Amount        int64
	Status        string
	Message       string
	CreatedAt     time.Time
}

// setupDB 测试前置设置函数（类似 Java 的 @Before）
// 负责初始化数据库连接、迁移表结构、重置测试数据
// 每个测试函数都应该在开始时调用此函数
func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	// 创建测试数据库连接
	db := testutil.NewTestDB(t, "transactions.db")

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&account{}, &transferRecord{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// 重置账户数据，确保每次测试都从干净的状态开始
	if err := resetAccounts(t, db); err != nil {
		t.Fatalf("reset accounts: %v", err)
	}

	// 注册清理函数，测试结束后自动清理（类似 Java 的 @After）
	t.Cleanup(func() {
		// 可以在这里添加测试后的清理逻辑
		// 例如：关闭连接、清理临时数据等
	})

	return db
}

// ============================================================================
// 知识点 1: 自动事务 - 自动提交
// ============================================================================

// TestTransactionAutoCommit 测试自动事务的正常提交
// 演示：使用 db.Transaction 自动管理事务，返回 nil 时自动提交
func TestTransactionAutoCommit(t *testing.T) {
	db := setupDB(t)

	// 使用自动事务执行转账操作
	// 特点：函数返回 nil 时，事务会自动提交
	err := db.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 扣减转出账户余额
		if err := tx.Model(&account{}).
			Where("id = ?", 1).
			Update("balance", gorm.Expr("balance - ?", 5000)).Error; err != nil {
			return fmt.Errorf("debit account: %w", err)
		}

		// 步骤2: 增加转入账户余额
		if err := tx.Model(&account{}).
			Where("id = ?", 2).
			Update("balance", gorm.Expr("balance + ?", 5000)).Error; err != nil {
			return fmt.Errorf("credit account: %w", err)
		}

		// 步骤3: 创建转账记录
		record := transferRecord{
			Reference:     "TX-001",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        5000,
			Status:        "SUCCESS",
			Message:       "自动事务测试",
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("create record: %w", err)
		}

		// 返回 nil，事务会自动提交
		return nil
	})

	if err != nil {
		t.Fatalf("transaction failed: %v", err)
	}

	// 验证转账结果
	var accounts []account
	if err := db.Order("id").Find(&accounts).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 验证账户余额变化
	if accounts[0].Balance != 95000 { // 100000 - 5000
		t.Errorf("expected account 1 balance 95000, got %d", accounts[0].Balance)
	}
	if accounts[1].Balance != 35000 { // 30000 + 5000
		t.Errorf("expected account 2 balance 35000, got %d", accounts[1].Balance)
	}

	// 验证转账记录已创建
	var record transferRecord
	if err := db.Where("reference = ?", "TX-001").First(&record).Error; err != nil {
		t.Errorf("transfer record should be created: %v", err)
	}
}

// ============================================================================
// 知识点 2: 自动事务 - 自动回滚
// ============================================================================

// TestTransactionAutoRollback 测试自动事务的回滚
// 演示：当返回 error 时，事务会自动回滚，所有已执行的操作都会被撤销
func TestTransactionAutoRollback(t *testing.T) {
	db := setupDB(t)

	// 记录转账前的账户余额
	var accountsBefore []account
	if err := db.Order("id").Find(&accountsBefore).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	transferAmount := int64(5000)

	// 使用自动事务执行转账操作
	// 特点：函数返回 error 时，事务会自动回滚，已执行的操作都会被撤销
	err := db.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 扣减转出账户余额（这个操作会成功执行）
		if err := tx.Model(&account{}).
			Where("id = ?", 1).
			Update("balance", gorm.Expr("balance - ?", transferAmount)).Error; err != nil {
			return fmt.Errorf("debit account: %w", err)
		}

		// 重要说明：accountsBefore 是在事务外部查询的，它存储的是查询时的快照值
		// 即使事务内部执行了 UPDATE，这个 Go 变量不会自动更新
		// 如果要在事务内部看到余额变化，需要使用事务的 tx 重新查询数据库
		fmt.Printf("事务外部查询的余额（不会变）: %d\n", accountsBefore[0].Balance)

		// 在事务内部重新查询，可以看到更新后的余额
		var accountInTx account
		if err := tx.First(&accountInTx, 1).Error; err != nil {
			return fmt.Errorf("query account in tx: %w", err)
		}
		fmt.Printf("事务内部查询的余额（已更新）: %d\n", accountInTx.Balance)

		if accountInTx.Balance < 100000 {
			fmt.Println("模拟报错，余额不足")
			return errInsufficientBalance
		}

		// 步骤2: 增加转入账户余额（这个操作也会成功执行）
		if err := tx.Model(&account{}).
			Where("id = ?", 2).
			Update("balance", gorm.Expr("balance + ?", transferAmount)).Error; err != nil {
			return fmt.Errorf("credit account: %w", err)
		}

		// 步骤3: 创建转账记录（模拟这里出错，比如违反唯一约束）
		// 使用一个会失败的 Reference，模拟业务逻辑错误
		record := transferRecord{
			Reference:     "TX-ROLLBACK-001",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        transferAmount,
			Status:        "SUCCESS",
			Message:       "回滚测试",
		}
		if err := tx.Create(&record).Error; err != nil {
			// 假设这里因为某种原因失败了（比如数据库约束、业务规则等）
			// 返回错误后，前面已执行的扣款和加款操作都会被回滚
			return fmt.Errorf("create record failed: %w", err)
		}

		// 步骤4: 模拟后续操作出错（比如业务校验失败）
		// 这里故意返回错误，演示事务回滚的效果
		// 即使前面的扣款和加款操作已经执行，但因为返回了错误，整个事务都会回滚
		return fmt.Errorf("simulated business error: %w", errors.New("business validation failed"))
	})

	// 应该返回错误
	if err == nil {
		t.Fatal("transaction should fail with error")
	}

	// 验证账户余额没有变化（事务已回滚，已执行的扣款和加款操作都被撤销）
	var accountsAfter []account
	if err := db.Order("id").Find(&accountsAfter).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 验证账户1的余额没有变化（扣款操作被回滚）
	if accountsAfter[0].Balance != accountsBefore[0].Balance {
		t.Errorf("account 1 balance should be rolled back, expected %d, got %d",
			accountsBefore[0].Balance, accountsAfter[0].Balance)
	}

	// 验证账户2的余额没有变化（加款操作被回滚）
	if accountsAfter[1].Balance != accountsBefore[1].Balance {
		t.Errorf("account 2 balance should be rolled back, expected %d, got %d",
			accountsBefore[1].Balance, accountsAfter[1].Balance)
	}

	// 验证转账记录没有被创建（因为事务回滚）
	var record transferRecord
	if err := db.Where("reference = ?", "TX-ROLLBACK-001").First(&record).Error; err == nil {
		t.Error("transfer record should not be created after rollback")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("unexpected error checking record: %v", err)
	}

	t.Log("事务回滚成功：即使扣款和加款操作已执行，但因为后续步骤出错，所有操作都被回滚")
}

// ============================================================================
// 知识点 3: 手动事务
// ============================================================================

// TestTransactionManual 测试手动事务
// 演示：手动控制事务的开始、提交和回滚，需要自己处理所有错误情况
func TestTransactionManual(t *testing.T) {
	db := setupDB(t)

	// 手动开始事务
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}

	// 使用 defer 确保在 panic 时回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出 panic
		}
	}()

	// 执行转账操作
	// 步骤1: 扣减转出账户余额
	if err := tx.Model(&account{}).
		Where("id = ?", 1).
		Update("balance", gorm.Expr("balance - ?", 3000)).Error; err != nil {
		tx.Rollback() // 手动回滚
		t.Fatalf("debit account: %v", err)
	}

	// 步骤2: 增加转入账户余额
	if err := tx.Model(&account{}).
		Where("id = ?", 2).
		Update("balance", gorm.Expr("balance + ?", 3000)).Error; err != nil {
		tx.Rollback() // 手动回滚
		t.Fatalf("credit account: %v", err)
	}

	// 步骤3: 创建转账记录
	record := transferRecord{
		Reference:     "TX-003",
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        3000,
		Status:        "SUCCESS",
		Message:       "手动事务测试",
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback() // 手动回滚
		t.Fatalf("create record: %v", err)
	}

	// 手动提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 提交失败时回滚
		t.Fatalf("commit transaction: %v", err)
	}

	// 验证转账结果
	var accounts []account
	if err := db.Order("id").Find(&accounts).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 验证账户余额变化
	if accounts[0].Balance != 97000 { // 100000 - 3000
		t.Errorf("expected account 1 balance 97000, got %d", accounts[0].Balance)
	}
	if accounts[1].Balance != 33000 { // 30000 + 3000
		t.Errorf("expected account 2 balance 33000, got %d", accounts[1].Balance)
	}

	// 验证转账记录已创建
	var createdRecord transferRecord
	if err := db.Where("reference = ?", "TX-003").First(&createdRecord).Error; err != nil {
		t.Errorf("transfer record should be created: %v", err)
	}
}

// ============================================================================
// 知识点 4: SavePoint（保存点）
// ============================================================================

// TestTransactionSavePoint 测试 SavePoint（保存点）
// 演示：在事务中创建检查点，可以回滚到特定点而不回滚整个事务
func TestTransactionSavePoint(t *testing.T) {
	db := setupDB(t)

	// 记录转账前的账户余额
	var accountsBefore []account
	if err := db.Order("id").Find(&accountsBefore).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 使用自动事务，内部使用 SavePoint
	err := db.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 扣减转出账户余额
		if err := tx.Model(&account{}).
			Where("id = ?", 1).
			Update("balance", gorm.Expr("balance - ?", 2000)).Error; err != nil {
			return fmt.Errorf("debit account: %w", err)
		}

		// 步骤2: 创建保存点（SavePoint）
		// SavePoint 允许在事务中创建检查点，可以回滚到特定点而不回滚整个事务
		if err := tx.SavePoint("after_debit").Error; err != nil {
			return fmt.Errorf("savepoint: %w", err)
		}

		// 步骤3: 尝试增加转入账户余额（模拟可能失败的操作）
		// 这里故意使用一个会失败的操作来演示回滚到保存点
		if err := tx.Model(&account{}).
			Where("id = ?", 999). // 不存在的账户ID，会失败
			Update("balance", gorm.Expr("balance + ?", 2000)).Error; err != nil {
			// 回滚到保存点（只回滚加款操作，扣款操作保留）
			if rollbackErr := tx.RollbackTo("after_debit").Error; rollbackErr != nil {
				return fmt.Errorf("rollback to savepoint: %w", rollbackErr)
			}
			// 注意：在实际业务中，转账操作通常要么全部成功要么全部失败
			// 这里仅演示 SavePoint 的用法
			t.Log("Rolled back to savepoint: after_debit")
			// 继续执行，不返回错误（演示部分回滚的效果）
		}

		// 步骤4: 创建转账记录（即使加款失败，记录也会创建）
		record := transferRecord{
			Reference:     "TX-SAVEPOINT-001",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        2000,
			Status:        "PARTIAL", // 部分成功
			Message:       "SavePoint 测试",
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("create record: %w", err)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("transaction failed: %v", err)
	}

	// 验证转账结果
	var accountsAfter []account
	if err := db.Order("id").Find(&accountsAfter).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 验证账户1的余额已扣减（扣款操作成功）
	if accountsAfter[0].Balance != accountsBefore[0].Balance-2000 {
		t.Errorf("account 1 balance should decrease by 2000, expected %d, got %d",
			accountsBefore[0].Balance-2000, accountsAfter[0].Balance)
	}

	// 验证账户2的余额没有变化（加款操作被回滚到保存点）
	if accountsAfter[1].Balance != accountsBefore[1].Balance {
		t.Errorf("account 2 balance should not change after rollback to savepoint, expected %d, got %d",
			accountsBefore[1].Balance, accountsAfter[1].Balance)
	}

	// 验证转账记录已创建
	var record transferRecord
	if err := db.Where("reference = ?", "TX-SAVEPOINT-001").First(&record).Error; err != nil {
		t.Errorf("transfer record should be created: %v", err)
	}
}

// ============================================================================
// 知识点 5: 嵌套事务
// ============================================================================

// TestTransactionNested 测试嵌套事务（成功场景）
// 演示：GORM 支持嵌套事务，内层事务实际上会使用 SavePoint 实现
//
// 为什么使用 SavePoint？
// - 数据库本身不支持真正的嵌套事务（大多数数据库只支持单个事务）
// - GORM 通过 SavePoint 机制模拟嵌套事务的行为
// - 提供统一的 Transaction() API，无论是否在事务中都可以使用
func TestTransactionNested(t *testing.T) {
	db := setupDB(t)

	// 执行嵌套事务（成功场景）
	err := db.Transaction(func(tx1 *gorm.DB) error {
		// 外层事务：创建第一个转账记录
		outerRecord := transferRecord{
			Reference:     "TX-NESTED-001",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        1000,
			Status:        "PENDING",
			Message:       "嵌套事务测试-外层",
		}
		if err := tx1.Create(&outerRecord).Error; err != nil {
			return fmt.Errorf("create outer record: %w", err)
		}

		// 内层事务：创建第二个转账记录
		// 内层事务实际上会使用 SavePoint 实现，原因如下：
		// 1. 数据库限制：大多数数据库（MySQL、PostgreSQL、SQLite等）不支持真正的嵌套事务
		//    它们只支持单个事务，不能在一个事务内部启动另一个独立的事务
		// 2. SavePoint 模拟：GORM 使用 SavePoint（保存点）来模拟嵌套事务的行为：
		//    - 当在已存在的事务中调用 Transaction() 时，GORM 会自动创建一个 SavePoint
		//    - 如果内层事务成功（返回 nil），GORM 会释放 SavePoint（相当于提交内层事务）
		//    - 如果内层事务失败（返回错误），GORM 会回滚到 SavePoint（相当于回滚内层事务）
		// 3. 行为一致性：这样设计提供了统一的 API，无论是否在事务中，都可以使用 Transaction() 方法
		// 4. 默认行为：内层事务失败时，GORM 默认会让外层事务也回滚（可以通过配置改变）
		return tx1.Transaction(func(tx2 *gorm.DB) error {
			innerRecord := transferRecord{
				Reference:     "TX-NESTED-002",
				FromAccountID: 2,
				ToAccountID:   1,
				Amount:        500,
				Status:        "PENDING",
				Message:       "嵌套事务测试-内层",
			}
			if err := tx2.Create(&innerRecord).Error; err != nil {
				// 内层事务返回错误，会导致外层事务也回滚（GORM 默认行为）
				return fmt.Errorf("create inner record: %w", err)
			}

			// 返回 nil，内层事务提交
			return nil
		})
		// 如果内层事务成功，外层事务也会提交
		// 如果内层事务失败，外层事务会回滚（因为内层返回了错误）
	})

	if err != nil {
		t.Fatalf("nested transaction failed: %v", err)
	}

	// 验证外层事务的记录已创建
	var outerRecord transferRecord
	if err := db.Where("reference = ?", "TX-NESTED-001").First(&outerRecord).Error; err != nil {
		t.Errorf("outer transaction record should be created: %v", err)
	}

	// 验证内层事务的记录已创建
	var innerRecord transferRecord
	if err := db.Where("reference = ?", "TX-NESTED-002").First(&innerRecord).Error; err != nil {
		t.Errorf("inner transaction record should be created: %v", err)
	}

	t.Log("嵌套事务成功：外层和内层事务都成功提交")
}

// TestTransactionNestedWithRollback 测试嵌套事务的回滚行为
// 演示：当内层事务失败时，外层事务也会回滚（GORM 的默认行为）
func TestTransactionNestedWithRollback(t *testing.T) {
	db := setupDB(t)

	// 记录操作前的转账记录数量
	var countBefore int64
	db.Model(&transferRecord{}).Count(&countBefore)

	// 执行一个会失败的嵌套事务（内层事务会失败）
	err := db.Transaction(func(tx1 *gorm.DB) error {
		// 外层事务：创建第一个记录
		outerRecord := transferRecord{
			Reference:     "TX-NESTED-ROLLBACK-001",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        1000,
			Status:        "PENDING",
			Message:       "嵌套事务测试-外层",
		}
		if err := tx1.Create(&outerRecord).Error; err != nil {
			return fmt.Errorf("create outer record: %w", err)
		}

		// 内层事务：尝试创建一个会失败的记录（使用已存在的 Reference）
		return tx1.Transaction(func(tx2 *gorm.DB) error {
			// 先创建一个记录
			innerRecord1 := transferRecord{
				Reference:     "TX-NESTED-ROLLBACK-002",
				FromAccountID: 2,
				ToAccountID:   1,
				Amount:        500,
				Status:        "PENDING",
				Message:       "嵌套事务测试-内层1",
			}
			if err := tx2.Create(&innerRecord1).Error; err != nil {
				return fmt.Errorf("create inner record 1: %w", err)
			}

			// 尝试创建一个会失败的记录（重复的 Reference，违反唯一约束）
			innerRecord2 := transferRecord{
				Reference:     "TX-NESTED-ROLLBACK-002", // 重复的 Reference，会失败
				FromAccountID: 2,
				ToAccountID:   1,
				Amount:        300,
				Status:        "PENDING",
				Message:       "嵌套事务测试-内层2（会失败）",
			}
			if err := tx2.Create(&innerRecord2).Error; err != nil {
				// 内层事务返回错误，会导致外层事务也回滚
				return fmt.Errorf("create inner record 2 (should fail): %w", err)
			}

			return nil
		})
	})

	// 应该返回错误（内层事务失败）
	if err == nil {
		t.Fatal("nested transaction should fail")
	}

	// 验证所有记录都没有创建（因为外层事务也回滚了）
	var countAfter int64
	db.Model(&transferRecord{}).Count(&countAfter)

	if countAfter != countBefore {
		t.Errorf("no records should be created after nested transaction rollback, expected %d, got %d",
			countBefore, countAfter)
	}

	t.Log("嵌套事务回滚：内层事务失败导致外层事务也回滚（GORM 默认行为）")
}

// ============================================================================
// 知识点 6: 幂等性设计
// ============================================================================

// TestTransactionIdempotency 测试事务的幂等性设计
// 演示：通过 Reference 字段防止重复操作，即使多次调用相同的转账请求，也只会执行一次
func TestTransactionIdempotency(t *testing.T) {
	db := setupDB(t)

	// 第一次转账，应该成功
	err := db.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 幂等性检查
		// 检查是否已存在相同的转账记录，防止重复操作
		var exists transferRecord
		if err := tx.Where("reference = ?", "TX-IDEMPOTENT-001").Take(&exists).Error; err == nil {
			// 已存在相同 Reference 的记录，返回错误（事务会自动回滚）
			return errDuplicateTransfer
		}

		// 步骤2: 执行转账操作
		if err := tx.Model(&account{}).
			Where("id = ?", 1).
			Update("balance", gorm.Expr("balance - ?", 5000)).Error; err != nil {
			return fmt.Errorf("debit account: %w", err)
		}

		if err := tx.Model(&account{}).
			Where("id = ?", 2).
			Update("balance", gorm.Expr("balance + ?", 5000)).Error; err != nil {
			return fmt.Errorf("credit account: %w", err)
		}

		// 步骤3: 创建转账记录
		record := transferRecord{
			Reference:     "TX-IDEMPOTENT-001",
			FromAccountID: 1,
			ToAccountID:   2,
			Amount:        5000,
			Status:        "SUCCESS",
			Message:       "第一次转账",
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("create record: %w", err)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("first transfer should succeed: %v", err)
	}

	// 记录第一次转账后的余额
	var accountsAfterFirst []account
	if err := db.Order("id").Find(&accountsAfterFirst).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 使用相同的 Reference 再次转账，应该被拒绝（幂等性保护）
	err = db.Transaction(func(tx *gorm.DB) error {
		// 幂等性检查：发现已存在相同的 Reference
		var exists transferRecord
		if err := tx.Where("reference = ?", "TX-IDEMPOTENT-001").Take(&exists).Error; err == nil {
			// 已存在，返回错误（事务会自动回滚）
			return errDuplicateTransfer
		}

		// 即使金额不同，也应该被拒绝（因为 Reference 相同）
		if err := tx.Model(&account{}).
			Where("id = ?", 1).
			Update("balance", gorm.Expr("balance - ?", 100)).Error; err != nil {
			return fmt.Errorf("debit account: %w", err)
		}

		return nil
	})

	// 应该返回重复转账的错误
	if !errors.Is(err, errDuplicateTransfer) {
		t.Fatalf("expected duplicate transfer error, got %v", err)
	}

	// 验证账户余额没有再次变化（第二次转账被拒绝）
	var accountsAfterSecond []account
	if err := db.Order("id").Find(&accountsAfterSecond).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	if accountsAfterSecond[0].Balance != accountsAfterFirst[0].Balance {
		t.Errorf("account balance should not change after duplicate transfer rejected, expected %d, got %d",
			accountsAfterFirst[0].Balance, accountsAfterSecond[0].Balance)
	}
}

// ============================================================================
// 知识点 7: 悲观锁（SELECT FOR UPDATE）
// ============================================================================

// TestTransactionPessimisticLocking 测试悲观锁
// 演示：使用 SELECT ... FOR UPDATE 锁定账户记录，防止并发修改
//
// ⚠️ 为什么不建议使用悲观锁？
//  1. 死锁风险：当多个事务以不同顺序锁定资源时，容易产生死锁
//     例如：事务A锁定账户1后尝试锁定账户2，事务B锁定账户2后尝试锁定账户1
//  2. 性能问题：
//     - 阻塞其他事务：被锁定的记录会阻塞其他需要修改该记录的事务
//     - 锁持有时间长：锁会持续到事务结束，如果事务中有慢操作，锁持有时间会更长
//     - 并发性能差：高并发场景下，大量事务会排队等待锁释放
//  3. 资源浪费：即使事务最终可能失败，锁也会一直持有直到事务结束
//  4. 扩展性差：随着并发量增加，性能会急剧下降
//
// 💡 建议：
// - 优先使用乐观锁（版本号机制），适合读多写少场景
// - 如果必须使用悲观锁，确保：
//   - 锁定顺序一致（避免死锁）
//   - 事务尽可能短（减少锁持有时间）
//   - 只锁定必要的记录（避免锁范围过大）
//   - 考虑使用超时机制（避免长时间等待）
func TestTransactionPessimisticLocking(t *testing.T) {
	db := setupDB(t)

	// 使用自动事务，内部使用悲观锁
	err := db.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 使用悲观锁查询转出账户
		// clause.Locking{Strength: "UPDATE"} 相当于 SQL 的 SELECT ... FOR UPDATE
		// 这会锁定查询到的记录（行锁），防止其他事务同时修改，直到事务结束
		// ⚠️ 注意：SELECT FOR UPDATE 是行锁，不是表锁，但如果锁定的行很多，影响范围也会很大
		// ⚠️ 死锁风险：如果多个事务以不同顺序锁定账户，可能产生死锁
		//    例如：事务A先锁账户1再锁账户2，事务B先锁账户2再锁账户1
		var from account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&from, 1).Error; err != nil {
			return fmt.Errorf("fetch from account: %w", err)
		}

		// 步骤2: 余额校验
		if from.Balance < 5000 {
			return errInsufficientBalance
		}

		// 步骤3: 扣减转出账户余额
		// 由于使用了悲观锁，其他尝试修改这个账户的事务会被阻塞，直到当前事务结束
		// ⚠️ 性能影响：如果有多个并发转账操作涉及同一个账户，它们会串行执行，严重影响性能
		if err := tx.Model(&account{}).
			Where("id = ?", from.ID).
			Update("balance", gorm.Expr("balance - ?", 5000)).Error; err != nil {
			return fmt.Errorf("debit account: %w", err)
		}

		// 步骤4: 使用悲观锁查询转入账户
		var to account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&to, 2).Error; err != nil {
			return fmt.Errorf("fetch to account: %w", err)
		}

		// 步骤5: 增加转入账户余额
		if err := tx.Model(&account{}).
			Where("id = ?", to.ID).
			Update("balance", gorm.Expr("balance + ?", 5000)).Error; err != nil {
			return fmt.Errorf("credit account: %w", err)
		}

		// 步骤6: 创建转账记录
		record := transferRecord{
			Reference:     "TX-LOCK-001",
			FromAccountID: from.ID,
			ToAccountID:   to.ID,
			Amount:        5000,
			Status:        "SUCCESS",
			Message:       "悲观锁测试",
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("create record: %w", err)
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		t.Fatalf("transaction failed: %v", err)
	}

	// 验证转账结果
	var accounts []account
	if err := db.Order("id").Find(&accounts).Error; err != nil {
		t.Fatalf("list accounts: %v", err)
	}

	// 验证账户余额变化
	if accounts[0].Balance != 95000 { // 100000 - 5000
		t.Errorf("expected account 1 balance 95000, got %d", accounts[0].Balance)
	}
	if accounts[1].Balance != 35000 { // 30000 + 5000
		t.Errorf("expected account 2 balance 35000, got %d", accounts[1].Balance)
	}

	// 验证转账记录已创建
	var record transferRecord
	if err := db.Where("reference = ?", "TX-LOCK-001").First(&record).Error; err != nil {
		t.Errorf("transfer record should be created: %v", err)
	}
}

// ============================================================================
// 辅助函数
// ============================================================================

// resetAccounts 重置账户数据，用于测试前的数据准备
// 删除所有转账记录和账户，然后创建初始测试账户
func resetAccounts(t *testing.T, db *gorm.DB) error {
	t.Helper()

	// 删除所有转账记录
	if err := db.Exec("DELETE FROM transfer_records").Error; err != nil {
		return err
	}

	// 删除所有账户
	if err := db.Exec("DELETE FROM accounts").Error; err != nil {
		return err
	}

	// 重置 SQLite 的 AUTOINCREMENT 序列（确保 ID 从 1 开始）
	if err := db.Exec("DELETE FROM sqlite_sequence WHERE name='accounts'").Error; err != nil {
		// 忽略错误，因为表可能还没有序列记录
		_ = err
	}

	// 创建初始测试账户，显式设置 ID 为 1 和 2，确保每次测试都使用相同的 ID
	accounts := []account{
		{ID: 1, Name: "Alice Corp", Balance: 100000}, // Alice 公司，初始余额 100000，ID=1
		{ID: 2, Name: "Bob Studio", Balance: 30000},  // Bob 工作室，初始余额 30000，ID=2
	}
	return db.Create(&accounts).Error
}
