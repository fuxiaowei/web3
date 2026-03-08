package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 创建上传目录
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		panic(err)
	}

	// ========== 单文件上传 ==========
	r.POST("/api/upload", uploadFile)

	// ========== 多文件上传 ==========
	r.POST("/api/upload-multiple", uploadFiles)

	// ========== 文件下载 ==========
	r.GET("/api/download/:filename", downloadFile)

	// ========== 静态文件服务 ==========
	// 提供静态文件服务
	r.Static("/static", "./static")

	// 提供文件系统
	r.StaticFS("/files", http.Dir("./uploads"))

	r.Run(":8080")
}

// ========== 单文件上传 ==========
func uploadFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存文件
	filename := filepath.Base(file.Filename)
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
		"size":     file.Size,
	})
}

// ========== 多文件上传 ==========
func uploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	var filenames []string
	for _, file := range files {
		filename := filepath.Join("uploads", filepath.Base(file.Filename))
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		filenames = append(filenames, filename)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"files":   filenames,
		"count":   len(filenames),
	})
}

// ========== 文件下载 ==========
func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join("uploads", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	// 返回文件
	c.File(filepath)
}

