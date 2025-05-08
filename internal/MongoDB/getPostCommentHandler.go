package MongoDB

import (
	"context"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
	"friend_graphql/pkg/MongoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostCommentHandler struct {
	storage *MongoDB.ClientMongo
}

func NewPostCommentHandler(client *MongoDB.ClientMongo) *PostCommentHandler {
	return &PostCommentHandler{storage: client}
}
func (h *PostCommentHandler) GetPostWithHashtag(hashtags []string, limit, offset int32) ([]*model.Post, []string, error) {
	posts := make([]*model.Post, 0, limit)
	users := make([]string, 0, limit)
	ctx := context.Background()
	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
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
func (h *PostCommentHandler) GetPostWithID(postIDs []string) ([]*model.Post, []string, error) {
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
