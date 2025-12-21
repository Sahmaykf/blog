package model

import (
	"gorm.io/gorm"
)

// Post文章模型
type Post struct {
	gorm.Model
	Title         string    `json:"title" binding:"required"`
	Content       string    `json:"content" binding:"required"`
	UserID        uint      `json:"user_id"` //uid
	User          User      `json:"user"`
	Tags          []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	TagNames      []string  `json:"tag_names" gorm:"-"` // 用于接收前端传来的标签名列表
	Comments      []Comment `json:"comments"`
	LikedBy       []*User   `gorm:"many2many:user_likes;" json:"liked_by,omitempty"`
	FavoritedBy   []*User   `gorm:"many2many:user_favorites;" json:"favorited_by,omitempty"`
	LikeCount     int       `json:"like_count" gorm:"default:0"`
	FavoriteCount int       `json:"favorite_count" gorm:"default:0"`
	CommentCount  int       `json:"comment_count" gorm:"default:0"`
	Status        string    `json:"status" gorm:"type:varchar(20);default:'published'"` // published, hidden
	IsTop         bool      `json:"is_top" gorm:"default:false"`                        // 个人主页置顶
	IsSystemTop   bool      `json:"is_system_top" gorm:"default:false"`                 // 全站首页置顶（仅管理员）
	ViewCount     int       `json:"view_count" gorm:"default:0"`
}
