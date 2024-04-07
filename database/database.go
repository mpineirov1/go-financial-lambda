package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mpineirov1/go-financial-lambda/utils"
	"sync"

	_ "github.com/lib/pq"
)

var (
	getDbOnce sync.Once
	db        *sql.DB
)

func connect(ctx context.Context) (*sql.DB, error) {
	dbHost := utils.GoDotEnvVariable("DB_HOST")
	dbPort := utils.GoDotEnvVariable("DB_PORT")
	dbUser := utils.GoDotEnvVariable("DB_USER")
	password := utils.GoDotEnvVariable("DB_PASSWORD")
	dbName := utils.GoDotEnvVariable("DB_NAME")

	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=require search_path=public`,
		dbHost,
		dbPort,
		dbUser,
		password,
		dbName,
	)

	return sql.Open("postgres", dsn)
}

func Connect(ctx context.Context) (*sql.DB, error) {
	var err error

	getDbOnce.Do(func() {
		db, err = connect(ctx)
	})

	return db, err
}
