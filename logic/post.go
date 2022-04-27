package logic

import (
	"fmt"
	"strconv"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.ParamPost, userID int64) (err error) {
	// 先校验 community_id的正确性
	communityDetail, err := mysql.GetCommunityDetailByID(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID error: ", zap.Error(err), zap.String("communityID is", strconv.Itoa(int(p.CommunityID))))
		return
	}
	post := &models.Post{
		PostID:      snowflake.GenID(),
		AuthorID:    userID,
		CommunityID: p.CommunityID,
		Title:       p.Content,
		Content:     p.Content,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	// 2. 保存到 mysql
	if err = mysql.CreatePost(post); err != nil {
		return
	}
	// 3. 保存进 redis(用做缓存)
	err = redis.CreatePost(post, communityDetail.CommunityName)
	return
}

// GetPostByID 查询帖子通过帖子ID
func GetPostByID(postID int64) (data *models.ApiPostDetails, err error) {
	var community *models.Community
	// 先查询redis缓存，如果没有再查询mysql
	post,  err := redis.GetPostByID(fmt.Sprint(postID))
	if err != nil {
		zap.L().Error("redis.GetPostByID failed", zap.Error(err))
		// 查询 数据库
		post, err = mysql.GetPostByID(postID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID failed", zap.Error(err))
			return
		}
		community, err = mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			zap.Int64("communityID", post.CommunityID)
			return
		}
		// 写入 redis
		redis.CreatePost(post, community.CommunityName)
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
		zap.Int64("AuthorID", post.AuthorID)
		return
	}
	// 根据社区id查询社区详细信息
	if community == nil {
		community, err = mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			zap.Int64("communityID", post.CommunityID)
			return
		}
	}

	data = &models.ApiPostDetails{
		Post:          post,
		AuthorName:    user.UserName,
		CommunityName: community.CommunityName,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetails, err error) {

	// 获取 这些 帖子 的 ID
	postIDs, err := redis.GetPostIDsInOrder("","time", page, size)
	if err != nil {
		return
	}

	return GetPostListByIDs(postIDs)
}

// GetPostList2 升级版获取帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	// 可以选择 order规则
	postIDs, err := redis.GetPostIDsInOrder("", p.Order, p.Page, p.Size)
	if err != nil {
		return data, err
	}
	return GetPostListByIDs(postIDs)
}

// GetCommunityPostList 通过社区ID来获取帖子列表
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	// 1. 去redis里面查询帖子id列表
	postIDs, err := redis.GetPostIDsByCommunityID(p)
	return GetPostListByIDs(postIDs)
}

func GetPostListByIDs(postIDs []string) (data []*models.ApiPostDetails, err error){
	data = make([]*models.ApiPostDetails, 0, len(postIDs))
	for _, postIDStr := range postIDs {
		postId, _ := strconv.Atoi(postIDStr)
		post, err := GetPostByID(int64(postId))
		if err != nil {
			return data, err
		}
		data = append(data, post)
	}
	return
}


// GetPostListNew 将两个查询帖子的接口合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	if p.CommunityName == "" {
		return GetPostList2(p)
	} else {
		return GetCommunityPostList(p)
	}

}
