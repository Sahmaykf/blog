package service

import (
	"errors"
	"simple-blog/internal/database"
	"simple-blog/internal/model"
)

type CommentService struct{}

func (s *CommentService) CreateComment(comment *model.Comment) error {
	if err := database.DB.Create(comment).Error; err != nil {
		return err
	}

	// 创建通知给文章作者
	var post model.Post
	if err := database.DB.First(&post, comment.PostID).Error; err == nil {
		// 如果评论者不是作者本人，则发送通知
		if post.UserID != comment.UserID {
			notification := &model.Notification{
				UserID:     post.UserID,
				Type:       "comment",
				Content:    comment.Content,
				FromUserID: comment.UserID,
				PostID:     comment.PostID,
				CommentID:  comment.ID,
			}
			database.DB.Create(notification)
		}
	}

	return nil
}

func (s *CommentService) GetCommentsOnUserPosts(userID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := database.DB.Joins("JOIN posts ON posts.id = comments.post_id").
		Where("posts.user_id = ?", userID).
		Preload("User").
		Preload("Post").
		Order("comments.created_at desc").
		Find(&comments).Error
	return comments, err
}

func (s *CommentService) GetCommentsByUser(userID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := database.DB.Where("user_id = ?", userID).
		Preload("Post").
		Order("created_at desc").
		Find(&comments).Error
	return comments, err
}

func (s *CommentService) DeleteComment(commentID uint, userID uint, userRole string) error {
	var comment model.Comment
	// Preload Post to check post author
	if err := database.DB.Preload("Post").First(&comment, commentID).Error; err != nil {
		return err
	}

	// Logic: Allow delete if userID matches comment author OR if userID matches the post author OR if user is admin
	if userRole == "admin" || comment.UserID == userID || (comment.Post != nil && comment.Post.UserID == userID) {
		return database.DB.Delete(&comment).Error
	}

	return errors.New("unauthorized to delete this comment")
}

func (s *CommentService) GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	// Preload User to show who commented
	err := database.DB.Where("post_id = ?", postID).Preload("User").Order("created_at desc").Find(&comments).Error
	return comments, err
}
