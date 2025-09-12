package graph

import (
	"log/slog"

	"github.com/daffadon/graphy/internal/domain/notes"
	"github.com/daffadon/graphy/internal/domain/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Ur users.UserRepository
	Nr notes.NoteRepository
	S  *slog.Logger
}
