package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const authBodyLimitBytes int64 = 1 << 20

func registerRoutes(router *gin.Engine, dependencies dependencies) {
	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authLimiter := newRateLimiter(20, time.Minute)

	router.POST("/api/users", limitBody(authBodyLimitBytes), authLimiter.Middleware(), dependencies.authHandler.CreateUser)
	router.POST("/api/auth/login", limitBody(authBodyLimitBytes), authLimiter.Middleware(), dependencies.authHandler.Login)
	router.POST("/api/auth/logout", dependencies.authHandler.Logout)
	router.GET("/api/auth/me", dependencies.authHandler.AuthMiddleware(), dependencies.authHandler.Me)

	authenticated := router.Group("/api", dependencies.authHandler.AuthMiddleware())
	authenticated.GET("/categories", dependencies.categoryHandler.ListCategories)
	authenticated.POST("/categories", limitBody(authBodyLimitBytes), dependencies.categoryHandler.CreateCategory)
	authenticated.PATCH("/categories/:categoryID", limitBody(authBodyLimitBytes), dependencies.categoryHandler.UpdateCategory)
	authenticated.DELETE("/categories/:categoryID", dependencies.categoryHandler.DeleteCategory)
	authenticated.GET("/accounts", dependencies.accountHandler.ListAccounts)
	authenticated.POST("/accounts", limitBody(authBodyLimitBytes), dependencies.accountHandler.CreateAccount)
	authenticated.GET("/accounts/total-balance", dependencies.accountHandler.GetTotalBalance)
	authenticated.PATCH("/accounts/:accountID", limitBody(authBodyLimitBytes), dependencies.accountHandler.UpdateAccount)
	authenticated.DELETE("/accounts/:accountID", dependencies.accountHandler.DeleteAccount)
	authenticated.GET("/transactions", dependencies.transactionHandler.ListTransactions)
	authenticated.POST("/transactions/income", limitBody(authBodyLimitBytes), dependencies.transactionHandler.CreateIncome)
	authenticated.POST("/transactions/expense", limitBody(authBodyLimitBytes), dependencies.transactionHandler.CreateExpense)
	authenticated.POST("/transactions/transfer", limitBody(authBodyLimitBytes), dependencies.transactionHandler.CreateTransfer)
	authenticated.PATCH("/transactions/:transactionID", limitBody(authBodyLimitBytes), dependencies.transactionHandler.UpdateTransaction)
	authenticated.DELETE("/transactions/:transactionID", dependencies.transactionHandler.DeleteTransaction)
	authenticated.GET("/credit-cards", dependencies.creditCardHandler.ListCreditCards)
	authenticated.POST("/credit-cards", limitBody(authBodyLimitBytes), dependencies.creditCardHandler.CreateCreditCard)
	authenticated.PATCH("/credit-cards/:cardID", limitBody(authBodyLimitBytes), dependencies.creditCardHandler.UpdateCreditCard)
	authenticated.PATCH("/credit-cards/:cardID/inactivate", dependencies.creditCardHandler.InactivateCreditCard)
	authenticated.GET("/credit-cards/:cardID/purchases", dependencies.creditCardHandler.ListPurchases)
	authenticated.POST("/credit-cards/:cardID/purchases", limitBody(authBodyLimitBytes), dependencies.creditCardHandler.CreatePurchase)
	authenticated.PATCH("/credit-card-purchases/:purchaseID/cancel", dependencies.creditCardHandler.CancelPurchase)
	authenticated.GET("/credit-cards/:cardID/invoices", dependencies.creditCardHandler.ListInvoices)
	authenticated.GET("/credit-card-invoices/:invoiceID", dependencies.creditCardHandler.GetInvoice)
	authenticated.PATCH("/credit-card-invoices/:invoiceID/close", dependencies.creditCardHandler.CloseInvoice)
	authenticated.PATCH("/credit-card-invoices/:invoiceID/pay", limitBody(authBodyLimitBytes), dependencies.creditCardHandler.PayInvoice)
	authenticated.GET("/dashboard/monthly", dependencies.dashboardHandler.GetMonthlyDashboard)
	authenticated.GET("/dashboard/monthly-series", dependencies.dashboardHandler.GetMonthlySeries)
	authenticated.GET("/reports/financial", dependencies.reportHandler.GetFinancialReport)
	authenticated.GET("/reports/financial/pdf", dependencies.reportHandler.DownloadFinancialPDF)
}
