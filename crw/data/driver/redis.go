package driver

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

// CreateDBConnection with postgres db
func CreateRedisConnection() (*redis.Client, context.Context) {
	ctx := context.Background()
	// create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	// ping the Redis server to check if it's running
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	// return the connection
	return rdb, ctx
}
