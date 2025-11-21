package tokenauthority

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Laelapa/CompanyRegistry/internal/config"
)

type TokenAuthority struct {
	cfg *config.AuthConfig
}

func New(cfg *config.AuthConfig) *TokenAuthority {
	return &TokenAuthority{
		cfg: cfg,
	}
}

func (t *TokenAuthority) IssueJWT(userID uuid.UUID) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    t.cfg.JwtIssuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(t.cfg.JwtLifetime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.cfg.JwtSecret))
	if err != nil {
		return "", errors.New("failed to sign JWT: " + err.Error())
	}

	return signedToken, nil
}
