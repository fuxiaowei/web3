package advanced

import (
	"testing"

	"lesson02examples/testutil"

	"gorm.io/gorm"
)

// TestAssociationsAPI demonstrates Association API operations for Many-to-Many relationships
// This test covers:
// 1. Count: Count associations
// 2. Find: Find associations
// 3. Replace: Replace all associations with new ones
// 4. Delete: Remove specific associations
// 5. Clear: Remove all associations
func TestAssociationsAPI(t *testing.T) {
	db := testutil.NewTestDB(t, "associations-api.db")

	// AutoMigrate creates all tables and their relationships
	if err := db.AutoMigrate(&user{}, &profile{}, &product{}, &order{}, &orderItem{}, &role{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Clean up existing data
	db.Exec("PRAGMA foreign_keys = OFF")
	db.Exec("DELETE FROM user_roles")
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM profiles")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM roles")
	db.Exec("PRAGMA foreign_keys = ON")

	// Create roles
	roles := []role{
		{Name: "admin", Description: "系统管理员，拥有所有权限"},
		{Name: "editor", Description: "编辑者，可以创建和编辑内容"},
		{Name: "viewer", Description: "查看者，只能查看内容"},
	}
	if err := db.Create(&roles).Error; err != nil {
		t.Fatalf("seed roles: %v", err)
	}

	// Create users
	alice := user{
		Name:  "Alice",
		Email: "alice@example.com",
		Profile: profile{
			Nickname: "阿狸",
			Phone:    "13800000000",
			Address:  "上海市徐汇区漕河泾开发区",
		},
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&alice).Error; err != nil {
		t.Fatalf("create alice: %v", err)
	}

	charlie := user{
		Name:  "Charlie",
		Email: "charlie@example.com",
	}
	if err := db.Create(&charlie).Error; err != nil {
		t.Fatalf("create charlie: %v", err)
	}

	// Associate roles with Alice
	if err := db.Model(&alice).Association("Roles").Append(&roles[0], &roles[1]); err != nil {
		t.Fatalf("append roles: %v", err)
	}

	// COUNT ASSOCIATIONS: How many roles does a user have?
	roleCount := db.Model(&alice).Association("Roles").Count()
	if roleCount != 2 {
		t.Fatalf("expected alice to have 2 roles, got %d", roleCount)
	}

	// FIND ASSOCIATIONS: Get all roles for a user
	var aliceRoles []role
	if err := db.Model(&alice).Association("Roles").Find(&aliceRoles); err != nil {
		t.Fatalf("find alice roles: %v", err)
	}
	if len(aliceRoles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(aliceRoles))
	}

	// REPLACE ASSOCIATIONS: Replace all roles with new ones
	// This removes existing associations and adds new ones
	if err := db.Model(&charlie).Association("Roles").Replace(&roles[2]); err != nil {
		t.Fatalf("replace charlie roles: %v", err)
	}
	// Verify Charlie now has only viewer role
	if db.Model(&charlie).Association("Roles").Count() != 1 {
		t.Fatalf("expected charlie to have 1 role")
	}

	// DELETE SPECIFIC ASSOCIATIONS: Remove one role from user
	// This removes the association but doesn't delete the role itself
	if err := db.Model(&alice).Association("Roles").Delete(&roles[1]); err != nil {
		t.Fatalf("delete alice editor role: %v", err)
	}
	// Verify Alice now has only admin role
	if db.Model(&alice).Association("Roles").Count() != 1 {
		t.Fatalf("expected alice to have 1 role after deletion")
	}

	// CLEAR ALL ASSOCIATIONS: Remove all roles from a user
	if err := db.Model(&charlie).Association("Roles").Clear(); err != nil {
		t.Fatalf("clear charlie roles: %v", err)
	}
	// Verify Charlie has no roles
	if db.Model(&charlie).Association("Roles").Count() != 0 {
		t.Fatalf("expected charlie to have 0 roles after clear")
	}
}
