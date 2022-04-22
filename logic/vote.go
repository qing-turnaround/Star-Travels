package logic

import (
	"strconv"
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

/* VoteForPost 为帖子投票
投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况
	1.之前没投过票，现在要投赞成票      -> +432
	2.之前投过反对票，现在要改为赞成票   -> +432*2
v=0时，有两种情况
    1.之前投过反对票，现在要取消		-> +432
	2.之前投过赞成票，现在要取消        -> -432
v=-1时，有两种情况
	1.之前没投过票，现在要投反对票      -> -432
	2.之前投过赞成票，现在要改为反对票   -> -432*2
*/

func VoteForPost(userID int64, p *models.ParamVote) error {
	zap.L().Debug("PostVote",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.String("direction", p.Direction))
	d, _ := strconv.ParseFloat(p.Direction, 64)
	return redis.PostVote(strconv.Itoa(int(userID)), p.PostID, d)
}
