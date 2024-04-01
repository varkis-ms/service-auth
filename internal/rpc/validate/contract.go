package validate

import (
	"context"

	"github.com/varkis-ms/service-auth/internal/model"
)

type Repository interface {
	GetUserById(ctx context.Context, id int64) (*model.User, error)
}
