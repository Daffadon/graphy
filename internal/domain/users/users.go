package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/daffadon/graphy/internal/domain/dto"
	"github.com/daffadon/graphy/internal/infrastructure/database"
	"github.com/jackc/pgx/v5"
)

type (
	UserRepository interface {
		GetUserByEmail(ctx context.Context, email string) (dto.User, error)
		CreateUser(ctx context.Context, input *dto.User) error
	}
	userRepository struct {
		q database.Querier
		l *slog.Logger
	}
)

func NewUserRepository(q database.Querier, l *slog.Logger) UserRepository {
	return &userRepository{
		q: q,
		l: l,
	}
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (dto.User, error) {
	query, args, err := sq.Select("id", "email", "fullname", "password").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		u.l.Error("failed to build sql", slog.Any("err", err))
		return dto.User{}, err
	}

	var user dto.User
	err = u.q.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Fullname, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.User{}, nil
		}
		u.l.Error(fmt.Sprintf("failed to execute query Err: %v", err))
		return dto.User{}, err
	}
	return user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, input *dto.User) error {
	query, args, err := sq.Insert("users").
		Columns("id", "email", "fullname", "password").
		Values(input.ID, input.Email, input.Fullname, input.Password).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		u.l.Error("failed to build insert sql", slog.Any("err", err))
		return err
	}

	_, err = u.q.Exec(ctx, query, args...)
	if err != nil {
		u.l.Error("failed to execute insert", slog.Any("err", err))
		return err
	}
	return nil
}
