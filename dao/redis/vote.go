package redis

import (
	"errors"
	"math"

	"github.com/go-redis/redis/v8"
)

const (
	scorePerVote    = 432 // 每一票值多少分数
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已经过期")
)


// PostVote 用户为帖子投票
func PostVote(userID, postID string, value float64) error {
	votedKey := KeyPostVotedZsetPF + postID
	postKey := KeyPostInfoHashPrefix + postID

	// 获取对应的客户端
	clientPostScore := rb.GetRandomConn(rb.Get(KeyPostScoreZSet))
	clientVote := rb.GetRandomConn(rb.Get(votedKey))
	clientPost := rb.GetRandomConn(rb.Get(postKey))

	// 获取用户关于帖子的状态（-1，0，1）
	ov := clientVote.ZScore(ctx, votedKey, userID).Val() // 获取当前分数
	// 计算差值的绝对值
	diffAbs := math.Abs(ov - value)

	clientVote.ZAdd(ctx, votedKey, &redis.Z{ // 记录已投票
		Score:  value,
		Member: userID,
	})

	clientPostScore.ZIncrBy(ctx, KeyPostScoreZSet, scorePerVote*diffAbs*value, postID) // 更新分数

	// 更新帖子的投票数
	switch math.Abs(ov) - math.Abs(value) {
	case 1:
		// 取消投票 ov=1/-1 v=0
		// 投票数-1
		clientPost.HIncrBy(ctx, postKey, "votes", -1)
	case 0:
		// 反转投票 ov=-1/1 v=1/-1
		// 投票数不用更新
	case -1:
		// 新增投票 ov=0 v=1/-1
		// 投票数+1
		clientPost.HIncrBy(ctx, KeyPostInfoHashPrefix+postID, "votes", 1)
	default:
		// 已经投过票了
		return errors.New("已经投过票了！")
	}
	return nil

}