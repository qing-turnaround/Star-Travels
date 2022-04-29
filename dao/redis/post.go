package redis

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis/v8"
)

//
// // GetPostVoteData 根据postID列表来查询对应帖子的投票数
// func GetPostVoteData(ids []string) (data []int64, err error) {
// 	data = make([]int64, 0, len(ids)) // 初始化返回结果
//
//
// 	postVotedZsetKey := getRedisKey(KeyPostVotedZsetPF)
// 	for _, id := range ids {
// 		key := postVotedZsetKey + id
// 		client := rb.Clients[rb.Get(key)]
// 		count := client.ZCount(ctx, key, "1", "1").Val() //计算Score在 min-max之间的Member数量
// 		data = append(data, count)
// 	}
//
//
// 	return data, nil
// }
//
// // GetCommunityPostIDsInOrder 按社区查询ids
// func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
// 	// 使用 ZInterStore 把分区的帖子和帖子的分数的 zset 生成一个新的 zset
// 	orderKey := getRedisKey(KeyPostTimeZSet) // 与帖子创建时间相关的Key
// 	if p.Order == models.OrderScore {
// 		orderKey = getRedisKey(KeyPostScoreZSet) // 与帖子分数相关的Key
// 	}
// 	// 社区的key
// 	cKey := getRedisKey(KeyCommunitySetPF + p.CommunityName)
//
// 	// 定义缓存键的名字
// 	key := orderKey + p.CommunityName
// 	if rdb.Exists(ctx, key).Val() < 1 {
// 		pipeline := rdb.TxPipeline()
// 		// 缓存键不存在，创建缓存键
// 		pipeline.ZInterStore(ctx, key, &redis.ZStore{
// 			Aggregate: "MAX", // 交集取Score的最大值
// 			Keys: []string{orderKey, cKey},
// 		}) //orderKey, cKey 需要操作的Key
// 		pipeline.Expire(ctx, key, 60*time.Second) // 设置过期时间
// 		_, err := pipeline.Exec(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return getIDsFormKey(key, p.Page, p.Size)
// }

func CreatePost(post *models.Post, communityName string) (err error) {
	postID, userID := fmt.Sprint(post.PostID), fmt.Sprint(post.AuthorID)
	now := float64(time.Now().Unix())
	postKey := KeyPostInfoHashPrefix + fmt.Sprint(post.PostID)
	votedKey := KeyPostVotedZsetPF + fmt.Sprint(post.PostID)
	communityKey := KeyCommunitySetPF + communityName // 社区Key
	postInfo := map[string]interface{}{
		"title":    post.Title,
		"summary":  post.Content,
		"post:id":  postID,
		"user:id":  userID,
		"community:id": fmt.Sprint(post.CommunityID),
		"create_time": 		fmt.Sprint(post.CreateTime.Format("2006-01-02 15:04:05")),

		"update_time": 		fmt.Sprint(post.UpdateTime.Format("2006-01-02 15:04:05")),

		"votes":    1,
	}
	// 获得对应键的 client
	clientPost := rb.GetRandomConn(rb.Get(postKey))
	clientVote:= rb.GetRandomConn(rb.Get(votedKey))
	clientCommunity := rb.GetRandomConn(rb.Get(communityKey))
	clientPostScore := rb.GetRandomConn(rb.Get(KeyPostScoreZSet))
	clientPostTime := rb.GetRandomConn(rb.Get(KeyPostTimeZSet))


	// 创建帖子
	clientPost.HMSet(ctx, postKey, postInfo)
	// 设置过期时间
	clientPost.Expire(ctx, postKey, RandomExpire())

	// 哪些人给这个帖子投了票
	clientVote.ZAdd(ctx, votedKey, &redis.Z{ // 作者默认投赞成票
		Score:  1,
		Member: userID,
	})
	clientVote.Expire(ctx, postKey, RandomExpire())

	// 投票分数 zset（永不过期key）
	clientPostScore.ZAdd(ctx, KeyPostScoreZSet, &redis.Z{
		Score:  scorePerVote,
		Member: postID,
	})

	// 时间分数 zset（永不过期key）
	clientPostTime.ZAdd(ctx, KeyPostTimeZSet, &redis.Z{ // 添加到时间的ZSet
		Score:  now,
		Member: postID,
	})

	// 更新：把帖子id加到社区的Set
	clientCommunity.SAdd(ctx, communityKey, postID)
	return err
}


// 通过 postID 获取帖子详细信息
func GetPostByID(postID string) (data *models.Post, err error) {
	// post Key
	postKey := KeyPostInfoHashPrefix + postID
	// 获取客户端
	clientPost := rb.GetRandomConn(rb.Get(postKey))

	// 获取结果
	postInfo, err := clientPost.HGetAll(ctx, postKey).Result()
	// 判断键是否还存在
	if postInfo["title"] == "" {
		return data, errors.New("键不存在或者过期")
	}

	hPostID , err := strconv.Atoi(postInfo["post:id"])
	hUserID , err := strconv.Atoi(postInfo["user:id"])
	hCommunityID , err := strconv.Atoi(postInfo["community:id"])
	hCreateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", postInfo["create_time"], time.Local)
	hUpdateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", postInfo["update_time"], time.Local)
	hVotes, _ := strconv.Atoi(postInfo["votes"])
	data = &models.Post{
		Content:   postInfo["title"],
		Title: postInfo["summary"],
		PostID:      int64(hPostID),
		AuthorID:    int64(hUserID),
		CommunityID: int64(hCommunityID),
		CreateTime:  hCreateTime,
		UpdateTime:  hUpdateTime,
		Votes: int64(hVotes),
	}
	return  data, err
}

func GetPostIDsByCommunityID(p *models.ParamPostList) (postIDs []string,err error) {
	orderKey := KeyPostTimeZSet
	if p.Order == models.OrderScore {
		orderKey = KeyPostScoreZSet
	}
	communityKey := KeyCommunitySetPF + p.CommunityName

	// 临时 key
	key := orderKey + communityKey
	client := rb.GetRandomConn(rb.Get(key))
	// 检查 key是否存在
	if client.Exists(ctx, key).Val() < 0 {
		// key 不存在，创建 key，设置过期时间，得出结果
		pipeline := client.TxPipeline()

		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys: []string{orderKey, communityKey},
			Aggregate: "max",
		})
		pipeline.Expire(ctx, key, RandomExpire())

		_, err = pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return GetPostIDsInOrder(key, p.Order, p.Page, p.Size)
}

func GetPostIDsInOrder(key, order string, page, size int64) ([]string, error) {
	if key == "" {
		key = KeyPostTimeZSet // 默认是时间key
		if order == models.OrderScore {
			key = KeyPostScoreZSet
		}
	}

	// 获取客户端
	clientPostTime := rb.GetRandomConn(rb.Get(key))
	start := (page-1) * size - 1
	end := start + size -1
	return clientPostTime.ZRevRange(ctx, key, start, end).Result()
}


