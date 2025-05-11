package runtime

import (
	"context"
	"errors"
)

func AuthorizationCheck(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", errors.New("Unauthorized")
	}
	return userID, nil
}
