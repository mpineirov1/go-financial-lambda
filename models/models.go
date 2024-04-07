package models

import "time"

type Transaction struct {
	ID          int
	Date        time.Time
	Transaction float64
	CreatedAt   time.Time
}

type MonthSummary struct {
	TransactionNumber int
	DebitAvg          float64
	CreditAvg         float64
}
