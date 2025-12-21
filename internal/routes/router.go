package routes

import (
	"simple-blog/internal/config"
	"simple-blog/internal/controller"
	"simple-blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回 Gin 引擎
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 增加文件上传大小限制 (默认 32MB，这里设为 100MB)
	r.MaxMultipartMemory = 100 << 20

	// 静态资源
	cfg := config.GetAppConfig()
	r.Static("/uploads", cfg.UploadDir)

	// API 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开路由
		v1.POST("/register", controller.Register)
		v1.POST("/login", controller.Login)
		v1.GET("/posts", middleware.SoftJWTAuth(), controller.GetPostList)
		v1.GET("/posts/hot", controller.GetHotPosts)
		v1.GET("/posts/:id", controller.GetPostDetail)
		v1.GET("/tags", controller.GetTags)
		searchCtrl := controller.NewSearchController()
		v1.GET("/search", searchCtrl.GlobalSearch)

		// Comments (Public Read)
		v1.GET("/posts/:id/comments", controller.GetComments)
		v1.GET("/posts/:id/like", middleware.SoftJWTAuth(), controller.GetPostLikeStatus)

		// User Profile & Relations (Public Read)
		v1.GET("/users/:id/followers", controller.GetFollowers)
		v1.GET("/users/:id/following", controller.GetFollowing)
		v1.GET("/users/:id/liked-posts", controller.GetLikedPosts)
		v1.GET("/users/:id/favorite-posts", controller.GetFavoritePosts)
		v1.GET("/users/:id/posts", controller.GetUserPosts)
		v1.GET("/users/:id", middleware.SoftJWTAuth(), controller.GetUserProfile)

		// 需要认证的路由组
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			auth.POST("/posts", controller.CreatePost)
			auth.PUT("/posts/:id", controller.UpdatePost)
			auth.DELETE("/posts/:id", controller.DeletePost)
			auth.POST("/posts/:id/like", controller.ToggleLike)
			auth.POST("/posts/:id/favorite", controller.ToggleFavorite)
			auth.POST("/posts/:id/top", controller.ToggleTop)
			auth.POST("/posts/:id/system-top", controller.ToggleSystemTop)

			// Comments
			auth.POST("/comments", controller.CreateComment)
			auth.DELETE("/comments/:id", controller.DeleteComment)
			auth.GET("/my/comments", controller.GetMyComments)
			auth.GET("/my/post-comments", controller.GetCommentsOnMyPosts)

			// Notifications
			auth.GET("/notifications", controller.GetNotifications)
			auth.GET("/notifications/unread-count", controller.GetUnreadNotificationCount)
			auth.PUT("/notifications/:id/read", controller.MarkNotificationRead)

			// User Relations
			auth.POST("/users/:id/follow", controller.FollowUser)
			auth.POST("/users/:id/unfollow", controller.UnfollowUser)

			// User Profile Update
			auth.PUT("/user/profile", controller.UpdateProfile)
			auth.POST("/user/avatar", controller.UploadAvatar)
			auth.POST("/upload/image", controller.UploadImage)
		}
	}

	return r
}
