package config

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

func Load() *Config {
	// 简化配置加载，实际应该使用 Viper
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "password",
			DBName:   "mydb",
		},
		JWT: JWTConfig{
			Secret: "your-secret-key-change-in-production",
			Expire: "24h",
		},
	}
}

