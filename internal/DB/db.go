package db

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	redisHost               = "redis"
	redisPort               = "6379"
	redisPasswordEnvVarName = "REDIS_DATABASE_PASSWORD"
)

type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type LoggerInterface interface {
	Print(v ...any)
	Fatal(v ...any)
}

// InitRedisDB connects to a redis server
func InitRedisDB(logger LoggerInterface) *redis.Client {
	ctx := context.Background()
	// For production usage, the password can alternatively be stored in AWS secrets manager and retrieved from there
	password := os.Getenv(redisPasswordEnvVarName)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: password,
		DB:       0})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Print("Successfully connected to the redis database")

	return rdb
}

// GetValue gets the long url that a short url maps to
func GetValue(ctx context.Context, key string, rdb RedisClientInterface, logger LoggerInterface) (*string, error) {
	val, err := rdb.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		logger.Print("getLongURL: key does not exist")
		return nil, err
	case err != nil:
		logger.Print("getLongURL: Get failed")
		return nil, err
	case val == "":
		logger.Print("getLongURL: value is empty")
		return &val, nil
	}

	return &val, nil
}
