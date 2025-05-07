package domain

import (
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
)

type Domain struct {
	StorageRedis    StorageRedisInterface
	StoragePostgres StoragePostgresInterface
}

type StorageRedisInterface interface {
	GetPostHashtagHash(hashtags []string, limit, offset string) ([]model.Post, error)
	GetPostHash(postIds []string) ([]model.Post, []string, error)
	CreatePopularPostHash(posts []model.Post) error
}

type StoragePostgresInterface interface {
	GetUserInfo(users []string) (map[string]model.UserInfo, error)
}

func NewPostCommentDomain(storage StorageRedisInterface) *Domain {
	return &Domain{StorageRedis: storage}
}

func (d *Domain) FeedGetPosts(interestPostIds []string) ([]model.Post, error) {
	op := "FeedGetPost"
	postLogger := logger.GetLogger().With("operation", op)
	posts, postIds, redisErr := d.StorageRedis.GetPostHash(interestPostIds)
	if redisErr != nil {
		postLogger.Error("Problem get post hash", "err", redisErr)

	}
	interestPostIds = subtractSlices(interestPostIds, postIds)
	var popularPost []model.Post
	var userInfo map[string]model.UserInfo
	if interestPostIds != nil {
		//TODO: написать обращение через MongoDB
		postgresPosts, users, err := d.PostgresPostComment.GetPosts(interestPostIds)
		users = uniqueSlice(users)
		if err != nil && redisErr != nil {
			localLogger.Error("Redis error", "err", redisErr)
			localLogger.Error("Postgres error", "err", err)
			return nil, err
		} else if err != nil {
			localLogger.Error("Postgres error", "err", err)
			return posts, nil
		}
		userInfo, err = d.StoragePostgres.GetUserInfo(users)
		if err != nil {
			localLogger.Error("Problem get user info postgres", "err", err)
			return posts, err
		}
		for i := 0; i < len(postgresPosts); i++ {
			author := postgresPosts[i].AuthorId
			postgresPosts[i].Author = userInfo[author].FirstName + " " + userInfo[author].SecondName
			postgresPosts[i].ImgPersonURL = userInfo[author].ImgUrl
			if postgresPosts[i].Watched > -1 {
				popularPost = append(popularPost, postgresPosts[i])
			}
			if i == (len(postgresPosts) - 1) {
				err = d.Redis.CreatePopularPostHash(popularPost)
				if err != nil {
					localLogger.Error("Problem to sent post in Redis", "err", err)
				}
			}
		}
		posts = append(posts, postgresPosts...)
	}
	return posts, nil
}
