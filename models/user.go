package models

import "time"

type User struct {
	ID         int       `gorm:"primary_key;type:int(11) AUTO_INCREMENT;comment:'id'"`
	UserID     int64     `gorm:"unique;type:bigint(20) NOT NULL;comment:'用户ID'"`
	UserName   string    `gorm:"unique;type:varchar(64) COLLATE utf8mb4_general_ci NOT NULL;comment:'用户名称'"`
	Password   string    `gorm:"type:varchar(64) COLLATE utf8mb4_general_ci NOT NULL;comment:'密码'"`
	Email      string    `gorm:"type:varchar(64) COLLATE utf8mb4_general_ci;comment:'邮箱'"`
	Gender     string    `gorm:"type:char(1) NOT NULL DEFAULT '男';comment:'性别'"`
	Token      string    `gorm:"type:varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL;comment:'token'"`
	CreateTime time.Time `gorm:"type:timestamp;not null;comment:'创建时间'"` //创建时间
	UpdateTime time.Time `gorm:"type:timestamp;not null;comment:'修改时间'"`
}
