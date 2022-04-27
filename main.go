package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"
	"web_app/viper"

	"go.uber.org/zap"
)

// @title           Star-Travels
// @version         1.0
// @description     基于gin框架的社区帖子项目
// @termsOfService  http://swagger.io/terms/
// @host      localhost:8080
// @BasePath  /api/v1/

func main() {
	// 1. 加载配置（configure）
	var configFile string
	//定义命令行参数
	flag.StringVar(&configFile, "f", "./conf/config.yaml", "配置文件的路径")
	//解析命令行参数
	flag.Parse()
	if err := settings.Init(configFile); err != nil {
		fmt.Printf("init settings failed: %v\n", err)
		return
	}
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed: %v\n", err)
		return
	}
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed: %v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed: %v\n", err)
		return
	}
	//把缓存区的日志追加到日志
	defer zap.L().Sync()
	zap.L().Debug("logger init success")
	// 3. 初始化Mysql
	if err := mysql.Init(settings.Conf.MysqlMasterConfig, settings.Conf.MysqlSlaveConfig); err != nil {
		fmt.Printf("init mysql failed: %v\n", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed: %v\n", err)
		return
	}
	defer redis.Close()
	// viper 热加载
	viper.Watch()
	// 5. 注册路由
	r := routes.SetUp(settings.Conf.Mode)
	// 6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown ", zap.Error(err))
	}

	zap.L().Info("Server exiting")

}
