package models

import "time"

// 内存对齐概念（相同的类型放在一起定义）
type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id" ` // author_id 是通过 GetCurrentUser 方法传递的
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	VoteScore        int64              `json:"vote_score"`
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区分类信息
}
