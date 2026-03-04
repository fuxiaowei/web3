package basics

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"lesson02examples/testutil"

	"gorm.io/gorm"
)

// TestCRUDDemo demonstrates the complete CRUD operations in GORM
// This test covers Create, Read, Update, and Delete operations with various patterns
func TestCRUDDemo(t *testing.T) {
	db := testutil.NewTestDB(t, "crud.db")

	// Define the User model
	// GORM will automatically map this struct to a "users" table
	type User struct {
		ID        uint      `gorm:"primaryKey"` // Primary key, auto-increment
		Name      string    // Regular field
		Email     string    `gorm:"uniqueIndex"` // Unique index for email
		Age       uint8     // Age field
		Status    string    // Status field
		CreatedAt time.Time // GORM will auto-populate on create
		UpdatedAt time.Time // GORM will auto-populate on create/update
	}

	// AutoMigrate creates the table if it doesn't exist
	// It will also add new columns if the struct has new fields
	// Note: It will NOT delete existing columns or modify existing data
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Seed initial data: batch insert using Create
	// Create can accept a single struct or a slice for batch insertion
	seed := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 28, Status: "active"},
		{Name: "Alice1", Email: "alice1@example.com", Age: 28, Status: "active"},
		{Name: "Alice2", Email: "alice2@example.com", Age: 28, Status: "inactive"},
		{Name: "Alice3", Email: "alice3@example.com", Age: 28, Status: "inactive"},
		{Name: "Bob", Email: "bob@example.com", Age: 31, Status: "active"},
		{Name: "Bob1", Email: "bob1@example.com", Age: 32, Status: "active"},
		{Name: "Bob2", Email: "bob2@example.com", Age: 33, Status: "inactive"},
		{Name: "Bob3", Email: "bob3example.com", Age: 34, Status: "inactive"},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}

	// CREATE: Single record insertion
	// Create inserts a new record and automatically populates:
	// - Primary key (ID) after insertion
	// - CreatedAt and UpdatedAt timestamps
	// You can also use Select/Omit to control which fields are inserted
	t.Run("create", func(t *testing.T) {
		u := User{Name: "Diane", Email: "diane@example.com", Age: 30, Status: "active"}
		// Create returns the inserted record with ID populated
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("create user: %v", err)
		}
		// After Create, u.ID is automatically populated by GORM
		t.Logf("new user id=%d", u.ID)
	})

	// // READ: Query operations
	// // GORM provides several query methods:
	// // - First: Get the first record matching conditions (returns error if not found)
	// // - Take: Get one record (doesn't require conditions)
	// // - Find: Get all matching records (returns empty slice if none found)
	// // - Scan: Scan results into a struct or map
	// // Always check for gorm.ErrRecordNotFound when using First

	t.Run("query/first", func(t *testing.T) {
		var user User
		// First: Get the first record matching conditions
		// Returns gorm.ErrRecordNotFound if no record found
		// Can use conditions: db.First(&user, "email = ?", "alice@example.com")
		// Or primary key: db.First(&user, 1)
		if err := db.Where("status = ?", "active").First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				t.Fatalf("no active user found")
			}
			t.Fatalf("query first active user: %v", err)
		}
		t.Logf("first active user: %+v", user)

		// First with primary key
		var userByID User
		if err := db.First(&userByID, 1).Error; err != nil {
			t.Fatalf("query user by ID: %v", err)
		}
		t.Logf("user by ID 1: %+v", userByID)
	})

	t.Run("query/take", func(t *testing.T) {
		var user User
		// Take: Get one record without requiring conditions
		// Doesn't return error if no record found (just doesn't populate the struct)
		// Useful when you just want any record from the table
		if err := db.Take(&user).Error; err != nil {
			t.Fatalf("take user: %v", err)
		}
		t.Logf("taken user: %+v", user)

		// Take with conditions
		var activeUser User
		if err := db.Where("status = ?", "inactive").Take(&activeUser).Error; err != nil {
			t.Fatalf("take active user: %v", err)
		}
		t.Logf("taken active user: %+v", activeUser)
	})

	t.Run("query/find", func(t *testing.T) {
		var actives []User
		// Find: Get all matching records
		// Returns empty slice if no records found (no error)
		// Where: Add conditions to the query
		// Order: Sort results (asc/desc)
		if err := db.Where("status = ?", "active").Order("created_at desc").Find(&actives).Error; err != nil {
			t.Fatalf("query actives: %v", err)
		}
		if len(actives) == 0 {
			t.Fatalf("expected at least one active user")
		}
		t.Logf("active users: %+v", actives)

		// Find all records
		var allUsers []User
		if err := db.Find(&allUsers).Error; err != nil {
			t.Fatalf("find all users: %v", err)
		}
		t.Logf("all users count: %d", len(allUsers))
	})

	t.Run("query/scan", func(t *testing.T) {
		// Scan: Scan results into a struct or map
		// Useful when you only need specific fields or want to scan into a different structure
		type UserSummary struct {
			Name   string
			Email  string
			Status string
		}
		var summaries []UserSummary
		// Select specific fields and scan into a different struct
		if err := db.Model(&User{}).Select("name", "email", "status").Where("status = ?", "active").Scan(&summaries).Error; err != nil {
			t.Fatalf("scan user summaries: %v", err)
		}
		t.Logf("user summaries: %+v", summaries)

		// Scan into a map
		var result map[string]interface{}
		if err := db.Model(&User{}).Select("name", "email", "age").Where("email = ?", "alice@example.com").Scan(&result).Error; err != nil {
			t.Fatalf("scan to map: %v", err)
		}
		t.Logf("user as map: %+v", result)

		// Scan into primitive values
		var count int64
		if err := db.Model(&User{}).Where("status = ?", "active").Count(&count).Error; err != nil {
			t.Fatalf("count active users: %v", err)
		}
		t.Logf("active users count: %d", count)
	})

	// // UPDATE: Update operations
	// // GORM provides different update methods:
	// // - Save: Updates all fields (including zero values)
	// // - Updates: Updates specified fields (ignores zero values by default)
	// // - Update: Updates a single field
	// // Use Select to specify which fields to update, or Omit to exclude fields
	// // Model(&user) is used to specify the model for the update operation
	t.Run("update", func(t *testing.T) {
		var user User
		// First: Get the first record matching the condition
		// Second parameter can be a condition string or primary key value
		if err := db.First(&user, "email = ?", "diane@example.com").Error; err != nil {
			t.Errorf("load user: %v", err)
		}
		fmt.Print(&user)
		// Select: Only update specified fields (Age and Status)
		// This prevents updating other fields and ignores zero values for non-selected fields
		if err := db.Model(&user).Select("Age", "Status").Where("email = ?", "alice@example.com").Updates(User{Age: 31, Status: "vip"}).Error; err != nil {
			t.Fatalf("update fields: %v", err)
		}
		// Reload the user to verify the update
		// First with ID: Query by primary key
		if err := db.First(&user, 1).Error; err != nil {
			t.Fatalf("reload user: %v", err)
		}
		if user.Age != 31 || user.Status != "vip" {
			t.Fatalf("unexpected updated values: %+v", user)
		}
	})

	// // BULK UPDATE: Update multiple records at once
	// // Use Model(&User{}) without a specific instance to perform bulk operations
	// // Updates can accept a struct or a map[string]any
	// // RowsAffected indicates how many rows were actually updated
	t.Run("bulk update", func(t *testing.T) {
		// Model(&User{}): Specify the model for bulk operation
		// Where: Add conditions to filter which records to update
		// Updates: Update all matching records
		// Using map[string]any allows updating specific fields without zero value issues
		res := db.Model(&User{}).Where("status = ?", "inactive").Updates(map[string]any{"status": "pending_review"})
		if res.Error != nil {
			t.Fatalf("bulk update: %v", res.Error)
		}
		// RowsAffected: Check how many rows were actually updated
		if res.RowsAffected == 0 {
			t.Fatalf("expected rows to be updated")
		}
	})

	// // DELETE: Delete operations
	// // Delete can be used with:
	// // - A specific instance: db.Delete(&user)
	// // - A model with conditions: db.Delete(&User{}, "id = ?", id)
	// // - Bulk delete: db.Where(...).Delete(&User{})
	// // Note: Soft delete will be covered in the advanced section
	// // After deletion, querying the record should return gorm.ErrRecordNotFound
	t.Run("delete", func(t *testing.T) {
		var user User
		// First: Load the user to delete
		if err := db.First(&user, "email = ?", "alice1@example.com").Error; err != nil {
			t.Fatalf("load user: %v", err)
		}
		// Delete: Delete by primary key
		// First parameter is the model type, second is the primary key value
		if err := db.Delete(&User{}, user.ID).Error; err != nil {
			t.Fatalf("delete: %v", err)
		}
		// Verify deletion: Query should return gorm.ErrRecordNotFound
		// Always use errors.Is to check for gorm.ErrRecordNotFound
		err := db.First(&User{}, user.ID).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("expected not found, got %v", err)
		}
	})
}
