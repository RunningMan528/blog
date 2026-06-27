package controllers

import (
	"blog/config"
	"blog/models"
	"blog/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthController struct{}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register 用户注册
func (au *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	// 获取参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(ctx, "Username already exists")
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.BadRequest(ctx, "Email already exsits")
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 密码会在BeforeCreate钩子中自动加密
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.BadRequest(ctx, "Failed to create user")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		logrus.Info(map[string]string{"generateErr": err.Error()})
		utils.InternalServerError(ctx, "Failed to generate token")
		return
	}
	utils.Success(ctx, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login 用户登录
func (au *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	// 获取登录参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).Find(&user).Error; err != nil {
		utils.BadRequest(ctx, "Account not exist")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		utils.Unauthorized(ctx, "Invalid username or password")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to generate token")
		return
	}

	utils.Success(ctx, AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile 获取用户信息
func (au *AuthController) GetProfile(ctx *gin.Context) {
	// 获取授权中间件gin.Context中保存的用户信息
	// 如果没有证明没有进行授权登录
	userID, exist := ctx.Get("user_id")
	if !exist {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	utils.Success(ctx, user)

}
