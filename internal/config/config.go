package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// 尝试加载 .env 文件，如果不存在则忽略
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

// JwtSecret 用于 JWT 签名和验证
var JwtSecret = []byte(getEnv("JWT_SECRET", "D7A91B3C-5E2F-48A0-9C1D-6E8F2B4A0C3D-STRONG-KEY-2025"))

// DBConfig 数据库配置
type DBConfig struct {
	DSN string
}

// GetDBConfig 获取数据库配置
func GetDBConfig() DBConfig {
	// 默认本地开发配置
	defaultDSN := "root:Hyfhgz87370376.@tcp(127.0.0.1:3306)/blog_system?charset=utf8mb4&parseTime=True&loc=Local"
	return DBConfig{
		DSN: getEnv("DB_DSN", defaultDSN),
	}
}

// AppConfig 应用配置
type AppConfig struct {
	Port      string
	UploadDir string
}

// GetAppConfig 获取应用配置
func GetAppConfig() AppConfig {
	return AppConfig{
		Port:      getEnv("APP_PORT", "8080"),
		UploadDir: getEnv("UPLOAD_DIR", "uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
