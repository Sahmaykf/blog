package database

import (
	"fmt"
	"log"
	"simple-blog/internal/config"
	"simple-blog/internal/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// 从配置获取 DSN
	cfg := config.GetDBConfig()
	dsn := cfg.DSN

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("连接数据库失败: ", err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("获取数据库失败: ", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	fmt.Println("MySQL started successfully.")
	err = DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Tag{}, &model.Comment{}, &model.Notification{})
	if err != nil {
		log.Fatal("数据库迁移失败: ", err)
	}
}
