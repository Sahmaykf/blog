package service

import (
	"errors"
	"simple-blog/internal/config"
	"simple-blog/internal/database"
	"simple-blog/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func (s *UserService) Register(username, password, email string) error {
	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	return database.DB.Create(&newUser).Error
}

func (s *UserService) Login(username, password string) (string, uint, string, error) {
	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", 0, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", 0, "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", 0, "", err
	}

	return tokenString, user.ID, user.Role, nil
}

func (s *UserService) FollowUser(followerID, followedID uint) error {
	if followerID == followedID {
		return errors.New("cannot follow yourself")
	}

	var follower model.User
	var followed model.User

	if err := database.DB.First(&follower, followerID).Error; err != nil {
		return err
	}
	if err := database.DB.First(&followed, followedID).Error; err != nil {
		return err
	}

	return database.DB.Model(&follower).Association("Following").Append(&followed)
}

func (s *UserService) UnfollowUser(followerID, followedID uint) error {
	var follower model.User
	var followed model.User

	if err := database.DB.First(&follower, followerID).Error; err != nil {
		return err
	}
	if err := database.DB.First(&followed, followedID).Error; err != nil {
		return err
	}

	return database.DB.Model(&follower).Association("Following").Delete(&followed)
}

func (s *UserService) GetFollowers(userID uint) ([]model.User, error) {
	var user model.User
	if err := database.DB.Preload("Followers").First(&user, userID).Error; err != nil {
		return nil, err
	}
	// Sanitize passwords
	for i := range user.Followers {
		user.Followers[i].Password = ""
	}
	// Convert []*User to []model.User
	followers := make([]model.User, len(user.Followers))
	for i, u := range user.Followers {
		followers[i] = *u
	}
	return followers, nil
}

func (s *UserService) GetFollowing(userID uint) ([]model.User, error) {
	var user model.User
	if err := database.DB.Preload("Following").First(&user, userID).Error; err != nil {
		return nil, err
	}
	// Sanitize passwords
	for i := range user.Following {
		user.Following[i].Password = ""
	}
	following := make([]model.User, len(user.Following))
	for i, u := range user.Following {
		following[i] = *u
	}
	return following, nil
}

func (s *UserService) GetUserProfile(targetUserID uint, currentUserID uint) (map[string]interface{}, error) {
	var user model.User
	if err := database.DB.First(&user, targetUserID).Error; err != nil {
		return nil, err
	}

	var postCount int64
	database.DB.Model(&model.Post{}).Where("user_id = ?", targetUserID).Count(&postCount)

	followerCount := database.DB.Model(&user).Association("Followers").Count()
	followingCount := database.DB.Model(&user).Association("Following").Count()

	isFollowing := false
	if currentUserID != 0 {
		var c int64
		database.DB.Table("user_followers").
			Where("follower_id = ? AND followed_id = ?", currentUserID, targetUserID).
			Count(&c)
		isFollowing = c > 0
	}

	return map[string]interface{}{
		"id":              user.ID,
		"username":        user.Username,
		"email":           user.Email,
		"avatar":          user.Avatar,
		"role":            user.Role,
		"created_at":      user.CreatedAt,
		"post_count":      postCount,
		"follower_count":  followerCount,
		"following_count": followingCount,
		"is_following":    isFollowing,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, email string, avatar string) error {
	updates := map[string]interface{}{
		"email": email,
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) UpdateAvatar(userID uint, avatarPath string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarPath).Error
}
