package http

import (
	"errors"
	"net/http"

	authhttp "contai/internal/auth/adapters/http"
	"contai/internal/dashboard/app/ports"
	"contai/internal/dashboard/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dashboardService ports.DashboardService
}

func NewHandler(dashboardService ports.DashboardService) Handler {
	return Handler{dashboardService: dashboardService}
}

func (handler Handler) GetMonthlyDashboard(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	period, err := parsePeriod(ctx.Query("startAt"), ctx.Query("endAt"))
	if err != nil {
		writeError(ctx, err)
		return
	}

	dashboard, err := handler.dashboardService.GetMonthlyDashboard(ctx.Request.Context(), ports.GetMonthlyDashboardInput{
		UserID: authenticatedUser.UserID,
		Period: period,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toMonthlyDashboardResponse(dashboard))
}

func parsePeriod(startAtValue, endAtValue string) (domain.Period, error) {
	if startAtValue == "" || endAtValue == "" {
		return domain.Period{}, domain.ErrDashboardInvalidPeriod
	}
	startAt, err := parseRFC3339(startAtValue)
	if err != nil {
		return domain.Period{}, domain.ErrDashboardInvalidPeriod
	}
	endAt, err := parseRFC3339(endAtValue)
	if err != nil {
		return domain.Period{}, domain.ErrDashboardInvalidPeriod
	}
	return domain.NewPeriod(startAt, endAt)
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrDashboardUserIDRequired),
		errors.Is(err, domain.ErrDashboardInvalidPeriod):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid dashboard request"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
