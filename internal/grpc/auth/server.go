package auth

import (
	"context"
	"errors"

	"github.com/varkis-ms/service-auth/internal/model"
	auth "github.com/varkis-ms/service-auth/internal/pkg/pb"

	"github.com/varkis-ms/service-auth/internal/rpc/login"
	"github.com/varkis-ms/service-auth/internal/rpc/signup"
	"github.com/varkis-ms/service-auth/internal/rpc/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	auth.UnimplementedAuthServer
	loginHandler    *login.Handler
	signupHandler   *signup.Handler
	validateHandler *validate.Handler
}

func Register(
	gRPCServer *grpc.Server,
	loginHandler *login.Handler,
	signupHandler *signup.Handler,
	validateHandler *validate.Handler,
) {
	auth.RegisterAuthServer(gRPCServer, &server{
		loginHandler:    loginHandler,
		signupHandler:   signupHandler,
		validateHandler: validateHandler,
	})
	reflection.Register(gRPCServer)
}

func (s *server) Login(
	ctx context.Context,
	in *auth.LoginRequest,
) (*auth.LoginResponse, error) {
	if in.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	out := &auth.LoginResponse{}
	if err := s.loginHandler.Handle(ctx, in, out); err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return out, nil
}

func (s *server) Signup(
	ctx context.Context,
	in *auth.SignupRequest,
) (*auth.SignupResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	out := &auth.SignupResponse{}
	if err := s.signupHandler.Handle(ctx, in, out); err != nil {
		if errors.Is(err, model.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to signup user")
	}

	return out, nil
}

func (s *server) Validate(
	ctx context.Context,
	in *auth.ValidateRequest,
) (*auth.ValidateResponse, error) {
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	out := &auth.ValidateResponse{}
	if err := s.validateHandler.Handle(ctx, in, out); err != nil {
		if errors.Is(err, model.ErrUnauthenticated) {
			return nil, status.Error(codes.Unauthenticated, "user unauthenticated")
		}
		return nil, status.Error(codes.Internal, "failed to validate user")
	}
	return out, nil
}
