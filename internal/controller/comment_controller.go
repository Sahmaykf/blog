package controller

import (
	"net/http"
	"simple-blog/internal/common"
	"simple-blog/internal/database"
	"simple-blog/internal/model"
	"simple-blog/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var commentService = new(service.CommentService)

func CreateComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	comment.UserID = userID.(uint)

	if err := commentService.CreateComment(&comment); err != nil {
		common.Error(c, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// Reload to get user info
	database.DB.Preload("User").First(&comment, comment.ID)

	common.Success(c, comment)
}

func DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if err := commentService.DeleteComment(uint(id), userID.(uint), role.(string)); err != nil {
		common.Error(c, http.StatusForbidden, err.Error())
		return
	}
	common.Success(c, nil)
}

func GetComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	comments, err := commentService.GetCommentsByPostID(uint(postID))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	common.Success(c, comments)
}

func GetMyComments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	comments, err := commentService.GetCommentsByUser(userID.(uint))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, comments)
}

func GetCommentsOnMyPosts(c *gin.Context) {
	userID, _ := c.Get("user_id")
	comments, err := commentService.GetCommentsOnUserPosts(userID.(uint))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, comments)
}
