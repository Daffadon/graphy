package cmd

import (
	"log/slog"

	_database "github.com/daffadon/graphy/config/database"
	"github.com/daffadon/graphy/config/env"
	"github.com/daffadon/graphy/config/logger"
	"github.com/daffadon/graphy/graph"
	"github.com/daffadon/graphy/internal/domain/notes"
	"github.com/daffadon/graphy/internal/domain/users"
	"github.com/daffadon/graphy/internal/infrastructure/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IBootstrap struct {
	G  *graph.Resolver
	S  *slog.Logger
	DB *pgxpool.Pool
}

func BootstrapRun() *IBootstrap {
	env.Load()

	slog := logger.NewSlog()
	dbcon := _database.NewSQLConn(slog)
	q := database.NewQuerier(dbcon)
	ur := users.NewUserRepository(q, slog)
	nr := notes.NewNoteRepository(q, slog)

	return &IBootstrap{
		G: &graph.Resolver{
			Ur: ur,
			Nr: nr,
			S:  slog,
		},
		S:  slog,
		DB: dbcon,
	}
}
