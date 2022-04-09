package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"web_app/settings"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	//初始化客户端连接
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,  //viper.GetString("redis.host"),
			cfg.Port), //viper.GetInt("redis.port")), //redis-serverip地址和端口号
		Password: cfg.Password, //viper.GetString("redis.password"), //密码设置
		DB:       cfg.DB,       //viper.GetInt("redis.db"), 		 //选择数据库
		PoolSize: cfg.PoolSize, //viper.GetInt("redis.poolsize"),    //连接池大小
	})

	_, err = rdb.Ping().Result()
	return
}

func Close() {
	_ = rdb.Close()
}
