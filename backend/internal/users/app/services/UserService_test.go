package services

import (
	"context"
	"errors"
	"testing"

	databaseports "contai/internal/database/ports"
	"contai/internal/users/app/ports"
	"contai/internal/users/domain"
)

func TestUserService_CreateUserNormalizesEmailAndHashesPassword(t *testing.T) {
	repository := &fakeUserRepository{}
	service := NewUserService(repository, fakeUserIDGenerator{}, fakePasswordHasher{}, nil, nil)

	user, err := service.CreateUser(context.Background(), createUserInput(" John@Example.COM "))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "john@example.com" {
		t.Fatalf("expected normalized email, got %s", user.Email)
	}
	if repository.created == nil || repository.created.PasswordHash != "hashed-secret123" {
		t.Fatalf("expected hashed password to be persisted, got %#v", repository.created)
	}
}

func TestUserService_CreateUserRejectsDuplicateEmail(t *testing.T) {
	service := NewUserService(&fakeUserRepository{emailExists: true}, fakeUserIDGenerator{}, fakePasswordHasher{}, nil, nil)

	_, err := service.CreateUser(context.Background(), createUserInput("john@example.com"))

	if !errors.Is(err, domain.ErrUserEmailAlreadyExists) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}

func TestUserService_CreateUserRejectsWeakPassword(t *testing.T) {
	repository := &fakeUserRepository{}
	service := NewUserService(repository, fakeUserIDGenerator{}, fakePasswordHasher{}, nil, nil)

	_, err := service.CreateUser(context.Background(), ports.CreateUserInput{
		Name:          "John Doe",
		Email:         "john@example.com",
		PlainPassword: "short",
	})

	if !errors.Is(err, domain.ErrUserPasswordTooWeak) {
		t.Fatalf("expected weak password error, got %v", err)
	}
	if repository.created != nil {
		t.Fatalf("expected no persisted user, got %#v", repository.created)
	}
}

func TestUserService_CreateUserUsesUnitOfWorkWhenConfigured(t *testing.T) {
	repository := &fakeUserRepository{}
	unitOfWork := &fakeUnitOfWork{}
	categoryCreator := &fakeDefaultCategoryCreator{}
	service := NewUserService(repository, fakeUserIDGenerator{}, fakePasswordHasher{}, categoryCreator, unitOfWork)

	_, err := service.CreateUser(context.Background(), createUserInput("john@example.com"))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !unitOfWork.called {
		t.Fatal("expected unit of work to be called")
	}
	if !repository.usedTx {
		t.Fatal("expected repository to receive transaction")
	}
	if !categoryCreator.usedTx {
		t.Fatal("expected default category creator to receive transaction")
	}
}

func TestUserService_CreateUserAllowsNilDefaultCategoryCreator(t *testing.T) {
	service := NewUserService(&fakeUserRepository{}, fakeUserIDGenerator{}, fakePasswordHasher{}, nil, &fakeUnitOfWork{})

	_, err := service.CreateUser(context.Background(), createUserInput("john@example.com"))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func createUserInput(email string) ports.CreateUserInput {
	return ports.CreateUserInput{
		Name:          "John Doe",
		Email:         email,
		PlainPassword: "secret123",
	}
}

type fakeUserRepository struct {
	emailExists bool
	created     *domain.User
	usedTx      bool
}

func (repository *fakeUserRepository) WithTx(tx databaseports.TxHandle) ports.UserRepository {
	repository.usedTx = true
	return repository
}

func (repository *fakeUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	repository.created = user
	return user, nil
}

func (repository *fakeUserRepository) FindUserById(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	return nil, nil
}

func (repository *fakeUserRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (repository *fakeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return repository.emailExists, nil
}

type fakeUserIDGenerator struct{}

func (generator fakeUserIDGenerator) NewUserID() domain.UserID {
	return domain.UserID("user-id")
}

type fakePasswordHasher struct{}

func (hasher fakePasswordHasher) HashPassword(ctx context.Context, plainPassword string) (string, error) {
	return "hashed-" + plainPassword, nil
}

func (hasher fakePasswordHasher) ComparePassword(ctx context.Context, passwordHash string, plainPassword string) error {
	return nil
}

type fakeUnitOfWork struct {
	called bool
}

func (unit *fakeUnitOfWork) WithinTx(ctx context.Context, fn func(ctx context.Context, tx databaseports.TxHandle) error) error {
	unit.called = true
	return fn(ctx, databaseports.NewTxHandle("tx"))
}

type fakeDefaultCategoryCreator struct {
	usedTx bool
}

func (creator *fakeDefaultCategoryCreator) WithTx(tx databaseports.TxHandle) ports.DefaultCategoryCreator {
	creator.usedTx = true
	return creator
}

func (creator *fakeDefaultCategoryCreator) EnsureDefaultCategories(ctx context.Context, userID domain.UserID) error {
	return nil
}
