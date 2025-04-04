package mysql

import (
	"blog_demo/models"
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
