package mysql

import (
	"web_app/models"
)


/**** 读操作 ****/

// GetPostByID 查询帖子返回数据
func GetPostByID(postID int64) (data *models.Post, err error) {
	data = new(models.Post) // 必要的一步，不然是空指针无法传入
	err = randomGetSlave().Where("post_id = ?", postID).Find(data).Error
	return
}

// GetPostList 查询帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	// 最新的帖子优先显示（offset 跳过 多少条记录，Limit 取多少条记录）
	err = randomGetSlave().Order("create_time desc").Limit(size).Offset((page - 1) * size).Find(&posts).Error
	return
}

// GetPostListByIDs 通过postID列表来查询帖子详细信息
func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	err = randomGetSlave().Where("post_id in (?)", ids).Order("post_id desc").Find(&posts).Error
	return
}


/* 写操作 */

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	err = master.Create(p).Error
	return
}