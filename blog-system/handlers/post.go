package handlers

import (
	"blog-system/config"
	"blog-system/models"
	"blog-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var post models.Post
	// 绑定参数
	if err := c.ShouldBindJSON(&post); err != nil {
		utils.Logger.Error("创建文章参数绑定失败: ", err)
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.ERROR, "参数格式错误"))
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	post.UserID = userID.(uint)

	exist, err := models.CheckUserExist(userID.(uint))
	if err != nil {
		utils.Logger.Error("检查用户是否存在失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "检查用户是否存在失败"))
		return
	}

	// 校验用户是否存在，防止用户登录token存在，但是表里用户数据被删的情况
	if !exist {
		utils.Logger.Error("用户不存在: ", userID.(uint))
		c.JSON(http.StatusNotFound, utils.ErrResponse(utils.UserNotExist, "用户不存在，请先注册用户！"))
		return
	}

	// 写入数据库
	if err := config.DB.Create(&post).Error; err != nil {
		utils.Logger.Error("创建文章失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "创建文章失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(post))
}

// GetPosts 获取所有文章
func GetPosts(c *gin.Context) {
	var posts []models.Post

	// 查询所有文章（关联用户信息）
	if err := config.DB.Preload("User").Find(&posts).Error; err != nil {
		utils.Logger.Error("查询文章列表失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "查询失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(posts))
}

// GetPost 获取单篇文章
func GetPost(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	utils.Logger.Info("获取文章详情，id: ", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.ERROR, "无效的文章ID"))
		return
	}

	var post models.Post

	// 查询文章（关联用户+评论）
	if err := config.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrResponse(utils.PostNotExist, "文章不存在"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(post))
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	utils.Logger.Info("更新文章，id: ", id)

	userID, _ := c.Get("userID")

	var post models.Post

	// 查询文章是否存在
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrResponse(utils.PostNotExist, "文章不存在"))
		return
	}

	// 验证权限（仅作者可更新）
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, utils.ErrResponse(utils.PermissionDeny, "无权限修改该文章"))
		return
	}

	// 绑定更新参数
	var updateData models.Post
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Logger.Error("更新文章参数绑定失败: ", err)
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.ERROR, "参数格式错误"))
		return
	}

	// 更新数据库
	if err := config.DB.Model(&post).Updates(models.Post{Title: updateData.Title, Content: updateData.Content}).Error; err != nil {
		utils.Logger.Error("更新文章失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "更新失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("更新成功"))
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	utils.Logger.Info("删除文章，id: ", id)

	userID, err := c.Get("userID")
	if !err {
		utils.Logger.Error("获取用户ID失败: userID 不存在于上下文中")
		c.JSON(http.StatusUnauthorized, utils.ErrResponse(utils.TokenEmpty, "Token 为空"))
		return
	}

	var post models.Post

	// 查询文章是否存在
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrResponse(utils.PostNotExist, "文章不存在"))
		return
	}

	// 验证权限
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, utils.ErrResponse(utils.PermissionDeny, "无权限删除该文章"))
		return
	}

	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		utils.Logger.Error("删除文章失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "删除失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("删除成功"))
}
