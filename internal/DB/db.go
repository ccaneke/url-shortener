package db

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

const (
	REDIS_HOST     = "localhost"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
)

// InitRedisDB connects to a redis server
func InitRedisDB(ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the redis database")

	return rdb
}
