package models

import "time"

//
// type Community struct {
// 	ID   int64  `json:"id,string" db:"community_id" gorm:"primary_key not null"`
// 	Name string `json:"name" db:"community_name"`
// }

type Community struct {
	ID            int64     `gorm:"primary_key;type:int(11) AUTO_INCREMENT;comment:'id'"`
	CommunityID   int64     `json:"community_id" gorm:"not null;unique;type:int(11) unsigned;comment:'社区ID'"`
	CommunityName string    `json:"community_name" gorm:"not null;unique;size:128;comment:'社区名字'"`
	Introduction  string    `json:"introduction" gorm:"not null;size:256;comment:'社区介绍'"`
	CreateTime    time.Time `gorm:"type:timestamp;not null;comment:'创建时间'"`
	UpdateTime    time.Time `gorm:"type:timestamp;not null;comment:'修改时间'"`
}

// 在 json tag后面加上string，可以让其传到前端时自动转换成string类型
type CommunityDetail struct {
	ID          int64     `json:"id,string" db:"community_id"`
	Name        string    `json:"name" db:"community_name"`
	Instruction string    `json:"introduction,omitempty" db:"introduction"` // omitempty表示如果为空，那么json字段就没必要展示出来
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}
