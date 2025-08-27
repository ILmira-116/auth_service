package validation

import (
	"github.com/ILmira-116/protos/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const emptyvalue = 0

func ValidateLoginRequest(req *auth.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyvalue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func ValidateRegisterRequest(req *auth.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func ValidateIsAdminRequest(req *auth.IsAdminRequest) error {
	if req.GetUserId() == emptyvalue {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}
