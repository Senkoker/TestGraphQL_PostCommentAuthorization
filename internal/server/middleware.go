package server

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	secret = "my_secret"
)

func AuthorizationMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		tokenString := authorization[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "token is invalid")
		}
		var userID string
		tokenClaims := token.Claims.(jwt.MapClaims)
		userID = tokenClaims["userID"].(string)
		ctx := context.WithValue(c.Request().Context(), "userID", userID)
		c.Request().WithContext(ctx)
		return next(c)
	}
}
func AuthorizationCheck(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", errors.New("Unauthorized")
	}
	return userID, nil
}
