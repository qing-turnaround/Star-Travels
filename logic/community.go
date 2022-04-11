package logic

import (
	"time"
	"web_app/dao/mysql"
	"web_app/models"
)

// CreateCommunity 创建社区
func CreateCommunity(p *models.ParamCreateCommunity) error {
	community := &models.Community{
		CommunityID:   p.CommunityID,
		CommunityName: p.CommunityName,
		Introduction:  p.Introduction,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	return mysql.CreateCommunity(community)
}

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库，查到所有community并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailByName(communityName string) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByName(communityName)
}
