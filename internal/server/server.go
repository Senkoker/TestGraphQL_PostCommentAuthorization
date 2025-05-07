package server

import (
	"context"
	"friend_graphql/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/vektah/gqlparser/v2/ast"
	"time"
)

type Server struct {
	router *echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	return &Server{router: e}
}

func (s *Server) QraphQLHandle(producer graph.ProducerKafkaInterface) {
	s.router.Use()
	s.router.Use()

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})

	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	s.router.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	s.router.GET("/query", echo.WrapHandler(srv))
}

func (s *Server) Start() {
	go func() {
		s.router.Start("localhost:8080")
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.router.Shutdown(ctx)
}
