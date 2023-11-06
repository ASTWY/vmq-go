package main

import (
	"fmt"
	"time"

	"sync"
	"vmq-go/config"
	"vmq-go/db"
	"vmq-go/logger"
	"vmq-go/router"
	"vmq-go/task"
)

func main() {
	// 初始化日志
	logger.InitLogger(config.Conf.Log.Level, config.Conf.Log.Path)
	// 初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.Database.User,
		config.Conf.Database.Password,
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.DBName,
	)
	if err := db.SetupDatabase(dsn); err != nil {
		logger.Logger.Fatal("Database connection failed!")
	}
	wg := sync.WaitGroup{} // 定义一个同步等待的组

	// 订单过期任务
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 每隔1秒执行一次
			time.Sleep(time.Second)
			go task.CheckOrderExpire()
		}
	}()

	// 启动Gin
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Logger.Info("Server startup...")
		server := router.SetupRouter()

		if err := server.Run(fmt.Sprintf("%s:%d", config.Conf.Host, config.Conf.Port)); err != nil {
			logger.Logger.Fatalf("Server startup failed! %s", err.Error())
		}
	}()

	wg.Wait()
}
