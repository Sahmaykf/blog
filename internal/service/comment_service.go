package service

import (
	"errors"
	"simple-blog/internal/database"
	"simple-blog/internal/model"

	"gorm.io/gorm"
)

type CommentService struct{}

func (s *CommentService) CreateComment(comment *model.Comment) error {
	if err := database.DB.Create(comment).Error; err != nil {
		return err
	}

	// 如果是回复评论
	if comment.ParentID != nil {
		var targetComment model.Comment
		// 查找被回复的评论
		// targetID := comment.ID
		var targetID uint
		if comment.ReplyToID != nil {
			targetID = *comment.ReplyToID
		} else {
			targetID = *comment.ParentID
		}

		if err := database.DB.Preload("User").First(&targetComment, targetID).Error; err == nil {
			// 通知被回复的人
			if targetComment.UserID != comment.UserID {
				notification := &model.Notification{
					UserID:     targetComment.UserID,
					Type:       "reply",
					Content:    comment.Content,
					FromUserID: comment.UserID,
					PostID:     comment.PostID,
					CommentID:  comment.ID,
				}
				database.DB.Create(notification)
			}
		}
	} else {
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
	// 只获取顶级评论（ParentID 为 nil），并预加载回复及其作者
	err := database.DB.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at asc")
		}).
		Preload("Replies.User").
		Order("created_at desc").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	// 手动填充回复的目标用户信息
	for i := range comments {
		for j := range comments[i].Replies {
			if comments[i].Replies[j].ReplyToID != nil {
				var replyToComment model.Comment
				if err := database.DB.Preload("User").First(&replyToComment, *comments[i].Replies[j].ReplyToID).Error; err == nil {
					comments[i].Replies[j].ReplyToUser = &replyToComment.User
				}
			}
		}
	}

	return comments, nil
}
