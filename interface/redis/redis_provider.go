package redis

import "github.com/go-redis/redis"

type redisProvider struct {}

type IRedisProvider interface {
	CreateClient() *redis.Client
}

func NewRedisProvider()  IRedisProvider {
	return &redisProvider{}
}

func (r *redisProvider) CreateClient() *redis.Client{
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}