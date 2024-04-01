package login

import (
	"context"

	"github.com/varkis-ms/service-auth/internal/model"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}
