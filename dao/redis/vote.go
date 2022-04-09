package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"time"
)

const (
	oneWeekInSecond = 7 * 24 * 3600
	scorePerVote    = 432 // 每一票值多少分数
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已经过期")
)

func CreatePost(postID int64, communityName string) error {
	pipeline := rdb.TxPipeline() // 事务
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postID,
	})
	// 更新：把帖子id加到社区的Set
	cKey := getRedisKey(KeyCommunitySetPF + communityName) // 社区Key
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec() //只有执行时才需要判断错误
	return err
}

// VoteForPost 用户为帖子投票
func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 取redis中的发帖时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val() // ZScore().Val() 根据 member取 value
	if float64(time.Now().Unix())-postTime > oneWeekInSecond {
		return ErrorVoteTimeExpire
	}

	// 第2步和第3步放进事务中执行
	// 2. 更新帖子的分数
	//  先查当前用户给帖子的投票记录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedZsetPF+postID), userID).Val()
	if ov == value { // 若重复投票，之间返回
		return nil
	}
	var op float64 // 方便重新赋值，1表示投票分数大于用户历史投票分数，-1 则表示小于
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(value - ov)                                                  // 计算两次投票的差值
	pipeline := rdb.TxPipeline()                                                  // 事务
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff*op*scorePerVote, postID) // 修改帖子投票分数
	// 3. 记录用户为帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPF+postID), redis.Z{
			Score:  value, // 投赞同还是反对
			Member: userID,
		})
	}
	_, err := pipeline.Exec() // 只有执行时才需要判断错误
	return err
}
