package handlers

import (
	"blog-system/config"
	"blog-system/models"
	"blog-system/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(c *gin.Context) {
	var user models.User
	// 绑定JSON参数
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Logger.Error("注册参数绑定失败: ", err)
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.ERROR, "参数格式错误"))
		return
	}

	// 检查用户是否已存在
	var existUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.UserExist, "用户名已存在"))
		return
	}

	// 密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error("密码加密失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "密码加密失败"))
		return
	}
	user.Password = string(hashedPwd)

	// 写入数据库
	if err := config.DB.Create(&user).Error; err != nil {
		utils.Logger.Error("用户创建失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "注册失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("注册成功"))
}

// Login 用户登录
func Login(c *gin.Context) {
	var user models.User
	// 绑定JSON参数
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Logger.Error("登录参数绑定失败: ", err)
		c.JSON(http.StatusBadRequest, utils.ErrResponse(utils.ERROR, "参数格式错误"))
		return
	}

	// 查询用户
	var storedUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrResponse(utils.UserNotExist, "用户名或密码错误"))
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrResponse(utils.PasswordError, "用户名或密码错误"))
		return
	}

	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		utils.Logger.Error("Token生成失败: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrResponse(utils.ERROR, "登录失败"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"token": tokenString}))
}
