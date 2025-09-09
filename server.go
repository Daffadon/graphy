package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daffadon/graphy/cmd"
	"github.com/daffadon/graphy/graph"
	database "github.com/daffadon/graphy/internal/pkg/database/postgresql"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := viper.GetString("app.port")
	if port == "" {
		port = defaultPort
	}

	b := cmd.BootstrapRun()

	database.MigratePostgres(b.DB, b.S)
	defer database.ClosePostgres(b.DB)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: b.G}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	b.R.Handle("/", playground.Handler("GraphQL playground", "/query"))
	b.R.Handle("/query", srv)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: b.R,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		b.S.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			b.S.Error(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()

	<-quit
	b.S.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		b.S.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}
	b.S.Info("Server exiting")
}
