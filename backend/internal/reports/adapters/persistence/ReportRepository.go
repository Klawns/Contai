package persistence

import (
	"context"
	"errors"

	accountpersistence "contai/internal/account/adapters/persistence"
	accountdomain "contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"
	transactionpersistence "contai/internal/transactions/adapters/persistence"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"

	"gorm.io/gorm"
)

var _ reportports.ReportRepository = ReportRepository{}

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return ReportRepository{db: db}
}

func (repository ReportRepository) FindAccountByID(
	ctx context.Context,
	userID userdomain.UserID,
	accountID accountdomain.AccountID,
) (*reportports.AccountReportRow, error) {
	var entity accountpersistence.AccountEntity
	err := repository.db.WithContext(ctx).
		First(&entity, "id = ? AND user_id = ?", string(accountID), string(userID)).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	row := accountEntityToReportRow(entity)
	return &row, nil
}

func (repository ReportRepository) ListAccounts(
	ctx context.Context,
	userID userdomain.UserID,
) ([]reportports.AccountReportRow, error) {
	var entities []accountpersistence.AccountEntity
	if err := repository.db.WithContext(ctx).
		Where("user_id = ?", string(userID)).
		Order("status ASC, name ASC").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	rows := make([]reportports.AccountReportRow, 0, len(entities))
	for _, entity := range entities {
		rows = append(rows, accountEntityToReportRow(entity))
	}
	return rows, nil
}

func (repository ReportRepository) ListTransactions(
	ctx context.Context,
	input reportports.ListReportTransactionsInput,
) ([]reportports.ReportTransactionRow, error) {
	query := repository.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND occurred_at >= ? AND occurred_at <= ?",
			string(input.UserID),
			string(transactiondomain.TransactionStatusActive),
			input.StartAt,
			input.EndAt,
		)

	if input.Type != nil {
		query = query.Where("type = ?", string(*input.Type))
	}
	if input.AccountID != nil {
		accountID := string(*input.AccountID)
		query = query.Where(
			"account_id = ? OR source_account_id = ? OR destination_account_id = ?",
			accountID,
			accountID,
			accountID,
		)
	}

	var entities []transactionpersistence.TransactionEntity
	if err := query.Order("occurred_at ASC, created_at ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	accountNames, err := repository.accountNameMap(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	rows := make([]reportports.ReportTransactionRow, 0, len(entities))
	for _, entity := range entities {
		row := transactionEntityToReportRow(entity)
		row.AccountName = transactionAccountName(row, accountNames)
		rows = append(rows, row)
	}
	return rows, nil
}

func (repository ReportRepository) accountNameMap(
	ctx context.Context,
	userID userdomain.UserID,
) (map[accountdomain.AccountID]string, error) {
	accounts, err := repository.ListAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}

	names := make(map[accountdomain.AccountID]string, len(accounts))
	for _, account := range accounts {
		names[account.ID] = account.Name
	}
	return names, nil
}

func accountEntityToReportRow(entity accountpersistence.AccountEntity) reportports.AccountReportRow {
	return reportports.AccountReportRow{
		ID:                      accountdomain.AccountID(entity.ID),
		Name:                    entity.Name,
		Type:                    accountdomain.AccountType(entity.Type),
		Status:                  accountdomain.AccountStatus(entity.Status),
		InitialBalance:          financedomain.NewMoney(entity.InitialBalance),
		CurrentBalance:          financedomain.NewMoney(entity.CurrentBalance),
		IncludeInDashboardTotal: entity.IncludeInDashboardTotal,
	}
}

func transactionEntityToReportRow(entity transactionpersistence.TransactionEntity) reportports.ReportTransactionRow {
	return reportports.ReportTransactionRow{
		ID:                   transactiondomain.TransactionID(entity.ID),
		Type:                 transactiondomain.TransactionType(entity.Type),
		Description:          entity.Description,
		Amount:               financedomain.NewMoney(entity.Amount),
		OccurredAt:           entity.OccurredAt,
		AccountID:            stringToAccountID(entity.AccountID),
		SourceAccountID:      stringToAccountID(entity.SourceAccountID),
		DestinationAccountID: stringToAccountID(entity.DestinationAccountID),
	}
}

func stringToAccountID(value *string) *accountdomain.AccountID {
	if value == nil {
		return nil
	}
	converted := accountdomain.AccountID(*value)
	return &converted
}

func transactionAccountName(
	transaction reportports.ReportTransactionRow,
	accountNames map[accountdomain.AccountID]string,
) string {
	if transaction.AccountID != nil {
		return accountNames[*transaction.AccountID]
	}
	if transaction.SourceAccountID != nil && transaction.DestinationAccountID != nil {
		source := accountNames[*transaction.SourceAccountID]
		destination := accountNames[*transaction.DestinationAccountID]
		if source != "" && destination != "" {
			return source + " -> " + destination
		}
	}
	return ""
}
