package grpcapi

import (
	"context"
	proto "github.com/legenda-hortici/hw-6-proto/gen/go/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"skillsRockAuthService/internal/storage/myerr"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type serverAPI struct {
	proto.UnimplementedAuthServiceServer
	auth AuthService
}

func NewAPI(gRPCServer *grpc.Server, auth AuthService) {
	proto.RegisterAuthServiceServer(
		gRPCServer,
		&serverAPI{
			auth: auth,
		})
}

func (s *serverAPI) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if request.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if request.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	err := s.auth.Register(ctx, request.Username, request.Password)
	if err != nil {
		if errors.Is(err, myerr.UserNotFoundErr) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RegisterResponse{
		Message: "success",
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	if request.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if request.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.Login(ctx, request.Username, request.Password)
	if err != nil {
		if errors.Is(err, myerr.UserNotFoundErr) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.LoginResponse{
		Token: token,
	}, nil
}
