package model

import (
	"gorm.io/gorm"
)

// Post文章模型
type Post struct {
	gorm.Model
	Title    string    `json:"title" binding:"required"`
	Content  string    `json:"content" binding:"required"`
	UserID   uint      `json:"user_id"` //uid
	User     User      `json:"user"`
	Tags     []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	TagNames []string  `json:"tag_names" gorm:"-"` // 用于接收前端传来的标签名列表
	Comments []Comment `json:"comments"`
	LikedBy  []*User   `gorm:"many2many:user_likes;" json:"liked_by,omitempty"`
	Status   string    `json:"status" gorm:"type:varchar(20);default:'published'"` // published, hidden
}
