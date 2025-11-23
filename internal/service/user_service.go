package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/domain"
	"github.com/Laelapa/CompanyRegistry/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type UserService struct {
	repo           UserRepository
	tokenAuthority *tokenauthority.TokenAuthority
	logger         *logging.Logger
	producer       EventProducer
	topic          string
}

func NewUserService(
	repo UserRepository,
	tokenAuthority *tokenauthority.TokenAuthority,
	logger *logging.Logger,
	producer EventProducer,
	topic string,
) *UserService {
	return &UserService{
		repo:           repo,
		tokenAuthority: tokenAuthority,
		logger:         logger,
		producer:       producer,
		topic:          topic,
	}
}

// Register creates a new user and returns a signed JWT
// effectively logging them in upon registration.
// It returns domain.ErrConflict if the username is already taken.
func (u *UserService) Register(
	ctx context.Context,
	username,
	password string,
) (signedJWT string, err error) {
	if username == "" {
		return "", fmt.Errorf("username is required: %w", domain.ErrBadCredentials)
	}
	if password == "" {
		return "", fmt.Errorf("password is required: %w", domain.ErrBadCredentials)
	}
	// Hash the password
	hashedPassword, pErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if pErr != nil {
		return "", fmt.Errorf("failed to hash password: %w", pErr)
	}

	hashedPasswordStr := string(hashedPassword)

	// Create the user
	user := &domain.User{
		Username:     &username,
		PasswordHash: &hashedPasswordStr,
	}
	dbUser, rErr := u.repo.Create(ctx, user)
	if rErr != nil {
		// Explicitly showing that it can return ErrConflict
		if errors.Is(rErr, domain.ErrConflict) {
			return "", domain.ErrConflict
		}
		return "", rErr
	}

	// Generate JWT
	jwt, jErr := u.tokenAuthority.IssueJWT(*dbUser.ID)
	if jErr != nil {
		return "", fmt.Errorf("failed to issue JWT: %w", jErr)
	}

	u.publishEvent(ctx, "SIGNUP", *dbUser.ID)
	return jwt, nil
}

func (u *UserService) Login(ctx context.Context, username, password string) (signedJWT string, err error) {
	// Retrieve user by username
	user, uErr := u.repo.GetByUsername(ctx, username)
	if uErr != nil {
		return "", domain.ErrBadCredentials
	}

	// Compare provided password with stored hash
	if pErr := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); pErr != nil {
		return "", domain.ErrBadCredentials
	}

	jwt, jErr := u.tokenAuthority.IssueJWT(*user.ID)
	if jErr != nil {
		return "", fmt.Errorf("failed to issue JWT: %w", jErr)
	}

	return jwt, nil
}

func (u *UserService) publishEvent(ctx context.Context, eventType string, id uuid.UUID) {
	// pub/sub not configured
	if u.producer == nil {
		return
	}

	eventData := map[string]any{
		"event": eventType,
	}

	marshalledEvent, err := json.Marshal(eventData)
	if err != nil {
		u.logger.Error("Failed to marshal event", zap.Error(err))
		return
	}

	if pErr := u.producer.Produce(ctx, u.topic, id.String(), marshalledEvent); pErr != nil {
		u.logger.Error("Failed to produce event", zap.Error(pErr))
		return
	}
}
