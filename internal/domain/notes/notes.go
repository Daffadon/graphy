package notes

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/daffadon/graphy/graph/model"
	"github.com/daffadon/graphy/internal/domain/dto"
	"github.com/daffadon/graphy/internal/infrastructure/database"
)

type (
	NoteRepository interface {
		CreateNewNote(ctx context.Context, note *dto.Note, userId string) error
		UpdateNote(ctx context.Context, note *model.Note, userId string) error
		DeleteNote(ctx context.Context, noteid, userid string) error
		GetAllNotes(ctx context.Context, userid string) ([]*model.Note, error)
		GetNote(ctx context.Context, noteid, userid string) (*model.Note, error)
	}
	noteRepository struct {
		q database.Querier
		l *slog.Logger
	}
)

func NewNoteRepository(q database.Querier, l *slog.Logger) NoteRepository {
	return &noteRepository{
		q: q,
		l: l,
	}
}

func (n *noteRepository) UpdateNote(ctx context.Context, note *model.Note, userId string) error {
	query, args, err := sq.
		Update("notes").
		Set("title", note.Title).
		Set("description", note.Description).
		Set("text", note.Text).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": note.ID, "user_id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		n.l.Error("failed to build update query", slog.Any("error", err))
		return err
	}
	_, err = n.q.Exec(ctx, query, args...)
	if err != nil {
		n.l.Error("failed to update note", slog.Any("error", err))
		return err
	}
	return nil
}

func (n *noteRepository) GetNote(ctx context.Context, noteid string, userid string) (*model.Note, error) {
	query, args, err := sq.
		Select("id", "title", "description", "text").
		From("notes").
		Where(sq.Eq{"id": noteid, "user_id": userid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		n.l.Error("failed to build select query", slog.Any("error", err))
		return nil, err
	}

	row := n.q.QueryRow(ctx, query, args...)
	var note model.Note
	if err := row.Scan(&note.ID, &note.Title, &note.Description, &note.Text); err != nil {
		n.l.Error("failed to scan note row", slog.Any("error", err))
		return nil, err
	}
	return &note, nil
}

func (n *noteRepository) GetAllNotes(ctx context.Context, userid string) ([]*model.Note, error) {
	query, args, err := sq.
		Select("id", "title", "description", "text").
		From("notes").
		Where(sq.Eq{"user_id": userid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		n.l.Error("failed to build select query", slog.Any("error", err))
		return nil, err
	}

	rows, err := n.q.Query(ctx, query, args...)
	if err != nil {
		n.l.Error("failed to query notes", slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var notes []*model.Note
	for rows.Next() {
		var note model.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Description, &note.Text); err != nil {
			n.l.Error("failed to scan note row", slog.Any("error", err))
			return nil, err
		}
		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		n.l.Error("rows error", slog.Any("error", err))
		return nil, err
	}
	return notes, nil
}

func (n *noteRepository) CreateNewNote(ctx context.Context, note *dto.Note, userId string) error {
	query, args, err := sq.
		Insert("notes").
		Columns("id", "title", "description", "text", "user_id").
		Values(note.ID, note.Title, note.Description, note.Text, userId).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		n.l.Error("failed to build insert query", slog.Any("error", err))
		return err
	}
	_, err = n.q.Exec(ctx, query, args...)
	if err != nil {
		n.l.Error("failed to insert note", slog.Any("error", err))
		return err
	}
	return nil
}

func (n *noteRepository) DeleteNote(ctx context.Context, noteid, userid string) error {
	query, args, err := sq.
		Delete("notes").
		Where(sq.Eq{"id": noteid, "user_id": userid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		n.l.Error("failed to build delete query", slog.Any("error", err))
		return err
	}
	_, err = n.q.Exec(ctx, query, args...)
	if err != nil {
		n.l.Error("failed to delete note", slog.Any("error", err))
		return err
	}
	return nil
}
