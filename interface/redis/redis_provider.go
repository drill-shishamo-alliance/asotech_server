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
		Addr:     "redis",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}