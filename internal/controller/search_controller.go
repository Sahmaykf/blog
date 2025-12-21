package controller

import (
	"net/http"
	"simple-blog/internal/common"
	"simple-blog/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	postService *service.PostService
	userService *service.UserService
}

func NewSearchController() *SearchController {
	return &SearchController{
		postService: &service.PostService{},
		userService: &service.UserService{},
	}
}

func (ctrl *SearchController) GlobalSearch(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		common.Error(c, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	// 搜索文章 (前10条)
	posts, _, err := ctrl.postService.GetPostList(1, 10, "", 0, "published", "created_at", keyword)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, "搜索文章失败")
		return
	}

	// 搜索用户 (前10条)
	users, err := ctrl.userService.SearchUsers(keyword, 10)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, "搜索用户失败")
		return
	}

	// 搜索标签 (前20条)
	tags, err := ctrl.postService.SearchTags(keyword, 20)
	if err != nil {
		common.Error(c, http.StatusInternalServerError, "搜索标签失败")
		return
	}

	common.Success(c, gin.H{
		"posts": posts,
		"users": users,
		"tags":  tags,
	})
}
