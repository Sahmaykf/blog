package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"simple-blog/internal/common"
	"simple-blog/internal/config"
	"simple-blog/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var userService = new(service.UserService)

// Register 用户注册
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := userService.Register(input.Username, input.Password, input.Email); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	common.Success(c, gin.H{"message": "Registration successful"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, userId, role, err := userService.Login(input.Username, input.Password)
	if err != nil {
		common.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	common.Success(c, gin.H{
		"token":   token,
		"user_id": userId,
		"role":    role,
	})
}

func FollowUser(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	currentUserID, _ := c.Get("user_id")

	if err := userService.FollowUser(currentUserID.(uint), uint(targetID)); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, nil)
}

func UnfollowUser(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	currentUserID, _ := c.Get("user_id")

	if err := userService.UnfollowUser(currentUserID.(uint), uint(targetID)); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, nil)
}

func GetFollowers(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	followers, err := userService.GetFollowers(uint(targetID))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, followers)
}

func GetFollowing(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	following, err := userService.GetFollowing(uint(targetID))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, following)
}

func GetUserProfile(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Optional Auth: Check if user is logged in to determine "IsFollowing"
	var currentUserID uint = 0
	if val, exists := c.Get("user_id"); exists {
		currentUserID = val.(uint)
	}

	profile, err := userService.GetUserProfile(uint(targetID), currentUserID)
	if err != nil {
		common.Error(c, http.StatusNotFound, "User not found")
		return
	}
	common.Success(c, profile)
}

func UpdateProfile(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Avatar   string `json:"avatar"`
		BlogName string `json:"blog_name"`
		Bio      string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	if err := userService.UpdateProfile(userID.(uint), input.Email, input.Avatar, input.BlogName, input.Bio); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.Success(c, nil)
}

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		common.Error(c, http.StatusBadRequest, "No file uploaded")
		return
	}

	userID, _ := c.Get("user_id")

	// Create directory if not exists
	cfg := config.GetAppConfig()
	uploadDir := filepath.Join(cfg.UploadDir, "avatars")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID.(uint), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		common.Error(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Return the URL but don't update DB yet
	// Ensure the URL starts with /uploads/
	avatarURL := "/" + filepath.ToSlash(filePath)
	common.Success(c, gin.H{"avatar_url": avatarURL})
}
