package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"simple-blog/internal/common"
	"simple-blog/internal/config"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// 增加恢复机制，防止 panic 导致连接重置
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic recovered in UploadImage: %v\n", err)
			common.Error(c, http.StatusInternalServerError, "服务器内部错误")
		}
	}()

	val, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	var userID uint
	switch v := val.(type) {
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		common.Error(c, http.StatusInternalServerError, fmt.Sprintf("用户信息类型异常: %T", val))
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		fmt.Printf("FormFile error: %v\n", err)
		common.Error(c, http.StatusBadRequest, "未接收到图片文件或文件过大")
		return
	}

	// Create directory if not exists
	cfg := config.GetAppConfig()
	uploadDir := filepath.Join(cfg.UploadDir, "posts")

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory: %v\n", err)
			common.Error(c, http.StatusInternalServerError, fmt.Sprintf("创建上传目录失败: %v", err))
			return
		}
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		common.Error(c, http.StatusInternalServerError, fmt.Sprintf("保存文件失败: %v", err))
		return
	}

	// Return the URL
	imageURL := "/" + filepath.ToSlash(filePath)
	fmt.Printf("Upload success: %s\n", imageURL)
	common.Success(c, gin.H{"url": imageURL})
}
