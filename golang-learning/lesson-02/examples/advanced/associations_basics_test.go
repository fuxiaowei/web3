package advanced

import (
	"testing"

	"lesson02examples/testutil"

	"gorm.io/gorm"
)

// TestAssociationsBasics demonstrates basic association relationships in GORM
// This test covers:
// 1. Defining associations (Has One, Has Many, Belongs To)
// 2. Creating records with associations
// 3. Basic preloading of associations
func TestAssociationsBasics(t *testing.T) {
	db := testutil.NewTestDB(t, "associations_basics.db")

	// AutoMigrate creates all tables and their relationships
	// GORM automatically creates foreign key constraints based on the struct relationships
	// Important: All related models must be migrated together
	if err := db.AutoMigrate(&user{}, &profile{}, &product{}, &order{}, &orderItem{}, &role{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Clean up existing data before seeding (for test isolation)
	// Disable foreign key constraints temporarily to allow deletion in any order
	db.Exec("PRAGMA foreign_keys = OFF")
	db.Exec("DELETE FROM user_roles")
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM profiles")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM roles")
	db.Exec("PRAGMA foreign_keys = ON")

	// Seed products first (independent entities, no foreign keys)
	// Products are referenced by order items, so they must exist before creating orders
	products := []product{
		{Name: "Go 语言圣经", Price: 10800, SKU: "BOOK-001"},
		{Name: "机械键盘", Price: 39900, SKU: "GEAR-201"},
		{Name: "人体工学椅", Price: 129900, SKU: "GEAR-301"},
	}
	if err := db.Create(&products).Error; err != nil {
		t.Fatalf("seed products: %v", err)
	}

	// CREATE WITH ASSOCIATIONS: Create a user with nested associations
	// When creating a record with associations, GORM will automatically:
	// 1. Create the main record (user)
	// 2. Create associated records (profile, orders, order items)
	// 3. Set foreign keys correctly
	//
	// FullSaveAssociations: When true, saves all associations even if they are zero values
	// This is useful when you want to create nested structures in one operation
	u := user{
		Name:  "Alice",
		Email: "alice@example.com",
		// Has One: Profile association
		Profile: profile{
			Nickname: "阿狸",
			Phone:    "13800000000",
			Address:  "上海市徐汇区漕河泾开发区",
		},
		// Has Many: Orders association (nested with Has Many OrderItems)
		Orders: []order{
			{
				OrderNo:    "ORDER-20241001-001",
				Status:     "PAID",
				TotalPrice: 10800 + 39900,
				// Has Many: OrderItems within Order
				Items: []orderItem{
					{ProductID: products[0].ID, Quantity: 1, UnitPrice: products[0].Price},
					{ProductID: products[1].ID, Quantity: 1, UnitPrice: products[1].Price},
				},
			},
			{
				OrderNo:    "ORDER-20241002-002",
				Status:     "PENDING",
				TotalPrice: 2*products[2].Price + products[0].Price,
				Items: []orderItem{
					{ProductID: products[2].ID, Quantity: 2, UnitPrice: products[2].Price},
					{ProductID: products[0].ID, Quantity: 1, UnitPrice: products[0].Price},
				},
			},
		},
	}

	// Session with FullSaveAssociations: Ensures all nested associations are saved
	// Without this, GORM might skip zero-value associations
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// BASIC PRELOAD: Eager loading associations
	// Preload loads associated records in separate queries (more efficient than lazy loading)
	// Lazy loading: Accessing association triggers additional SQL queries (N+1 problem)
	// Eager loading: All associations loaded upfront in optimized queries
	var result user
	// Preload Profile: Load the user's profile (Has One)
	// Preload Orders: Load all orders for the user (Has Many)
	if err := db.Preload("Profile").
		Preload("Orders").
		Preload("Orders.Items").
		First(&result, "email = ?", u.Email).Error; err != nil {
		t.Fatalf("preload query: %v", err)
	}

	// Verify associations are loaded
	if result.Profile.Nickname == "" {
		t.Fatal("profile should be loaded")
	}
	if len(result.Orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(result.Orders))
	}
}
