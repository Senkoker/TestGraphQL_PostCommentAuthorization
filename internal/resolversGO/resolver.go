package runtime

import (
	"context"
	"friend_graphql/graph/model"
)

type PostDomainInterface interface {
	UploadPostKafka(input *model.NewPost, userID string) (string, error)
	UploadCommentKafka(input *model.NewComment, userID string) (string, error)
	FeedGetPosts(interestPostIds []string) ([]*model.Post, error)
	FeedGetPostsWithHashtag(hashtags *string, limit, offset *int32, redisStatus string) ([]*model.Post, error)
	GetUserPosts(ctx context.Context, userID string, limit, offset int32) ([]*model.Post, error)
}

type CommentDomainInterface interface {
	GetComment(replyToID string, limit, offset int32) ([]*model.Comment, error)
}

type UserDomainInterface interface {
	GetUserInfo(userID string) (*model.User, error)
	GetUserFriendSubscriber(user *model.User, ctx context.Context,
		limit, offset int32, friendStatus bool) ([]*model.User, error)
}

type Resolver struct {
	PostDomain    PostDomainInterface
	CommentDomain CommentDomainInterface
	UserDomain    UserDomainInterface
}
