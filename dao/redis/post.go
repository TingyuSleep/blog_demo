package redis

import (
	"blog_demo/models"

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
