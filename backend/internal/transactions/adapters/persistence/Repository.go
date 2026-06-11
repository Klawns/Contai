package persistence

import (
	"context"
	"errors"

	databaseports "contai/internal/database/ports"
	"contai/internal/transactions/app/ports"
	"contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.TransactionRepository = TransactionRepository{}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return TransactionRepository{db: db}
}

func (repository TransactionRepository) WithTx(tx databaseports.TxHandle) ports.TransactionRepository {
	if db, ok := tx.Value().(*gorm.DB); ok && db != nil {
		return TransactionRepository{db: db}
	}
	return repository
}

func (repository TransactionRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	entity := toTransactionEntity(*transaction)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}
	created, err := toDomainTransaction(entity)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (repository TransactionRepository) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	entity := toTransactionEntity(*transaction)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrTransactionNotFound
	}
	updated, err := toDomainTransaction(entity)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (repository TransactionRepository) FindTransactionByID(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID) (*domain.Transaction, error) {
	return repository.findTransactionByID(ctx, transactionID, userID, false)
}

func (repository TransactionRepository) FindTransactionByIDForUpdate(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID) (*domain.Transaction, error) {
	return repository.findTransactionByID(ctx, transactionID, userID, true)
}

func (repository TransactionRepository) FindTransactionsByUserID(ctx context.Context, input ports.ListTransactionsInput) ([]domain.Transaction, error) {
	query := repository.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", string(input.UserID), string(domain.TransactionStatusActive))
	if input.StartAt != nil {
		query = query.Where("occurred_at >= ?", *input.StartAt)
	}
	if input.EndAt != nil {
		query = query.Where("occurred_at <= ?", *input.EndAt)
	}
	if input.AccountID != nil {
		accountID := string(*input.AccountID)
		query = query.Where("account_id = ? OR source_account_id = ? OR destination_account_id = ?", accountID, accountID, accountID)
	}
	if input.AccountIDNone {
		query = query.Where("type IN ? AND account_id IS NULL", []string{string(domain.TransactionTypeIncome), string(domain.TransactionTypeExpense)})
	}
	if input.CategoryID != nil {
		query = query.Where("category_id = ?", string(*input.CategoryID))
	}
	if input.Type != nil {
		query = query.Where("type = ?", string(*input.Type))
	}
	if input.SettlementStatus != nil {
		query = query.Where("settlement_status = ?", string(*input.SettlementStatus))
	}
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}
	if input.Offset > 0 {
		query = query.Offset(input.Offset)
	}

	var entities []TransactionEntity
	if err := query.Order("occurred_at DESC, created_at DESC").Find(&entities).Error; err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, 0, len(entities))
	for _, entity := range entities {
		transaction, err := toDomainTransaction(entity)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (repository TransactionRepository) findTransactionByID(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID, lock bool) (*domain.Transaction, error) {
	query := repository.db.WithContext(ctx)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var entity TransactionEntity
	err := query.First(&entity, "id = ? AND user_id = ?", string(transactionID), string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	transaction, err := toDomainTransaction(entity)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
