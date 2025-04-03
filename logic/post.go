package logic

import (
	"blog_demo/dao/mysql"
	"blog_demo/models"
	"blog_demo/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()

	// 2. 保存进数据库
	return mysql.CreatePost(p)
}
