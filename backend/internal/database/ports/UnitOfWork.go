package ports

import "context"

type TxHandle struct {
	value any
}

func NewTxHandle(value any) TxHandle {
	return TxHandle{value: value}
}

func (handle TxHandle) Value() any {
	return handle.value
}

type UnitOfWork interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, tx TxHandle) error) error
}
