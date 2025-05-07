package RedisHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
	"friend_graphql/pkg/DBRedis"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type HandlerRedis struct {
	storage *DBRedis.RedisStorage
	ctxTime time.Duration
}

func NewRedisHandler(storage *DBRedis.RedisStorage, ctxTime time.Duration) *HandlerRedis {
	return &HandlerRedis{storage: storage, ctxTime: ctxTime}
}

func (r *HandlerRedis) GetPostHashtagHash(hashtags []string, limit, offset string) ([]model.Post, error) {
	var allQuery string
	for i, hashtag := range hashtags {
		text := fmt.Sprintf("@tags_ids:{%s}", hashtag)
		if i == 0 {
			allQuery = text
		} else {
			allQuery = allQuery + " " + text
		}
	}
	ctx := context.Background()
	res, err := r.storage.RStorage.Do(ctx,
		"FT.SEARCH",
		"Post_index",
		allQuery,
		"LIMIT", offset, limit,
	).Result()
	if err != nil {
		return nil, fmt.Errorf("Problem to return data from redis:%v", err)
	}
	massive, _ := strconv.Atoi(limit)
	posts := make([]model.Post, 0, massive)
	var post model.Post
	result := res.([]interface{})
	for i := 1; i < len(result); i++ {
		if i%2 == 0 {
			postString := result[i].([]interface{})[1].(string)
			err = json.Unmarshal([]byte(postString), &post)
			if err != nil {
				fmt.Println(err)
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *HandlerRedis) GetPostHash(postIds []string) ([]model.Post, []string, error) {
	ctx := context.Background()
	pipe := r.storage.RStorage.Pipeline()
	cmds := make([]*redis.Cmd, len(postIds))
	posts := make([]model.Post, 0, len(postIds))
	existPostIds := make([]string, 0, len(postIds))
	for i := 0; i < len(postIds); i++ {
		cmds[i] = pipe.Do(ctx, "JSON.GET", fmt.Sprintf("Post_id:%s", postIds[i]))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		if err != redis.Nil {
			logger.GetLogger().Error("Problem to pipe redis data", "err", err)
			return nil, nil, fmt.Errorf("problem to connect any statement for send query to redis:%w", err)
		}
	}
	for i := 0; i < len(postIds); i++ {
		data, err := cmds[i].Result()
		if err != nil {
			logger.GetLogger().Error("Problem to Json.Get", "err", err)
			continue
		}
		var post model.Post
		err = json.Unmarshal([]byte(data.(string)), &post)
		if err != nil {
			continue
		}
		existPostIds = append(existPostIds, post.ID)
		posts = append(posts, post)
	}
	fmt.Println(existPostIds, "найденные id")
	return posts, existPostIds, nil
}
func (r *HandlerRedis) CreatePopularPostHash(posts []model.Post) error {
	pipe := r.storage.RStorage.Pipeline()
	for _, post := range posts {
		jsonPost, err := json.Marshal(post)
		if err != nil {
			return fmt.Errorf("problem to marshal post to redis:%w", err)
		}
		ctx := context.Background()
		pipe.Do(ctx, "JSON.SET", fmt.Sprintf("Post_id:%s", post.PostID), "$", jsonPost)

	}
	ctx := context.Background()
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("problem to connect any statement for send query to redis:%w", err)
	}
	return nil
}
