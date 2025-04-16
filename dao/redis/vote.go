package redis

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

// 本项目使用简化版的投票分数，投一票加432分 86400/200  -> 200张赞成票可以给你的帖子续一天  ->《redis实战》

/*
投票的几种情况：
direction=1时，有两种情况：
	1.之前没有投票，现在投赞成票	----> 更新分数和投票记录	差值的绝对值：1	+432
	2.之前投反对票，现在投赞成票	----> 更新分数和投票记录	差值的绝对值：2	+432*2
direction=0是，有两种情况：
	1.之前投赞成票，现在取消投票	----> 更新分数和投票记录	差值的绝对值：1	-432
	2.之前投反对票，现在取消投票	----> 更新分数和投票记录	差值的绝对值：1	+432
direction=-1是，有两种情况：
	1.之前投赞成票，现在投反对票	----> 更新分数和投票记录	差值的绝对值：2	-432*2
	2.之前没有投票，现在投反对票	----> 更新分数和投票记录	差值的绝对值：1	-432

投票的限制：
每个帖子自发表之日起，一个星期之内允许用户投票，超过一个星期就不允许投票了
	1.到期之后，将redis中保存的赞成票数和反对票数存储到mysql表中
	2.到期之后，删除 KeyPostVoteType

*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票的分数
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
)

func VoteForPost(c *gin.Context, userID, postID string, direction float64) (err error) {
	// 1.判断投票限制，如果过期，则不能投票
	// 去redis获取发帖时间
	postTime := rdb.ZScore(c, getRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}

	// 注意，2和3要放到一个pipeline事务中执行

	// 2.更新帖子的分数
	// 先查当前用户给当前帖子的投票记录，拿到以前投票（赞成1/反对-1/未投票0）
	ov, err := rdb.ZScore(c, getRedisKey(KeyPostVoteType+postID), userID).Result()

	//计算两次投票差值（可以保证用户不能重复投票）
	diff := direction - ov
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(c, getRedisKey(KeyPostScore), diff*scorePerVote, postID)

	// 3.记录用户为该帖子投票的数据
	if direction == 0 { //如果为0，则删除这条记录
		pipeline.ZRem(c, getRedisKey(KeyPostVoteType+postID), userID)
	}

	pipeline.ZAdd(c, getRedisKey(KeyPostVoteType+postID), redis.Z{
		Score:  direction,
		Member: userID,
	})
	_, err = pipeline.Exec(c)
	return err
}

// CreatePost 在redis中记录发帖时间以及分数
func CreatePost(c *gin.Context, postID int64) error {
	pipeline := rdb.TxPipeline() //事务
	// 发帖时间
	pipeline.ZAdd(c, getRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//// 帖子分数
	//pipeline.ZAdd(c, getRedisKey(KeyPostScore), redis.Z{
	//	Score:  float64(time.Now().Unix()), // 不是很明白?
	//	Member: postID,
	//})
	_, err := pipeline.Exec(c)
	return err
}
