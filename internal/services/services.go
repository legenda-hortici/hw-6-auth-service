package services

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"skillsRockAuthService/internal/domain"
	"skillsRockAuthService/internal/storage/myerr"
)

type AuthService struct {
	log          *zap.SugaredLogger
	authProvider AuthRepository
}

type AuthRepository interface {
	Register(ctx context.Context, username string, passHash []byte) error
	Login(ctx context.Context, email string) (*domain.User, error)
	CheckUser(ctx context.Context, username string) (bool, error)
}

func NewAuthService(
	log *zap.SugaredLogger,
	authProvider AuthRepository,
) *AuthService {
	return &AuthService{
		log:          log,
		authProvider: authProvider,
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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	const op = "services.Login"

	user, err := s.authProvider.Login(ctx, email)
	if err != nil {
		if errors.Is(err, myerr.UserNotFoundErr) {
			s.log.Errorf("%s: %v", op, err)
			return "", errors.Wrap(err, ":"+op)
		}

		s.log.Errorf("%s: %v", op, err)
		return "", errors.Wrap(err, ":"+op)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Errorf("%s: %v", op, err)
		return "", errors.Wrap(err, ":"+op)
	}

	return user.Username, nil
}
