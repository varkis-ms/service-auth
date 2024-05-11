package storage

import (
	"context"

	"github.com/varkis-ms/service-auth/internal/model"
)

// Repository описывает операции на уровне хранилища
type Repository interface {
	SignupToDb(ctx context.Context, email string, passHash []byte) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
}
