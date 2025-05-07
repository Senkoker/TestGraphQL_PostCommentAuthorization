package graph

import (
	"context"
	"errors"
)

func checkAuthorization(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", errors.New("user unauthorized")
	}
	return userID, nil
}
