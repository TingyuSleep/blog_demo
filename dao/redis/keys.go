package redis

// redis key
// 注意使用命名空间
const (
	Prefix          = "blog:"
	KeyPostTime     = "post:time"   //zset，帖子及发帖时间
	KeyPostScore    = "post:score"  //zset，帖子及投票的分数
	KeyPostVoteType = "post:voted:" //zset，记录用户及投票类型（赞同/反对/未投票），参数是post id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
