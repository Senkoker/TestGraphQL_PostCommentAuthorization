package runtime

import (
	"context"
	"errors"
	"fmt"
)

func AuthorizationCheck(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	fmt.Println(userID, "ID после контекста")
	if !ok || userID == "" {
		return "", errors.New("Unauthorized")
	}
	return userID, nil
}
