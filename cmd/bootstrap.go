package cmd

import (
	"log/slog"

	_database "github.com/daffadon/graphy/config/database"
	"github.com/daffadon/graphy/config/env"
	"github.com/daffadon/graphy/config/logger"
	"github.com/daffadon/graphy/config/router"
	"github.com/daffadon/graphy/graph"
	"github.com/daffadon/graphy/internal/domain/users"
	"github.com/daffadon/graphy/internal/infrastructure/database"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IBootstrap struct {
	G  *graph.Resolver
	S  *slog.Logger
	DB *pgxpool.Pool
	R  *chi.Mux
}

func BootstrapRun() *IBootstrap {
	env.Load()

	slog := logger.NewSlog()
	dbcon := _database.NewSQLConn(slog)
	q := database.NewQuerier(dbcon)
	ur := users.NewUserRepository(q, slog)
	r := router.NewHTTPRouter()

	return &IBootstrap{
		G: &graph.Resolver{
			Ur: ur,
			S:  slog,
		},
		S:  slog,
		DB: dbcon,
		R:  r,
	}
}
