package MongoDB

import (
	"friend_graphql/graph/model"
	"friend_graphql/pkg/MongoDB"
)

type PostCommentHandler struct {
	storage *MongoDB.ClientMongo
}

func NewPostCommentHandler(client *MongoDB.ClientMongo) *PostCommentHandler {
	return &PostCommentHandler{storage: client}
}
func (h *PostCommentHandler) GetPostWithHashtag(hashtags []string) {
	post := new(model.Post)
	h.storage.Client.Find()
}
func (h *PostCommentHandler) GetPostWithID(postIDs []string) {

}
