package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content" binding:"required"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user"`
	PostID  uint   `json:"post_id" binding:"required"`
	Post    *Post  `json:"-"`
}
