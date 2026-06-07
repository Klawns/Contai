package persistence

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	creditcarddomain "contai/internal/creditcards/domain"
	"contai/internal/dashboard/app/ports"
	dashboarddomain "contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"

	"gorm.io/gorm"
)

var _ ports.DashboardRepository = Repository{}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (repository Repository) FindActiveAccountBalances(ctx context.Context, userID userdomain.UserID) ([]ports.AccountBalanceDTO, error) {
	var rows []accountBalanceRow
	if err := repository.db.WithContext(ctx).
		Table("accounts").
		Select("id, name, type, current_balance, bank_icon_id, include_in_dashboard_total").
		Where("user_id = ? AND status = ?", string(userID), string(accountdomain.AccountStatusActive)).
		Order("name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	balances := make([]ports.AccountBalanceDTO, 0, len(rows))
	for _, row := range rows {
		balances = append(balances, ports.AccountBalanceDTO{
			AccountID:               accountdomain.AccountID(row.ID),
			Name:                    row.Name,
			Type:                    accountdomain.AccountType(row.Type),
			Balance:                 financedomain.NewMoney(row.CurrentBalance),
			BankIconID:              row.BankIconID,
			IncludeInDashboardTotal: row.IncludeInDashboardTotal,
		})
	}
	return balances, nil
}

func (repository Repository) FindCreditCards(ctx context.Context, userID userdomain.UserID, now time.Time) ([]ports.CreditCardDashboardDTO, error) {
	var cards []creditCardRow
	if err := repository.db.WithContext(ctx).
		Table("credit_cards").
		Select("id, name, linked_account_id, limit_total, closing_day, due_day").
		Where("user_id = ? AND status = ?", string(userID), string(creditcarddomain.CreditCardStatusActive)).
		Order("name ASC").
		Scan(&cards).Error; err != nil {
		return nil, err
	}

	values := make([]ports.CreditCardDashboardDTO, 0, len(cards))
	for _, card := range cards {
		limitUsed, err := repository.sumCreditCardLimitUsed(ctx, userID, card.ID)
		if err != nil {
			return nil, err
		}
		referenceMonth, _, _ := creditcarddomain.CycleForPurchase(now, card.ClosingDay, card.DueDay)
		invoice, err := repository.findCurrentInvoice(ctx, userID, card.ID, referenceMonth)
		if err != nil {
			return nil, err
		}
		value := ports.CreditCardDashboardDTO{
			CardID:          creditcarddomain.CreditCardID(card.ID),
			Name:            card.Name,
			LinkedAccountID: accountdomain.AccountID(card.LinkedAccountID),
			LimitTotal:      financedomain.NewMoney(card.LimitTotal),
			LimitUsed:       limitUsed,
			LimitAvailable:  financedomain.NewMoney(card.LimitTotal).Sub(limitUsed),
		}
		if invoice != nil {
			invoiceID := creditcarddomain.InvoiceID(invoice.ID)
			value.CurrentInvoiceID = &invoiceID
			value.CurrentInvoiceAmount = invoice.Amount
			value.CurrentInvoiceDueAt = &invoice.DueAt
			value.CurrentInvoiceEffectiveStatus = creditcarddomain.Invoice{
				ID:             creditcarddomain.InvoiceID(invoice.ID),
				UserID:         userID,
				CardID:         creditcarddomain.CreditCardID(card.ID),
				ReferenceMonth: invoice.ReferenceMonth,
				DueAt:          invoice.DueAt,
				Status:         creditcarddomain.InvoiceStatus(invoice.Status),
				PaidAt:         invoice.PaidAt,
				CreatedAt:      now,
				UpdatedAt:      now,
			}.EffectiveStatus(now)
		}
		values = append(values, value)
	}
	return values, nil
}

func (repository Repository) SumIncome(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) (financedomain.Money, error) {
	return repository.sumByType(ctx, userID, period, transactiondomain.TransactionTypeIncome)
}

func (repository Repository) SumExpense(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) (financedomain.Money, error) {
	return repository.sumByType(ctx, userID, period, transactiondomain.TransactionTypeExpense)
}

func (repository Repository) FindMonthlyIncomeExpense(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]ports.MonthlyIncomeExpenseDTO, error) {
	var rows []monthlyIncomeExpenseRow
	if err := repository.db.WithContext(ctx).
		Table("transactions").
		Select(`
			date_trunc('month', occurred_at) AS month_start_at,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS expense
		`, string(transactiondomain.TransactionTypeIncome), string(transactiondomain.TransactionTypeExpense)).
		Where("user_id = ?", string(userID)).
		Where("status = ? AND removed_at IS NULL", string(transactiondomain.TransactionStatusActive)).
		Where("type IN ?", []string{string(transactiondomain.TransactionTypeIncome), string(transactiondomain.TransactionTypeExpense)}).
		Where("occurred_at >= ? AND occurred_at <= ?", period.StartAt, period.EndAt).
		Group("date_trunc('month', occurred_at)").
		Order("month_start_at ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	values := make([]ports.MonthlyIncomeExpenseDTO, 0, len(rows))
	for _, row := range rows {
		values = append(values, ports.MonthlyIncomeExpenseDTO{
			MonthStartAt: row.MonthStartAt,
			Income:       financedomain.NewMoney(row.Income),
			Expense:      financedomain.NewMoney(row.Expense),
		})
	}
	return values, nil
}

func (repository Repository) FindMonthlyBalances(ctx context.Context, userID userdomain.UserID, monthEnds []time.Time) ([]ports.MonthlyBalanceDTO, error) {
	balances := make([]ports.MonthlyBalanceDTO, 0, len(monthEnds))
	for _, monthEnd := range monthEnds {
		var balance int64
		if err := repository.db.WithContext(ctx).Raw(`
			SELECT COALESCE(SUM(account_balance), 0)
			FROM (
				SELECT
					accounts.initial_balance + COALESCE(SUM(CASE
						WHEN transactions.type = ? AND transactions.account_id = accounts.id THEN transactions.amount
						WHEN transactions.type = ? AND transactions.account_id = accounts.id THEN -transactions.amount
						WHEN transactions.type = ? AND transactions.source_account_id = accounts.id THEN -transactions.amount
						WHEN transactions.type = ? AND transactions.destination_account_id = accounts.id THEN transactions.amount
						ELSE 0
					END), 0) AS account_balance
				FROM accounts
				LEFT JOIN transactions ON transactions.user_id = accounts.user_id
					AND transactions.status = ?
					AND transactions.removed_at IS NULL
					AND transactions.occurred_at <= ?
					AND (
						transactions.account_id = accounts.id
						OR transactions.source_account_id = accounts.id
						OR transactions.destination_account_id = accounts.id
					)
				WHERE accounts.user_id = ?
					AND accounts.status = ?
					AND accounts.include_in_dashboard_total = ?
				GROUP BY accounts.id, accounts.initial_balance
			) account_balances
		`,
			string(transactiondomain.TransactionTypeIncome),
			string(transactiondomain.TransactionTypeExpense),
			string(transactiondomain.TransactionTypeTransfer),
			string(transactiondomain.TransactionTypeTransfer),
			string(transactiondomain.TransactionStatusActive),
			monthEnd,
			string(userID),
			string(accountdomain.AccountStatusActive),
			true,
		).Scan(&balance).Error; err != nil {
			return nil, err
		}
		balances = append(balances, ports.MonthlyBalanceDTO{
			MonthEndAt: monthEnd,
			Balance:    financedomain.NewMoney(balance),
		})
	}
	return balances, nil
}

func (repository Repository) FindTransactionsByPeriod(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]transactiondomain.Transaction, error) {
	var rows []recentTransactionRow
	if err := repository.db.WithContext(ctx).
		Table("transactions").
		Select("id, user_id, type, description, amount, occurred_at, account_id, source_account_id, destination_account_id, category_id, status, note, removed_at, created_at, updated_at").
		Where("user_id = ?", string(userID)).
		Where("status = ? AND removed_at IS NULL", string(transactiondomain.TransactionStatusActive)).
		Where("occurred_at >= ? AND occurred_at <= ?", period.StartAt, period.EndAt).
		Order("occurred_at DESC, created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	transactions := make([]transactiondomain.Transaction, 0, len(rows))
	for _, row := range rows {
		transaction, err := toDomainTransaction(row)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (repository Repository) FindExpensesByCategory(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]ports.CategoryExpenseDTO, error) {
	var rows []expenseByCategoryRow
	if err := repository.db.WithContext(ctx).
		Table("transactions").
		Select("categories.id AS category_id, categories.name, categories.color, categories.icon, COALESCE(SUM(transactions.amount), 0) AS total").
		Joins("JOIN categories ON categories.id = transactions.category_id AND categories.user_id = transactions.user_id").
		Where("transactions.user_id = ?", string(userID)).
		Where("transactions.type = ?", string(transactiondomain.TransactionTypeExpense)).
		Where("transactions.status = ? AND transactions.removed_at IS NULL", string(transactiondomain.TransactionStatusActive)).
		Where("transactions.occurred_at >= ? AND transactions.occurred_at <= ?", period.StartAt, period.EndAt).
		Group("categories.id, categories.name, categories.color, categories.icon").
		Order("total DESC, categories.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	expenses := make([]ports.CategoryExpenseDTO, 0, len(rows))
	for _, row := range rows {
		expenses = append(expenses, ports.CategoryExpenseDTO{
			CategoryID: categorydomain.CategoryID(row.CategoryID),
			Name:       row.Name,
			Color:      row.Color,
			Icon:       row.Icon,
			Total:      financedomain.NewMoney(row.Total),
		})
	}
	return expenses, nil
}

func (repository Repository) FindRecentTransactions(ctx context.Context, userID userdomain.UserID, limit int) ([]ports.TransactionDTO, error) {
	var rows []recentTransactionRow
	query := repository.db.WithContext(ctx).
		Table("transactions").
		Select("id, user_id, type, description, amount, occurred_at, account_id, source_account_id, destination_account_id, category_id, status, note, removed_at, created_at, updated_at").
		Where("user_id = ?", string(userID)).
		Where("status = ? AND removed_at IS NULL", string(transactiondomain.TransactionStatusActive)).
		Order("occurred_at DESC, created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	transactions := make([]ports.TransactionDTO, 0, len(rows))
	for _, row := range rows {
		transactions = append(transactions, ports.TransactionDTO{
			ID:                   transactiondomain.TransactionID(row.ID),
			UserID:               userdomain.UserID(row.UserID),
			Type:                 transactiondomain.TransactionType(row.Type),
			Description:          row.Description,
			Amount:               financedomain.NewMoney(row.Amount),
			OccurredAt:           row.OccurredAt,
			AccountID:            toAccountID(row.AccountID),
			SourceAccountID:      toAccountID(row.SourceAccountID),
			DestinationAccountID: toAccountID(row.DestinationAccountID),
			CategoryID:           toCategoryID(row.CategoryID),
			Status:               transactiondomain.TransactionStatus(row.Status),
			Note:                 row.Note,
			RemovedAt:            row.RemovedAt,
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            row.UpdatedAt,
		})
	}
	return transactions, nil
}

func (repository Repository) sumCreditCardLimitUsed(ctx context.Context, userID userdomain.UserID, cardID string) (financedomain.Money, error) {
	var total int64
	if err := repository.db.WithContext(ctx).
		Table("card_installments").
		Select("COALESCE(SUM(card_installments.amount), 0)").
		Joins("JOIN card_invoices ON card_invoices.id = card_installments.invoice_id").
		Where("card_installments.user_id = ? AND card_installments.card_id = ?", string(userID), cardID).
		Where("card_installments.status = ?", string(creditcarddomain.PurchaseStatusActive)).
		Where("card_invoices.status NOT IN ?", []string{string(creditcarddomain.InvoiceStatusPaid), string(creditcarddomain.InvoiceStatusCanceled)}).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return financedomain.NewMoney(total), nil
}

func (repository Repository) findCurrentInvoice(ctx context.Context, userID userdomain.UserID, cardID string, referenceMonth time.Time) (*currentInvoiceRow, error) {
	var row currentInvoiceRow
	err := repository.db.WithContext(ctx).
		Table("card_invoices").
		Select("card_invoices.id, card_invoices.reference_month, card_invoices.due_at, card_invoices.status, card_invoices.paid_at, COALESCE(SUM(card_installments.amount), 0) AS amount").
		Joins("LEFT JOIN card_installments ON card_installments.invoice_id = card_invoices.id AND card_installments.status = ?", string(creditcarddomain.PurchaseStatusActive)).
		Where("card_invoices.user_id = ? AND card_invoices.card_id = ? AND card_invoices.reference_month = ?", string(userID), cardID, creditcarddomain.FirstDayOfMonth(referenceMonth)).
		Group("card_invoices.id, card_invoices.reference_month, card_invoices.due_at, card_invoices.status, card_invoices.paid_at").
		First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (repository Repository) sumByType(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period, transactionType transactiondomain.TransactionType) (financedomain.Money, error) {
	var total int64
	if err := repository.db.WithContext(ctx).
		Table("transactions").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", string(userID)).
		Where("type = ?", string(transactionType)).
		Where("status = ? AND removed_at IS NULL", string(transactiondomain.TransactionStatusActive)).
		Where("occurred_at >= ? AND occurred_at <= ?", period.StartAt, period.EndAt).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return financedomain.NewMoney(total), nil
}

type accountBalanceRow struct {
	ID                      string
	Name                    string
	Type                    string
	CurrentBalance          int64
	BankIconID              string
	IncludeInDashboardTotal bool
}

type creditCardRow struct {
	ID              string
	Name            string
	LinkedAccountID string
	LimitTotal      int64
	ClosingDay      int
	DueDay          int
}

type currentInvoiceRow struct {
	ID             string
	ReferenceMonth time.Time
	DueAt          time.Time
	Status         string
	PaidAt         *time.Time
	Amount         financedomain.Money
}

type expenseByCategoryRow struct {
	CategoryID string
	Name       string
	Color      string
	Icon       string
	Total      int64
}

type monthlyIncomeExpenseRow struct {
	MonthStartAt time.Time
	Income       int64
	Expense      int64
}

type recentTransactionRow struct {
	ID                   string
	UserID               string
	Type                 string
	Description          string
	Amount               int64
	OccurredAt           time.Time
	AccountID            *string
	SourceAccountID      *string
	DestinationAccountID *string
	CategoryID           *string
	Status               string
	Note                 string
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func toAccountID(value *string) *accountdomain.AccountID {
	if value == nil {
		return nil
	}
	converted := accountdomain.AccountID(*value)
	return &converted
}

func toCategoryID(value *string) *categorydomain.CategoryID {
	if value == nil {
		return nil
	}
	converted := categorydomain.CategoryID(*value)
	return &converted
}

func toDomainTransaction(row recentTransactionRow) (transactiondomain.Transaction, error) {
	return transactiondomain.RehydrateTransaction(
		transactiondomain.TransactionID(row.ID),
		userdomain.UserID(row.UserID),
		transactiondomain.TransactionType(row.Type),
		row.Description,
		financedomain.NewMoney(row.Amount),
		row.OccurredAt,
		toAccountID(row.AccountID),
		toAccountID(row.SourceAccountID),
		toAccountID(row.DestinationAccountID),
		toCategoryID(row.CategoryID),
		transactiondomain.TransactionStatus(row.Status),
		transactiondomain.TransactionOriginTypeManual,
		nil,
		row.Note,
		row.RemovedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
