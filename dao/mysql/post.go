package mysql

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"web_app/models"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) error {
	sqlStr := `insert into post(post_id, author_id, title, content, community_id)
		values(?, ?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, p.PostID, p.AuthorID, p.Title, p.Content, p.CommunityID)
	return err
}

// GetPostByID 查询帖子返回数据
func GetPostByID(postID int64) (data *models.Post, err error) {
	data = new(models.Post) // 必要的一步，不然是空指针无法传入db.Get
	sqlStr := `select post_id, author_id, title, content, community_id, create_time 
		from post where post_id=?`
	err = db.Get(data, sqlStr, postID)
	return
}

// GetPostList 查询帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	// 最新的帖子优先显示
	sqlStr := `select
	post_id, author_id, title, content, community_id, create_time
	from post
	order by create_time
	desc
	limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size) //取那一页，每一页多少个帖子
	return posts, err
}

// GetPostListByIDs 通过postID列表来查询帖子详细信息
func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`
	// 动态填充id（https://www.liwenzhou.com/posts/Go/sqlx/#autoid-0-4-1）
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
