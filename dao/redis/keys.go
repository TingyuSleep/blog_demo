package redis

// redis key
// 注意使用命名空间
const (
	Prefix          = "blog:"
	KeyPostTime     = "post:time"   //zset，存储帖子id及发帖时间
	KeyPostScore    = "post:score"  //zset，存储帖子id及投票的分数
	KeyPostVoteType = "post:voted:" //zset，存储用户id及投票类型（赞同/反对/未投票），参数是post id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
