package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/internal/domain"
	"github.com/legenda-hortici/hw-6-auth-service/internal/storage/myerr"
	"github.com/legenda-hortici/hw-6-auth-service/pkg/jwt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg           config.Config
	log           *zap.SugaredLogger
	authProvider  AuthRepository
	tokenProvider TokenRepository
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=AuthRepository
type AuthRepository interface {
	Register(ctx context.Context, username string, passHash []byte) error
	Login(ctx context.Context, email string) (*domain.Users, error)
	CheckUser(ctx context.Context, username string) (bool, error)
}

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, refresh domain.RefreshToken) error
	RefreshTokenCheck(ctx context.Context, refreshID uuid.UUID) (bool, error)
	RefreshTokenUpdate(ctx context.Context, refresh domain.RefreshToken) error
	UserByID(ctx context.Context, tokenHash uuid.UUID) (domain.Users, error)
}

func NewAuthService(
	cfg config.Config,
	log *zap.SugaredLogger,
	authProvider AuthRepository,
	tokenProvider TokenRepository,
) *AuthService {
	return &AuthService{
		cfg:           cfg,
		log:           log,
		authProvider:  authProvider,
		tokenProvider: tokenProvider,
	}
}

func (s *AuthService) Register(ctx context.Context, username string, password string) error {
	const op = "services.Register"

	exists, err := s.authProvider.CheckUser(ctx, username)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)
		return errors.Wrap(err, ":"+op)
	}

	if exists {
		return errors.New("User already exists")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)
		return errors.Wrap(err, ":"+op)
	}

	err = s.authProvider.Register(ctx, username, passHash)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)
		return errors.Wrap(err, ":"+op)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	const op = "services.Login"

	user, err := s.authProvider.Login(ctx, email)
	if err != nil {
		if errors.Is(err, myerr.UserNotFoundErr) {
			s.log.Errorf("%s: %v", op, err)
			return "", "", errors.Wrap(err, ":"+op)
		}

		s.log.Errorf("%s: %v", op, err)
		return "", "", errors.Wrap(err, ":"+op)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		s.log.Errorf("%s: %v", op, err)
		return "", "", errors.Wrap(err, ":"+op)
	}

	accessToken, err := jwt.NewAccessToken(s.cfg.TokenJWT, *user)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)
		return "", "", errors.Wrap(err, ":"+op)
	}

	refreshTokenStr, refreshToken, err := jwt.NewRefreshToken(s.cfg.TokenJWT, user.ID)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)
		return "", "", errors.Wrap(err, ":"+op)
	}

	if err := s.tokenProvider.SaveRefreshToken(ctx, domain.RefreshToken{
		UserID:    user.ID,
		Hash:      refreshToken.Hash,
		ExpiresAt: refreshToken.ExpireAt,
		CreatedAt: refreshToken.CreatedAt,
	}); err != nil {
		s.log.Errorf("login failed: %v", err)

		return "", "", errors.Wrap(err, ":"+op)
	}

	return accessToken, refreshTokenStr, nil
}

func (s *AuthService) CheckToken(ctx context.Context, token string) (string, string, error) {
	const op = "services.CheckToken"

	tokenID, err := jwt.ParseRefreshToken(token, s.cfg.TokenJWT.Secret)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.Wrap(err, "failed to refresh token:")
	}

	exist, err := s.tokenProvider.RefreshTokenCheck(ctx, tokenID)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.Wrap(err, "failed to check token:")
	}

	if !exist {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.New("token not found")
	}

	user, err := s.tokenProvider.UserByID(ctx, tokenID)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.Wrap(err, "failed to check user:")
	}

	tokenAccess, err := jwt.NewAccessToken(s.cfg.TokenJWT, user)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.Wrap(err, "failed to refresh token:")
	}

	refreshTokenStr, refreshToken, err := jwt.NewRefreshToken(s.cfg.TokenJWT, user.ID)
	if err != nil {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.Wrap(err, "failed to refresh token:")
	}

	if err := s.tokenProvider.RefreshTokenUpdate(ctx, domain.RefreshToken{
		UserID:    user.ID,
		Hash:      refreshToken.Hash,
		ExpiresAt: refreshToken.ExpireAt,
		CreatedAt: refreshToken.CreatedAt,
	}); err != nil {
		s.log.Errorf("%s: %v", op, err)

		return "", "", errors.Wrap(err, "failed to refresh token:")
	}

	return tokenAccess, refreshTokenStr, nil
}
