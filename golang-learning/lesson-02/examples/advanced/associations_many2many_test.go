package advanced

import (
	"fmt"
	"testing"

	"lesson02examples/testutil"

	"gorm.io/gorm"
)

// TestAssociationsManyToMany demonstrates Many-to-Many relationships in GORM
// This test covers:
// 1. Creating Many-to-Many relationships
// 2. Associating roles with users using Association API
// 3. Associating roles when creating users
// 4. Preloading Many-to-Many associations
func TestAssociationsManyToMany(t *testing.T) {
	db := testutil.NewTestDB(t, "associations-many2many.db")

	// AutoMigrate creates all tables and their relationships
	// For Many-to-Many: GORM automatically creates the join table (user_roles)
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

	// MANY TO MANY: Create roles and associate them with users
	// Many-to-many relationships require a join table (automatically created by GORM)
	// The join table "user_roles" contains user_id and role_id columns
	//
	// Steps for Many-to-Many:
	// 1. Create the roles first (independent entities)
	// 2. Create users (or use existing users)
	// 3. Associate roles with users using Association API or by setting the Roles field
	roles := []role{
		{Name: "admin", Description: "系统管理员，拥有所有权限"},
		{Name: "editor", Description: "编辑者，可以创建和编辑内容"},
		{Name: "viewer", Description: "查看者，只能查看内容"},
	}
	// if err := db.Create(&roles).Error; err != nil {
	// 	t.Fatalf("seed roles: %v", err)
	// }

	// Create a user first
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

	// METHOD 1: Associate roles using Association API
	// Association API provides fine-grained control over many-to-many relationships
	// This is useful when you want to add/remove associations without loading the full user record
	//
	// Available methods:
	// - Append: Add associations (doesn't remove existing ones)
	// - Replace: Replace all associations with new ones
	// - Delete: Remove associations
	// - Clear: Remove all associations
	// - Count: Count associations
	// - Find: Find associations

	// Append roles to user: Alice gets admin and editor roles
	// This adds entries to the user_roles join table
	if err := db.Model(&alice).Association("Roles").Append(&roles[0], &roles[1]); err != nil {
		t.Fatalf("append roles: %v", err)
	}

	// Create another user and assign roles directly
	// METHOD 2: Set roles when creating user (roles must exist first)
	bob := user{
		Name:  "Bob",
		Email: "bob@example.com",
		Profile: profile{
			Nickname: "鲍勃",
			Phone:    "13900000000",
			Address:  "北京市海淀区中关村",
		},
		Roles: []role{roles[1], roles[2]}, // Bob gets editor and viewer roles
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&bob).Error; err != nil {
		t.Fatalf("create bob: %v", err)
	}

	// PRELOAD MANY TO MANY: Load users with their roles
	// Preloading many-to-many works the same way as other associations
	var usersWithRoles []user
	// Preload Roles: Load all roles for each user
	// GORM will execute:
	// 1. SELECT * FROM users
	// 2. SELECT * FROM roles INNER JOIN user_roles ON roles.id = user_roles.role_id WHERE user_roles.user_id IN (...)
	if err := db.Preload("Roles").Find(&usersWithRoles).Error; err != nil {
		t.Fatalf("preload users with roles: %v", err)
	}
	fmt.Println(usersWithRoles)

	// Verify Alice has 2 roles (admin, editor)
	var aliceWithRoles user
	if err := db.Preload("Roles").First(&aliceWithRoles, "email = ?", "alice@example.com").Error; err != nil {
		t.Fatalf("preload alice roles: %v", err)
	}
	if len(aliceWithRoles.Roles) != 2 {
		t.Fatalf("expected alice to have 2 roles, got %d", len(aliceWithRoles.Roles))
	}

	// CONDITIONAL PRELOAD MANY TO MANY: Load only specific roles
	// Example: Load only admin roles for users
	var adminUsers []user
	// Preload Roles with condition: Only load roles where name = 'admin'
	if err := db.Preload("Roles", "name = ?", "admin").Find(&adminUsers).Error; err != nil {
		t.Fatalf("preload admin roles: %v", err)
	}

	// REVERSE PRELOAD: Load roles with their users
	// You can also preload the reverse side of many-to-many
	var rolesWithUsers []role
	// Preload Users: Load all users that have each role
	// This demonstrates bidirectional many-to-many preloading
	if err := db.Preload("Users").Find(&rolesWithUsers).Error; err != nil {
		t.Fatalf("preload roles with users: %v", err)
	}
	// Verify admin role has at least one user (Alice)
	var adminRole role
	if err := db.Preload("Users").First(&adminRole, "name = ?", "admin").Error; err != nil {
		t.Fatalf("preload admin role users: %v", err)
	}
	if len(adminRole.Users) == 0 {
		t.Fatalf("expected admin role to have at least one user")
	}
}
