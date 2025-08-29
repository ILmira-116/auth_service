package authgrpc

import (
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/internal/validation"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/ILmira-116/protos/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	auth.UnimplementedAuthServer
	auth Auth
	log  *slog.Logger
}

// регистрация обработчика
func Register(gRPC *grpc.Server, authSvc *service.Auth, logger *slog.Logger) {
	auth.RegisterAuthServer(gRPC, &serverAPI{
		auth: authSvc,
		log:  logger,
	})
}

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := validation.ValidateLoginRequest(req); err != nil {
		s.log.Warn("login request validation failed", "email", req.GetEmail(), "err", err)
		return nil, err
	}

	s.log.Info("attempting to login user", "email", req.GetEmail(), "app_id", req.GetAppId())

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		errText := err.Error()
		if errors.Is(err, repository.ErrInvalidCredentials) || strings.Contains(errText, "invalid credentials") {
			s.log.Warn("login failed: invalid credentials", "email", req.GetEmail(), "err", err)
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		s.log.Error("login failed: internal error", "email", req.GetEmail(), "err", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	s.log.Info("user logged in successfully", "email", req.GetEmail(), "app_id", req.GetAppId())
	return &auth.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := validation.ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		errText := err.Error()
		s.log.Info("Register error received in handler", "err", errText)

		if errors.Is(err, repository.ErrUserExists) || strings.Contains(errText, "duplicate key") {
			s.log.Warn("duplicate registration attempt", "email", req.GetEmail(), "err", err)
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		s.log.Error("registration failed", "err", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &auth.RegisterResponse{
		UserId: userID,
	}, nil

}

func (s *serverAPI) IsAdmin(ctx context.Context, req *auth.IsAdminRequest) (*auth.IsAdminResponse, error) {
	if err := validation.ValidateIsAdminRequest(req); err != nil {
		s.log.Warn("IsAdmin request validation failed", "user_id", req.GetUserId(), "err", err)
		return nil, err
	}

	s.log.Info("checking if user is admin", "user_id", req.GetUserId())
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found in IsAdmin", "user_id", req.GetUserId(), "err", err)
			return nil, status.Error(codes.NotFound, "user not found")
		}
		s.log.Error("IsAdmin failed: internal error", "user_id", req.GetUserId(), "err", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	s.log.Info("checked if user is admin", "user_id", req.GetUserId(), "is_admin", isAdmin)
	return &auth.IsAdminResponse{IsAdmin: isAdmin}, nil
}
