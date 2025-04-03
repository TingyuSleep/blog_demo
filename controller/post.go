package controller

import (
	"blog_demo/logic"
	"blog_demo/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON error", zap.Error(err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从c中取到当前用户的ID
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedAuth)
		return
	}
	p.AuthorID = userID

	// 2. 创建帖子当前用户的用户ID
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)

}
