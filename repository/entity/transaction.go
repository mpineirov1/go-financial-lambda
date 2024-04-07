package entity

import "time"

type Transaction struct {
	ID          int       `db:"id"`
	Date        time.Time `db:"date"`
	Transaction float64   `db:"transaction"`
	CreatedAt   time.Time `db:"created_at"`
}
