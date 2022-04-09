package models

import "time"

type Community struct {
	ID   int64  `json:"id,string" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"id,string" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Instructions string    `json:"introduction,omitempty" db:"introduction"` // omitempty表示如果为空，那么json字段就没必要展示出来
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
