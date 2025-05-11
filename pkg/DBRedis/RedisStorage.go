package DBRedis

import (
	"context"
	"friend_graphql/internal/config"
	"friend_graphql/internal/logger"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisStorage struct {
	RStorage *redis.Client
}

func NewRedis(cfg *config.Cfg) *RedisStorage {
	redisStorage := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := redisStorage.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	_, err = redisStorage.Do(ctx,
		"FT.CREATE", "Post_index",
		"ON", "JSON",
		"PREFIX", "1", "Post_id:",
		"SCHEMA",
		"$.post_id", "AS", "id", "TEXT",
		"$.img_person_url", "AS", "name", "TEXT",
		"$.author", "AS", "author", "TEXT",
		"$.author_id", "AS", "author_id", "TEXT",
		"$.tags_ids", "AS", "tags_ids", "TAG",
		"$.content", "AS", "content", "TEXT",
		"$.created_at", "AS", "created_at", "TEXT",
		"$.watched", "AS", "watched", "NUMERIC",
		"$.likes", "AS", "likes", "NUMERIC",
	).Result()
	if err != nil {
		logger.GetLogger().Error("Error creating post index: %v", err)
	}
	return &RedisStorage{redisStorage}
}
