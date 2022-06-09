package db

import (
	"context"
	"time"

	"github.com/inoth/ino-gathere/components/config"

	"github.com/go-redis/redis/v8"
)

var RedisConnect *redis.ClusterClient

type RedisConnectCluster struct{}

func (RedisConnectCluster) Init() error {

	hosts := config.Cfg.GetStringSlice("Redis.Host")
	password := config.Cfg.GetString("Redis.Passwd")
	poolSize := config.Cfg.GetInt("Redis.PoolSize")

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    hosts,
		Password: password,
		PoolSize: poolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.Cfg.GetDuration("Redis.PoolTimeout"))
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}
	RedisConnect = client
	return nil
}
