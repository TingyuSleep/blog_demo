package mysql

import (
	"blog_demo/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id,title,content,author_id,community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id=?`
	err = db.Get(post, sqlStr, pid)
	return
}

func GetPostList(page, size int64) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time 
	from post limit ?,?`
	postList = make([]*models.Post, 0, 2)
	err = db.Select(&postList, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time 
	from post 
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}

	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return
}
