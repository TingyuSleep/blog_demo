package logic

import (
	"blog_demo/dao/redis"
	"blog_demo/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// 投票功能

// VoteForPost 为帖子投票
func VoteForPost(c *gin.Context, userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(c, strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
