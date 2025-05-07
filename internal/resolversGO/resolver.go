package runtime

import "friend_graphql/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type PostCommentDomainInterface interface {
	UploadPostKafka(input *model.NewPost, userID string) (string, error)
	UploadCommentKafka(input *model.NewComment, userID string) (string, error)
	FindPost(hashtags string, postIDS []string) (*model.Post, error)
}

type Resolver struct {
	PostCommentDomain PostCommentDomainInterface
}
