package middleware

import (
	"fmt"
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
			fmt.Println("JWTAuth: Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("JWTAuth: Invalid header format: %s\n", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("JWTAuth: Token error: %v\n", err)
			// 区分过期和其他错误
			msg := "Invalid or expired token"
			if err != nil && strings.Contains(err.Error(), "expired") {
				msg = "token_expired"
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
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
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
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
