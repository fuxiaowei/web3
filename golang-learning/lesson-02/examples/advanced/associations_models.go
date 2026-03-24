package advanced

import "time"

// ASSOCIATION RELATIONSHIPS IN GORM
// GORM supports three types of associations:
// 1. Has One / Belongs To: One-to-one relationship
// 2. Has Many: One-to-many relationship
// 3. Many to Many: Many-to-many relationship (with join table)

// user represents a user with multiple associations
// - Has One Profile: Each user has one profile (one-to-one)
// - Has Many Orders: Each user can have multiple orders (one-to-many)
// - Many to Many Roles: Each user can have multiple roles, each role can belong to multiple users
type user struct {
	ID        uint
	Name      string
	Email     string
	Profile   profile // Has One: One user has one profile
	Orders    []order // Has Many: One user has many orders
	Roles     []role  `gorm:"many2many:user_roles;"` // Many to Many: User has many roles through user_roles join table
	CreatedAt time.Time
	UpdatedAt time.Time
}

// profile represents user profile information
// Belongs To User: Profile belongs to one user (inverse of Has One)
// UserID is the foreign key that references user.ID
// uniqueIndex ensures one-to-one relationship (one profile per user)
type profile struct {
	ID        uint
	UserID    uint `gorm:"uniqueIndex"` // Foreign key to user, unique to enforce one-to-one
	Nickname  string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// order represents an order placed by a user
// Belongs To User: Order belongs to one user
// Has Many OrderItems: One order has many order items
type order struct {
	ID         uint
	OrderNo    string      `gorm:"uniqueIndex"` // Unique order number
	UserID     uint        // Foreign key to user
	Items      []orderItem // Has Many: One order has many items
	TotalPrice int64
	Status     string
	CreatedAt  time.Time
}

// orderItem represents an item in an order
// Belongs To Order: OrderItem belongs to one order
// Belongs To Product: OrderItem references one product
type orderItem struct {
	ID        uint
	OrderID   uint    // Foreign key to order
	ProductID uint    // Foreign key to product
	Product   product // Belongs To: OrderItem belongs to one product
	Quantity  int
	UnitPrice int64
	CreatedAt time.Time
}

// product represents a product in the system
// Referenced by OrderItem (many order items can reference one product)
type product struct {
	ID        uint
	Name      string
	Price     int64
	SKU       string `gorm:"uniqueIndex"` // Stock Keeping Unit, unique identifier
	CreatedAt time.Time
	UpdatedAt time.Time
}

// role represents a role in the system
// Many to Many Users: Each role can belong to multiple users, each user can have multiple roles
// GORM automatically creates a join table "user_roles" with user_id and role_id columns
type role struct {
	ID          uint
	Name        string `gorm:"uniqueIndex"` // Role name must be unique (e.g., "admin", "user", "editor")
	Description string
	Users       []user `gorm:"many2many:user_roles;"` // Many to Many: Role belongs to many users
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

