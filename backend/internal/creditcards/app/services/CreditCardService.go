package services

import (
	"context"
	"fmt"
	"time"

	accountports "contai/internal/account/app/ports"
	accountdomain "contai/internal/account/domain"
	categoryports "contai/internal/category/app/ports"
	categorydomain "contai/internal/category/domain"
	"contai/internal/creditcards/app/ports"
	"contai/internal/creditcards/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	transactionports "contai/internal/transactions/app/ports"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

var _ ports.CreditCardService = CreditCardService{}

type CreditCardService struct {
	repository             ports.CreditCardRepository
	accountRepository      accountports.AccountRepository
	categoryRepository     categoryports.CategoryRepository
	transactionRepository  transactionports.TransactionRepository
	idGenerator            ports.CreditCardIDGenerator
	transactionIDGenerator transactionports.TransactionIDGenerator
	unitOfWork             databaseports.UnitOfWork
}

func NewCreditCardService(repository ports.CreditCardRepository, accountRepository accountports.AccountRepository, categoryRepository categoryports.CategoryRepository, transactionRepository transactionports.TransactionRepository, idGenerator ports.CreditCardIDGenerator, transactionIDGenerator transactionports.TransactionIDGenerator, unitOfWork databaseports.UnitOfWork) CreditCardService {
	return CreditCardService{
		repository:             repository,
		accountRepository:      accountRepository,
		categoryRepository:     categoryRepository,
		transactionRepository:  transactionRepository,
		idGenerator:            idGenerator,
		transactionIDGenerator: transactionIDGenerator,
		unitOfWork:             unitOfWork,
	}
}

func (service CreditCardService) ListCreditCards(ctx context.Context, userID userdomain.UserID) ([]ports.CreditCardDTO, error) {
	if userID == "" {
		return nil, domain.ErrCreditCardUserIDRequired
	}
	cards, err := service.repository.FindCreditCardsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return service.toCreditCardDTOs(ctx, cards)
}

func (service CreditCardService) CreateCreditCard(ctx context.Context, input ports.CreateCreditCardInput) (ports.CreditCardDTO, error) {
	card, err := domain.NewCreditCard(service.idGenerator.NewCreditCardID(), input.UserID, input.Name, input.LinkedAccountID, input.LimitTotal, input.ClosingDay, input.DueDay)
	if err != nil {
		return ports.CreditCardDTO{}, err
	}
	if err := service.validateAccount(ctx, service.accountRepository, card.UserID, card.LinkedAccountID); err != nil {
		return ports.CreditCardDTO{}, err
	}
	created, err := service.repository.CreateCreditCard(ctx, &card)
	if err != nil {
		return ports.CreditCardDTO{}, err
	}
	return service.toCreditCardDTO(ctx, *created)
}

func (service CreditCardService) UpdateCreditCard(ctx context.Context, input ports.UpdateCreditCardInput) (ports.CreditCardDTO, error) {
	var dto ports.CreditCardDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		repository := service.repository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)

		card, err := repository.FindCreditCardByIDForUpdate(txCtx, input.CardID, input.UserID)
		if err != nil {
			return err
		}
		if card == nil {
			return domain.ErrCreditCardNotFound
		}
		status := input.Status
		if status == "" {
			status = card.Status
		}
		if err := card.Edit(input.Name, input.LinkedAccountID, input.LimitTotal, input.ClosingDay, input.DueDay, status); err != nil {
			return err
		}
		if err := service.validateAccount(txCtx, accountRepository, card.UserID, card.LinkedAccountID); err != nil {
			return err
		}
		used, err := repository.SumLimitUsed(txCtx, card.ID, card.UserID)
		if err != nil {
			return err
		}
		if card.LimitTotal.Cents() < used.Cents() {
			return domain.ErrCreditCardLimitExceeded
		}
		updated, err := repository.UpdateCreditCard(txCtx, card)
		if err != nil {
			return err
		}
		dto = toCreditCardDTO(*updated, used)
		return nil
	})
	if err != nil {
		return ports.CreditCardDTO{}, err
	}
	return dto, nil
}

func (service CreditCardService) InactivateCreditCard(ctx context.Context, input ports.CardIDInput) (ports.CreditCardDTO, error) {
	var dto ports.CreditCardDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		repository := service.repository.WithTx(tx)
		card, err := repository.FindCreditCardByIDForUpdate(txCtx, input.CardID, input.UserID)
		if err != nil {
			return err
		}
		if card == nil {
			return domain.ErrCreditCardNotFound
		}
		if err := card.Inactivate(); err != nil {
			return err
		}
		updated, err := repository.UpdateCreditCard(txCtx, card)
		if err != nil {
			return err
		}
		used, err := repository.SumLimitUsed(txCtx, updated.ID, updated.UserID)
		if err != nil {
			return err
		}
		dto = toCreditCardDTO(*updated, used)
		return nil
	})
	if err != nil {
		return ports.CreditCardDTO{}, err
	}
	return dto, nil
}

func (service CreditCardService) ListPurchases(ctx context.Context, input ports.CardIDInput) ([]ports.PurchaseDTO, error) {
	if _, err := service.findCard(ctx, service.repository, input.CardID, input.UserID); err != nil {
		return nil, err
	}
	purchases, err := service.repository.FindPurchasesByCardID(ctx, input.CardID, input.UserID)
	if err != nil {
		return nil, err
	}
	return toPurchaseDTOs(purchases), nil
}

func (service CreditCardService) CreatePurchase(ctx context.Context, input ports.CreatePurchaseInput) (ports.PurchaseDTO, error) {
	var dto ports.PurchaseDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		repository := service.repository.WithTx(tx)
		categoryRepository := service.categoryRepository.WithTx(tx)

		card, err := repository.FindCreditCardByIDForUpdate(txCtx, input.CardID, input.UserID)
		if err != nil {
			return err
		}
		if card == nil || card.Status != domain.CreditCardStatusActive {
			return domain.ErrCreditCardNotFound
		}
		if err := service.validateExpenseCategory(txCtx, categoryRepository, input.UserID, input.CategoryID); err != nil {
			return err
		}
		purchase, err := domain.NewPurchase(service.idGenerator.NewPurchaseID(), input.UserID, input.CardID, input.CategoryID, input.Description, input.TotalAmount, input.PurchaseDate, input.InstallmentCount, input.Note)
		if err != nil {
			return err
		}
		used, err := repository.SumLimitUsed(txCtx, card.ID, card.UserID)
		if err != nil {
			return err
		}
		if used.Add(purchase.TotalAmount).Cents() > card.LimitTotal.Cents() {
			return domain.ErrCreditCardLimitExceeded
		}
		createdPurchase, err := repository.CreatePurchase(txCtx, &purchase)
		if err != nil {
			return err
		}
		installments, err := service.buildInstallments(txCtx, repository, *card, *createdPurchase)
		if err != nil {
			return err
		}
		if _, err := repository.CreateInstallments(txCtx, installments); err != nil {
			return err
		}
		dto = toPurchaseDTO(*createdPurchase)
		return nil
	})
	if err != nil {
		return ports.PurchaseDTO{}, err
	}
	return dto, nil
}

func (service CreditCardService) CancelPurchase(ctx context.Context, input ports.PurchaseIDInput) (ports.PurchaseDTO, error) {
	var dto ports.PurchaseDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		repository := service.repository.WithTx(tx)
		purchase, err := repository.FindPurchaseByIDForUpdate(txCtx, input.PurchaseID, input.UserID)
		if err != nil {
			return err
		}
		if purchase == nil {
			return domain.ErrPurchaseNotFound
		}
		invoices, err := repository.FindInvoicesByPurchaseID(txCtx, purchase.ID, purchase.UserID)
		if err != nil {
			return err
		}
		for _, invoice := range invoices {
			if invoice.Status == domain.InvoiceStatusPaid {
				return domain.ErrPurchaseInvoiceAlreadyPaid
			}
		}
		if err := purchase.Cancel(); err != nil {
			return err
		}
		installments, err := repository.FindInstallmentsByPurchaseID(txCtx, purchase.ID, purchase.UserID)
		if err != nil {
			return err
		}
		for index := range installments {
			if err := installments[index].Cancel(); err != nil {
				return err
			}
		}
		if err := repository.UpdateInstallments(txCtx, installments); err != nil {
			return err
		}
		updated, err := repository.UpdatePurchase(txCtx, purchase)
		if err != nil {
			return err
		}
		dto = toPurchaseDTO(*updated)
		return nil
	})
	if err != nil {
		return ports.PurchaseDTO{}, err
	}
	return dto, nil
}

func (service CreditCardService) ListInvoices(ctx context.Context, input ports.CardIDInput) ([]ports.InvoiceDTO, error) {
	if _, err := service.findCard(ctx, service.repository, input.CardID, input.UserID); err != nil {
		return nil, err
	}
	invoices, err := service.repository.FindInvoicesByCardID(ctx, input.CardID, input.UserID)
	if err != nil {
		return nil, err
	}
	return service.toInvoiceDTOs(ctx, invoices)
}

func (service CreditCardService) GetInvoice(ctx context.Context, input ports.InvoiceIDInput) (ports.InvoiceDTO, error) {
	invoice, err := service.repository.FindInvoiceByID(ctx, input.InvoiceID, input.UserID)
	if err != nil {
		return ports.InvoiceDTO{}, err
	}
	if invoice == nil {
		return ports.InvoiceDTO{}, domain.ErrInvoiceNotFound
	}
	return service.toInvoiceDTO(ctx, *invoice)
}

func (service CreditCardService) CloseInvoice(ctx context.Context, input ports.InvoiceIDInput) (ports.InvoiceDTO, error) {
	var dto ports.InvoiceDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		repository := service.repository.WithTx(tx)
		invoice, err := repository.FindInvoiceByIDForUpdate(txCtx, input.InvoiceID, input.UserID)
		if err != nil {
			return err
		}
		if invoice == nil {
			return domain.ErrInvoiceNotFound
		}
		if err := invoice.Close(); err != nil {
			return err
		}
		updated, err := repository.UpdateInvoice(txCtx, invoice)
		if err != nil {
			return err
		}
		dto, err = service.toInvoiceDTOWithRepository(txCtx, repository, *updated)
		return err
	})
	if err != nil {
		return ports.InvoiceDTO{}, err
	}
	return dto, nil
}

func (service CreditCardService) PayInvoice(ctx context.Context, input ports.PayInvoiceInput) (ports.InvoiceDTO, error) {
	var dto ports.InvoiceDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		repository := service.repository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)
		categoryRepository := service.categoryRepository.WithTx(tx)
		transactionRepository := service.transactionRepository.WithTx(tx)

		invoice, err := repository.FindInvoiceByIDForUpdate(txCtx, input.InvoiceID, input.UserID)
		if err != nil {
			return err
		}
		if invoice == nil {
			return domain.ErrInvoiceNotFound
		}
		card, err := repository.FindCreditCardByIDForUpdate(txCtx, invoice.CardID, invoice.UserID)
		if err != nil {
			return err
		}
		if card == nil {
			return domain.ErrCreditCardNotFound
		}
		if err := service.validateExpenseCategory(txCtx, categoryRepository, input.UserID, input.CategoryID); err != nil {
			return err
		}
		amount, err := repository.SumInvoiceAmount(txCtx, invoice.ID, invoice.UserID)
		if err != nil {
			return err
		}
		if !amount.IsPositive() {
			return domain.ErrInvoiceNotPayable
		}
		account, err := accountRepository.FindAccountByIDForUpdate(txCtx, card.LinkedAccountID, card.UserID)
		if err != nil {
			return err
		}
		if account == nil || account.Status != accountdomain.AccountStatusActive {
			return domain.ErrCreditCardAccountNotFound
		}
		transaction, err := transactiondomain.NewExpense(service.transactionIDGenerator.NewTransactionID(), input.UserID, fmt.Sprintf("Fatura %s", card.Name), amount, input.OccurredAt, card.LinkedAccountID, input.CategoryID, input.Note)
		if err != nil {
			return err
		}
		if err := transaction.SetOrigin(transactiondomain.TransactionOriginTypeCreditCardInvoice, string(invoice.ID)); err != nil {
			return err
		}
		if err := account.DecreaseBalance(amount); err != nil {
			return err
		}
		if _, err := accountRepository.UpdateAccount(txCtx, account); err != nil {
			return err
		}
		createdTransaction, err := transactionRepository.CreateTransaction(txCtx, &transaction)
		if err != nil {
			return err
		}
		if err := invoice.MarkPaid(createdTransaction.ID, input.OccurredAt); err != nil {
			return err
		}
		updated, err := repository.UpdateInvoice(txCtx, invoice)
		if err != nil {
			return err
		}
		dto, err = service.toInvoiceDTOWithRepository(txCtx, repository, *updated)
		return err
	})
	if err != nil {
		return ports.InvoiceDTO{}, err
	}
	return dto, nil
}

func (service CreditCardService) buildInstallments(ctx context.Context, repository ports.CreditCardRepository, card domain.CreditCard, purchase domain.Purchase) ([]domain.Installment, error) {
	amounts := domain.SplitInstallments(purchase.TotalAmount, purchase.InstallmentCount)
	installments := make([]domain.Installment, 0, len(amounts))
	for index, amount := range amounts {
		installmentDate := purchase.PurchaseDate.AddDate(0, index, 0)
		referenceMonth, closingAt, dueAt := domain.CycleForPurchase(installmentDate, card.ClosingDay, card.DueDay)
		invoice, err := repository.FindInvoiceByCardAndReferenceMonth(ctx, card.ID, card.UserID, referenceMonth)
		if err != nil {
			return nil, err
		}
		if invoice == nil {
			newInvoice, err := domain.NewInvoice(service.idGenerator.NewInvoiceID(), card.UserID, card.ID, referenceMonth, closingAt, dueAt)
			if err != nil {
				return nil, err
			}
			invoice, err = repository.CreateInvoice(ctx, &newInvoice)
			if err != nil {
				return nil, err
			}
		}
		installment, err := domain.NewInstallment(service.idGenerator.NewInstallmentID(), purchase.UserID, purchase.CardID, purchase.ID, invoice.ID, index+1, amount, referenceMonth)
		if err != nil {
			return nil, err
		}
		installments = append(installments, installment)
	}
	return installments, nil
}

func (service CreditCardService) validateAccount(ctx context.Context, repository accountports.AccountRepository, userID userdomain.UserID, accountID accountdomain.AccountID) error {
	account, err := repository.FindAccountByID(ctx, accountID, userID)
	if err != nil {
		return err
	}
	if account == nil || account.Status != accountdomain.AccountStatusActive {
		return domain.ErrCreditCardAccountNotFound
	}
	return nil
}

func (service CreditCardService) validateExpenseCategory(ctx context.Context, repository categoryports.CategoryRepository, userID userdomain.UserID, categoryID categorydomain.CategoryID) error {
	category, err := repository.FindCategoryByID(ctx, categoryID, userID)
	if err != nil {
		return err
	}
	if category == nil || category.Status != categorydomain.CategoryStatusActive {
		return domain.ErrCreditCardCategoryNotFound
	}
	if category.Type != categorydomain.CategoryTypeExpense {
		return domain.ErrCreditCardCategoryTypeMismatch
	}
	return nil
}

func (service CreditCardService) findCard(ctx context.Context, repository ports.CreditCardRepository, cardID domain.CreditCardID, userID userdomain.UserID) (*domain.CreditCard, error) {
	card, err := repository.FindCreditCardByID(ctx, cardID, userID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, domain.ErrCreditCardNotFound
	}
	return card, nil
}

func (service CreditCardService) toCreditCardDTOs(ctx context.Context, cards []domain.CreditCard) ([]ports.CreditCardDTO, error) {
	dtos := make([]ports.CreditCardDTO, 0, len(cards))
	for _, card := range cards {
		dto, err := service.toCreditCardDTO(ctx, card)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (service CreditCardService) toCreditCardDTO(ctx context.Context, card domain.CreditCard) (ports.CreditCardDTO, error) {
	used, err := service.repository.SumLimitUsed(ctx, card.ID, card.UserID)
	if err != nil {
		return ports.CreditCardDTO{}, err
	}
	return toCreditCardDTO(card, used), nil
}

func (service CreditCardService) toInvoiceDTOs(ctx context.Context, invoices []domain.Invoice) ([]ports.InvoiceDTO, error) {
	dtos := make([]ports.InvoiceDTO, 0, len(invoices))
	for _, invoice := range invoices {
		dto, err := service.toInvoiceDTO(ctx, invoice)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (service CreditCardService) toInvoiceDTO(ctx context.Context, invoice domain.Invoice) (ports.InvoiceDTO, error) {
	return service.toInvoiceDTOWithRepository(ctx, service.repository, invoice)
}

func (service CreditCardService) toInvoiceDTOWithRepository(ctx context.Context, repository ports.CreditCardRepository, invoice domain.Invoice) (ports.InvoiceDTO, error) {
	amount, err := repository.SumInvoiceAmount(ctx, invoice.ID, invoice.UserID)
	if err != nil {
		return ports.InvoiceDTO{}, err
	}
	installments, err := repository.FindInstallmentsByInvoiceID(ctx, invoice.ID, invoice.UserID)
	if err != nil {
		return ports.InvoiceDTO{}, err
	}
	return toInvoiceDTO(invoice, amount, installments, time.Now()), nil
}

func (service CreditCardService) withinTx(ctx context.Context, fn func(context.Context, databaseports.TxHandle) error) error {
	if service.unitOfWork == nil {
		return fn(ctx, databaseports.NewTxHandle(nil))
	}
	return service.unitOfWork.WithinTx(ctx, fn)
}

func toCreditCardDTO(card domain.CreditCard, used financedomain.Money) ports.CreditCardDTO {
	return ports.CreditCardDTO{
		ID:              card.ID,
		UserID:          card.UserID,
		Name:            card.Name,
		LinkedAccountID: card.LinkedAccountID,
		LimitTotal:      card.LimitTotal,
		LimitUsed:       used,
		LimitAvailable:  card.LimitTotal.Sub(used),
		ClosingDay:      card.ClosingDay,
		DueDay:          card.DueDay,
		Status:          card.Status,
		CreatedAt:       card.CreatedAt,
		UpdatedAt:       card.UpdatedAt,
	}
}

func toPurchaseDTO(purchase domain.Purchase) ports.PurchaseDTO {
	return ports.PurchaseDTO{
		ID:               purchase.ID,
		UserID:           purchase.UserID,
		CardID:           purchase.CardID,
		CategoryID:       purchase.CategoryID,
		Description:      purchase.Description,
		TotalAmount:      purchase.TotalAmount,
		PurchaseDate:     purchase.PurchaseDate,
		InstallmentCount: purchase.InstallmentCount,
		Note:             purchase.Note,
		Status:           purchase.Status,
		CanceledAt:       purchase.CanceledAt,
		CreatedAt:        purchase.CreatedAt,
		UpdatedAt:        purchase.UpdatedAt,
	}
}

func toPurchaseDTOs(purchases []domain.Purchase) []ports.PurchaseDTO {
	dtos := make([]ports.PurchaseDTO, 0, len(purchases))
	for _, purchase := range purchases {
		dtos = append(dtos, toPurchaseDTO(purchase))
	}
	return dtos
}

func toInvoiceDTO(invoice domain.Invoice, amount financedomain.Money, installments []domain.Installment, now time.Time) ports.InvoiceDTO {
	return ports.InvoiceDTO{
		ID:                   invoice.ID,
		UserID:               invoice.UserID,
		CardID:               invoice.CardID,
		ReferenceMonth:       invoice.ReferenceMonth,
		ClosingAt:            invoice.ClosingAt,
		DueAt:                invoice.DueAt,
		Amount:               amount,
		Status:               invoice.Status,
		EffectiveStatus:      invoice.EffectiveStatus(now),
		PaidAt:               invoice.PaidAt,
		PaymentTransactionID: invoice.PaymentTransactionID,
		Installments:         toInstallmentDTOs(installments),
		CreatedAt:            invoice.CreatedAt,
		UpdatedAt:            invoice.UpdatedAt,
	}
}

func toInstallmentDTOs(installments []domain.Installment) []ports.InstallmentDTO {
	dtos := make([]ports.InstallmentDTO, 0, len(installments))
	for _, installment := range installments {
		dtos = append(dtos, ports.InstallmentDTO{
			ID:             installment.ID,
			UserID:         installment.UserID,
			CardID:         installment.CardID,
			PurchaseID:     installment.PurchaseID,
			InvoiceID:      installment.InvoiceID,
			Number:         installment.Number,
			Amount:         installment.Amount,
			Status:         installment.Status,
			ReferenceMonth: installment.ReferenceMonth,
			CreatedAt:      installment.CreatedAt,
			UpdatedAt:      installment.UpdatedAt,
		})
	}
	return dtos
}
