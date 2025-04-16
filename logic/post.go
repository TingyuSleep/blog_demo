package logic

import (
	"blog_demo/dao/mysql"
	"blog_demo/dao/redis"
	"blog_demo/models"
	"blog_demo/pkg/snowflake"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(c *gin.Context, p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()

	// 2. 保存进数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(c, p.ID) //把帖子发布时间记录到redis数据库中
	return
}

// GetPostByID 根据id查询单个帖子的详情数据
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
	// 初始化指针变量data,拼接接口数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	//初始化data
	data = make([]*models.ApiPostDetail, 0, len(postList))

	// 查询多个帖子信息
	for _, post := range postList {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		// 根据社区id查询社区详情
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID failed",
				zap.Int64("community_id",
					post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(c *gin.Context, p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2.去redis中查询id列表
	ids, err := redis.GetPostIdInOrder(c, p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIdInOrder return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))

	// 3.根据id去mysql中查询帖子详细信息
	// 从redis中按照什么顺序拿到的id，在mysql中还要按照该顺序返回对应的数据
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("postList", postList))
	// 提前查询每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(c, ids)
	if err != nil {
		return
	}

	// 查询多个帖子信息
	for idx, post := range postList {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		// 根据社区id查询社区详情
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("GetCommunityDetailByID failed",
				zap.Int64("community_id",
					post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteScore:       voteData[idx],
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}
	return
}
