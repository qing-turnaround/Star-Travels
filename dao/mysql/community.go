package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"web_app/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&communityList, sqlStr); err != nil { // select语句的dest应该为切片指针
		if err == sql.ErrNoRows { //如果是空行错误
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityNameByID(id int64) (name string) {
	sqlStr := `select community_name from community
		where community_id=?`
	_ = db.Get(&name, sqlStr, id)
	return
}

// GetCommunityDetailByName 提供查询社区详情
func GetCommunityDetailByName(communityName string) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select 
			community_id, community_name, introduction, create_time 
			from community
			where community_name = ?`
	if err = db.Get(communityDetail, sqlStr, communityName); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidName
		}
	}
	return communityDetail, err
}

func GetCommunityDetailByID(communityID int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `select 
			community_id, community_name, introduction, create_time 
			from community
			where community_id = ?`
	if err = db.Get(communityDetail, sqlStr, communityID); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return communityDetail, err
}

// CreateCommunity 在mysql数据库中Community表中 add 一行
func CreateCommunity(p *models.ParamCreateCommunity) (err error) {
	sqlStr := `insert into community(community_id, community_name, introduction)values(?,?,?)`
	_, err = db.Exec(sqlStr, p.CommunityID, p.CommunityName, p.Introduction)
	return
}
