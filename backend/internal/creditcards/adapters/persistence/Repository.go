package persistence

import (
	"context"
	"errors"
	"time"

	"contai/internal/creditcards/app/ports"
	"contai/internal/creditcards/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CreditCardRepository = Repository{}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (repository Repository) WithTx(tx databaseports.TxHandle) ports.CreditCardRepository {
	if db, ok := tx.Value().(*gorm.DB); ok && db != nil {
		return Repository{db: db}
	}
	return repository
}

func (repository Repository) CreateCreditCard(ctx context.Context, card *domain.CreditCard) (*domain.CreditCard, error) {
	entity := toCreditCardEntity(*card)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}
	created, err := toDomainCreditCard(entity)
	return &created, err
}

func (repository Repository) UpdateCreditCard(ctx context.Context, card *domain.CreditCard) (*domain.CreditCard, error) {
	entity := toCreditCardEntity(*card)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrCreditCardNotFound
	}
	updated, err := toDomainCreditCard(entity)
	return &updated, err
}

func (repository Repository) FindCreditCardByID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) (*domain.CreditCard, error) {
	return repository.findCreditCardByID(ctx, cardID, userID, false)
}

func (repository Repository) FindCreditCardByIDForUpdate(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) (*domain.CreditCard, error) {
	return repository.findCreditCardByID(ctx, cardID, userID, true)
}

func (repository Repository) FindCreditCardsByUserID(ctx context.Context, userID userdomain.UserID) ([]domain.CreditCard, error) {
	var entities []CreditCardEntity
	if err := repository.db.WithContext(ctx).Where("user_id = ?", string(userID)).Order("name ASC").Find(&entities).Error; err != nil {
		return nil, err
	}
	cards := make([]domain.CreditCard, 0, len(entities))
	for _, entity := range entities {
		card, err := toDomainCreditCard(entity)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (repository Repository) CreatePurchase(ctx context.Context, purchase *domain.Purchase) (*domain.Purchase, error) {
	entity := toPurchaseEntity(*purchase)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}
	created, err := toDomainPurchase(entity)
	return &created, err
}

func (repository Repository) UpdatePurchase(ctx context.Context, purchase *domain.Purchase) (*domain.Purchase, error) {
	entity := toPurchaseEntity(*purchase)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrPurchaseNotFound
	}
	updated, err := toDomainPurchase(entity)
	return &updated, err
}

func (repository Repository) FindPurchaseByIDForUpdate(ctx context.Context, purchaseID domain.PurchaseID, userID userdomain.UserID) (*domain.Purchase, error) {
	var entity CardPurchaseEntity
	err := repository.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&entity, "id = ? AND user_id = ?", string(purchaseID), string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	purchase, err := toDomainPurchase(entity)
	return &purchase, err
}

func (repository Repository) FindPurchasesByCardID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) ([]domain.Purchase, error) {
	var entities []CardPurchaseEntity
	if err := repository.db.WithContext(ctx).Where("card_id = ? AND user_id = ?", string(cardID), string(userID)).Order("purchase_date DESC, created_at DESC").Find(&entities).Error; err != nil {
		return nil, err
	}
	purchases := make([]domain.Purchase, 0, len(entities))
	for _, entity := range entities {
		purchase, err := toDomainPurchase(entity)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	return purchases, nil
}

func (repository Repository) CreateInstallments(ctx context.Context, installments []domain.Installment) ([]domain.Installment, error) {
	entities := make([]CardInstallmentEntity, 0, len(installments))
	for _, installment := range installments {
		entities = append(entities, toInstallmentEntity(installment))
	}
	if len(entities) > 0 {
		if err := repository.db.WithContext(ctx).Create(&entities).Error; err != nil {
			return nil, err
		}
	}
	created := make([]domain.Installment, 0, len(entities))
	for _, entity := range entities {
		installment, err := toDomainInstallment(entity)
		if err != nil {
			return nil, err
		}
		created = append(created, installment)
	}
	return created, nil
}

func (repository Repository) UpdateInstallments(ctx context.Context, installments []domain.Installment) error {
	for _, installment := range installments {
		entity := toInstallmentEntity(installment)
		result := repository.db.WithContext(ctx).Save(&entity)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrInstallmentInvalid
		}
	}
	return nil
}

func (repository Repository) FindInstallmentsByPurchaseID(ctx context.Context, purchaseID domain.PurchaseID, userID userdomain.UserID) ([]domain.Installment, error) {
	var entities []CardInstallmentEntity
	if err := repository.db.WithContext(ctx).Where("purchase_id = ? AND user_id = ?", string(purchaseID), string(userID)).Order("number ASC").Find(&entities).Error; err != nil {
		return nil, err
	}
	return toDomainInstallments(entities)
}

func (repository Repository) FindInstallmentsByInvoiceID(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) ([]domain.Installment, error) {
	var entities []CardInstallmentEntity
	if err := repository.db.WithContext(ctx).Where("invoice_id = ? AND user_id = ?", string(invoiceID), string(userID)).Order("reference_month ASC, number ASC").Find(&entities).Error; err != nil {
		return nil, err
	}
	return toDomainInstallments(entities)
}

func (repository Repository) FindInvoicesByPurchaseID(ctx context.Context, purchaseID domain.PurchaseID, userID userdomain.UserID) ([]domain.Invoice, error) {
	var entities []CardInvoiceEntity
	if err := repository.db.WithContext(ctx).
		Table("card_invoices").
		Select("DISTINCT card_invoices.*").
		Joins("JOIN card_installments ON card_installments.invoice_id = card_invoices.id").
		Where("card_installments.purchase_id = ? AND card_installments.user_id = ?", string(purchaseID), string(userID)).
		Find(&entities).Error; err != nil {
		return nil, err
	}
	return toDomainInvoices(entities)
}

func (repository Repository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error) {
	entity := toInvoiceEntity(*invoice)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}
	created, err := toDomainInvoice(entity)
	return &created, err
}

func (repository Repository) UpdateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error) {
	entity := toInvoiceEntity(*invoice)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrInvoiceNotFound
	}
	updated, err := toDomainInvoice(entity)
	return &updated, err
}

func (repository Repository) FindInvoiceByID(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) (*domain.Invoice, error) {
	return repository.findInvoiceByID(ctx, invoiceID, userID, false)
}

func (repository Repository) FindInvoiceByIDForUpdate(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) (*domain.Invoice, error) {
	return repository.findInvoiceByID(ctx, invoiceID, userID, true)
}

func (repository Repository) FindInvoiceByCardAndReferenceMonth(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID, referenceMonth time.Time) (*domain.Invoice, error) {
	var entity CardInvoiceEntity
	err := repository.db.WithContext(ctx).Where("card_id = ? AND user_id = ? AND reference_month = ?", string(cardID), string(userID), domain.FirstDayOfMonth(referenceMonth)).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	invoice, err := toDomainInvoice(entity)
	return &invoice, err
}

func (repository Repository) FindInvoicesByCardID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) ([]domain.Invoice, error) {
	var entities []CardInvoiceEntity
	if err := repository.db.WithContext(ctx).Where("card_id = ? AND user_id = ?", string(cardID), string(userID)).Order("reference_month DESC").Find(&entities).Error; err != nil {
		return nil, err
	}
	return toDomainInvoices(entities)
}

func (repository Repository) SumInvoiceAmount(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) (financedomain.Money, error) {
	var total int64
	if err := repository.db.WithContext(ctx).
		Table("card_installments").
		Select("COALESCE(SUM(amount), 0)").
		Where("invoice_id = ? AND user_id = ? AND status = ?", string(invoiceID), string(userID), string(domain.PurchaseStatusActive)).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return financedomain.NewMoney(total), nil
}

func (repository Repository) SumLimitUsed(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) (financedomain.Money, error) {
	var total int64
	if err := repository.db.WithContext(ctx).
		Table("card_installments").
		Select("COALESCE(SUM(card_installments.amount), 0)").
		Joins("JOIN card_invoices ON card_invoices.id = card_installments.invoice_id").
		Where("card_installments.card_id = ? AND card_installments.user_id = ?", string(cardID), string(userID)).
		Where("card_installments.status = ?", string(domain.PurchaseStatusActive)).
		Where("card_invoices.status NOT IN ?", []string{string(domain.InvoiceStatusPaid), string(domain.InvoiceStatusCanceled)}).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return financedomain.NewMoney(total), nil
}

func (repository Repository) findCreditCardByID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID, lock bool) (*domain.CreditCard, error) {
	query := repository.db.WithContext(ctx)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var entity CreditCardEntity
	err := query.First(&entity, "id = ? AND user_id = ?", string(cardID), string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	card, err := toDomainCreditCard(entity)
	return &card, err
}

func (repository Repository) findInvoiceByID(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID, lock bool) (*domain.Invoice, error) {
	query := repository.db.WithContext(ctx)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var entity CardInvoiceEntity
	err := query.First(&entity, "id = ? AND user_id = ?", string(invoiceID), string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	invoice, err := toDomainInvoice(entity)
	return &invoice, err
}

func toDomainInstallments(entities []CardInstallmentEntity) ([]domain.Installment, error) {
	installments := make([]domain.Installment, 0, len(entities))
	for _, entity := range entities {
		installment, err := toDomainInstallment(entity)
		if err != nil {
			return nil, err
		}
		installments = append(installments, installment)
	}
	return installments, nil
}

func toDomainInvoices(entities []CardInvoiceEntity) ([]domain.Invoice, error) {
	invoices := make([]domain.Invoice, 0, len(entities))
	for _, entity := range entities {
		invoice, err := toDomainInvoice(entity)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}
