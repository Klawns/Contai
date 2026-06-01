package persistence

import (
	"context"
	"errors"

	databaseports "contai/internal/database/ports"
	"contai/internal/users/app/ports"
	"contai/internal/users/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var _ ports.UserRepository = UserRepository{}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (repository UserRepository) WithTx(tx databaseports.TxHandle) ports.UserRepository {
	if db, ok := tx.Value().(*gorm.DB); ok && db != nil {
		return UserRepository{db: db}
	}

	return repository
}

func (repository UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	entity := toUserEntity(*user)

	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrUserEmailAlreadyExists
		}
		return nil, err
	}

	created, err := toDomainUser(entity)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (repository UserRepository) FindUserById(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var entity UserEntity
	err := repository.db.WithContext(ctx).First(&entity, "id = ?", string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user, err := toDomainUser(entity)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository UserRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var entity UserEntity
	err := repository.db.WithContext(ctx).First(&entity, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user, err := toDomainUser(entity)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := repository.db.WithContext(ctx).Model(&UserEntity{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func isUniqueViolation(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
