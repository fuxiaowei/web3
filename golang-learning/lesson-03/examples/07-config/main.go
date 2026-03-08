package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

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

func init() {
	// 设置配置文件名称（不含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 添加配置文件搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.app")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "debug")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v", err)
		log.Println("Using default values and environment variables")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

// LoadConfig 加载并解析配置文件，将 Viper 中的配置数据映射到 Config 结构体
// 返回解析后的配置对象指针和可能的错误
// 配置数据来源包括：配置文件、环境变量和默认值（在 init 函数中已设置）
func LoadConfig() (*Config, error) {
	var config Config

	// 使用 viper.Unmarshal 将配置数据解析到 config 结构体中
	// 如果解析失败，返回错误
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 返回解析成功的配置对象
	return &config, nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(config.Server.Mode)

	r := gin.Default()

	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"server": config.Server,
			"database": gin.H{
				"host":     config.Database.Host,
				"port":     config.Database.Port,
				"username": config.Database.Username,
				"dbname":   config.Database.DBName,
				// 不返回密码
			},
			"jwt": gin.H{
				"expire": config.JWT.Expire,
				// 不返回密钥
			},
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"config_file": viper.ConfigFileUsed(),
		})
	})

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Server starting on %s", addr)
	r.Run(addr)
}
