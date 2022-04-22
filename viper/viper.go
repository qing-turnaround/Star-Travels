package viper

import (
	"fmt"
	"web_app/dao/redis"
	"web_app/settings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 监控，热加载配置文件
func Watch() {
	viper.WatchConfig() //实时监控配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		//当配置文件信息发生变化 就修改 Conf 变量
		if err := viper.Unmarshal(settings.Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v\n", err)
		}
		// 配置文件发生变化，也变更一下一致性哈希
		redis.Update(settings.Conf.RedisConfig)
	})
}