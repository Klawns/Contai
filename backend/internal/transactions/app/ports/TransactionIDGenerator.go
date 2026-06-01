package ports

import "contai/internal/transactions/domain"

type TransactionIDGenerator interface {
	NewTransactionID() domain.TransactionID
}
