package redis

import (
	"blog_demo/models"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func GetPostIdInOrder(c *gin.Context, p *models.ParamPostList) ([]string, error) {
	// 从redis获取帖子id
	// 1.根据用户请求中的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTime)
	// 根据时间排序
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}
	// 2.确定帖子每一页列表索引起止位置
	start := (p.Page - 1) * p.Size
	stop := start + p.Size - 1
	// 3.查询帖子id, 按照 时间/分数 从大到小获取某一页中指定数量的帖子id
	return rdb.ZRevRange(c, key, start, stop).Result() // 返回值为[]string类型的切片和err
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(c *gin.Context, ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVoteType + id)
	//	//ZCOUNT key min max：统计score值在给定范围内的所有元素的个数
	//	val := rdb.ZCount(c, key, "1", "1").Val()
	//	data = append(data, val)
	//}

	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVoteType + id)
		//ZCOUNT key min max：统计score值在给定范围内的所有元素的个数
		pipeline.ZCount(c, key, "1", "1")
	}
	cmders, err := pipeline.Exec(c)
	if err != nil {
		return
	}
	for _, cmder := range cmders {
		val := cmder.(*redis.IntCmd).Val()
		data = append(data, val)
	}
	return
}
