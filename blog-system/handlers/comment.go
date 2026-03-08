package handlers

import (
	"blog-system/config"
	"blog-system/models"
	"blog-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	var comment models.Comment
	// 绑定参数
	if err := c.ShouldBindJSON(&comment); err != nil {
		utils.Logger.Error("创建评论参数绑定失败: ", err)
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.ERROR, "参数格式错误"))
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, comment.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrResponse(utils.PostNotExist, "文章不存在"))
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	comment.UserID = userID.(uint)

	// 写入数据库
	if err := config.DB.Create(&comment).Error; err != nil {
		utils.Logger.Error("创建评论失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "创建评论失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(comment))
}

// GetComments 获取文章评论
func GetComments(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []models.Comment
	// 查询指定文章的所有评论（关联用户）
	if err := config.DB.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		utils.Logger.Error("查询评论失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "查询评论失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(comments))
}
