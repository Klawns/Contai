package services

import (
	"context"

	accountports "contai/internal/account/app/ports"
	accountdomain "contai/internal/account/domain"
	categoryports "contai/internal/category/app/ports"
	categorydomain "contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	"contai/internal/transactions/app/ports"
	"contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

var _ ports.TransactionService = TransactionService{}

type TransactionService struct {
	transactionRepository ports.TransactionRepository
	accountRepository     accountports.AccountRepository
	categoryRepository    categoryports.CategoryRepository
	idGenerator           ports.TransactionIDGenerator
	unitOfWork            databaseports.UnitOfWork
}

func NewTransactionService(transactionRepository ports.TransactionRepository, accountRepository accountports.AccountRepository, categoryRepository categoryports.CategoryRepository, idGenerator ports.TransactionIDGenerator, unitOfWork databaseports.UnitOfWork) TransactionService {
	return TransactionService{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		categoryRepository:    categoryRepository,
		idGenerator:           idGenerator,
		unitOfWork:            unitOfWork,
	}
}

func (service TransactionService) CreateIncome(ctx context.Context, input ports.CreateIncomeInput) (ports.TransactionDTO, error) {
	transaction, err := domain.NewIncome(service.idGenerator.NewTransactionID(), input.UserID, input.Description, input.Amount, input.OccurredAt, input.AccountID, input.CategoryID, input.Note)
	if err != nil {
		return ports.TransactionDTO{}, err
	}
	return service.create(ctx, transaction)
}

func (service TransactionService) CreateExpense(ctx context.Context, input ports.CreateExpenseInput) (ports.TransactionDTO, error) {
	transaction, err := domain.NewExpense(service.idGenerator.NewTransactionID(), input.UserID, input.Description, input.Amount, input.OccurredAt, input.AccountID, input.CategoryID, input.Note)
	if err != nil {
		return ports.TransactionDTO{}, err
	}
	return service.create(ctx, transaction)
}

func (service TransactionService) CreateTransfer(ctx context.Context, input ports.CreateTransferInput) (ports.TransactionDTO, error) {
	transaction, err := domain.NewTransfer(service.idGenerator.NewTransactionID(), input.UserID, input.Description, input.Amount, input.OccurredAt, input.SourceAccountID, input.DestinationAccountID, input.Note)
	if err != nil {
		return ports.TransactionDTO{}, err
	}
	return service.create(ctx, transaction)
}

func (service TransactionService) ListTransactions(ctx context.Context, input ports.ListTransactionsInput) ([]ports.TransactionDTO, error) {
	if input.UserID == "" {
		return nil, domain.ErrTransactionUserIDRequired
	}
	if input.Type != nil && !isValidTransactionType(*input.Type) {
		return nil, domain.ErrTransactionInvalidType
	}
	if input.Limit < 0 || input.Offset < 0 {
		return nil, domain.ErrTransactionInvalidType
	}

	transactions, err := service.transactionRepository.FindTransactionsByUserID(ctx, input)
	if err != nil {
		return nil, err
	}
	return toTransactionDTOs(transactions), nil
}

func (service TransactionService) UpdateTransaction(ctx context.Context, input ports.UpdateTransactionInput) (ports.TransactionDTO, error) {
	var dto ports.TransactionDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		transactionRepository := service.transactionRepository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)
		categoryRepository := service.categoryRepository.WithTx(tx)

		transaction, err := transactionRepository.FindTransactionByIDForUpdate(txCtx, input.TransactionID, input.UserID)
		if err != nil {
			return err
		}
		if transaction == nil {
			return domain.ErrTransactionNotFound
		}
		if transaction.Status == domain.TransactionStatusRemoved {
			return domain.ErrTransactionRemoved
		}

		if err := service.applyEffects(txCtx, accountRepository, transaction.UserID, transaction.ReversalEffects()); err != nil {
			return err
		}

		switch transaction.Type {
		case domain.TransactionTypeIncome:
			if err := transaction.EditIncome(input.Description, input.Amount, input.OccurredAt, input.AccountID, input.CategoryID, input.Note); err != nil {
				return err
			}
		case domain.TransactionTypeExpense:
			if err := transaction.EditExpense(input.Description, input.Amount, input.OccurredAt, input.AccountID, input.CategoryID, input.Note); err != nil {
				return err
			}
		case domain.TransactionTypeTransfer:
			if err := transaction.EditTransfer(input.Description, input.Amount, input.OccurredAt, input.SourceAccountID, input.DestinationAccountID, input.Note); err != nil {
				return err
			}
		default:
			return domain.ErrTransactionInvalidType
		}

		if err := service.validateReferences(txCtx, accountRepository, categoryRepository, *transaction); err != nil {
			return err
		}
		if err := service.applyEffects(txCtx, accountRepository, transaction.UserID, transaction.BalanceEffects()); err != nil {
			return err
		}

		updated, err := transactionRepository.UpdateTransaction(txCtx, transaction)
		if err != nil {
			return err
		}
		dto = toTransactionDTO(*updated)
		return nil
	})
	if err != nil {
		return ports.TransactionDTO{}, err
	}
	return dto, nil
}

func (service TransactionService) DeleteTransaction(ctx context.Context, input ports.DeleteTransactionInput) error {
	return service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		transactionRepository := service.transactionRepository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)

		transaction, err := transactionRepository.FindTransactionByIDForUpdate(txCtx, input.TransactionID, input.UserID)
		if err != nil {
			return err
		}
		if transaction == nil {
			return domain.ErrTransactionNotFound
		}
		if transaction.Status == domain.TransactionStatusRemoved {
			return nil
		}
		if err := service.applyEffects(txCtx, accountRepository, transaction.UserID, transaction.ReversalEffects()); err != nil {
			return err
		}
		if err := transaction.MarkRemoved(); err != nil {
			return err
		}
		_, err = transactionRepository.UpdateTransaction(txCtx, transaction)
		return err
	})
}

func (service TransactionService) create(ctx context.Context, transaction domain.Transaction) (ports.TransactionDTO, error) {
	var dto ports.TransactionDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		transactionRepository := service.transactionRepository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)
		categoryRepository := service.categoryRepository.WithTx(tx)

		if err := service.validateReferences(txCtx, accountRepository, categoryRepository, transaction); err != nil {
			return err
		}
		if err := service.applyEffects(txCtx, accountRepository, transaction.UserID, transaction.BalanceEffects()); err != nil {
			return err
		}
		created, err := transactionRepository.CreateTransaction(txCtx, &transaction)
		if err != nil {
			return err
		}
		dto = toTransactionDTO(*created)
		return nil
	})
	if err != nil {
		return ports.TransactionDTO{}, err
	}
	return dto, nil
}

func (service TransactionService) validateReferences(ctx context.Context, accountRepository accountports.AccountRepository, categoryRepository categoryports.CategoryRepository, transaction domain.Transaction) error {
	for _, accountID := range accountIDs(transaction) {
		account, err := accountRepository.FindAccountByID(ctx, accountID, transaction.UserID)
		if err != nil {
			return err
		}
		if account == nil || account.Status != accountdomain.AccountStatusActive {
			return domain.ErrTransactionAccountNotFound
		}
	}

	if transaction.CategoryID != nil {
		category, err := categoryRepository.FindCategoryByID(ctx, *transaction.CategoryID, transaction.UserID)
		if err != nil {
			return err
		}
		if category == nil || category.Status != categorydomain.CategoryStatusActive {
			return domain.ErrTransactionCategoryNotFound
		}
		if (transaction.Type == domain.TransactionTypeIncome && category.Type != categorydomain.CategoryTypeIncome) ||
			(transaction.Type == domain.TransactionTypeExpense && category.Type != categorydomain.CategoryTypeExpense) {
			return domain.ErrTransactionCategoryTypeMismatch
		}
	}
	return nil
}

func (service TransactionService) applyEffects(ctx context.Context, accountRepository accountports.AccountRepository, userID userdomain.UserID, effects []domain.BalanceEffect) error {
	for _, effect := range effects {
		account, err := accountRepository.FindAccountByIDForUpdate(ctx, effect.AccountID, userID)
		if err != nil {
			return err
		}
		if account == nil {
			return domain.ErrTransactionAccountNotFound
		}
		if effect.Amount.IsPositive() {
			if err := account.IncreaseBalance(effect.Amount); err != nil {
				return err
			}
		} else {
			if err := account.DecreaseBalance(effect.Amount.Neg()); err != nil {
				return err
			}
		}
		if _, err := accountRepository.UpdateAccount(ctx, account); err != nil {
			return err
		}
	}
	return nil
}

func (service TransactionService) withinTx(ctx context.Context, fn func(context.Context, databaseports.TxHandle) error) error {
	if service.unitOfWork == nil {
		return fn(ctx, databaseports.NewTxHandle(nil))
	}
	return service.unitOfWork.WithinTx(ctx, fn)
}

func accountIDs(transaction domain.Transaction) []accountdomain.AccountID {
	if transaction.Type == domain.TransactionTypeTransfer {
		return []accountdomain.AccountID{*transaction.SourceAccountID, *transaction.DestinationAccountID}
	}
	return []accountdomain.AccountID{*transaction.AccountID}
}

func isValidTransactionType(transactionType domain.TransactionType) bool {
	return transactionType == domain.TransactionTypeIncome ||
		transactionType == domain.TransactionTypeExpense ||
		transactionType == domain.TransactionTypeTransfer
}

func toTransactionDTO(transaction domain.Transaction) ports.TransactionDTO {
	return ports.TransactionDTO{
		ID:                   transaction.ID,
		UserID:               transaction.UserID,
		Type:                 transaction.Type,
		Description:          transaction.Description,
		Amount:               transaction.Amount,
		OccurredAt:           transaction.OccurredAt,
		AccountID:            transaction.AccountID,
		SourceAccountID:      transaction.SourceAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		CategoryID:           transaction.CategoryID,
		Status:               transaction.Status,
		Note:                 transaction.Note,
		RemovedAt:            transaction.RemovedAt,
		CreatedAt:            transaction.CreatedAt,
		UpdatedAt:            transaction.UpdatedAt,
	}
}

func toTransactionDTOs(transactions []domain.Transaction) []ports.TransactionDTO {
	dtos := make([]ports.TransactionDTO, 0, len(transactions))
	for _, transaction := range transactions {
		dtos = append(dtos, toTransactionDTO(transaction))
	}
	return dtos
}
