package cache

import "github.com/go-redis/redis/v7"

type CacheAccess struct {
	Redis *redis.Client
}

func InitCacheAccess(client *redis.Client) (*CacheAccess, error) {
	return &CacheAccess{Redis: client}, nil
}

func GetRedisConn() (*redis.Client, error) {
	dsn := "localhost:6379"

	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
