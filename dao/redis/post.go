package redis

import (
	"github.com/go-redis/redis"
	"time"
	"web_app/models"
)

// 获取PostID列表
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 2. 确定查询的索引的起点
	start := (page - 1) * size // 从第（page-1）* size 条记录开始
	stop := start + size - 1   // 到 start + size - 1结束

	// 3. ZrevRange 查询（从大到小查询，查询结果为Member）
	return rdb.ZRevRange(key, start, stop).Result()
}

// GetPostIDsInOrder 得到postID列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1. 根据用户请求携带的order参数来确定要查询的 redis key
	key := getRedisKey(KeyPostTimeZSet) // 与帖子创建时间相关的Key
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet) // 与帖子分数相关的Key
	}

	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据postID列表（）来查询对应帖子的投票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	pipeline := rdb.TxPipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZsetPF) + id
		pipeline.ZCount(key, "1", "1") //计算Score在 min-max之间的Member数量
	}

	cmders, err := pipeline.Exec()
	if err != nil {
		return
	}
	data = make([]int64, 0, len(ids)) // 初始化返回结果
	for _, cmder := range cmders {
		data = append(data, cmder.(*redis.IntCmd).Val())
	}
	return data, nil
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 使用 ZInterStore 把分区的帖子和帖子的分数的 zset 生成一个新的 zset
	orderKey := getRedisKey(KeyPostTimeZSet) // 与帖子创建时间相关的Key
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet) // 与帖子分数相关的Key
	}
	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + p.CommunityName)

	// 定义缓存键的名字
	key := orderKey + p.CommunityName
	if rdb.Exists(key).Val() < 1 {
		pipeline := rdb.TxPipeline()
		// 缓存键不存在，创建缓存键
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", // 交集取Score的最大值
		}, orderKey, cKey) //orderKey, cKey 需要操作的Key
		pipeline.Expire(key, 60*time.Second) // 设置过期时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	return getIDsFormKey(key, p.Page, p.Size)
}
