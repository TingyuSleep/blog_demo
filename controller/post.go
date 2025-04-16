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
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 2. 创建帖子当前用户的用户ID
	if err := logic.CreatePost(c, p); err != nil {
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

// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口 ，根据发布时间 或者 分数 动态获取帖子列表
func GetPostListHandler2(c *gin.Context) {
	/*
		   1. 获取参数
		   2. 去redis查询id列表
		   3. 根据id去数据库查询帖子详细信息
		get请求参数(query string)：/api/v1/postList2?page=2&size=10&order=time
	*/

	// 1. 获取参数
	// 初始化结构体时,指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	//c.ShouldBind()	根据请求的参数类型选择相应的方法获取数据
	//c.ShouldBindJSON()	如果前端的请求中鞋带的是json格式的数据，才能使用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2.获取数据
	data, err := logic.GetPostList2(c, p)
	if err != nil {
		zap.L().Error("GetPostList2 failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}
