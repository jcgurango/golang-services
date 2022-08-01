package main

import (
	"context"
	"ledger/cmd/migrate/migrations"
	"os"

	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"jcgurango.com/ledger/dbmodel"
)

func main() {
	dbmodel.SetupDB()
	db := dbmodel.GetDB()

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv(""),
	))

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	ctx := context.Background()

	if len(os.Args) > 1 && os.Args[1] == "rollback" {
		_, err := migrator.Rollback(ctx)

		if err != nil {
			panic(err)
		}
	} else {
		migrator.Init(ctx)
		_, err := migrator.Migrate(ctx)

		if err != nil {
			panic(err)
		}
	}
}
