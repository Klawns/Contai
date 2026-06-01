package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"contai/internal/auth/app/contracts"
	authdomain "contai/internal/auth/domain"
	databaseports "contai/internal/database/ports"
	userports "contai/internal/users/app/ports"
	userdomain "contai/internal/users/domain"
)

func TestAuthService_Login(t *testing.T) {
	activeUser := mustUser(t, "john@example.com", "hash", userdomain.UserStatusActive)

	t.Run("should login active user with valid password", func(t *testing.T) {
		service := NewAuthService(
			fakeUserRepository{byEmail: map[string]*userdomain.User{"john@example.com": &activeUser}},
			fakePasswordHasher{},
			fakeJWTService{},
		)

		authenticatedUser, tokens, err := service.Login(context.Background(), contracts.LoginInput{
			Email:         " John@Example.COM ",
			PlainPassword: "secret",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if authenticatedUser.UserID != activeUser.ID {
			t.Fatalf("expected authenticated user %s, got %s", activeUser.ID, authenticatedUser.UserID)
		}
		if tokens.AccessToken != "access-token" {
			t.Fatalf("expected access token, got %s", tokens.AccessToken)
		}
	})

	t.Run("should reject missing user", func(t *testing.T) {
		service := NewAuthService(fakeUserRepository{}, fakePasswordHasher{}, fakeJWTService{})

		_, _, err := service.Login(context.Background(), contracts.LoginInput{Email: "missing@example.com", PlainPassword: "secret"})

		if !errors.Is(err, authdomain.ErrInvalidCredentials) {
			t.Fatalf("expected invalid credentials, got %v", err)
		}
	})

	t.Run("should reject invalid password", func(t *testing.T) {
		service := NewAuthService(
			fakeUserRepository{byEmail: map[string]*userdomain.User{"john@example.com": &activeUser}},
			fakePasswordHasher{compareErr: errors.New("bad password")},
			fakeJWTService{},
		)

		_, _, err := service.Login(context.Background(), contracts.LoginInput{Email: "john@example.com", PlainPassword: "wrong"})

		if !errors.Is(err, authdomain.ErrInvalidCredentials) {
			t.Fatalf("expected invalid credentials, got %v", err)
		}
	})

	t.Run("should reject inactive user", func(t *testing.T) {
		inactiveUser := mustUser(t, "inactive@example.com", "hash", userdomain.UserStatusInactive)
		service := NewAuthService(
			fakeUserRepository{byEmail: map[string]*userdomain.User{"inactive@example.com": &inactiveUser}},
			fakePasswordHasher{},
			fakeJWTService{},
		)

		_, _, err := service.Login(context.Background(), contracts.LoginInput{Email: "inactive@example.com", PlainPassword: "secret"})

		if !errors.Is(err, userdomain.ErrUserInactive) {
			t.Fatalf("expected inactive user, got %v", err)
		}
	})
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	activeUser := mustUser(t, "john@example.com", "hash", userdomain.UserStatusActive)

	t.Run("should validate active user token", func(t *testing.T) {
		service := NewAuthService(
			fakeUserRepository{byID: map[userdomain.UserID]*userdomain.User{activeUser.ID: &activeUser}},
			fakePasswordHasher{},
			fakeJWTService{claims: authdomain.AuthClaims{Type: authdomain.AuthTokenTypeAccess, UserID: activeUser.ID}},
		)

		authenticatedUser, err := service.ValidateAccessToken(context.Background(), "valid")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if authenticatedUser.UserID != activeUser.ID {
			t.Fatalf("expected user %s, got %s", activeUser.ID, authenticatedUser.UserID)
		}
	})

	t.Run("should reject invalid jwt", func(t *testing.T) {
		service := NewAuthService(fakeUserRepository{}, fakePasswordHasher{}, fakeJWTService{validateErr: authdomain.ErrInvalidToken})

		_, err := service.ValidateAccessToken(context.Background(), "invalid")

		if !errors.Is(err, authdomain.ErrInvalidToken) {
			t.Fatalf("expected invalid token, got %v", err)
		}
	})

	t.Run("should reject inactive token user", func(t *testing.T) {
		inactiveUser := mustUser(t, "inactive@example.com", "hash", userdomain.UserStatusInactive)
		service := NewAuthService(
			fakeUserRepository{byID: map[userdomain.UserID]*userdomain.User{inactiveUser.ID: &inactiveUser}},
			fakePasswordHasher{},
			fakeJWTService{claims: authdomain.AuthClaims{Type: authdomain.AuthTokenTypeAccess, UserID: inactiveUser.ID}},
		)

		_, err := service.ValidateAccessToken(context.Background(), "valid")

		if !errors.Is(err, userdomain.ErrUserInactive) {
			t.Fatalf("expected inactive user, got %v", err)
		}
	})
}

type fakeUserRepository struct {
	byEmail map[string]*userdomain.User
	byID    map[userdomain.UserID]*userdomain.User
}

func (repository fakeUserRepository) WithTx(tx databaseports.TxHandle) userports.UserRepository {
	return repository
}

func (repository fakeUserRepository) CreateUser(ctx context.Context, user *userdomain.User) (*userdomain.User, error) {
	return user, nil
}

func (repository fakeUserRepository) FindUserById(ctx context.Context, userID userdomain.UserID) (*userdomain.User, error) {
	return repository.byID[userID], nil
}

func (repository fakeUserRepository) FindUserByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	return repository.byEmail[email], nil
}

func (repository fakeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return repository.byEmail[email] != nil, nil
}

type fakePasswordHasher struct {
	compareErr error
}

func (hasher fakePasswordHasher) HashPassword(ctx context.Context, plainPassword string) (string, error) {
	return "hash", nil
}

func (hasher fakePasswordHasher) ComparePassword(ctx context.Context, passwordHash string, plainPassword string) error {
	return hasher.compareErr
}

type fakeJWTService struct {
	claims      authdomain.AuthClaims
	validateErr error
}

func (service fakeJWTService) IssueAccessToken(ctx context.Context, user userdomain.User) (string, authdomain.AuthClaims, error) {
	return "access-token", authdomain.AuthClaims{
		Type:      authdomain.AuthTokenTypeAccess,
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}, nil
}

func (service fakeJWTService) ValidateAccessToken(ctx context.Context, token string) (authdomain.AuthClaims, error) {
	if service.validateErr != nil {
		return authdomain.AuthClaims{}, service.validateErr
	}

	return service.claims, nil
}

func mustUser(t *testing.T, email string, passwordHash string, status userdomain.UserStatus) userdomain.User {
	t.Helper()

	user, err := userdomain.NewUser(userdomain.UserID(email), "John Doe", email, passwordHash)
	if err != nil {
		t.Fatalf("expected user, got %v", err)
	}
	if status == userdomain.UserStatusInactive {
		if err := user.Deactivate(); err != nil {
			t.Fatalf("expected inactive user, got %v", err)
		}
	}

	return user
}
