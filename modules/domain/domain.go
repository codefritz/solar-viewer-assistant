package domain

import (
	"time"
)

type DayReport struct {
	ReportDate time.Time
	Energy     float64
}
