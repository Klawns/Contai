package domain

import "errors"

var (
	ErrDashboardUserIDRequired = errors.New("dashboard user id required")
	ErrDashboardInvalidPeriod  = errors.New("dashboard invalid period")
)
