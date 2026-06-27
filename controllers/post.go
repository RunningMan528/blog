package controllers

import (
	"blog/config"
	"blog/models"
	"blog/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
	Url     string `json:"url"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
	Url     string `json:"url"`
}

type GetPostsRequest struct {
	Page uint `json:"page"`
	Size uint `json:"size"`
}

// CreatePost 创建文章
func (ps *PostController) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest
	// 获取参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 从context获取用户ID
	uID, exist := ctx.Get("user_id")
	if !exist {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	userID, idErr := strconv.Atoi(fmt.Sprintf("%v", uID))
	// 校验参数
	if req.Content == "" || req.Title == "" || idErr != nil {
		utils.BadRequest(ctx, "Invalid arguments")
		return
	}

	post := models.Post{
		UserID:  uint(userID),
		Content: req.Content,
		Title:   req.Title,
		Url:     req.Url,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, post)
}

// GetPosts 获取文章列表
func (ps *PostController) GetPosts(ctx *gin.Context) {
	var req GetPostsRequest
	// 获取query参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	// 分页查询
	var posts []models.Post
	if err := config.DB.Preload("User").
		Scopes(utils.Paginate(int(req.Page), int(req.Size))).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}
	// 文章总数量
	var postCount int64
	if err := config.DB.Model(&models.Post{}).Count(&postCount).Error; err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}
	utils.Success(ctx, gin.H{
		"post_list": posts,
		"total":     postCount,
		"page":      req.Page,
		"size":      req.Size,
	})
}

// GetPost 获取单个文章详情
func (ps *PostController) GetPost(ctx *gin.Context) {
	// 获取参数
	postID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid post ID")
		return
	}

	// 查询文章详情
	var post models.Post
	if err := config.DB.Model(&models.Post{}).
		Preload("Comments.User").
		Preload("User").
		Find(&post, postID).Error; err != nil {
		utils.NotFound(ctx, "Not found this post")
		return
	}

	utils.Success(ctx, post)
}

// UpdatePost 更新文章
func (ps *PostController) UpdatePost(ctx *gin.Context) {
	// path参数
	postID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid post ID")
		return
	}
	var req UpdatePostRequest
	// 获取body参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	// 校验是否授权
	uID, exist := ctx.Get("user_id")
	if !exist {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(ctx, "post not found")
		return
	}

	// 检查是否是文章作者
	if post.UserID != uID.(uint) {
		utils.Forbidden(ctx, "You can only update your own posts")
		return
	}

	// 校验参数
	_, idErr := strconv.Atoi(fmt.Sprintf("%v", uID))
	if req.Content == "" || req.Title == "" || idErr != nil {
		utils.BadRequest(ctx, "Invalid arguments")
		return
	}

	// 更新数据库
	// 第一种方式
	/*
		if err := config.DB.Model(&models.Post{}).
			Where("id = ?", postID).
			Updates(map[string]any{"content": req.Content, "title": req.Title}).Error; err != nil {
			utils.InternalServerError(ctx, err.Error())
			return
		}
	*/
	// 第二种方式
	post.Content = req.Content
	post.Title = req.Title
	post.Url = req.Url

	if err := config.DB.Save(&post).Error; err != nil {
		utils.InternalServerError(ctx, "Failed to update post")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&post, postID)
	utils.Success(ctx, post)

	//utils.Success(ctx, fmt.Sprintf("Update post success post_id:%v", postID))
}

// DeletePost 删除文章
func (ps *PostController) DeletePost(ctx *gin.Context) {
	// 获取参数
	postID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid post ID")
		return
	}

	// 校验是否登录
	userID, exist := ctx.Get("user_id")
	if !exist {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	var existPost models.Post
	if err := config.DB.First(&existPost, postID).Error; err != nil {
		utils.NotFound(ctx, "Post not found")
		return
	}

	// 检查是否是文章作者
	if existPost.UserID != userID.(uint) {
		utils.Forbidden(ctx, "You can only delete your own posts")
		return
	}

	// 执行软删除
	if err := config.DB.Delete(&existPost).Error; err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"message": "Post deleted successfully"})
}
