package model

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID     uint   `json:"user_id"` // 接收通知的用户ID
	Type       string `json:"type"`    // comment, like, etc.
	Content    string `json:"content"` // 通知内容摘要
	IsRead     bool   `json:"is_read" gorm:"default:false"`
	FromUserID uint   `json:"from_user_id"` // 触发通知的用户ID
	FromUser   User   `json:"from_user"`
	PostID     uint   `json:"post_id"`    // 相关文章ID
	CommentID  uint   `json:"comment_id"` // 相关评论ID
}
