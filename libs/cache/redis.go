package cache

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type CacheAccess struct {
	Redis *redis.Client
}

func InitCacheAccess(client *redis.Client) (*CacheAccess, error) {
	return &CacheAccess{Redis: client}, nil
}

func GetRedisConn(host string, port string) (*redis.Client, error) {
	dsn := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
