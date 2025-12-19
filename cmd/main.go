package main

import (
	"fmt"
	"simple-blog/internal/config"
	"simple-blog/internal/database"
	"simple-blog/internal/model"
	"simple-blog/internal/routes"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. 初始化数据库
	database.InitDB()

	// 初始化管理员账号
	initAdmin()

	// 2. 初始化路由
	r := routes.SetupRouter()

	// 3. 启动服务
	cfg := config.GetAppConfig()
	fmt.Printf("Server is running on port %s...\n", cfg.Port)
	r.Run(":" + cfg.Port)
}

func initAdmin() {
	var admin model.User
	if err := database.DB.Where("role = ?", "admin").First(&admin).Error; err != nil {
		// 创建默认管理员
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin = model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Role:     "admin",
		}
		if err := database.DB.Create(&admin).Error; err == nil {
			fmt.Println("Created default admin user: admin / admin123")
		} else {
			fmt.Printf("Failed to create admin user: %v\n", err)
		}
	}
}
