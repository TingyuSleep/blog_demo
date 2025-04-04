package logic

import (
	"blog_demo/dao/mysql"
	"blog_demo/models"
	"blog_demo/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()

	// 2. 保存进数据库
	return mysql.CreatePost(p)
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	//data = new(models.ApiPostDetail) //如果指针未初始化，会由于空指针引用问题而报错
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("GetPostByID failed", zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}

	// 根据社区id查询社区详情
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("GetCommunityDetailByID failed",
			zap.Int64("community_id",
				post.CommunityID), zap.Error(err))
		return
	}
	// 初始化指针变量data
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	return
}
