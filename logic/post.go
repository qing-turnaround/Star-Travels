package logic

import (
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
	_, err = mysql.GetCommunityDetailByID(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID error: ", zap.Error(err))
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

	// 2. 保存到 mysql 数据库
	if err = mysql.CreatePost(post); err != nil {
		return
	}
	// 3. 在redis中创建 Key
	err = redis.CreatePost(post.PostID, mysql.GetCommunityNameByID(post.CommunityID))
	return
}

// GetPostByID 查询帖子通过帖子ID
func GetPostByID(postID int64) (data *models.ApiPostDetails, err error) {
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
		zap.Int64("AuthorID", post.AuthorID)
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
		zap.Int64("communityID", post.CommunityID)
		return
	}
	data = &models.ApiPostDetails{
		Post:          post,
		AuthorName:    user.Username,
		CommunityName: community.Name,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetails, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetails, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
			zap.Int64("AuthorID", post.AuthorID)
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			zap.Int64("communityID", post.CommunityID)
			continue
		}
		postDetail := &models.ApiPostDetails{
			Post:          post,
			AuthorName:    user.Username,
			CommunityName: community.Name,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

// GetPostList2 升级版获取帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	// 1. 去redis里面查询帖子id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	// 2. 根据id去mysql数据库中查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 将帖子的作者信息及其分区信息查询出来（还需填充帖子的票数）
	data = make([]*models.ApiPostDetails, 0, len(posts))
	for index, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
			zap.Int64("AuthorID", post.AuthorID)
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			zap.Int64("communityID", post.CommunityID)
			continue
		}
		postDetail := &models.ApiPostDetails{
			Post:          post,
			VoteNum:       voteData[index],
			AuthorName:    user.Username,
			CommunityName: community.Name,
		}
		data = append(data, postDetail)
	}
	return data, nil

}

// GetCommunityPostList 通过社区ID来获取帖子列表
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	// 1. 去redis里面查询帖子id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	// 2. 根据id去mysql数据库中查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 将帖子的作者信息及其分区信息查询出来（还需填充帖子的票数）
	data = make([]*models.ApiPostDetails, 0, len(posts))
	for index, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Error(err))
			zap.Int64("AuthorID", post.AuthorID)
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			zap.Int64("communityID", post.CommunityID)
			continue
		}
		postDetail := &models.ApiPostDetails{
			Post:          post,
			VoteNum:       voteData[index],
			AuthorName:    user.Username,
			CommunityName: community.Name,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

// GetPostListNew 将两个查询帖子的接口合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetails, err error) {
	if p.CommunityName == "" {
		return GetPostList2(p)
	} else {
		return GetCommunityPostList(p)
	}

}
