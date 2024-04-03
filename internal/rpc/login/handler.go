package login

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/varkis-ms/service-auth/internal/model"
	"github.com/varkis-ms/service-auth/internal/pkg/logger/sl"
	pb "github.com/varkis-ms/service-auth/internal/pkg/pb"
	"github.com/varkis-ms/service-auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserSaver

const ttlToken = time.Hour * 4

type Handler struct {
	repo      Repository
	secretKey string
	log       *slog.Logger
}

func New(
	repo storage.Repository,
	secretKey string,
	log *slog.Logger,
) *Handler {
	return &Handler{
		repo:      repo,
		secretKey: secretKey,
		log:       log,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.LoginRequest, out *pb.LoginResponse) error {
	log := h.log.With(slog.String("email", in.Email))

	user, err := h.repo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		log.Error("repo.GetUserByEmail failed", sl.Err(err))

		return err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(in.Password)); err != nil {
		log.Info("bcrypt.CompareHashAndPassword failed", sl.Err(err))

		return model.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(ttlToken).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		log.Info("failed to generate token", sl.Err(err))

		return err
	}

	out.Token = tokenString
	return nil
}
