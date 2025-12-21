package controller

import (
	"net/http"
	"simple-blog/internal/common"
	"simple-blog/internal/model"
	"simple-blog/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var postService = new(service.PostService)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var post model.Post

	// 1. 绑定 JSON 参数到结构体
	if err := c.ShouldBindJSON(&post); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 从上下文获取 UserID (由 JWT 中间件设置)
	userId, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "User not found in context")
		return
	}
	post.UserID = userId.(uint)

	// 2. 调用 Service
	if err := postService.CreatePost(&post); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. 返回结果
	common.Success(c, post)
}

// GetPostDetail 获取文章详情
func GetPostDetail(c *gin.Context) {
	id := c.Param("id")

	post, err := postService.GetPostDetail(id)
	if err != nil {
		common.Error(c, http.StatusNotFound, "Post not found")
		return
	}

	common.Success(c, post)
}

// GetPostList 获取文章列表（分页）
func GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	tagName := c.Query("tag")
	userID, _ := strconv.Atoi(c.Query("user_id"))
	status := c.DefaultQuery("status", "published")
	orderBy := c.DefaultQuery("order_by", "created_at")
	keyword := c.Query("keyword")

	// 如果是管理员或者查询自己的文章，可以查看所有状态
	currentUserID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	if status == "all" && (userRole == "admin" || (currentUserID != nil && currentUserID.(uint) == uint(userID))) {
		status = "" // 不过滤状态
	}

	posts, total, err := postService.GetPostList(page, pageSize, tagName, uint(userID), status, orderBy, keyword)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.Success(c, gin.H{
		"list": posts,
		"meta": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total":        total,
		},
	})
}

// GetHotPosts 获取热门文章
func GetHotPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	posts, err := postService.GetHotPosts(limit)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, posts)
}

// GetUserPosts 获取指定用户的文章列表
func GetUserPosts(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 个人主页只显示已发布的文章
	posts, _, err := postService.GetPostList(page, pageSize, "", uint(userID), "published", "created_at", "")
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.Success(c, posts)
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	userId, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "User not found in context")
		return
	}

	userRole, _ := c.Get("user_role")

	var input model.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		common.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := postService.UpdatePost(id, userId.(uint), userRole.(string), &input); err != nil {
		if err.Error() == "unauthorized" {
			common.Error(c, http.StatusForbidden, "You are not authorized to update this post")
		} else if err.Error() == "post not found" {
			common.Error(c, http.StatusNotFound, "Post not found")
		} else {
			common.Error(c, http.StatusInternalServerError, "Failed to update post")
		}
		return
	}

	common.Success(c, nil)
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")

	userId, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "User not found in context")
		return
	}
	userRole, _ := c.Get("user_role")

	if err := postService.DeletePost(id, userId.(uint), userRole.(string)); err != nil {
		if err.Error() == "unauthorized" {
			common.Error(c, http.StatusForbidden, "You are not authorized to delete this post")
		} else if err.Error() == "post not found" {
			common.Error(c, http.StatusNotFound, "Post not found")
		} else {
			common.Error(c, http.StatusInternalServerError, "Failed to delete post")
		}
		return
	}

	common.Success(c, "Post deleted successfully")
}

// GetTags 获取所有标签
func GetTags(c *gin.Context) {
	tags, err := postService.GetAllTags()
	if err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, tags)
}

// ToggleLike 处理点赞/取消点赞
func ToggleLike(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.Atoi(postIDStr)
	userID, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 简单的逻辑：先检查是否已点赞，反转状态
	isLiked, _ := postService.IsPostLiked(userID.(uint), uint(postID))

	var err error
	if isLiked {
		err = postService.UnlikePost(userID.(uint), uint(postID))
	} else {
		err = postService.LikePost(userID.(uint), uint(postID))
	}

	if err != nil {
		common.Error(c, http.StatusInternalServerError, "Operation failed")
		return
	}

	common.Success(c, !isLiked)
}

// GetLikedPosts 获取某用户的点赞列表
func GetLikedPosts(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	posts, err := postService.GetLikedPosts(uint(userID))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, "Failed to fetch liked posts")
		return
	}
	common.Success(c, posts)
}

// GetPostLikeStatus 获取文章的点赞状态和数量
func GetPostLikeStatus(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.Atoi(postIDStr)

	likeCount, _ := postService.GetPostLikeCount(uint(postID))
	favoriteCount, _ := postService.GetPostFavoriteCount(uint(postID))

	isLiked := false
	isFavorited := false
	// 如果用户登录了，检查是否点赞和收藏
	if userID, exists := c.Get("user_id"); exists {
		isLiked, _ = postService.IsPostLiked(userID.(uint), uint(postID))
		isFavorited, _ = postService.IsPostFavorited(userID.(uint), uint(postID))
	}

	common.Success(c, gin.H{
		"count":          likeCount,
		"is_liked":       isLiked,
		"favorite_count": favoriteCount,
		"is_favorited":   isFavorited,
	})
}

// ToggleFavorite 处理收藏/取消收藏
func ToggleFavorite(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.Atoi(postIDStr)
	userID, exists := c.Get("user_id")
	if !exists {
		common.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	isFavorited, _ := postService.IsPostFavorited(userID.(uint), uint(postID))

	var err error
	if isFavorited {
		err = postService.UnfavoritePost(userID.(uint), uint(postID))
	} else {
		err = postService.FavoritePost(userID.(uint), uint(postID))
	}

	if err != nil {
		common.Error(c, http.StatusInternalServerError, "Operation failed")
		return
	}

	common.Success(c, !isFavorited)
}

// GetFavoritePosts 获取某用户的收藏列表
func GetFavoritePosts(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		common.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	posts, err := postService.GetFavoritePosts(uint(userID))
	if err != nil {
		common.Error(c, http.StatusInternalServerError, "Failed to fetch favorite posts")
		return
	}
	common.Success(c, posts)
}

// ToggleTop 切换个人置顶状态
func ToggleTop(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	if err := postService.ToggleTop(id, userID.(uint), userRole.(string)); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.Success(c, nil)
}

// ToggleSystemTop 切换全站置顶状态（仅管理员）
func ToggleSystemTop(c *gin.Context) {
	id := c.Param("id")
	userRole, _ := c.Get("user_role")

	if err := postService.ToggleSystemTop(id, userRole.(string)); err != nil {
		common.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.Success(c, nil)
}
