package models

import (
	"blog-system/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;size:50" json:"username"` // 用户名
	Password string `gorm:"not null;size:100" json:"password"`       // 加密后的密码
	Email    string `gorm:"unique;not null;size:100" json:"email"`   // 邮箱
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title   string `gorm:"not null;size:200" json:"title"`    // 文章标题
	Content string `gorm:"not null;type:text" json:"content"` // 文章内容
	UserID  uint   `gorm:"not null" json:"user_id"`           // 关联用户ID
	User    User   `gorm:"foreignKey:UserID" json:"user"`     // 用户关联
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;type:text" json:"content"` // 评论内容
	UserID  uint   `gorm:"not null" json:"user_id"`           // 关联用户ID
	PostID  uint   `gorm:"not null" json:"post_id"`           // 关联文章ID
	User    User   `gorm:"foreignKey:UserID" json:"user"`     // 用户关联
	Post    Post   `gorm:"foreignKey:PostID" json:"post"`     // 文章关联
}

// 初始化数据库表
func InitDB() {
	var err error
	// 方式1：使用SQLite（无需额外部署，适合测试）
	// config.DB, err = gorm.Open(sqlite.Open(config.SQLiteDSN), &gorm.Config{})

	// 方式2：使用MySQL（生产环境推荐）
	config.DB, err = gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移表结构（创建/更新表）
	err = config.DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("表结构迁移失败: " + err.Error())
	}
}
