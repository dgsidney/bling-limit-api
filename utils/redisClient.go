package utils

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var Ctx = context.Background()

func init() {
	_ = godotenv.Load()
}

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, _ = client.Ping(Ctx).Result()

	return client
}
