package ports

import "context"

type PDFRenderer interface {
	RenderAccountsReport(ctx context.Context, report AccountsReportDTO) ([]byte, error)
	RenderFinancialReport(ctx context.Context, report FinancialReportDTO) ([]byte, error)
}
