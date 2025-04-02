package logic

import (
	"blog_demo/dao/mysql"
	"blog_demo/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	// 查询数据库 找到所有community 并返回
	return mysql.GetCommunityList()
}
