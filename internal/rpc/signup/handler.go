package signup

import (
	"context"
	"log/slog"

	"github.com/varkis-ms/service-auth/internal/pkg/logger/sl"
	pb "github.com/varkis-ms/service-auth/internal/pkg/pb"
	"github.com/varkis-ms/service-auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserSaver

type Handler struct {
	repo Repository
	log  *slog.Logger
}

func New(
	repo storage.Repository,
	log *slog.Logger,
) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

// Handle registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (h *Handler) Handle(ctx context.Context, in *pb.SignupRequest, out *pb.SignupResponse) error {
	log := h.log.With(slog.String("email", in.Email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Info("failed to generate password hash", sl.Err(err))

		return err
	}

	if err := h.repo.SignupToDb(ctx, in.Email, passHash); err != nil {
		log.Info("failed to save user", sl.Err(err))

		return err
	}

	out.Ok = true
	return nil
}
