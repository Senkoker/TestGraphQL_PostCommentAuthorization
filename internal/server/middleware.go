package server

import (
	"bytes"
	"context"
	"friend_graphql/internal/logger"
	"io"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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
			userID := ""
			ctx := context.WithValue(c.Request().Context(), "userID", userID)
			c.Request().WithContext(ctx)
			return next(c)
		}
		var userID string
		tokenClaims := token.Claims.(jwt.MapClaims)
		userID = tokenClaims["userID"].(string)
		ctx := context.WithValue(c.Request().Context(), "userID", userID)
		c.Request().WithContext(ctx)
		return next(c)
	}
}

func RequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return err
		}

		if len(body) > 0 {
			logger.GetLogger().Info(string(body))
		}

		// Восстанавливаем тело запроса
		c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

		return next(c)
	}
}
