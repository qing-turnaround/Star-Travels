package redis

import (
	"context"
	"web_app/settings"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rb *ConsistentHashBalance
)

func Init(cfg *settings.RedisConfig) (err error) {
	// 初始化一致性hash
	rb = NewConsistentHashBalance(nil, cfg.Replicas, cfg)

	// 添加一个主sentinel
	rb.SentinelClient = redis.NewSentinelClient(&redis.Options{
		Addr: cfg.Sentinels[0], // 随机取一个
		Password: cfg.Password,
	})

	// 测试连接
	for _, client := range rb.Clients {
		err = client.Ping(ctx).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() {
	for _, client := range rb.Clients {
		client.Close()
	}
}

func Update(conf *settings.RedisConfig) {
	rb.Update(conf)
}