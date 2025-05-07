package domain

import (
	"encoding/json"
	"fmt"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
	"github.com/99designs/gqlgen/graphql"
)

type Domain struct {
	StorageRedis    StorageRedisInterface
	StoragePostgres StoragePostgresInterface
	Producer        ProducerKafkaInterface
	amazonS3        AmazonS3Interface
}

type StorageRedisInterface interface {
	GetPostHashtagHash(hashtags []string, limit, offset string) ([]model.Post, error)
	GetPostHash(postIds []string) ([]model.Post, []string, error)
	CreatePopularPostHash(posts []model.Post) error
}
type AmazonS3Interface interface {
	UploadFile(file graphql.Upload) (string, error)
}
type ProducerKafkaInterface interface {
	Produce(msg []byte) error
}

type StoragePostgresInterface interface {
	GetUserInfo(users []string) (map[string]model.UserInfo, error)
}

func NewPostCommentDomain(storageRedis StorageRedisInterface, storagePostgres StoragePostgresInterface) *Domain {
	return &Domain{StorageRedis: storageRedis, StoragePostgres: storagePostgres}
}

func (d *Domain) UploadPostKafka(input *model.NewPost, userID string) (string, error) {
	op := "uploadPostKafka"
	loggerUpload := logger.GetLogger().With(op)
	imgUrl, err := d.amazonS3.UploadFile(input.File)
	if err != nil {
		loggerUpload.Error("Error uploading file", "err", err)
	}
	input.AuthorID = userID
	input.ImgUrl = imgUrl
	postMsgKafka, err := json.Marshal(input)
	if err != nil {
		loggerUpload.Error("Error marshalling post json", "err", err)
		return "", fmt.Errorf("error marshalling post json: %v", err)
	}
	err = d.Producer.Produce(postMsgKafka)
	if err != nil {
		loggerUpload.Error("problem sent file to broker", "err", err)
		return "", fmt.Errorf("problem sent file to broker:%v", err)
	}
	return "OK", nil
}

func (d *Domain) UploadCommentKafka(input *model.NewComment, userID string) (string, error) {
	op := "uploadCommentKafka"
	loggerUpload := logger.GetLogger().With(op)
	input.AuthorID = userID
	postMsgKafka, err := json.Marshal(input)
	if err != nil {
		loggerUpload.Error("Error marshalling post json", "err", err)
		return "", fmt.Errorf("error marshalling post json: %v", err)
	}
	err = d.Producer.Produce(postMsgKafka)
	if err != nil {
		loggerUpload.Error("problem sent file to broker", "err", err)
		return "", fmt.Errorf("problem sent file to broker:%v", err)
	}
	return "OK", nil
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
