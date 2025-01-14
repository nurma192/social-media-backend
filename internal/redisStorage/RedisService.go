package redisStorage

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func CreateClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

type RedisService struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
		ctx:    context.Background(),
	}
}
