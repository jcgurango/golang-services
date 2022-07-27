package migrations

import (
	"context"

	"github.com/uptrace/bun"
	"jcgurango.com/ledger/dbmodel"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model((*dbmodel.Account)(nil)).
			WithForeignKeys().
			ForeignKey(`("user") REFERENCES "users" ("id") ON DELETE CASCADE`).
			Exec(ctx)

		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().Model((*dbmodel.Account)(nil)).Exec(ctx)

		if err != nil {
			return err
		}

		return nil
	})
}
