package domain

import "time"

type Period struct {
	StartAt time.Time
	EndAt   time.Time
}

func NewPeriod(startAt, endAt time.Time) (Period, error) {
	period := Period{StartAt: startAt, EndAt: endAt}
	if err := period.Validate(); err != nil {
		return Period{}, err
	}
	return period, nil
}

func (period Period) Validate() error {
	if period.StartAt.IsZero() || period.EndAt.IsZero() {
		return ErrDashboardInvalidPeriod
	}
	if period.StartAt.After(period.EndAt) {
		return ErrDashboardInvalidPeriod
	}
	return nil
}
