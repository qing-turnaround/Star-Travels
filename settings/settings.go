package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify" //用于监控的包
	"github.com/spf13/viper"
)

//全局配置结构体
var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"` //mapstructure：通用structTag
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int    `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"` //tag需要与配置文件中的名字对应
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(configFile string) (err error) {
	viper.SetConfigFile(configFile)
	// viper.SetConfigFile("./config.yaml") //指定配置文件（带后缀，可写绝对路径和相对路径两种）
	//viper.SetConfigName("config") //指定配置文件的名字（不带后缀）
	// 基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	viper.SetConfigType("yaml") //远程配置文件传输 确定配置文件的格式
	viper.AddConfigPath(".")    //指定配置文件的一个寻找路径
	err = viper.ReadInConfig()  //读取配置信息
	if err != nil {
		//读取配置信息错误
		fmt.Printf("viper.ReadInConfig() failed: %v\n", err)
		return
	}
	//把读取到的信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed: %v\n", err)
	}
	viper.WatchConfig() //实时监控配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		//当配置文件信息发生变化 就修改 Conf 变量
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed: %v\n", err)
		}
	})
	return
}
