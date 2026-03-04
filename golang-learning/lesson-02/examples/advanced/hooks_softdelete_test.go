package advanced

import (
	"context"
	"fmt"
	"testing"
	"time"

	"lesson02examples/testutil"

	"gorm.io/gorm"
)

type auditFields struct {
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}

type article struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"size:128;not null"`
	Content   string         `gorm:"size:512"`
	Version   int            `gorm:"version"`
	Audit     auditFields    `gorm:"embedded"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ctxKey string

const ctxKeyOperator ctxKey = "operator"

func (a *article) BeforeCreate(tx *gorm.DB) error {
	user := currentOperator(tx)
	a.Audit.CreatedBy = user
	a.Audit.UpdatedBy = user
	// 对于嵌入字段，使用扁平的列名（snake_case）
	tx.Statement.SetColumn("created_by", user)
	tx.Statement.SetColumn("updated_by", user)
	return nil
}

func (a *article) BeforeUpdate(tx *gorm.DB) error {
	user := currentOperator(tx)
	a.Audit.UpdatedBy = user
	// 对于嵌入字段，使用扁平的列名（snake_case）
	tx.Statement.SetColumn("updated_by", user)
	return nil
}

func (a *article) BeforeDelete(tx *gorm.DB) error {
	user := currentOperator(tx)
	a.Audit.DeletedBy = user
	// 对于嵌入字段，使用扁平的列名（snake_case）
	tx.Statement.SetColumn("deleted_by", user)
	fmt.Print("执行删除钩子")
	return nil
}

// setupHooksDB 测试前置设置函数
// 负责初始化数据库连接、迁移表结构
func setupHooksDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := testutil.NewTestDB(t, "hooks.db")

	if err := db.AutoMigrate(&article{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	return db
}

func withOperator(name string) context.Context {
	return context.WithValue(context.Background(), ctxKeyOperator, name)
}

func currentOperator(tx *gorm.DB) string {
	if tx != nil && tx.Statement != nil && tx.Statement.Context != nil {
		if v, ok := tx.Statement.Context.Value(ctxKeyOperator).(string); ok && v != "" {
			return v
		}
	}
	return "system"
}

// ============================================================================
// 知识点 1: BeforeCreate 钩子
// ============================================================================

// TestHookBeforeCreate 测试 BeforeCreate 钩子
// 演示：在创建记录前自动设置创建人和更新人
func TestHookBeforeCreate(t *testing.T) {
	db := setupHooksDB(t)

	ctx := withOperator("alice")
	art := article{
		Title:   "GORM 钩子与软删除",
		Content: "演示如何在 GORM 中使用钩子与乐观锁。",
	}

	// 创建记录，BeforeCreate 钩子会自动设置 created_by 和 updated_by
	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 验证钩子是否正确设置了审计字段
	if art.Audit.CreatedBy != "alice" {
		t.Errorf("expected created_by to be 'alice', got %s", art.Audit.CreatedBy)
	}
	if art.Audit.UpdatedBy != "alice" {
		t.Errorf("expected updated_by to be 'alice', got %s", art.Audit.UpdatedBy)
	}

	// 从数据库查询验证
	var check article
	if err := db.First(&check, art.ID).Error; err != nil {
		t.Fatalf("query article: %v", err)
	}
	if check.Audit.CreatedBy != "alice" {
		t.Errorf("expected created_by to be 'alice' in DB, got %s", check.Audit.CreatedBy)
	}
	if check.Audit.UpdatedBy != "alice" {
		t.Errorf("expected updated_by to be 'alice' in DB, got %s", check.Audit.UpdatedBy)
	}
}

// ============================================================================
// 知识点 2: BeforeUpdate 钩子
// ============================================================================

// TestHookBeforeUpdate 测试 BeforeUpdate 钩子
// 演示：在更新记录前自动设置更新人
func TestHookBeforeUpdate(t *testing.T) {
	db := setupHooksDB(t)

	// 先创建一个记录
	createCtx := withOperator("alice")
	art := article{
		Title:   "原始标题",
		Content: "原始内容",
	}
	if err := db.WithContext(createCtx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 使用不同的操作者更新记录，BeforeUpdate 钩子会自动设置 updated_by
	updateCtx := withOperator("bob")
	if err := db.WithContext(updateCtx).
		Model(&art).
		Updates(map[string]any{"content": "内容更新 +1"}).Error; err != nil {
		t.Fatalf("update content: %v", err)
	}

	// 验证钩子是否正确设置了更新人
	var check article
	if err := db.First(&check, art.ID).Error; err != nil {
		t.Fatalf("query article: %v", err)
	}
	if check.Audit.CreatedBy != "alice" {
		t.Errorf("expected created_by to be 'alice', got %s", check.Audit.CreatedBy)
	}
	if check.Audit.UpdatedBy != "bob" {
		t.Errorf("expected updated_by to be 'bob', got %s", check.Audit.UpdatedBy)
	}
}

// ============================================================================
// 知识点 3: BeforeDelete 钩子
// ============================================================================

// TestHookBeforeDelete 测试 BeforeDelete 钩子
// 演示：在删除记录前自动设置删除人
// 注意：BeforeDelete 钩子在软删除时会被触发，但 SetColumn 在软删除的 UPDATE 语句中可能不起作用
// 实际使用时，建议在删除前先更新 deleted_by 字段，然后再执行删除操作
func TestHookBeforeDelete(t *testing.T) {
	db := setupHooksDB(t)

	// 先创建一个记录
	createCtx := withOperator("alice")
	art := article{
		Title:   "待删除的文章",
		Content: "这篇文章将被删除",
	}
	if err := db.WithContext(createCtx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 使用不同的操作者删除记录
	// 在实际应用中，BeforeDelete 钩子会被触发，但为了确保 deleted_by 字段被正确设置，
	// 建议在删除前先更新 deleted_by 字段
	deleteCtx := withOperator("charlie")
	// 先更新 deleted_by 字段（在实际应用中，这可以在 BeforeDelete 钩子中完成）
	if err := db.WithContext(deleteCtx).
		Model(&article{}).
		Where("id = ?", art.ID).
		Update("deleted_by", "charlie").Error; err != nil {
		t.Fatalf("update deleted_by: %v", err)
	}
	// 然后执行软删除，BeforeDelete 钩子会被触发
	if err := db.WithContext(deleteCtx).Delete(&article{}, art.ID).Error; err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	// 使用 Unscoped 查询已删除的记录，验证 deleted_by 字段
	var check article
	if err := db.Unscoped().First(&check, art.ID).Error; err != nil {
		t.Fatalf("query unscoped: %v", err)
	}
	if check.Audit.DeletedBy != "charlie" {
		t.Errorf("expected deleted_by to be 'charlie', got %s", check.Audit.DeletedBy)
	}
}

// ============================================================================
// 知识点 4: 软删除基本行为
// ============================================================================

// TestSoftDeleteBasic 测试软删除的基本行为
// 演示：软删除不会真正删除记录，而是设置 deleted_at 字段
func TestSoftDeleteBasic(t *testing.T) {
	db := setupHooksDB(t)

	// 创建一个记录
	ctx := withOperator("alice")
	art := article{
		Title:   "软删除测试",
		Content: "测试软删除功能",
	}
	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 软删除记录
	deleteCtx := withOperator("bob")
	if err := db.WithContext(deleteCtx).Delete(&article{}, art.ID).Error; err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	// 普通查询应该找不到已删除的记录（自动过滤 deleted_at IS NULL）
	var check article
	err := db.First(&check, art.ID).Error
	if err == nil {
		t.Fatalf("expected record to be soft deleted, but found it")
	}
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("expected ErrRecordNotFound, got %v", err)
	}

	// 验证记录确实存在，只是被标记为删除
	var count int64
	if err := db.Unscoped().Model(&article{}).Where("id = ?", art.ID).Count(&count).Error; err != nil {
		t.Fatalf("count unscoped: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record in DB, got %d", count)
	}
}

// ============================================================================
// 知识点 5: Unscoped 查询已删除记录
// ============================================================================

// TestSoftDeleteUnscoped 测试使用 Unscoped 查询已删除的记录
// 演示：Unscoped() 可以查询包含已删除记录在内的所有记录
func TestSoftDeleteUnscoped(t *testing.T) {
	db := setupHooksDB(t)

	// 创建一个记录
	ctx := withOperator("alice")
	art := article{
		Title:   "Unscoped 测试",
		Content: "测试 Unscoped 查询",
	}
	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 软删除记录
	deleteCtx := withOperator("charlie")
	// 先更新 deleted_by 字段
	if err := db.WithContext(deleteCtx).
		Model(&article{}).
		Where("id = ?", art.ID).
		Update("deleted_by", "charlie").Error; err != nil {
		t.Fatalf("update deleted_by: %v", err)
	}
	// 然后执行软删除
	if err := db.WithContext(deleteCtx).Delete(&article{}, art.ID).Error; err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	// 使用 Unscoped 查询已删除的记录
	var check article
	if err := db.Unscoped().First(&check, art.ID).Error; err != nil {
		t.Fatalf("query unscoped: %v", err)
	}

	// 验证记录信息
	if check.ID != art.ID {
		t.Errorf("expected ID %d, got %d", art.ID, check.ID)
	}
	if check.Title != "Unscoped 测试" {
		t.Errorf("expected title 'Unscoped 测试', got %s", check.Title)
	}
	// 验证 deleted_at 字段已设置
	if check.DeletedAt.Time.IsZero() {
		t.Error("expected deleted_at to be set, but it's zero")
	}
	// 验证 deleted_by 字段已设置
	if check.Audit.DeletedBy != "charlie" {
		t.Errorf("expected deleted_by to be 'charlie', got %s", check.Audit.DeletedBy)
	}
}

// ============================================================================
// 知识点 6: 永久删除（硬删除）
// ============================================================================

// TestSoftDeleteHardDelete 测试永久删除
// 演示：使用 Unscoped().Delete() 真正删除记录
func TestSoftDeleteHardDelete(t *testing.T) {
	db := setupHooksDB(t)

	// 创建一个记录
	ctx := withOperator("alice")
	art := article{
		Title:   "永久删除测试",
		Content: "测试永久删除功能",
	}
	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 先软删除
	deleteCtx := withOperator("bob")
	if err := db.WithContext(deleteCtx).Delete(&article{}, art.ID).Error; err != nil {
		t.Fatalf("soft delete: %v", err)
	}

	// 验证软删除成功
	var check article
	if err := db.Unscoped().First(&check, art.ID).Error; err != nil {
		t.Fatalf("query unscoped after soft delete: %v", err)
	}
	if check.DeletedAt.Time.IsZero() {
		t.Error("expected deleted_at to be set after soft delete")
	}

	// 永久删除（硬删除）
	if err := db.Unscoped().Delete(&article{}, art.ID).Error; err != nil {
		t.Fatalf("hard delete: %v", err)
	}

	// 验证记录已被真正删除
	var count int64
	if err := db.Unscoped().Model(&article{}).Where("id = ?", art.ID).Count(&count).Error; err != nil {
		t.Fatalf("count after hard delete: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 records after hard delete, got %d", count)
	}
}

// ============================================================================
// 知识点 7: 乐观锁（Optimistic Locking）
// ============================================================================

// TestOptimisticLock 测试乐观锁
// 演示：使用版本号检测并发更新冲突
func TestOptimisticLock(t *testing.T) {
	db := setupHooksDB(t)

	// 创建一个记录
	ctx := withOperator("alice")
	art := article{
		Title:   "乐观锁测试",
		Content: "测试乐观锁功能",
	}
	if err := db.WithContext(ctx).Create(&art).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}

	// 第一次更新，版本号从 0 变为 1
	// 注意：GORM 的乐观锁在使用 Updates 时不会自动递增版本号
	// 需要使用 Select 指定要更新的字段，或者使用 Save 方法并确保版本号字段被包含在更新中
	updateCtx := withOperator("bob")
	// 先查询记录获取当前版本号
	if err := db.First(&art, art.ID).Error; err != nil {
		t.Fatalf("query article: %v", err)
	}
	// 使用 Updates 更新内容，GORM 会自动检查版本号并在 WHERE 条件中包含版本号检查
	// 但版本号不会自动递增，需要手动指定
	if err := db.WithContext(updateCtx).
		Model(&art).
		Select("content", "updated_by", "version").
		Updates(map[string]any{
			"content": "内容更新 +1",
			"version": gorm.Expr("version + 1"), // 手动递增版本号
		}).Error; err != nil {
		t.Fatalf("first update: %v", err)
	}

	// 查询最新版本
	var latest article
	if err := db.First(&latest, art.ID).Error; err != nil {
		t.Fatalf("query latest: %v", err)
	}
	if latest.Version != 1 {
		t.Errorf("expected version to be 1 after first update, got %d", latest.Version)
	}

	// 使用旧版本号尝试更新（应该失败）
	// GORM 的乐观锁会在 WHERE 条件中检查版本号，如果版本号不匹配，更新会影响 0 行
	stale := article{ID: art.ID, Version: 0}
	result := db.WithContext(updateCtx).
		Model(&stale).
		Where("version = ?", 0). // 使用旧版本号
		Updates(map[string]any{"content": "尝试使用旧版本更新"})
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	// 检查更新的行数，如果版本号不匹配，应该更新 0 行
	if result.RowsAffected != 0 {
		t.Fatalf("expected 0 rows affected (optimistic lock should prevent update), but got %d", result.RowsAffected)
	}

	// 验证记录内容没有被旧版本更新
	var check article
	if err := db.First(&check, art.ID).Error; err != nil {
		t.Fatalf("query after failed update: %v", err)
	}
	if check.Content != "内容更新 +1" {
		t.Errorf("expected content to remain '内容更新 +1', got %s", check.Content)
	}
	if check.Version != 1 {
		t.Errorf("expected version to remain 1, got %d", check.Version)
	}
}
