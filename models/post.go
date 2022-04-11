package models

import "time"

// 内存对齐（cpu每次取8个字节的块大小）
//json string解决传给前端数字id失真问题，js数字最大值为1<<53-1,而Go确实1<<63-1
type Post struct {
	ID          int       `gorm:"primary_key;type:int(11) AUTO_INCREMENT;comment:'id'"`
	PostID      int64     `json:"post_id,string" gorm:"unique;type:bigint(20) NOT NULL;comment:'帖子id'"`                 // 帖子ID，无须传入
	AuthorID    int64     `json:"author_id,string" gorm:"index;type:bigint(20) NOT NULL;comment:'作者的用户id'"`             // 作者ID
	CommunityID int64     `json:"community_id" binding:"required" gorm:"index;type:bigint(20) NOT NULL;comment:'所属社区'"` // 社区ID
	Status      int32     `json:"status" gorm:"type:tinyint(4) NOT NULL DEFAULT '1';comment:'帖子状态'"`
	Title       string    `json:"title" binding:"required" gorm:"type:varchar(128) COLLATE utf8mb4_general_ci NOT NULL;comment:'帖子的标题'"`    // 标题
	Content     string    `json:"content" binding:"required" gorm:"type:varchar(8192) COLLATE utf8mb4_general_ci NOT NULL;comment:'帖子的内容'"` // 内容
	CreateTime  time.Time `json:"create_time" gorm:"type:timestamp;not null;comment:'创建时间'"`                                                //创建时间
	UpdateTime  time.Time `json:"update_time" gorm:"type:timestamp;not null;comment:'修改时间'"`
}

// ApiPostDetails 帖子详情接口的结构体
type ApiPostDetails struct {
	*Post
	VoteNum       int64  `json:"vote_num"` // 帖子的投票数量
	AuthorName    string `json:"author_name" db:"author_name"`
	CommunityName string `json:"community_name" db:"community_name"`
}
