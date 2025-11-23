package adapters

import (
	"context"
	"errors"

	"github.com/Laelapa/CompanyRegistry/internal/domain"
	"github.com/Laelapa/CompanyRegistry/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGUserRepoAdapter struct {
	q *repository.Queries
}

func NewPGUserRepoAdapter(q *repository.Queries) *PGUserRepoAdapter {
	return &PGUserRepoAdapter{q: q}
}

// Create creates a new user.
// If uniqueness constraints are violated, it returns domain.ErrConflict.
// It propagates other errors from the database layer.
func (p *PGUserRepoAdapter) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	if u.Username == nil {
		return nil, errors.New("username is required")
	}
	if u.PasswordHash == nil {
		return nil, errors.New("password hash is required")
	}
	params := repository.CreateUserParams{
		Username:     *u.Username,
		PasswordHash: *u.PasswordHash,
	}
	dbUser, err := p.q.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // pg:unique_violation
			return nil, domain.ErrConflict
		}
		return nil, err
	}

	return p.toDomainType(&dbUser), nil
}

func (p *PGUserRepoAdapter) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	dbUser, err := p.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return p.toDomainType(&dbUser), nil
}

// toDomain converts from DB model to domain model
func (p *PGUserRepoAdapter) toDomainType(u *repository.User) *domain.User {
	return &domain.User{
		ID:           &u.ID,
		Username:     &u.Username,
		PasswordHash: &u.PasswordHash,
	}
}
