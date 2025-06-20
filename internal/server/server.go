package server

import (
	"context"

	"github.com/labstack/echo/v4/middleware"

	runtime "friend_graphql/internal/resolversGO"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/vektah/gqlparser/v2/ast"
)

type Server struct {
	router *echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	return &Server{router: e}
}

func (s *Server) QraphQLHandle(postDomain runtime.PostDomainInterface, commentDomain runtime.CommentDomainInterface,
	userDomain runtime.UserDomainInterface) {
	s.router.Use(middleware.RequestID())
	s.router.Use(AuthorizationMiddleWare)
	s.router.Use(middleware.Logger())
	config := runtime.Config{Resolvers: &runtime.Resolver{PostDomain: postDomain,
		CommentDomain: commentDomain, UserDomain: userDomain}}
	config.Directives.InputUnion = runtime.NewInputUnionDirective()
	srv := handler.New(runtime.NewExecutableSchema(config))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(transport.Websocket{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	s.router.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	s.router.POST("/query", echo.WrapHandler(srv))
}

func (s *Server) Start() {
	go func() {
		s.router.Start("localhost:8085")
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.router.Shutdown(ctx)
}
