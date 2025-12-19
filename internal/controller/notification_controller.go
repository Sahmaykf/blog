package controller

import (
	"net/http"
	"simple-blog/internal/common"
	"simple-blog/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var notificationService = new(service.NotificationService)

func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notifs, err := notificationService.GetUserNotifications(userID.(uint))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, notifs)
}

func MarkNotificationRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notifIDStr := c.Param("id")
	notifID, _ := strconv.Atoi(notifIDStr)

	if err := notificationService.MarkAsRead(userID.(uint), uint(notifID)); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, nil)
}

func GetUnreadNotificationCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count, err := notificationService.GetUnreadCount(userID.(uint))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, count)
}
