package middleware

import (
	"net/http"
	"simple-blog/internal/config"
	"simple-blog/internal/database"
	"simple-blog/internal/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		//将UserID解析出来并存入上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userIdFloat, ok := claims["user_id"].(float64); ok {
				userId := uint(userIdFloat)
				c.Set("user_id", userId)
				c.Set("userID", userId)

				// 获取用户角色
				var user model.User
				if err := database.DB.First(&user, userId).Error; err == nil {
					c.Set("user_role", user.Role)
					c.Set("role", user.Role) // Alias for compatibility
				}
			}
		}

		c.Next()
	}
}

func SoftJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return config.JwtSecret, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userIdFloat, ok := claims["user_id"].(float64); ok {
						userId := uint(userIdFloat)
						c.Set("user_id", userId)
						c.Set("userID", userId)

						// 获取用户角色
						var user model.User
						if err := database.DB.First(&user, userId).Error; err == nil {
							c.Set("user_role", user.Role)
							c.Set("role", user.Role)
						}
					}
				}
			}
		}
		c.Next()
	}
}
