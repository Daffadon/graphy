package users

import (
	"log/slog"

	"github.com/daffadon/graphy/internal/infrastructure/database"
)

type (
	UserRepository interface{}
	userRepository struct {
		q database.Querier
		l *slog.Logger
	}
)

func NewUserRepository(q database.Querier, l *slog.Logger) UserRepository {
	return userRepository{
		q: q,
		l: l,
	}
}
