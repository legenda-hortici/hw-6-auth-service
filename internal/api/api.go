package grpcapi

import (
	"context"
	"github.com/legenda-hortici/hw-6-auth-service/internal/storage/myerr"
	"github.com/legenda-hortici/hw-6-auth-service/pkg/validator"
	proto "github.com/legenda-hortici/hw-6-proto/gen/go/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=AuthService
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
	if err := validator.ValidateRegister(request); err != nil {
	}

	err := s.auth.Register(ctx, request.Username, request.Password)
	if err != nil {
		// Если это статусная ошибка — возвращаем как есть
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}

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
	if err := validator.ValidateLogin(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
