package service

import (
	"simple-blog/internal/database"
	"simple-blog/internal/model"
)

type NotificationService struct{}

func (s *NotificationService) CreateNotification(notif *model.Notification) error {
	return database.DB.Create(notif).Error
}

func (s *NotificationService) GetUserNotifications(userID uint) ([]model.Notification, error) {
	var notifs []model.Notification
	err := database.DB.Preload("FromUser").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&notifs).Error
	return notifs, err
}

func (s *NotificationService) MarkAsRead(userID uint, notifID uint) error {
	return database.DB.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true).Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
