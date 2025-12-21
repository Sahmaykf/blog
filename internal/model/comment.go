package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content     string    `json:"content" binding:"required"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user"`
	PostID      uint      `json:"post_id" binding:"required"`
	Post        *Post     `json:"-"`
	ParentID    *uint     `json:"parent_id"`              // 父评论ID
	ReplyToID   *uint     `json:"reply_to_id"`            // 被回复的评论ID（用于多级回复显示）
	ReplyToUser *User     `json:"reply_to_user" gorm:"-"` // 仅用于前端显示，不直接关联
	Replies     []Comment `json:"replies" gorm:"foreignKey:ParentID"`
}
