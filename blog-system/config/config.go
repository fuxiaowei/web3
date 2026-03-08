package config

import (
	"gorm.io/gorm"
)

// 全局数据库连接
var DB *gorm.DB

// JWT 秘钥（生产环境建议通过环境变量读取）
const JWTSecret = "your_secure_jwt_secret_key_2024"

// 数据库配置（支持MySQL/SQLite，按需切换）
const (
	// SQLite 配置
	SQLiteDSN = "blog.db"

	// MySQL 配置（需提前创建数据库 blog）
	MySQLDSN = "root:root@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
)
