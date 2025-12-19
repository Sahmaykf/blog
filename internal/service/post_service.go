package service

import (
	"errors"
	"simple-blog/internal/database"
	"simple-blog/internal/model"

	"gorm.io/gorm"
)

type PostService struct{}

func (s *PostService) CreatePost(post *model.Post) error {
	// 处理标签
	var tags []model.Tag
	for _, tagName := range post.TagNames {
		var tag model.Tag
		// 查找或创建标签
		if err := database.DB.FirstOrCreate(&tag, model.Tag{Name: tagName}).Error; err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	post.Tags = tags

	return database.DB.Create(post).Error
}

func (s *PostService) GetPostDetail(id string) (*model.Post, error) {
	var post model.Post
	if err := database.DB.Preload("User").Preload("Tags").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetPostList(page, pageSize int, tagName string, userID uint, status string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.Post{})

	// 如果有标签筛选
	if tagName != "" {
		db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", tagName)
	}

	if userID > 0 {
		db = db.Where("posts.user_id = ?", userID)
	}

	if status != "" {
		db = db.Where("posts.status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Preload("User").Preload("Tags").Offset(offset).Limit(pageSize).Order("posts.created_at desc").Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *PostService) UpdatePost(id string, userId uint, userRole string, updatedPost *model.Post) error {
	var post model.Post
	if err := database.DB.Preload("Tags").First(&post, id).Error; err != nil {
		return errors.New("post not found")
	}

	if post.UserID != userId && userRole != "admin" {
		return errors.New("unauthorized")
	}

	// 处理标签更新
	var tags []model.Tag
	for _, tagName := range updatedPost.TagNames {
		var tag model.Tag
		if err := database.DB.FirstOrCreate(&tag, model.Tag{Name: tagName}).Error; err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	// 使用 Association 替换标签
	if err := database.DB.Model(&post).Association("Tags").Replace(tags); err != nil {
		return err
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	if updatedPost.Status != "" {
		post.Status = updatedPost.Status
	}
	return database.DB.Save(&post).Error
}

func (s *PostService) DeletePost(id string, userId uint, userRole string) error {
	var post model.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		return errors.New("post not found")
	}

	if post.UserID != userId && userRole != "admin" {
		return errors.New("unauthorized")
	}

	return database.DB.Delete(&post).Error
}

func (s *PostService) GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	// 获取所有标签，可以按使用频率排序，这里简单起见直接获取所有
	if err := database.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// LikePost 点赞文章
func (s *PostService) LikePost(userID uint, postID uint) error {
	return database.DB.Model(&model.User{Model: gorm.Model{ID: userID}}).
		Association("LikedPosts").
		Append(&model.Post{Model: gorm.Model{ID: postID}})
}

// UnlikePost 取消点赞
func (s *PostService) UnlikePost(userID uint, postID uint) error {
	return database.DB.Model(&model.User{Model: gorm.Model{ID: userID}}).
		Association("LikedPosts").
		Delete(&model.Post{Model: gorm.Model{ID: postID}})
}

// GetLikedPosts 获取用户点赞的文章列表
func (s *PostService) GetLikedPosts(userID uint) ([]*model.Post, error) {
	var user model.User
	// 预加载 LikedPosts 及其关联的 User 和 Tags 信息，以便前端展示
	err := database.DB.Preload("LikedPosts.User").Preload("LikedPosts.Tags").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.LikedPosts, nil
}

// IsPostLiked 检查用户是否点赞了某文章
func (s *PostService) IsPostLiked(userID uint, postID uint) (bool, error) {
	var count int64
	err := database.DB.Table("user_likes").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

// GetPostLikeCount 获取文章点赞数
func (s *PostService) GetPostLikeCount(postID uint) (int64, error) {
	return database.DB.Model(&model.Post{Model: gorm.Model{ID: postID}}).Association("LikedBy").Count(), nil
}
