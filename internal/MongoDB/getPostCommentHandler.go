package MongoDB

import (
	"context"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
	"friend_graphql/pkg/MongoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type PostCommentHandler struct {
	storage *MongoDB.ClientMongo
}

func NewPostCommentHandler(client *MongoDB.ClientMongo) *PostCommentHandler {
	return &PostCommentHandler{storage: client}
}
func (h *PostCommentHandler) StorageGetPostWithHashtag(hashtagsPointer *string, limit, offset *int32) ([]*model.Post, []string, error) {
	hashtags := strings.Split(*hashtagsPointer, "#")
	posts := make([]*model.Post, 0, *limit)
	users := make([]string, 0, *limit)
	ctx := context.Background()
	findOptions := options.Find()
	findOptions.SetSkip(int64(*offset))
	findOptions.SetLimit(int64(*limit))
	findOptions.SetSort(bson.D{{"likes", 1}})
	cursor, err := h.storage.Client.Find(ctx, bson.M{"TagIDS": bson.M{"$all": hashtags}}, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, nil, err
	}
	for cursor.Next(ctx) {
		post := new(model.Post)
		err = cursor.Decode(post)
		if err != nil {
			logger.GetLogger().With("MongoDBError").Error("Cursor Decode Error")
			continue
		}
		users = append(users, post.AuthorID)
		posts = append(posts, post)
	}
	return posts, users, nil
}
func (h *PostCommentHandler) StorageGetPostWithID(postIDs []string) ([]*model.Post, []string, error) {
	ctx := context.Background()
	posts := make([]*model.Post, 0, len(postIDs))
	post := new(model.Post)
	users := make([]string, 0, len(postIDs))
	cursor, err := h.storage.Client.Find(ctx, bson.M{"_id": bson.M{"$in": postIDs}})
	defer cursor.Close(ctx)
	if err != nil {
		return nil, nil, err
	}
	for cursor.Next(ctx) {
		err = cursor.Decode(post)
		if err != nil {
			logger.GetLogger().With("MongoDBError").Error("Cursor Decode Error")
			continue
		}
		users = append(users, post.AuthorID)
		posts = append(posts, post)
	}
	return posts, users, nil
}

func (h *PostCommentHandler) StorageGetComment(replyToID string, limit, offset int32) ([]*model.Comment, error) {
	ctx := context.Background()
	comments := make([]*model.Comment, 0, limit)
	findConfig := options.Find()
	findConfig.SetSkip(int64(offset))
	findConfig.SetLimit(int64(limit))
	cursor, err := h.storage.Client.Find(ctx, bson.M{"replyto": replyToID}, findConfig)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		comment := new(model.Comment)
		err = cursor.Decode(comment)
		if err != nil {
			logger.GetLogger().With("MongoDBError").Error("Cursor comment decode Error")
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (h *PostCommentHandler) StorageGetUserPosts(ctx context.Context, userID string, limit, offset int32) ([]*model.Post, error) {
	findOprtions := options.Find()
	findOprtions.SetSkip(int64(offset))
	findOprtions.SetLimit(int64(limit))
	cursor, err := h.storage.Client.Find(ctx, bson.M{"authorid": userID}, findOprtions)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	posts := []*model.Post{}
	for cursor.Next(ctx) {
		post := new(model.Post)
		err = cursor.Decode(post)
		if err != nil {
			logger.GetLogger().Error("cursor decode error ", "StorageGetUserPosts error", err.Error())
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}
