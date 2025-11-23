package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type UserService struct {
	repo           UserRepository
	tokenAuthority *tokenauthority.TokenAuthority
}

func NewUserService(
	repo UserRepository,
	tokenAuthority *tokenauthority.TokenAuthority,
) *UserService {
	return &UserService{
		repo:           repo,
		tokenAuthority: tokenAuthority,
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
