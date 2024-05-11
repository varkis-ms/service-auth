package validate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/varkis-ms/service-auth/internal/model"
	"github.com/varkis-ms/service-auth/internal/pkg/logger/sl"
	pb "github.com/varkis-ms/service-auth/internal/pkg/pb"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserSaver

type Handler struct {
	repo      Repository
	secretKey string
	log       *slog.Logger
}

func New(
	repo Repository,
	secretKey string,
	log *slog.Logger,
) *Handler {
	return &Handler{
		repo:      repo,
		secretKey: secretKey,
		log:       log,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.ValidateRequest, out *pb.ValidateResponse) error {
	//TODO: подумать над ошибками, мб стоит просто отправлять -> ok: false
	token, err := jwt.Parse(in.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}

		return []byte(h.secretKey), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !h.isValid(claims) {
			return model.ErrUnauthenticated
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return model.ErrUnauthenticated
		}

		userID := int64(claims["uid"].(float64))
		email := claims["email"].(string)

		user, err := h.repo.GetUserById(ctx, userID)
		if err != nil {
			if !errors.Is(err, model.ErrUserNotFound) {
				h.log.Error("repo.GetUserById", sl.Err(err))

				return err
			}

			return model.ErrUnauthenticated
		}

		if user.ID != userID || user.Email != email {
			return model.ErrUnauthenticated
		}

		out.UserID = user.ID
	}

	return nil
}

func (h *Handler) isValid(claims jwt.MapClaims) bool {
	if _, ok := claims["exp"].(float64); !ok {
		return false
	}

	if _, ok := claims["uid"].(float64); !ok {
		return false
	}

	if _, ok := claims["email"].(string); !ok {
		return false
	}

	return true
}
