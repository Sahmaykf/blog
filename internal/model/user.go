package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string    `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Password   string    `gorm:"not null" json:"-"` // 存储哈希后的密码
	Email      string    `gorm:"uniqueIndex;size:255" json:"email"`
	Avatar     string    `json:"avatar"`
	Role       string    `gorm:"default:'user'" json:"role"` // user or admin
	Posts      []Post    `json:"posts,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
	Followers  []*User   `gorm:"many2many:user_followers;joinForeignKey:followed_id;joinReferences:follower_id" json:"-"`
	Following  []*User   `gorm:"many2many:user_followers;joinForeignKey:follower_id;joinReferences:followed_id" json:"-"`
	LikedPosts []*Post   `gorm:"many2many:user_likes;" json:"liked_posts,omitempty"`
}
