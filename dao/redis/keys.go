package redis

import (
	"math/rand"
	"time"
)

// redis keys

// redis key注意使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix          = "Stars-Travel:"
	KeyPostInfoHashPrefix = KeyPrefix + "bluebell:post:"
	KeyPostTimeZSet    = KeyPrefix + "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZSet   = KeyPrefix + "post:score"  // zset;帖子及投票的分数
	KeyPostVotedZsetPF = KeyPrefix + "post:voted:" // zset;记录用户及投票类型，参数是post id
	KeyCommunitySetPF  = KeyPrefix + "community:"  // 保存每个分区下帖子的id
)

var (
	// 过期时间
	oneHour = 3600 // 一小时
	threeDaysInSecond          = 3 * 24 * oneHour // 三天
	oneWeekInSecond = 7 * 24 * oneHour    // 一周
)

// 随机初始化过期时间（至少3天过期，然后加上 随机一周大小的时间）
func RandomExpire() time.Duration {
	num := 	threeDaysInSecond + rand.Intn(oneWeekInSecond)
	return time.Duration(num) * time.Second

}