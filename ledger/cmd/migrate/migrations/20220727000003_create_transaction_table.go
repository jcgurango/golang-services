package migrations

import (
	"context"

	"github.com/uptrace/bun"
	"jcgurango.com/ledger/dbmodel"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model((*dbmodel.Transaction)(nil)).WithForeignKeys().
			ForeignKey(`("credit_account") REFERENCES "accounts" ("id") ON DELETE CASCADE`).
			ForeignKey(`("debit_account") REFERENCES "accounts" ("id") ON DELETE CASCADE`).
			Exec(ctx)

		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().Model((*dbmodel.Transaction)(nil)).Exec(ctx)

		if err != nil {
			return err
		}

		return nil
	})
}
