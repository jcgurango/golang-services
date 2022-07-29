package service

import (
	"context"

	"ledger/pkg/service/types"

	"github.com/uptrace/bun"
	"jcgurango.com/ledger/dbmodel"
)

type LedgerServiceDatabaseConnection struct {
	LedgerServiceConnection
}

func (l *LedgerServiceDatabaseConnection) UserExists(ctx context.Context, username string) (bool, error) {
	db := dbmodel.GetDB()
	count, err := db.NewSelect().Model((*dbmodel.User)(nil)).Where("Username = ?", username).Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (l *LedgerServiceDatabaseConnection) CreateUser(ctx context.Context, user *dbmodel.User) error {
	db := dbmodel.GetDB()
	_, err := db.NewInsert().Model(user).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (l *LedgerServiceDatabaseConnection) GetUser(ctx context.Context, username string) (dbmodel.User, error) {
	db := dbmodel.GetDB()
	dbUser := dbmodel.User{}
	err := db.NewSelect().Model(&dbUser).Where("Username = ?", username).Scan(ctx, &dbUser)

	if err != nil {
		return dbUser, err
	}

	return dbUser, nil
}

func (l *LedgerServiceDatabaseConnection) AccountExists(ctx context.Context, userId int64, accountName string) (bool, error) {
	db := dbmodel.GetDB()
	count, err := db.NewSelect().Model((*dbmodel.Account)(nil)).Where("? = ?", bun.Ident("user"), userId).Where("Name = ?", accountName).Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (l *LedgerServiceDatabaseConnection) AccountIdExists(ctx context.Context, userId int64, accountId int64) (bool, error) {
	db := dbmodel.GetDB()
	count, err := db.NewSelect().Model((*dbmodel.Account)(nil)).Where("? = ?", bun.Ident("user"), userId).Where("ID = ?", accountId).Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (l *LedgerServiceDatabaseConnection) CreateAccount(ctx context.Context, account *dbmodel.Account) error {
	db := dbmodel.GetDB()
	_, err := db.NewInsert().Model(account).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (l *LedgerServiceDatabaseConnection) CreateTransaction(ctx context.Context, transaction *dbmodel.Transaction) error {
	db := dbmodel.GetDB()
	_, err := db.NewInsert().Model(transaction).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (l *LedgerServiceDatabaseConnection) GetUserAccounts(ctx context.Context, userId int64) ([]dbmodel.Account, error) {
	db := dbmodel.GetDB()
	var accounts []dbmodel.Account = []dbmodel.Account{}
	err := db.NewSelect().Model(&accounts).Scan(ctx, &accounts)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (l *LedgerServiceDatabaseConnection) GetBalance(ctx context.Context, userId int64) ([]types.AccountBalance, error) {
	db := dbmodel.GetDB()
	var balances []types.AccountBalance = []types.AccountBalance{}
	err := db.
		NewSelect().
		Model((*dbmodel.Account)(nil)).
		ColumnExpr("ID as AccountID").
		ColumnExpr(`"name" as AccountName`).
		ColumnExpr(`COALESCE((SELECT SUM(amount) FROM transactions WHERE credit_account = account.id), '0'::money) - COALESCE((SELECT SUM(amount) FROM transactions WHERE debit_account = account.id), '0'::money) AS Balance`).
		Where("? = ?", bun.Ident("user"), userId).
		Scan(ctx, &balances)

	if err != nil {
		return nil, err
	}

	return balances, err
}

func (l *LedgerServiceDatabaseConnection) GetUserTransactions(ctx context.Context, userId int64) ([]dbmodel.Transaction, error) {
	db := dbmodel.GetDB()
	var transactions []dbmodel.Transaction = []dbmodel.Transaction{}
	err := db.NewSelect().
		Model(&transactions).
		Join(`INNER JOIN "accounts" AS ac ON "transactions".credit_account = "accounts".id`).
		Join(`INNER JOIN "accounts" AS ad ON "transactions".debit_account = "accounts".id`).
		Where(`"ac".user = ? OR "ad".user = ?`, userId).
		Scan(ctx, &transactions)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
