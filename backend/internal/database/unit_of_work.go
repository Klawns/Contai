package database

import (
	"context"

	"contai/internal/database/ports"

	"gorm.io/gorm"
)

var _ ports.UnitOfWork = UnitOfWork{}

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return UnitOfWork{db: db}
}

func (unit UnitOfWork) WithinTx(ctx context.Context, fn func(ctx context.Context, tx ports.TxHandle) error) error {
	return unit.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, ports.NewTxHandle(tx))
	})
}
