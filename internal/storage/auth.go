package storage

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/varkis-ms/service-auth/internal/model"
)

func (s *Storage) SignupToDb(ctx context.Context, email string, passHash []byte) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := s.db.Builder.
		Insert("auth_user").
		Columns("email", "pass_hash").
		Values(email, passHash).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		if err.Error() == errDuplicate {
			return model.ErrUserExists
		}
		// TODO: подумать над ситуацией, когда пользователь уже существует
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := s.db.Builder.
		Select("id, email, pass_hash").
		From("auth_user").
		Where(sq.Eq{"email": email}).
		ToSql()

	var user model.User
	if err := tx.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	return &user, tx.Commit(ctx)
}

func (s *Storage) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := s.db.Builder.
		Select("id, email").
		From("auth_user").
		Where(sq.Eq{"id": id}).
		ToSql()

	var user model.User
	if err := tx.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	return &user, tx.Commit(ctx)
}
