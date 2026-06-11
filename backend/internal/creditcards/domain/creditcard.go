package domain

import (
	"strings"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type CreditCardID string
type PurchaseID string
type InstallmentID string
type InvoiceID string

type CreditCardStatus string

const (
	CreditCardStatusActive   CreditCardStatus = "active"
	CreditCardStatusInactive CreditCardStatus = "inactive"
)

type PurchaseStatus string

const (
	PurchaseStatusActive   PurchaseStatus = "active"
	PurchaseStatusCanceled PurchaseStatus = "canceled"
)

type PurchaseType string

const (
	PurchaseTypeSingle      PurchaseType = "single"
	PurchaseTypeInstallment PurchaseType = "installment"
	PurchaseTypeFixed       PurchaseType = "fixed"
)

type InvoiceStatus string

const (
	InvoiceStatusOpen     InvoiceStatus = "open"
	InvoiceStatusClosed   InvoiceStatus = "closed"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusCanceled InvoiceStatus = "canceled"
)

type InvoiceEffectiveStatus string

const (
	InvoiceEffectiveStatusOpen     InvoiceEffectiveStatus = "open"
	InvoiceEffectiveStatusClosed   InvoiceEffectiveStatus = "closed"
	InvoiceEffectiveStatusOverdue  InvoiceEffectiveStatus = "overdue"
	InvoiceEffectiveStatusPaid     InvoiceEffectiveStatus = "paid"
	InvoiceEffectiveStatusCanceled InvoiceEffectiveStatus = "canceled"
)

type CreditCard struct {
	ID              CreditCardID
	UserID          userdomain.UserID
	Name            string
	LinkedAccountID accountdomain.AccountID
	LimitTotal      financedomain.Money
	ClosingDay      int
	DueDay          int
	Status          CreditCardStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Purchase struct {
	ID                PurchaseID
	UserID            userdomain.UserID
	CardID            CreditCardID
	CategoryID        categorydomain.CategoryID
	Description       string
	TotalAmount       financedomain.Money
	PurchaseDate      time.Time
	PurchaseType      PurchaseType
	InstallmentCount  int
	FirstInvoiceMonth time.Time
	Note              string
	Status            PurchaseStatus
	CanceledAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Installment struct {
	ID             InstallmentID
	UserID         userdomain.UserID
	CardID         CreditCardID
	PurchaseID     PurchaseID
	InvoiceID      InvoiceID
	Number         int
	Amount         financedomain.Money
	Status         PurchaseStatus
	ReferenceMonth time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Invoice struct {
	ID                   InvoiceID
	UserID               userdomain.UserID
	CardID               CreditCardID
	ReferenceMonth       time.Time
	ClosingAt            time.Time
	DueAt                time.Time
	Status               InvoiceStatus
	PaidAt               *time.Time
	PaymentTransactionID *transactiondomain.TransactionID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewCreditCard(id CreditCardID, userID userdomain.UserID, name string, linkedAccountID accountdomain.AccountID, limitTotal financedomain.Money, closingDay, dueDay int) (CreditCard, error) {
	now := time.Now()
	card := CreditCard{
		ID:              CreditCardID(strings.TrimSpace(string(id))),
		UserID:          userdomain.UserID(strings.TrimSpace(string(userID))),
		Name:            strings.TrimSpace(name),
		LinkedAccountID: accountdomain.AccountID(strings.TrimSpace(string(linkedAccountID))),
		LimitTotal:      limitTotal,
		ClosingDay:      closingDay,
		DueDay:          dueDay,
		Status:          CreditCardStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	return card, card.validate()
}

func RehydrateCreditCard(id CreditCardID, userID userdomain.UserID, name string, linkedAccountID accountdomain.AccountID, limitTotal financedomain.Money, closingDay, dueDay int, status CreditCardStatus, createdAt, updatedAt time.Time) (CreditCard, error) {
	card := CreditCard{
		ID:              CreditCardID(strings.TrimSpace(string(id))),
		UserID:          userdomain.UserID(strings.TrimSpace(string(userID))),
		Name:            strings.TrimSpace(name),
		LinkedAccountID: accountdomain.AccountID(strings.TrimSpace(string(linkedAccountID))),
		LimitTotal:      limitTotal,
		ClosingDay:      closingDay,
		DueDay:          dueDay,
		Status:          status,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	return card, card.validate()
}

func (card *CreditCard) Edit(name string, linkedAccountID accountdomain.AccountID, limitTotal financedomain.Money, closingDay, dueDay int, status CreditCardStatus) error {
	card.Name = strings.TrimSpace(name)
	card.LinkedAccountID = accountdomain.AccountID(strings.TrimSpace(string(linkedAccountID)))
	card.LimitTotal = limitTotal
	card.ClosingDay = closingDay
	card.DueDay = dueDay
	card.Status = status
	card.UpdatedAt = time.Now()
	return card.validate()
}

func (card *CreditCard) Inactivate() error {
	card.Status = CreditCardStatusInactive
	card.UpdatedAt = time.Now()
	return card.validate()
}

func (card CreditCard) validate() error {
	if strings.TrimSpace(string(card.ID)) == "" {
		return ErrCreditCardIDRequired
	}
	if strings.TrimSpace(string(card.UserID)) == "" {
		return ErrCreditCardUserIDRequired
	}
	if strings.TrimSpace(card.Name) == "" {
		return ErrCreditCardNameRequired
	}
	if strings.TrimSpace(string(card.LinkedAccountID)) == "" {
		return ErrCreditCardAccountIDRequired
	}
	if !card.LimitTotal.IsPositive() {
		return ErrCreditCardLimitInvalid
	}
	if card.ClosingDay < 1 || card.ClosingDay > 31 {
		return ErrCreditCardClosingDayInvalid
	}
	if card.DueDay < 1 || card.DueDay > 31 {
		return ErrCreditCardDueDayInvalid
	}
	if card.Status != CreditCardStatusActive && card.Status != CreditCardStatusInactive {
		return ErrCreditCardInvalidStatus
	}
	return nil
}

func NewPurchase(id PurchaseID, userID userdomain.UserID, cardID CreditCardID, categoryID categorydomain.CategoryID, description string, totalAmount financedomain.Money, purchaseDate time.Time, purchaseType PurchaseType, installmentCount int, firstInvoiceMonth time.Time, note string) (Purchase, error) {
	now := time.Now()
	purchase := Purchase{
		ID:                PurchaseID(strings.TrimSpace(string(id))),
		UserID:            userdomain.UserID(strings.TrimSpace(string(userID))),
		CardID:            CreditCardID(strings.TrimSpace(string(cardID))),
		CategoryID:        categorydomain.CategoryID(strings.TrimSpace(string(categoryID))),
		Description:       strings.TrimSpace(description),
		TotalAmount:       totalAmount,
		PurchaseDate:      purchaseDate,
		PurchaseType:      purchaseType,
		InstallmentCount:  installmentCount,
		FirstInvoiceMonth: FirstDayOfMonth(firstInvoiceMonth),
		Note:              strings.TrimSpace(note),
		Status:            PurchaseStatusActive,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	return purchase, purchase.validate()
}

func RehydratePurchase(id PurchaseID, userID userdomain.UserID, cardID CreditCardID, categoryID categorydomain.CategoryID, description string, totalAmount financedomain.Money, purchaseDate time.Time, purchaseType PurchaseType, installmentCount int, firstInvoiceMonth time.Time, note string, status PurchaseStatus, canceledAt *time.Time, createdAt, updatedAt time.Time) (Purchase, error) {
	purchase := Purchase{
		ID:                PurchaseID(strings.TrimSpace(string(id))),
		UserID:            userdomain.UserID(strings.TrimSpace(string(userID))),
		CardID:            CreditCardID(strings.TrimSpace(string(cardID))),
		CategoryID:        categorydomain.CategoryID(strings.TrimSpace(string(categoryID))),
		Description:       strings.TrimSpace(description),
		TotalAmount:       totalAmount,
		PurchaseDate:      purchaseDate,
		PurchaseType:      purchaseType,
		InstallmentCount:  installmentCount,
		FirstInvoiceMonth: FirstDayOfMonth(firstInvoiceMonth),
		Note:              strings.TrimSpace(note),
		Status:            status,
		CanceledAt:        canceledAt,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
	return purchase, purchase.validate()
}

func (purchase *Purchase) Cancel() error {
	if purchase.Status == PurchaseStatusCanceled {
		return nil
	}
	now := time.Now()
	purchase.Status = PurchaseStatusCanceled
	purchase.CanceledAt = &now
	purchase.UpdatedAt = now
	return purchase.validate()
}

func (purchase Purchase) validate() error {
	if strings.TrimSpace(string(purchase.ID)) == "" {
		return ErrPurchaseIDRequired
	}
	if strings.TrimSpace(string(purchase.UserID)) == "" {
		return ErrCreditCardUserIDRequired
	}
	if strings.TrimSpace(string(purchase.CardID)) == "" {
		return ErrCreditCardIDRequired
	}
	if strings.TrimSpace(string(purchase.CategoryID)) == "" {
		return ErrCreditCardCategoryNotFound
	}
	if strings.TrimSpace(purchase.Description) == "" {
		return ErrPurchaseDescriptionRequired
	}
	if !purchase.TotalAmount.IsPositive() {
		return ErrPurchaseAmountInvalid
	}
	if purchase.PurchaseDate.IsZero() {
		return ErrPurchaseDateRequired
	}
	if purchase.FirstInvoiceMonth.IsZero() {
		return ErrPurchaseFirstInvoiceMonthRequired
	}
	switch purchase.PurchaseType {
	case "":
		return ErrPurchaseTypeInvalid
	case PurchaseTypeSingle:
		if purchase.InstallmentCount != 1 {
			return ErrPurchaseInstallmentCountInvalid
		}
	case PurchaseTypeInstallment:
		if purchase.InstallmentCount < 2 || purchase.InstallmentCount > 12 {
			return ErrPurchaseInstallmentCountInvalid
		}
	case PurchaseTypeFixed:
		if purchase.InstallmentCount != 1 {
			return ErrPurchaseInstallmentCountInvalid
		}
	default:
		return ErrPurchaseTypeInvalid
	}
	if purchase.InstallmentCount < 1 || purchase.InstallmentCount > 12 {
		return ErrPurchaseInstallmentCountInvalid
	}
	if purchase.Status != PurchaseStatusActive && purchase.Status != PurchaseStatusCanceled {
		return ErrPurchaseInvalidStatus
	}
	if purchase.Status == PurchaseStatusActive && purchase.CanceledAt != nil {
		return ErrPurchaseInvalidStatus
	}
	if purchase.Status == PurchaseStatusCanceled && purchase.CanceledAt == nil {
		return ErrPurchaseInvalidStatus
	}
	return nil
}

func NewInstallment(id InstallmentID, userID userdomain.UserID, cardID CreditCardID, purchaseID PurchaseID, invoiceID InvoiceID, number int, amount financedomain.Money, referenceMonth time.Time) (Installment, error) {
	now := time.Now()
	installment := Installment{
		ID:             InstallmentID(strings.TrimSpace(string(id))),
		UserID:         userdomain.UserID(strings.TrimSpace(string(userID))),
		CardID:         CreditCardID(strings.TrimSpace(string(cardID))),
		PurchaseID:     PurchaseID(strings.TrimSpace(string(purchaseID))),
		InvoiceID:      InvoiceID(strings.TrimSpace(string(invoiceID))),
		Number:         number,
		Amount:         amount,
		Status:         PurchaseStatusActive,
		ReferenceMonth: FirstDayOfMonth(referenceMonth),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return installment, installment.validate()
}

func RehydrateInstallment(id InstallmentID, userID userdomain.UserID, cardID CreditCardID, purchaseID PurchaseID, invoiceID InvoiceID, number int, amount financedomain.Money, status PurchaseStatus, referenceMonth, createdAt, updatedAt time.Time) (Installment, error) {
	installment := Installment{
		ID:             InstallmentID(strings.TrimSpace(string(id))),
		UserID:         userdomain.UserID(strings.TrimSpace(string(userID))),
		CardID:         CreditCardID(strings.TrimSpace(string(cardID))),
		PurchaseID:     PurchaseID(strings.TrimSpace(string(purchaseID))),
		InvoiceID:      InvoiceID(strings.TrimSpace(string(invoiceID))),
		Number:         number,
		Amount:         amount,
		Status:         status,
		ReferenceMonth: FirstDayOfMonth(referenceMonth),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	return installment, installment.validate()
}

func (installment *Installment) Cancel() error {
	installment.Status = PurchaseStatusCanceled
	installment.UpdatedAt = time.Now()
	return installment.validate()
}

func (installment Installment) validate() error {
	if strings.TrimSpace(string(installment.ID)) == "" ||
		strings.TrimSpace(string(installment.UserID)) == "" ||
		strings.TrimSpace(string(installment.CardID)) == "" ||
		strings.TrimSpace(string(installment.PurchaseID)) == "" ||
		strings.TrimSpace(string(installment.InvoiceID)) == "" ||
		installment.Number < 1 ||
		!installment.Amount.IsPositive() ||
		installment.ReferenceMonth.IsZero() ||
		(installment.Status != PurchaseStatusActive && installment.Status != PurchaseStatusCanceled) {
		return ErrInstallmentInvalid
	}
	return nil
}

func NewInvoice(id InvoiceID, userID userdomain.UserID, cardID CreditCardID, referenceMonth, closingAt, dueAt time.Time) (Invoice, error) {
	now := time.Now()
	invoice := Invoice{
		ID:             InvoiceID(strings.TrimSpace(string(id))),
		UserID:         userdomain.UserID(strings.TrimSpace(string(userID))),
		CardID:         CreditCardID(strings.TrimSpace(string(cardID))),
		ReferenceMonth: FirstDayOfMonth(referenceMonth),
		ClosingAt:      closingAt,
		DueAt:          dueAt,
		Status:         InvoiceStatusOpen,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return invoice, invoice.validate()
}

func RehydrateInvoice(id InvoiceID, userID userdomain.UserID, cardID CreditCardID, referenceMonth, closingAt, dueAt time.Time, status InvoiceStatus, paidAt *time.Time, paymentTransactionID *transactiondomain.TransactionID, createdAt, updatedAt time.Time) (Invoice, error) {
	invoice := Invoice{
		ID:                   InvoiceID(strings.TrimSpace(string(id))),
		UserID:               userdomain.UserID(strings.TrimSpace(string(userID))),
		CardID:               CreditCardID(strings.TrimSpace(string(cardID))),
		ReferenceMonth:       FirstDayOfMonth(referenceMonth),
		ClosingAt:            closingAt,
		DueAt:                dueAt,
		Status:               status,
		PaidAt:               paidAt,
		PaymentTransactionID: trimTransactionID(paymentTransactionID),
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
	return invoice, invoice.validate()
}

func (invoice *Invoice) Close() error {
	if invoice.Status == InvoiceStatusClosed {
		return ErrInvoiceAlreadyClosed
	}
	if invoice.Status == InvoiceStatusPaid {
		return ErrInvoiceAlreadyPaid
	}
	if invoice.Status == InvoiceStatusCanceled {
		return ErrInvoiceAlreadyCanceled
	}
	invoice.Status = InvoiceStatusClosed
	invoice.UpdatedAt = time.Now()
	return invoice.validate()
}

func (invoice *Invoice) MarkPaid(transactionID transactiondomain.TransactionID, paidAt time.Time) error {
	if invoice.Status == InvoiceStatusPaid {
		return ErrInvoiceAlreadyPaid
	}
	if invoice.Status == InvoiceStatusCanceled {
		return ErrInvoiceAlreadyCanceled
	}
	if invoice.Status != InvoiceStatusOpen && invoice.Status != InvoiceStatusClosed {
		return ErrInvoiceNotPayable
	}
	if strings.TrimSpace(string(transactionID)) == "" || paidAt.IsZero() {
		return ErrInvoiceInvalidStatus
	}
	trimmed := transactiondomain.TransactionID(strings.TrimSpace(string(transactionID)))
	invoice.Status = InvoiceStatusPaid
	invoice.PaidAt = &paidAt
	invoice.PaymentTransactionID = &trimmed
	invoice.UpdatedAt = time.Now()
	return invoice.validate()
}

func (invoice Invoice) EffectiveStatus(now time.Time) InvoiceEffectiveStatus {
	switch invoice.Status {
	case InvoiceStatusPaid:
		return InvoiceEffectiveStatusPaid
	case InvoiceStatusCanceled:
		return InvoiceEffectiveStatusCanceled
	case InvoiceStatusOpen:
		if invoice.DueAt.Before(now) {
			return InvoiceEffectiveStatusOverdue
		}
		return InvoiceEffectiveStatusOpen
	case InvoiceStatusClosed:
		if invoice.DueAt.Before(now) {
			return InvoiceEffectiveStatusOverdue
		}
		return InvoiceEffectiveStatusClosed
	default:
		return InvoiceEffectiveStatusOpen
	}
}

func (invoice Invoice) validate() error {
	if strings.TrimSpace(string(invoice.ID)) == "" {
		return ErrInvoiceIDRequired
	}
	if strings.TrimSpace(string(invoice.UserID)) == "" {
		return ErrCreditCardUserIDRequired
	}
	if strings.TrimSpace(string(invoice.CardID)) == "" {
		return ErrCreditCardIDRequired
	}
	if invoice.ReferenceMonth.IsZero() {
		return ErrInvoiceReferenceMonthRequired
	}
	if invoice.DueAt.IsZero() {
		return ErrInvoiceDueAtRequired
	}
	switch invoice.Status {
	case InvoiceStatusOpen, InvoiceStatusClosed, InvoiceStatusCanceled:
		if invoice.PaidAt != nil || invoice.PaymentTransactionID != nil {
			return ErrInvoiceInvalidStatus
		}
	case InvoiceStatusPaid:
		if invoice.PaidAt == nil || invoice.PaymentTransactionID == nil {
			return ErrInvoiceInvalidStatus
		}
	default:
		return ErrInvoiceInvalidStatus
	}
	return nil
}

func FirstDayOfMonth(value time.Time) time.Time {
	location := value.Location()
	return time.Date(value.In(location).Year(), value.In(location).Month(), 1, 0, 0, 0, 0, location)
}

func CycleForPurchase(purchaseDate time.Time, closingDay, dueDay int) (referenceMonth time.Time, closingAt time.Time, dueAt time.Time) {
	location := purchaseDate.Location()
	purchaseLocal := purchaseDate.In(location)
	closingMonth := time.Date(purchaseLocal.Year(), purchaseLocal.Month(), 1, 0, 0, 0, 0, location)
	closingAt = dateWithClampedDay(closingMonth, closingDay)
	if purchaseLocal.Day() > closingAt.Day() {
		closingMonth = closingMonth.AddDate(0, 1, 0)
		closingAt = dateWithClampedDay(closingMonth, closingDay)
	}
	referenceMonth = FirstDayOfMonth(closingAt)
	dueMonth := referenceMonth
	if dueDay <= closingDay {
		dueMonth = dueMonth.AddDate(0, 1, 0)
	}
	dueAt = dateWithClampedDay(dueMonth, dueDay)
	return referenceMonth, closingAt, dueAt
}

func dateWithClampedDay(month time.Time, day int) time.Time {
	location := month.Location()
	first := time.Date(month.In(location).Year(), month.In(location).Month(), 1, 0, 0, 0, 0, location)
	lastDay := first.AddDate(0, 1, -1).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(first.Year(), first.Month(), day, 0, 0, 0, 0, location)
}

func SplitInstallments(total financedomain.Money, count int) []financedomain.Money {
	if count <= 0 {
		return nil
	}
	base := total.Cents() / int64(count)
	values := make([]financedomain.Money, 0, count)
	for index := 1; index <= count; index++ {
		amount := base
		if index == count {
			amount = total.Cents() - base*int64(count-1)
		}
		values = append(values, financedomain.NewMoney(amount))
	}
	return values
}

func trimTransactionID(value *transactiondomain.TransactionID) *transactiondomain.TransactionID {
	if value == nil {
		return nil
	}
	trimmed := transactiondomain.TransactionID(strings.TrimSpace(string(*value)))
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
