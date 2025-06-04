package grpcapi

import (
	"context"
	"github.com/legenda-hortici/hw-6-auth-service/internal/mocks"
	proto "github.com/legenda-hortici/hw-6-proto/gen/go/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServerAPI_Register(t *testing.T) {
	test := []struct {
		name      string
		request   proto.RegisterRequest
		setupMock func(m *mocks.AuthService)
		wantErr   codes.Code
	}{
		{
			name: "success register user",
			request: proto.RegisterRequest{
				Username: "user1",
				Password: "password1",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Register", mock.Anything, "user1", "password1").Return(nil)
			},
			wantErr: codes.OK,
		},
		{
			name: "fail register user",
			request: proto.RegisterRequest{
				Username: "user1",
				Password: "password1",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Register", mock.Anything, "user1", "password1").Return(status.Error(codes.Internal, "failed to register user"))
			},
			wantErr: codes.Internal,
		},
		{
			name: "already exists user",
			request: proto.RegisterRequest{
				Username: "user1",
				Password: "password1",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Register", mock.Anything, "user1", "password1").Return(status.Error(codes.AlreadyExists, "user already exists"))
			},
			wantErr: codes.AlreadyExists,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := mocks.NewAuthService(t)
			tt.setupMock(mockAuth)

			api := &serverAPI{
				auth: mockAuth,
			}
			_, err := api.Register(context.Background(), &tt.request)
			if tt.wantErr == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantErr, status.Code(err))
			} else {
				assert.Equal(t, tt.wantErr, status.Code(err))
			}
		})
	}
}
