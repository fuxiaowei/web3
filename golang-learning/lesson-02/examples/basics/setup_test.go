package basics

import (
	"testing"
	"time"

	"lesson02examples/testutil"
)

func TestSetupDemo(t *testing.T) {
	db := testutil.NewTestDB(t, "setup.db")

	type User struct {
		ID        uint      `gorm:"primaryKey"`
		Name      string    `gorm:"size:64;not null"`
		Email     string    `gorm:"size:128;uniqueIndex;not null"`
		Age       uint8     `gorm:"not null"`
		Status    string    `gorm:"size:16;default:active;index"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	// AutoMigrate automatically creates or updates the database table structure
	// based on the Go struct definition. It will:
	// - Create the table if it doesn't exist
	// - Add new columns if the struct has new fields
	// - Add indexes based on struct tags (e.g., uniqueIndex, index)
	// Note: It will NOT delete existing columns or modify existing data
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// Clear the users table before inserting new data
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		t.Fatalf("clear users table: %v", err)
	}

	// Prepare test data: a slice of User structs to be inserted
	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 28},
		{Name: "Bob", Email: "bob@example.com", Age: 32},
		{Name: "Celine", Email: "celine@example.com", Age: 25, Status: "inactive"},
	}

	// CreateInBatches inserts records in batches to optimize performance
	// The second parameter (2) specifies the batch size - how many records to insert per batch
	// In this case, with 3 users and batch size of 2:
	// - First batch: inserts 2 users (Alice and Bob)
	// - Second batch: inserts 1 user (Celine)
	// This approach is useful for large datasets to avoid SQL statement length limits
	// and reduce memory usage during bulk inserts
	if err := db.CreateInBatches(users, 3).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}

	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}

	if count != int64(len(users)) {
		t.Fatalf("expected %d users, got %d", len(users), count)
	}

	t.Logf("created %d users", count)
}
