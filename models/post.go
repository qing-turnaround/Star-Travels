package models

import "time"

// 内存对齐（cpu每次取8个字节的块大小）
//json string解决传给前端数字id失真问题，js数字最大值为1<<53-1,而Go确实1<<63-1

type Post struct {
	PostID      int64  `json:"post_id,string" db:"post_id"`     // 帖子ID，无须传入
	AuthorID    int64  `json:"author_id,string" db:"author_id"` // 作者ID
	Status      int32  `json:"status" db:"status"`
	CommunityID int64  `json:"community_id" db:"community_id" binding:"required"` // 社区ID
	Title       string `json:"title" db:"title" binding:"required"`               // 标题
	Content     string `json:"content" db:"content" binding:"required"`           // 内容

	CreateTime time.Time `json:"create_time" db:"create_time"` //创建时间
}

// ApiPostDetails 帖子详情接口的结构体
type ApiPostDetails struct {
	*Post
	VoteNum       int64  `json:"vote_num"` // 帖子的投票数量
	AuthorName    string `json:"author_name" db:"author_name"`
	CommunityName string `json:"community_name" db:"community_name"`
}
