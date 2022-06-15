package main

import (
	"fmt"
	"os"

	"github.com/inoth/ino-gathere/components/cache"
	"github.com/inoth/ino-gathere/components/config"
	"github.com/inoth/ino-gathere/components/logger"
	"github.com/inoth/ino-gathere/register"
	"github.com/inoth/ino-gathere/src/agent"
	_ "github.com/inoth/ino-gathere/src/plugins/inputs/all"
	_ "github.com/inoth/ino-gathere/src/plugins/outputs/all"
)

func main() {
	ag := &agent.AgentConfig{}
	// 注册组件
	err := register.Register(
		&cache.CacheComponents{}, // 本地缓存
		config.Instance(),        // 配置文件
		&logger.LogrusConfig{},   // logrus日志
		// &db.RedisConnectCluster{}, // redis 数据库
		// &db.MysqlConnect{},        // mysql 数据库
		ag,
	).Init().Run(ag)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
