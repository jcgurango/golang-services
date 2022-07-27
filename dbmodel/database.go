package dbmodel

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var globalDb *bun.DB = nil

func GetDB() *bun.DB {
	return globalDb
}

func SetupDB() {
	sqldb, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONNECTION_STRING"))

	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	globalDb = db
}
