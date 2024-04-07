package repository

import (
	"context"
	"database/sql"
	"github.com/mpineirov1/go-financial-lambda/repository/entity"
)

type Transaction interface {
	Insert(ctx context.Context, transaction *entity.Transaction) (err error)
}

type transaction struct {
	Conn *sql.DB
}

func NewTransaction(conn *sql.DB) transaction {
	return transaction{
		Conn: conn,
	}
}

func (m transaction) Insert(ctx context.Context, transaction *entity.Transaction) (err error) {
	query := `
		INSERT INTO transactions (date, transaction, created_at)
		VALUES ($1, $2, $3)
	`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		transaction.Date,
		transaction.Transaction,
		transaction.CreatedAt,
	).Err()

	return err
}
