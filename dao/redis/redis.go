package redis

import (
	"context"
	"web_app/settings"
)

var (
	ctx = context.Background()
	rb *ConsistentHashBalance
)

func Init(cfg *settings.RedisConfig) (err error) {
	// 初始化一致性hash
	rb = NewConsistentHashBalance(nil, cfg.Replicas, cfg)

	// 测试连接
	for _, clients := range rb.Clients {
		for _, conn := range clients {
			if err = conn.Ping(ctx).Err(); err != nil {
				return err
			}
		}
	}

	return
}

// 关闭实例连接
func Close() {
	for _, clients := range rb.Clients {
		for _, conn := range clients {
			conn.Close()
		}
	}

	rb.SentinelClient.Close()
}

func Update(conf *settings.RedisConfig) {
	rb.Update(conf)
}