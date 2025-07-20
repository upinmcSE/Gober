package cache

import (
	"Gober/configs"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func NewRedisClient(config *configs.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Username:     config.Redis.Username,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("🍺 Connected Redis")

	return client
}
