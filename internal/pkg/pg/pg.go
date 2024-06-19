package pg

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewConnection(dsn string) (*bun.DB, error) {
	//dsn2 := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", "postgres", "1234", "5432", "tgstat")
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
