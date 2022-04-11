package mysql

import (
	"fmt"
	"math/rand"
	"web_app/models"
	"web_app/settings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

// 定义一个全局对象db
var master *gorm.DB
var slave []*gorm.DB

// Init 定义一个初始化数据库的函数
func Init(masterCfg *settings.MysqlMasterConfig, slaveCfg *settings.MysqlSlaveConfig) (err error) {
	// DSN:Data Source Name
	//DSN格式为：[username[:password]@][protocol[(host:port)]]/dbname[?param1=value1&...&paramN=valueN]

	// 配置主节点
	masterDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		masterCfg.User,
		masterCfg.Password,
		masterCfg.Host,
		masterCfg.Port,
		masterCfg.Dbname,
	)
	master, err = gorm.Open("mysql", masterDsn)
	if err != nil {
		zap.L().Error("connect master failed", zap.Error(err))
		return
	}
	// 单数，表名不加s
	master.SingularTable(true)
	// 生成表
	master.Set("gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci").
		AutoMigrate(&models.Community{}, &models.Post{}, &models.User{})
	master.DB().SetMaxIdleConns(masterCfg.MaxIdleConns)
	master.DB().SetMaxOpenConns(masterCfg.MaxOpenConns)


	// 配置 所有 slave节点
	for i := 0; i < slaveCfg.Count; i++ {
		slaveDsn := getDsn(slaveCfg.User[i], slaveCfg.Password[i], slaveCfg.Host[i],
			slaveCfg.Dbname[i], slaveCfg.Port[i])
		slaveNode, err := gorm.Open("mysql", slaveDsn)
		slave = append(slave, slaveNode) // 加入到 从节点集合中
		if err != nil {
			zap.L().Error("connect slaveNode failed", zap.Error(err),
				zap.String("slaveNode Dsn is", "slaveDsn"))
			return err
		}
		// 设置最大空闲连接和最大连接数
		slaveNode.DB().SetMaxIdleConns(slaveCfg.MaxIdleConns)
		slaveNode.DB().SetMaxOpenConns(slaveCfg.MaxOpenConns)
	}

	return
}

func Close() {
	// 关闭master节点连接
	_ = master.Close()
	// 关闭slave 节点连接
	for i := range slave {
		slave[i].Close()
	}
}

func getDsn(user, password, host,  dbname string, port int) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, password, host, port, dbname,
	)
}

// 随机获取一个从节点 连接
func randomGetSlave() *gorm.DB {
	randInt := rand.Int31n(int32(settings.Conf.MysqlSlaveConfig.Count))
	return slave[randInt]
}