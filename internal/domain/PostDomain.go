package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"

	"github.com/99designs/gqlgen/graphql"
)

type Domain struct {
	StorageRedis    StorageRedisInterface
	StoragePostgres StoragePostgresInterface
	Producer        ProducerKafkaInterface
	AmazonS3        AmazonS3Interface
	StorageMongoDb  StorageMongoDBInterface
}
type StorageMongoDBInterface interface {
	StorageGetPostWithHashtag(hashtags *string, limit, offset *int32) ([]*model.Post, []string, error)
	StorageGetPostWithID(postIDs []string) ([]*model.Post, []string, error)
	StorageGetUserPosts(ctx context.Context, userID string, limit, offset int32) ([]*model.Post, error)
}

type StorageRedisInterface interface {
	GetPostHashtagHash(hashtagsPointer *string, limit, offset *int32) ([]*model.Post, error)
	GetPostHash(postIds []string) ([]*model.Post, []string, error)
	CreatePopularPostHash(posts []*model.Post) error
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

func NewPostDomain(storageRedis StorageRedisInterface, storagePostgres StoragePostgresInterface,
	amazon AmazonS3Interface, producer ProducerKafkaInterface, mongo StorageMongoDBInterface) *Domain {
	return &Domain{
		StorageRedis:    storageRedis,
		StoragePostgres: storagePostgres,
		AmazonS3:        amazon,
		Producer:        producer,
		StorageMongoDb:  mongo,
	}
}

func (d *Domain) UploadPostKafka(input *model.NewPost, userID string) (string, error) {
	op := "uploadPostKafka"
	loggerUpload := logger.GetLogger().With(op)
	imgUrl, err := d.AmazonS3.UploadFile(input.File)
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

func (d *Domain) FeedGetPosts(interestPostIds []string) ([]*model.Post, error) {
	op := "FeedGetPost"
	postLogger := logger.GetLogger().With("operation", op)
	posts, postIds, redisErr := d.StorageRedis.GetPostHash(interestPostIds)
	if redisErr != nil {
		postLogger.Error("Problem get post hash", "err", redisErr)
	}
	interestPostIds = subtractSlices(interestPostIds, postIds)
	fmt.Println("проверка", len(interestPostIds))
	var popularPost []*model.Post
	var userInfo map[string]model.UserInfo
	if len(interestPostIds) != 0 {
		mongoPosts, users, err := d.StorageMongoDb.StorageGetPostWithID(interestPostIds)
		users = uniqueSlice(users)
		if err != nil && redisErr != nil {
			postLogger.Error("Redis error", "err", redisErr)
			postLogger.Error("Postgres error", "err", err)
			return nil, err
		} else if err != nil {
			postLogger.Error("Postgres error", "err", err)
			return posts, nil
		}
		userInfo, err = d.StoragePostgres.GetUserInfo(users)
		if err != nil {
			postLogger.Error("Problem get user info postgres", "err", err)
			return posts, err
		}
		for i := 0; i < len(mongoPosts); i++ {
			author := mongoPosts[i].AuthorID
			mongoPosts[i].Author = userInfo[author].FirstName + " " + userInfo[author].SecondName
			mongoPosts[i].ImgPersonURL = userInfo[author].ImgUrl
			if *mongoPosts[i].Watched > -1 {
				popularPost = append(popularPost, mongoPosts[i])
			}
			if i == (len(mongoPosts) - 1) {
				err = d.StorageRedis.CreatePopularPostHash(popularPost)
				if err != nil {
					postLogger.Error("Problem to sent post in Redis", "err", err)
				}
			}
		}
		posts = append(posts, mongoPosts...)
	}
	fmt.Println("посты после обработки", posts)
	return posts, nil
}

func (d *Domain) FeedGetPostsWithHashtag(hashtags *string, limit, offset *int32, redisStatus string) ([]*model.Post, error) {
	op := "Domain GetPosts"
	logger := logger.GetLogger().With("op", op)
	if redisStatus == "true" {
		postWithHashtagsRedis, err := d.StorageRedis.GetPostHashtagHash(hashtags, limit, offset)
		return postWithHashtagsRedis, err
	}

	postHashtags, users, err := d.StorageMongoDb.StorageGetPostWithHashtag(hashtags, limit, offset)
	if err != nil {
		if errors.Is(err, errors.New("DoesNotExist")) {
			logger.Error("Problem get Post with hashtags", "err", err.Error())
			return nil, err
		}
		logger.Error("Problem get Post with hashtags", "err", err.Error())
		return nil, err
	}
	usersInfo, err := d.StoragePostgres.GetUserInfo(users)
	if err != nil {
		logger.Error("Problem get user info postgres", "err", err.Error())
		return nil, err
	}
	for i := 0; i < len(postHashtags); i++ {
		authorID := postHashtags[i].AuthorID
		info := usersInfo[authorID]
		postHashtags[i].Author = info.FirstName + " " + info.SecondName
		postHashtags[i].ImgPersonURL = info.ImgUrl
	}
	return postHashtags, nil
}

func (d *Domain) GetUserPosts(ctx context.Context, userID string, limit, offset int32) ([]*model.Post, error) {
	op := "GetUserPosts"
	loggerUserPosts := logger.GetLogger().With(op)
	posts, err := d.StorageMongoDb.StorageGetUserPosts(ctx, userID, limit, offset)
	if err != nil {
		loggerUserPosts.Error("Problem get posts", "err", err.Error())
		return nil, err
	}
	return posts, nil
}
