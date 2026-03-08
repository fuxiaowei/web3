package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// 全局日志对象
var Logger *logrus.Logger

// InitLogger 初始化日志
func InitLogger() {
	Logger = logrus.New()
	// 设置日志输出（文件+控制台）
	file, _ := os.OpenFile("blog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	Logger.SetOutput(file)
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// 设置日志级别
	Logger.SetLevel(logrus.InfoLevel)
}
