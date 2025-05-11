package server

import (
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
		buff := make([]byte, 1024) // Создаем буфер размером 1024 байта
		n, err := c.Request().Body.Read(buff)
		if err != nil && err != io.EOF {
			logger.GetLogger().Error(err.Error())
		}
		if n > 0 {
			logger.GetLogger().Info(string(buff[:n]))
		}
		return next(c)

	}
}
