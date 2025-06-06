package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/legenda-hortici/hw-6-auth-service/internal/storage/myerr"
	"github.com/pkg/errors"
	"time"

	"github.com/legenda-hortici/hw-6-auth-service/internal/config"
	"github.com/legenda-hortici/hw-6-auth-service/internal/domain"
)

func NewAccessToken(cfg config.TokenJWT, user domain.Users) (string, error) {
	if cfg.Secret == "" {
		return "", errors.New("jwt secret is required")
	}

	if cfg.AccessTTL < time.Minute*15 {
		return "", errors.New("jwt token ttl is less than 3600")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["uuid"] = user.ID
	claims["exp"] = time.Now().Add(cfg.AccessTTL).Unix()

	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewRefreshToken(cfg config.TokenJWT, id uuid.UUID) (string, domain.RefreshTokenClaims, error) {
	tokenID := uuid.New()
	createdAt := time.Now()
	expireAt := createdAt.Add(cfg.RefreshTTL).Unix()

	claims := domain.RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expireAt, 0)),
			IssuedAt:  jwt.NewNumericDate(createdAt),
		},
		Hash:      tokenID,
		ID:        id,
		CreatedAt: createdAt,
		ExpireAt:  expireAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", domain.RefreshTokenClaims{}, err
	}

	return tokenString, claims, nil
}

func ParseRefreshToken(tokenString, secret string) (uuid.UUID, error) {
	token, err := jwt.NewParser().ParseWithClaims(tokenString, &domain.RefreshTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*domain.RefreshTokenClaims); ok && token.Valid {
		if time.Now().Unix() > claims.ExpireAt {
			return uuid.Nil, errors.New("token expired")
		}
		return claims.Hash, nil
	}

	return uuid.Nil, myerr.ErrInvalidToken
}
