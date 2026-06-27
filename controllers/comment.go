package controllers

import (
	"blog/config"
	"blog/models"
	"blog/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
}

type GetCommentsRequest struct {
	Page uint `json:"page"`
	Size uint `json:"size"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(ctx *gin.Context) {
	// 获取path参数
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid post ID")
		return
	}
	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(ctx, "post not found")
		return
	}
	// 校验是否授权
	userID, exist := ctx.Get("user_id")
	if !exist {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	comment := models.Comment{
		UserID:  userID.(uint),
		PostID:  uint(postID),
		Content: req.Content,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		utils.InternalServerError(ctx, "Failed to create comment")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&comment, comment.ID)

	utils.Success(ctx, comment)
}

// GetComments 获取文章的评论列表
func (cc *CommentController) GetComments(ctx *gin.Context) {

	// 获取path参数
	postID, err := strconv.ParseUint(ctx.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "Invalid post ID")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(ctx, "post not found")
		return
	}

	// 获取query分页参数
	var req GetCommentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "Invalid arguments")
		return
	}

	var comments []models.Comment
	page, pageErr := strconv.Atoi(fmt.Sprintf("%v", req.Page))
	size, sizeErr := strconv.Atoi(fmt.Sprintf("%v", req.Size))
	if pageErr != nil || sizeErr != nil {
		utils.BadRequest(ctx, "Invalid arguments")
		return
	}

	if err := config.DB.Preload("User").
		Where("post_id = ?", postID).
		Scopes(utils.Paginate(page, size)).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		utils.InternalServerError(ctx, "Failed to get comments")
		return
	}

	var totalCount int64
	config.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&totalCount)

	utils.Success(ctx, gin.H{
		"comments": comments,
		"page":     page,
		"size":     size,
		"total":    totalCount,
	})

}
