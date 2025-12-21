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
	if err := database.DB.Model(&model.Post{}).
		Select("posts.*, "+
			"(SELECT COUNT(*) FROM user_likes WHERE user_likes.post_id = posts.id) as like_count, "+
			"(SELECT COUNT(*) FROM user_favorites WHERE user_favorites.post_id = posts.id) as favorite_count, "+
			"(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count").
		Preload("User").Preload("Tags").First(&post, id).Error; err != nil {
		return nil, err
	}
	// 增加阅读量
	database.DB.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	return &post, nil
}

func (s *PostService) GetPostList(page, pageSize int, tagName string, userID uint, status string, orderBy string, keyword string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.Post{})

	// 如果有关键词搜索
	if keyword != "" {
		k := "%" + keyword + "%"
		db = db.Joins("LEFT JOIN users ON users.id = posts.user_id").
			Where("(posts.title LIKE ? OR users.username LIKE ? OR posts.id IN (SELECT post_id FROM post_tags pt JOIN tags t ON t.id = pt.tag_id WHERE t.name LIKE ?))", k, k, k)
	}

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

	// 始终查询点赞数、收藏数和评论数
	db = db.Select("posts.*, " +
		"(SELECT COUNT(*) FROM user_likes WHERE user_likes.post_id = posts.id) as like_count, " +
		"(SELECT COUNT(*) FROM user_favorites WHERE user_favorites.post_id = posts.id) as favorite_count, " +
		"(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count")

	// 排序逻辑
	orderClause := ""
	if userID > 0 {
		// 个人主页：按个人置顶排序
		orderClause = "is_top DESC, "
	} else if orderBy == "created_at" {
		// 首页且按时间排序：按全站置顶排序
		orderClause = "is_system_top DESC, "
	}

	switch orderBy {
	case "views":
		orderClause += "view_count DESC"
	case "likes":
		orderClause += "like_count DESC"
	default:
		orderClause += "posts.created_at DESC"
	}

	if err := db.Preload("User").Preload("Tags").Offset(offset).Limit(pageSize).Order(orderClause).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *PostService) GetHotPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	// 热门文章逻辑：按点赞数排序，如果点赞数相同按评论数排序
	// 这里需要关联查询点赞数和评论数
	err := database.DB.Model(&model.Post{}).
		Select("posts.*, "+
			"(SELECT COUNT(*) FROM user_likes WHERE user_likes.post_id = posts.id) as like_count, "+
			"(SELECT COUNT(*) FROM user_favorites WHERE user_favorites.post_id = posts.id) as favorite_count, "+
			"(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count").
		Where("status = ?", "published").
		Order("like_count DESC, comment_count DESC").
		Limit(limit).
		Preload("User").
		Find(&posts).Error
	return posts, err
}

func (s *PostService) ToggleTop(id string, userId uint, userRole string) error {
	var post model.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		return errors.New("post not found")
	}

	// 只有作者或管理员可以操作个人置顶
	if post.UserID != userId && userRole != "admin" {
		return errors.New("unauthorized")
	}

	return database.DB.Model(&post).Update("is_top", !post.IsTop).Error
}

func (s *PostService) ToggleSystemTop(id string, userRole string) error {
	if userRole != "admin" {
		return errors.New("only admin can toggle system top")
	}

	var post model.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		return errors.New("post not found")
	}

	return database.DB.Model(&post).Update("is_system_top", !post.IsSystemTop).Error
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
	// 预加载 LikedPosts 及其关联的 User 和 Tags 信息，并包含点赞数和评论数
	err := database.DB.Preload("LikedPosts", func(db *gorm.DB) *gorm.DB {
		return db.Select("posts.*, " +
			"(SELECT COUNT(*) FROM user_likes WHERE user_likes.post_id = posts.id) as like_count, " +
			"(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count")
	}).Preload("LikedPosts.User").Preload("LikedPosts.Tags").First(&user, userID).Error
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
	var count int64
	err := database.DB.Table("user_likes").Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// FavoritePost 收藏文章
func (s *PostService) FavoritePost(userID uint, postID uint) error {
	return database.DB.Model(&model.User{Model: gorm.Model{ID: userID}}).
		Association("FavoritePosts").
		Append(&model.Post{Model: gorm.Model{ID: postID}})
}

// UnfavoritePost 取消收藏
func (s *PostService) UnfavoritePost(userID uint, postID uint) error {
	return database.DB.Model(&model.User{Model: gorm.Model{ID: userID}}).
		Association("FavoritePosts").
		Delete(&model.Post{Model: gorm.Model{ID: postID}})
}

// GetFavoritePosts 获取用户收藏的文章列表
func (s *PostService) GetFavoritePosts(userID uint) ([]*model.Post, error) {
	var user model.User
	err := database.DB.Preload("FavoritePosts", func(db *gorm.DB) *gorm.DB {
		return db.Select("posts.*, " +
			"(SELECT COUNT(*) FROM user_likes WHERE user_likes.post_id = posts.id) as like_count, " +
			"(SELECT COUNT(*) FROM user_favorites WHERE user_favorites.post_id = posts.id) as favorite_count, " +
			"(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count")
	}).Preload("FavoritePosts.User").Preload("FavoritePosts.Tags").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.FavoritePosts, nil
}

// IsPostFavorited 检查用户是否收藏了某文章
func (s *PostService) IsPostFavorited(userID uint, postID uint) (bool, error) {
	var count int64
	err := database.DB.Table("user_favorites").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

// GetPostFavoriteCount 获取文章收藏数
func (s *PostService) GetPostFavoriteCount(postID uint) (int64, error) {
	var count int64
	err := database.DB.Table("user_favorites").Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (s *PostService) SearchTags(keyword string, limit int) ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Where("name LIKE ?", "%"+keyword+"%").
		Limit(limit).
		Find(&tags).Error
	return tags, err
}
