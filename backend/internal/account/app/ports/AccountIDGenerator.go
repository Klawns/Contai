package ports

import "contai/internal/account/domain"

type AccountIDGenerator interface {
	NewAccountID() domain.AccountID
}
