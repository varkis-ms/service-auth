package signup

import (
	"context"
	"service-auth/internal/storage"
	"service-auth/pkg/logging"
	pb "service-auth/pkg/pb"

	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserSaver

type Handler struct {
	repo Repository
	log  *logging.Logger
}

func New(
	repo storage.Repository,
	log *logging.Logger,
) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

// Handle registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (h *Handler) Handle(ctx context.Context, in *pb.SignupRequest, out *pb.SignupResponse) error {
	log := h.log.WithField("email", in.Email)

	passHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Info("failed to generate password hash")

		return err
	}

	if err := h.repo.SignupToDb(ctx, in.Email, passHash); err != nil {
		log.WithError(err).Info("failed to save user")

		return err
	}

	out.Ok = true
	return nil
}
