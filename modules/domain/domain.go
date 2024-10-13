package domain

import (
	"time"
)

type DayReport struct {
	ReportDate time.Time
	Energy     float64
}

type Weather struct {
	ReportDate time.Time
	Cloudiness int
}
