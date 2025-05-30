package validator

import (
	proto "github.com/legenda-hortici/hw-6-proto/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UsernameError = status.Error(codes.InvalidArgument, "username is required")
	PasswordError = status.Error(codes.InvalidArgument, "password is required")
)

func ValidateLogin(req *proto.LoginRequest) error {
	if req.GetUsername() == "" {
		return UsernameError
	}

	if req.GetPassword() == "" {
		return PasswordError
	}

	return nil
}

func ValidateRegister(req *proto.RegisterRequest) error {
	if req.GetUsername() == "" {
		return UsernameError
	}

	if req.GetPassword() == "" {
		return PasswordError
	}

	return nil
}
