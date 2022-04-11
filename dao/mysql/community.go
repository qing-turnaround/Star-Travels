package mysql

import (
	"web_app/models"
)

/**** 读操作 ****/

func GetCommunityList() (communityList []*models.Community, err error) {
	err = randomGetSlave().Find(&communityList).Error
	return
}

func GetCommunityNameByID(communityID int64) (name string) {
	community := &models.Community{}
	randomGetSlave().Select("community_name").Where("community_id = ?", communityID).Find(community)
	return community.CommunityName
}

// GetCommunityDetailByName 提供查询社区详情
func GetCommunityDetailByName(communityName string) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	rowsAffected := randomGetSlave().Where("community_name = ?", communityName).Find(communityDetail).RowsAffected
	if rowsAffected == 0 {
		err = ErrorInvalidName
	}
	return
}

func GetCommunityDetailByID(communityID int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	rowsAffected := randomGetSlave().Where("community_id = ?", communityID).Find(communityDetail).RowsAffected
	if rowsAffected == 0 {
		err = ErrorInvalidID
	}
	return
}


/* 写操作 */

// CreateCommunity 在mysql数据库中Community表中 add 一行
func CreateCommunity(community *models.Community) (err error) {
	err = master.Create(community).Error
	return
}
