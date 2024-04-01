package login

import (
	"context"
	"service-auth/internal/model"
	"service-auth/internal/storage"
	"service-auth/pkg/logging"
	pb "service-auth/pkg/pb"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserSaver

const ttlToken = time.Hour * 4

type Handler struct {
	repo      Repository
	secretKey string
	log       *logging.Logger
}

func New(
	repo storage.Repository,
	secretKey string,
	log *logging.Logger,
) *Handler {
	return &Handler{
		repo:      repo,
		secretKey: secretKey,
		log:       log,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.LoginRequest, out *pb.LoginResponse) error {
	log := h.log.WithField("email", in.Email)

	user, err := h.repo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		h.log.WithError(err).Error("repo.GetUserByEmail failed")

		return err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(in.Password)); err != nil {
		log.WithError(err).Info("bcrypt.CompareHashAndPassword failed")

		return model.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(ttlToken).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		log.WithError(err).Info("failed to generate token")

		return err
	}

	out.Token = tokenString
	return nil
}
