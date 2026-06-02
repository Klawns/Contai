package persistence

import (
	"context"
	"errors"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.AccountRepository = AccountRepository{}

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return AccountRepository{db: db}
}

func (repository AccountRepository) WithTx(tx databaseports.TxHandle) ports.AccountRepository {
	if db, ok := tx.Value().(*gorm.DB); ok && db != nil {
		return AccountRepository{db: db}
	}

	return repository
}

func (repository AccountRepository) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	entity := toAccountEntity(*account)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	created, err := toDomainAccount(entity)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (repository AccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	entity := toAccountEntity(*account)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrAccountNotFound
	}

	updated, err := toDomainAccount(entity)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (repository AccountRepository) FindAccountByID(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID) (*domain.Account, error) {
	return repository.findAccountByID(ctx, accountID, userID, false)
}

func (repository AccountRepository) FindAccountByIDForUpdate(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID) (*domain.Account, error) {
	return repository.findAccountByID(ctx, accountID, userID, true)
}

func (repository AccountRepository) FindAccountsByUserID(ctx context.Context, input ports.ListAccountsInput) ([]domain.Account, error) {
	query := repository.db.WithContext(ctx).Where("user_id = ?", string(input.UserID))
	if input.Status != nil {
		query = query.Where("status = ?", string(*input.Status))
	}

	var entities []AccountEntity
	if err := query.Order("status ASC, name ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	accounts := make([]domain.Account, 0, len(entities))
	for _, entity := range entities {
		account, err := toDomainAccount(entity)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (repository AccountRepository) SumActiveAccountBalances(ctx context.Context, userID userdomain.UserID) (int64, error) {
	var total int64
	if err := repository.db.WithContext(ctx).
		Model(&AccountEntity{}).
		Where("user_id = ? AND status = ? AND include_in_dashboard_total = ?", string(userID), string(domain.AccountStatusActive), true).
		Select("COALESCE(SUM(current_balance), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (repository AccountRepository) findAccountByID(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID, lock bool) (*domain.Account, error) {
	query := repository.db.WithContext(ctx)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var entity AccountEntity
	err := query.First(&entity, "id = ? AND user_id = ?", string(accountID), string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	account, err := toDomainAccount(entity)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
