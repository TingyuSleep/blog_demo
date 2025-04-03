package controller

import (
	"blog_demo/logic"
	"blog_demo/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
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

// GetPostDetailHandler 获取帖子详细信息
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数及参数校验
	//帖子id从url中传进来的,用c.Param方法获取
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2.根据id查询数据库
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, data)
}
