package persistence

import (
	"context"
	"sort"
	"strings"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	creditcarddomain "contai/internal/creditcards/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"

	"gorm.io/gorm"
)

var _ reportports.ReportRepository = ReportRepository{}

type ReportRepository struct {
	db *gorm.DB
}

type movementRow struct {
	ID               string
	Source           string
	Type             string
	Description      string
	Amount           int64
	OccurredAt       time.Time
	CategoryID       *string
	CategoryName     *string
	AccountID        *string
	AccountName      *string
	SettlementStatus string
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return ReportRepository{db: db}
}

func (repository ReportRepository) ListFinancialMovements(
	ctx context.Context,
	input reportports.ListFinancialMovementsInput,
) ([]reportports.FinancialMovementDTO, error) {
	movements := make([]reportports.FinancialMovementDTO, 0)

	if input.MovementType == "" || input.MovementType == reportports.MovementTypeAll ||
		input.MovementType == reportports.MovementTypeIncome ||
		input.MovementType == reportports.MovementTypeExpense ||
		input.MovementType == reportports.MovementTypeTransfer {
		transactionMovements, err := repository.listTransactionMovements(ctx, input)
		if err != nil {
			return nil, err
		}
		movements = append(movements, transactionMovements...)
	}

	if input.MovementType == "" || input.MovementType == reportports.MovementTypeAll ||
		input.MovementType == reportports.MovementTypeCreditCardExpense {
		cardMovements, err := repository.listCardInstallmentMovements(ctx, input)
		if err != nil {
			return nil, err
		}
		movements = append(movements, cardMovements...)
	}

	sort.SliceStable(movements, func(i, j int) bool {
		if movements[i].OccurredAt.Equal(movements[j].OccurredAt) {
			return movements[i].ID < movements[j].ID
		}
		return movements[i].OccurredAt.Before(movements[j].OccurredAt)
	})
	return movements, nil
}

func (repository ReportRepository) listTransactionMovements(
	ctx context.Context,
	input reportports.ListFinancialMovementsInput,
) ([]reportports.FinancialMovementDTO, error) {
	query := repository.db.WithContext(ctx).
		Table("transactions AS t").
		Select(strings.Join([]string{
			"t.id AS id",
			"t.type AS source",
			"t.type AS type",
			"t.description AS description",
			"t.amount AS amount",
			"t.occurred_at AS occurred_at",
			"t.category_id AS category_id",
			"cat.name AS category_name",
			"COALESCE(t.account_id, t.source_account_id) AS account_id",
			"CASE WHEN t.type = ? THEN source_account.name || ' -> ' || destination_account.name ELSE account.name END AS account_name",
			"t.settlement_status AS settlement_status",
		}, ", "), string(transactiondomain.TransactionTypeTransfer)).
		Joins("LEFT JOIN categories AS cat ON cat.id = t.category_id AND cat.user_id = t.user_id").
		Joins("LEFT JOIN accounts AS account ON account.id = t.account_id AND account.user_id = t.user_id").
		Joins("LEFT JOIN accounts AS source_account ON source_account.id = t.source_account_id AND source_account.user_id = t.user_id").
		Joins("LEFT JOIN accounts AS destination_account ON destination_account.id = t.destination_account_id AND destination_account.user_id = t.user_id").
		Where("t.user_id = ? AND t.status = ? AND t.occurred_at >= ? AND t.occurred_at <= ?",
			string(input.UserID),
			string(transactiondomain.TransactionStatusActive),
			input.StartAt,
			input.EndAt,
		).
		Where("t.origin_type <> ?", string(transactiondomain.TransactionOriginTypeCreditCardInvoice))

	if input.MovementType != "" && input.MovementType != reportports.MovementTypeAll && input.MovementType != reportports.MovementTypeCreditCardExpense {
		query = query.Where("t.type = ?", string(input.MovementType))
	}
	if input.CategoryID != nil {
		query = query.Where("t.category_id = ?", string(*input.CategoryID))
	}
	if input.AccountID != nil {
		accountID := string(*input.AccountID)
		query = query.Where("t.account_id = ? OR t.source_account_id = ? OR t.destination_account_id = ?", accountID, accountID, accountID)
	}
	if input.SettlementStatus != "" && input.SettlementStatus != reportports.SettlementStatusAll {
		query = query.Where("t.settlement_status = ?", string(input.SettlementStatus))
	}

	var rows []movementRow
	if err := query.Order("t.occurred_at ASC, t.created_at ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return movementRowsToDTO(rows), nil
}

func (repository ReportRepository) listCardInstallmentMovements(
	ctx context.Context,
	input reportports.ListFinancialMovementsInput,
) ([]reportports.FinancialMovementDTO, error) {
	query := repository.db.WithContext(ctx).
		Table("card_installments AS installment").
		Select(strings.Join([]string{
			"installment.id AS id",
			"? AS source",
			"? AS type",
			"purchase.description AS description",
			"installment.amount AS amount",
			"purchase.purchase_date AS occurred_at",
			"purchase.category_id AS category_id",
			"cat.name AS category_name",
			"card.linked_account_id AS account_id",
			"account.name AS account_name",
			"CASE WHEN invoice.status = ? THEN ? ELSE ? END AS settlement_status",
		}, ", "),
			string(reportports.MovementTypeCreditCardExpense),
			string(reportports.MovementTypeCreditCardExpense),
			string(creditcarddomain.InvoiceStatusPaid),
			string(transactiondomain.SettlementStatusSettled),
			string(transactiondomain.SettlementStatusPending),
		).
		Joins("JOIN card_purchases AS purchase ON purchase.id = installment.purchase_id AND purchase.user_id = installment.user_id").
		Joins("JOIN card_invoices AS invoice ON invoice.id = installment.invoice_id AND invoice.user_id = installment.user_id").
		Joins("JOIN credit_cards AS card ON card.id = installment.card_id AND card.user_id = installment.user_id").
		Joins("LEFT JOIN categories AS cat ON cat.id = purchase.category_id AND cat.user_id = installment.user_id").
		Joins("LEFT JOIN accounts AS account ON account.id = card.linked_account_id AND account.user_id = installment.user_id").
		Where("installment.user_id = ? AND installment.status = ? AND purchase.status = ? AND invoice.status <> ? AND purchase.purchase_date >= ? AND purchase.purchase_date <= ?",
			string(input.UserID),
			string(creditcarddomain.PurchaseStatusActive),
			string(creditcarddomain.PurchaseStatusActive),
			string(creditcarddomain.InvoiceStatusCanceled),
			input.StartAt,
			input.EndAt,
		)

	if input.CategoryID != nil {
		query = query.Where("purchase.category_id = ?", string(*input.CategoryID))
	}
	if input.AccountID != nil {
		query = query.Where("card.linked_account_id = ?", string(*input.AccountID))
	}
	if input.SettlementStatus == reportports.SettlementStatusSettled {
		query = query.Where("invoice.status = ?", string(creditcarddomain.InvoiceStatusPaid))
	} else if input.SettlementStatus == reportports.SettlementStatusPending {
		query = query.Where("invoice.status <> ?", string(creditcarddomain.InvoiceStatusPaid))
	}

	var rows []movementRow
	if err := query.Order("purchase.purchase_date ASC, installment.number ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return movementRowsToDTO(rows), nil
}

func movementRowsToDTO(rows []movementRow) []reportports.FinancialMovementDTO {
	movements := make([]reportports.FinancialMovementDTO, 0, len(rows))
	for _, row := range rows {
		movements = append(movements, reportports.FinancialMovementDTO{
			ID:               row.ID,
			Source:           reportports.MovementType(row.Source),
			Type:             reportports.MovementType(row.Type),
			Description:      row.Description,
			Amount:           financedomain.NewMoney(row.Amount),
			OccurredAt:       row.OccurredAt,
			CategoryID:       stringToCategoryID(row.CategoryID),
			CategoryName:     stringValue(row.CategoryName),
			AccountID:        stringToAccountID(row.AccountID),
			AccountName:      stringValue(row.AccountName),
			SettlementStatus: transactiondomain.SettlementStatus(row.SettlementStatus),
		})
	}
	return movements
}

func stringToAccountID(value *string) *accountdomain.AccountID {
	if value == nil {
		return nil
	}
	converted := accountdomain.AccountID(*value)
	return &converted
}

func stringToCategoryID(value *string) *categorydomain.CategoryID {
	if value == nil {
		return nil
	}
	converted := categorydomain.CategoryID(*value)
	return &converted
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
