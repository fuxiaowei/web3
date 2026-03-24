package advanced

import (
	"testing"

	"lesson02examples/testutil"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TestAssociationsPreload demonstrates advanced Preload operations in GORM
// This test covers:
// 1. Conditional preloading (preload with conditions)
// 2. Nested preloading (preload associations of associations)
// 3. Preload all associations using clause.Associations
func TestAssociationsPreload(t *testing.T) {
	db := testutil.NewTestDB(t, "associations.db")

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

	// Seed products
	products := []product{
		{Name: "Go 语言圣经", Price: 10800, SKU: "BOOK-001"},
		{Name: "机械键盘", Price: 39900, SKU: "GEAR-201"},
		{Name: "人体工学椅", Price: 129900, SKU: "GEAR-301"},
	}
	if err := db.Create(&products).Error; err != nil {
		t.Fatalf("seed products: %v", err)
	}

	// Create user with orders
	u := user{
		Name:  "Alice",
		Email: "alice@example.com",
		Profile: profile{
			Nickname: "阿狸",
			Phone:    "13800000000",
			Address:  "上海市徐汇区漕河泾开发区",
		},
		Orders: []order{
			{
				OrderNo:    "ORDER-20241001-001",
				Status:     "PAID",
				TotalPrice: 10800 + 39900,
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

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// CONDITIONAL PRELOAD: Preload associations with conditions
	// This allows you to filter which associated records are loaded
	// Example: Load only PAID orders
	var result user
	// Preload Profile: Load the user's profile (Has One)
	// Preload Orders with condition: Load only PAID orders (Has Many with condition)
	// Preload Orders.Items.Product: Load nested associations (Order -> OrderItem -> Product)
	if err := db.Preload("Profile").
		Preload("Orders", "status = ?", "PAID"). // Conditional preload: only PAID orders
		Preload("Orders.Items.Product").         // Nested preload: OrderItems and their Products
		First(&result, "email = ?", u.Email).Error; err != nil {
		t.Fatalf("preload query: %v", err)
	}
	// After preload, result.Profile and result.Orders are populated
	// Since we filtered for PAID orders, only one order should be loaded
	if len(result.Orders) != 1 {
		t.Fatalf("expected 1 paid order, got %d", len(result.Orders))
	}

	// PRELOAD ALL ASSOCIATIONS: Using clause.Associations
	// clause.Associations automatically preloads all associations of the model
	// This is useful when you want to load all related data without specifying each association
	// Note: This only preloads direct associations, not nested ones
	var orders []order
	// Preload(clause.Associations): Preload all associations of order (Items, User, etc.)
	// This is equivalent to: Preload("Items").Preload("User") (if User association exists)
	if err := db.Model(&order{}).Preload(clause.Associations).
		Where("status = ?", "PAID").
		Find(&orders).Error; err != nil {
		t.Fatalf("load paid orders: %v", err)
	}
	if len(orders) == 0 {
		t.Fatalf("expected paid orders")
	}

	// CONDITIONAL PRELOAD WITH NESTED: Preload associations with conditions and nested associations
	// This demonstrates how to combine conditional preload with nested preload
	var filtered []order
	// Preload Items with condition: Only load items where quantity >= 2
	// Preload Items.Product: Also load the Product for each filtered item
	// This demonstrates how to combine conditional preload with nested preload
	if err := db.Preload("Items", "quantity >= ?", 2). // Conditional: quantity >= 2
								Preload("Items.Product"). // Nested: Load Product for each Item
								Find(&filtered).Error; err != nil {
		t.Fatalf("conditional preload: %v", err)
	}
	if len(filtered) == 0 {
		t.Fatalf("expected filtered orders")
	}
}

