package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Post Post `json:"post,omitempty" gorm:"foreignKey:PostID"`
}

// 钩子函数创建之后对应的文章评论数+1
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	result := tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		UpdateColumn("comment_count", gorm.Expr("COALESCE(comment_count,0) + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
