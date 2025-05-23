package models

import "time"

// 社区分类结构体
type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

// 社区详情结构体
type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"` //omitempty,如果内容为空，字段不展示
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
