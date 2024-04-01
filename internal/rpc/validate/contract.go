package validate

import (
	"context"
	"service-auth/internal/model"
)

type Repository interface {
	GetUserById(ctx context.Context, id int64) (*model.User, error)
}
