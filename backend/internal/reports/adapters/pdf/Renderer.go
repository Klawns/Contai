package pdf

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	accountdomain "contai/internal/account/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"

	"github.com/phpdave11/gofpdf"
)

const (
	pageWidth    = 210.0
	pageHeight   = 297.0
	marginLeft   = 12.0
	marginTop    = 12.0
	marginRight  = 12.0
	marginBottom = 14.0
)

var _ reportports.PDFRenderer = Renderer{}

type Renderer struct{}

type pdfDocument struct {
	pdf       *gofpdf.Fpdf
	translate func(string) string
	title     string
	subtitle  string
}

type tableColumn struct {
	Label string
	Width float64
	Align string
}

func NewRenderer() (Renderer, error) {
	return Renderer{}, nil
}

func (renderer Renderer) RenderAccountsReport(ctx context.Context, report reportports.AccountsReportDTO) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	doc := newPDFDocument("Relatorio de contas", "")
	doc.addPage()
	doc.addGeneratedAt(report.GeneratedAt)
	doc.addMoneySummary([]summaryItem{
		{Label: "Saldo total", Value: formatMoney(report.TotalBalance)},
		{Label: "Total no dashboard", Value: formatMoney(report.DashboardTotal)},
		{Label: "Contas", Value: fmt.Sprintf("%d", len(report.Accounts))},
	})

	columns := []tableColumn{
		{Label: "Nome", Width: 44, Align: "L"},
		{Label: "Tipo", Width: 32, Align: "L"},
		{Label: "Status", Width: 20, Align: "L"},
		{Label: "Inicial", Width: 28, Align: "R"},
		{Label: "Atual", Width: 28, Align: "R"},
		{Label: "Dashboard", Width: 30, Align: "L"},
	}
	doc.addTableHeader(columns)
	if len(report.Accounts) == 0 {
		doc.addEmptyRow(columns, "Nenhuma conta encontrada.")
	} else {
		for _, account := range report.Accounts {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			doc.addTableRow(columns, []string{
				account.Name,
				accountType(account.Type),
				accountStatus(account.Status),
				formatMoney(account.InitialBalance),
				formatMoney(account.CurrentBalance),
				yesNo(account.IncludeInDashboardTotal),
			})
		}
	}

	return doc.output()
}

func (renderer Renderer) RenderFinancialReport(ctx context.Context, report reportports.FinancialReportDTO) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	doc := newPDFDocument(report.Title, report.Subtitle)
	doc.addPage()
	doc.addGeneratedAt(report.GeneratedAt)
	doc.addPeriod(report)
	doc.addMoneySummary([]summaryItem{
		{Label: "Receitas", Value: formatMoney(report.IncomeTotal)},
		{Label: "Despesas", Value: formatMoney(report.ExpenseTotal)},
		{Label: "Transf. entrada", Value: formatMoney(report.TransferInTotal)},
		{Label: "Transf. saida", Value: formatMoney(report.TransferOutTotal)},
		{Label: "Resultado", Value: formatMoney(report.NetTotal)},
	})

	columns := []tableColumn{
		{Label: "Data", Width: 24, Align: "L"},
		{Label: "Tipo", Width: 27, Align: "L"},
		{Label: "Descricao", Width: 61, Align: "L"},
		{Label: "Conta", Width: 42, Align: "L"},
		{Label: "Valor", Width: 28, Align: "R"},
	}
	doc.addTableHeader(columns)
	if len(report.Transactions) == 0 {
		doc.addEmptyRow(columns, "Nenhuma transacao encontrada.")
	} else {
		for _, transaction := range report.Transactions {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			doc.addTableRow(columns, []string{
				formatDate(transaction.OccurredAt),
				transactionType(transaction.Type),
				transaction.Description,
				transactionAccountName(transaction),
				formatMoney(transaction.Amount),
			})
		}
	}

	return doc.output()
}

type summaryItem struct {
	Label string
	Value string
}

func newPDFDocument(title, subtitle string) *pdfDocument {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(false, marginBottom)
	pdf.SetCompression(true)
	pdf.SetTitle(title, false)
	pdf.SetCreator("Contai", false)

	doc := &pdfDocument{
		pdf:       pdf,
		translate: pdf.UnicodeTranslatorFromDescriptor(""),
		title:     title,
		subtitle:  subtitle,
	}
	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(120, 120, 120)
		pdf.CellFormat(0, 6, doc.text(fmt.Sprintf("Pagina %d", pdf.PageNo())), "", 0, "R", false, 0, "")
	})
	return doc
}

func (doc *pdfDocument) output() ([]byte, error) {
	var output bytes.Buffer
	if err := doc.pdf.Output(&output); err != nil {
		return nil, fmt.Errorf("write pdf: %w", err)
	}
	if output.Len() == 0 {
		return nil, fmt.Errorf("empty pdf output")
	}
	return output.Bytes(), nil
}

func (doc *pdfDocument) addPage() {
	doc.pdf.AddPage()
	doc.pdf.SetDrawColor(220, 225, 232)
	doc.pdf.SetFillColor(248, 250, 252)
	doc.pdf.Rect(0, 0, pageWidth, 28, "F")
	doc.pdf.SetTextColor(24, 31, 42)
	doc.pdf.SetFont("Arial", "B", 17)
	doc.pdf.SetXY(marginLeft, 9)
	doc.pdf.CellFormat(0, 8, doc.text(doc.title), "", 1, "L", false, 0, "")
	if strings.TrimSpace(doc.subtitle) != "" {
		doc.pdf.SetFont("Arial", "", 10)
		doc.pdf.SetTextColor(82, 92, 108)
		doc.pdf.SetX(marginLeft)
		doc.pdf.CellFormat(0, 6, doc.text(doc.subtitle), "", 1, "L", false, 0, "")
	}
	doc.pdf.SetY(34)
}

func (doc *pdfDocument) addGeneratedAt(generatedAt time.Time) {
	doc.pdf.SetFont("Arial", "", 9)
	doc.pdf.SetTextColor(82, 92, 108)
	doc.pdf.CellFormat(0, 5, doc.text("Gerado em "+formatDateTime(generatedAt)), "", 1, "L", false, 0, "")
	doc.pdf.Ln(3)
}

func (doc *pdfDocument) addPeriod(report reportports.FinancialReportDTO) {
	parts := []string{fmt.Sprintf("Periodo: %s a %s", formatDate(report.StartAt), formatDate(report.EndAt))}
	if strings.TrimSpace(report.AccountName) != "" {
		parts = append(parts, "Conta: "+report.AccountName)
	}
	doc.pdf.SetFont("Arial", "", 9)
	doc.pdf.SetTextColor(82, 92, 108)
	doc.pdf.CellFormat(0, 5, doc.text(strings.Join(parts, " | ")), "", 1, "L", false, 0, "")
	doc.pdf.Ln(3)
}

func (doc *pdfDocument) addMoneySummary(items []summaryItem) {
	doc.ensureSpace(22)
	itemWidth := (pageWidth - marginLeft - marginRight) / float64(len(items))
	startY := doc.pdf.GetY()
	for index, item := range items {
		x := marginLeft + float64(index)*itemWidth
		doc.pdf.SetXY(x, startY)
		doc.pdf.SetFillColor(245, 247, 250)
		doc.pdf.SetDrawColor(220, 225, 232)
		doc.pdf.Rect(x, startY, itemWidth-2, 18, "FD")
		doc.pdf.SetTextColor(82, 92, 108)
		doc.pdf.SetFont("Arial", "", 8)
		doc.pdf.SetXY(x+2, startY+2)
		doc.pdf.CellFormat(itemWidth-6, 5, doc.text(item.Label), "", 1, "L", false, 0, "")
		doc.pdf.SetTextColor(24, 31, 42)
		doc.pdf.SetFont("Arial", "B", 10)
		doc.pdf.SetXY(x+2, startY+9)
		doc.pdf.CellFormat(itemWidth-6, 6, doc.text(item.Value), "", 1, "L", false, 0, "")
	}
	doc.pdf.SetY(startY + 23)
}

func (doc *pdfDocument) addTableHeader(columns []tableColumn) {
	doc.ensureSpace(14)
	doc.pdf.SetFillColor(39, 50, 68)
	doc.pdf.SetTextColor(255, 255, 255)
	doc.pdf.SetDrawColor(39, 50, 68)
	doc.pdf.SetFont("Arial", "B", 8)
	for _, column := range columns {
		doc.pdf.CellFormat(column.Width, 8, doc.text(column.Label), "1", 0, "L", true, 0, "")
	}
	doc.pdf.Ln(-1)
}

func (doc *pdfDocument) addTableRow(columns []tableColumn, values []string) {
	if doc.pdf.GetY()+8 > pageHeight-marginBottom {
		doc.addPage()
		doc.addTableHeader(columns)
	}
	doc.pdf.SetFont("Arial", "", 8)
	doc.pdf.SetTextColor(35, 43, 55)
	doc.pdf.SetDrawColor(226, 232, 240)
	for index, column := range columns {
		value := ""
		if index < len(values) {
			value = values[index]
		}
		doc.pdf.CellFormat(column.Width, 7, doc.text(fitText(doc.pdf, value, column.Width-3)), "1", 0, column.Align, false, 0, "")
	}
	doc.pdf.Ln(-1)
}

func (doc *pdfDocument) addEmptyRow(columns []tableColumn, message string) {
	doc.ensureSpace(8)
	width := 0.0
	for _, column := range columns {
		width += column.Width
	}
	doc.pdf.SetFont("Arial", "", 8)
	doc.pdf.SetTextColor(82, 92, 108)
	doc.pdf.SetDrawColor(226, 232, 240)
	doc.pdf.CellFormat(width, 8, doc.text(message), "1", 1, "L", false, 0, "")
}

func (doc *pdfDocument) ensureSpace(required float64) {
	if doc.pdf.GetY()+required <= pageHeight-marginBottom {
		return
	}
	doc.addPage()
}

func (doc *pdfDocument) text(value string) string {
	return doc.translate(value)
}

func fitText(pdf *gofpdf.Fpdf, value string, width float64) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\n", " "))
	if pdf.GetStringWidth(value) <= width {
		return value
	}
	const suffix = "..."
	runes := []rune(value)
	for len(runes) > 0 && pdf.GetStringWidth(string(runes)+suffix) > width {
		runes = runes[:len(runes)-1]
	}
	value = strings.TrimSpace(string(runes))
	if value == "" {
		return suffix
	}
	return value + suffix
}

func formatDateTime(value time.Time) string {
	if value.IsZero() {
		return "-"
	}
	return value.Format("02/01/2006 15:04")
}

func formatDate(value time.Time) string {
	if value.IsZero() {
		return "-"
	}
	return value.Format("02/01/2006")
}

func formatMoney(value interface{ Cents() int64 }) string {
	cents := value.Cents()
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%sR$ %d,%02d", sign, cents/100, cents%100)
}

func accountType(value accountdomain.AccountType) string {
	switch value {
	case accountdomain.AccountTypeChecking:
		return "Conta corrente"
	case accountdomain.AccountTypeSavings:
		return "Poupanca"
	case accountdomain.AccountTypeDigital:
		return "Conta digital"
	case accountdomain.AccountTypeCash:
		return "Dinheiro"
	case accountdomain.AccountTypeSalary:
		return "Conta salario"
	case accountdomain.AccountTypeInvestment:
		return "Investimento"
	default:
		return "Outra"
	}
}

func accountStatus(value accountdomain.AccountStatus) string {
	if value == accountdomain.AccountStatusActive {
		return "Ativa"
	}
	return "Inativa"
}

func yesNo(value bool) string {
	if value {
		return "Sim"
	}
	return "Nao"
}

func transactionType(value transactiondomain.TransactionType) string {
	switch value {
	case transactiondomain.TransactionTypeIncome:
		return "Receita"
	case transactiondomain.TransactionTypeExpense:
		return "Despesa"
	case transactiondomain.TransactionTypeTransfer:
		return "Transferencia"
	default:
		return "Outro"
	}
}

func transactionAccountName(transaction reportports.ReportTransactionRow) string {
	if strings.TrimSpace(transaction.AccountName) != "" {
		return transaction.AccountName
	}
	if transaction.Type == transactiondomain.TransactionTypeTransfer {
		source := accountIDText(transaction.SourceAccountID)
		destination := accountIDText(transaction.DestinationAccountID)
		if source != "" && destination != "" {
			return source + " -> " + destination
		}
	}
	return accountIDText(transaction.AccountID)
}

func accountIDText(value *accountdomain.AccountID) string {
	if value == nil {
		return ""
	}
	return string(*value)
}
