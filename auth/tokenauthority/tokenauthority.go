package tokenauthority

import (
	"fmt"
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

func (t *TokenAuthority) IssueJWT(userID uuid.UUID) (signedToken string, err error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    t.cfg.JwtIssuer,
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(t.cfg.JwtLifetime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(t.cfg.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

func (t *TokenAuthority) ValidateJWT(tokenString string) (subjectUUID uuid.UUID, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
			}
			return []byte(t.cfg.JwtSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to extract subject from JWT: %w", err)
	}

	subjectUUID, err = uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("token subject is not a valid UUID: %w", err)
	}

	return subjectUUID, nil
}
