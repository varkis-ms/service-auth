package signup

import "context"

type Repository interface {
	SignupToDb(ctx context.Context, email string, passHash []byte) error
}
