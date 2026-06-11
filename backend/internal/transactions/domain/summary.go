package domain

import financedomain "contai/internal/finance/domain"

type TransactionTotals struct {
	IncomeTotal      financedomain.Money
	ExpenseTotal     financedomain.Money
	TransferInTotal  financedomain.Money
	TransferOutTotal financedomain.Money
}

func CalculateTransactionTotals(transactions []Transaction) TransactionTotals {
	totals := TransactionTotals{
		IncomeTotal:      financedomain.NewMoney(0),
		ExpenseTotal:     financedomain.NewMoney(0),
		TransferInTotal:  financedomain.NewMoney(0),
		TransferOutTotal: financedomain.NewMoney(0),
	}

	for _, transaction := range transactions {
		if transaction.Status != TransactionStatusActive || transaction.RemovedAt != nil {
			continue
		}
		if transaction.Type != TransactionTypeTransfer && transaction.SettlementStatus != SettlementStatusSettled {
			continue
		}

		switch transaction.Type {
		case TransactionTypeIncome:
			totals.IncomeTotal = totals.IncomeTotal.Add(transaction.Amount)
		case TransactionTypeExpense:
			totals.ExpenseTotal = totals.ExpenseTotal.Add(transaction.Amount)
		case TransactionTypeTransfer:
			totals.TransferInTotal = totals.TransferInTotal.Add(transaction.Amount)
			totals.TransferOutTotal = totals.TransferOutTotal.Add(transaction.Amount)
		}
	}

	return totals
}
