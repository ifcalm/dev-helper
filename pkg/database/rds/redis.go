package rds

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// 封装Redis客户端，用于实际项目
type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("无法连接到Redis: %v", err)
	} else {
		fmt.Println("Redis连接成功:", pong)
	}
}

// redis 客户端
func NewRedisClient() *RedisClient {
	return &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}
}

// 设置key
func SetKey(key, value string) error {
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Fatalf("could not get key: %v", err)
		return err
	}
	fmt.Println("set key successful")
	return nil
}

// 获取key
func GetKey(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Fatalf("could not get key: %v", err)
		return "", err
	}

	return val, nil
}
