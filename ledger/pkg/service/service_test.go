package service

import (
	"context"
	"errors"
	"ledger/pkg/service/types"
	stubs "ledger/pkg/stubs"
	"strconv"
	"testing"

	"jcgurango.com/ledger/dbmodel"
)

type TestServiceDatabaseConnection struct {
	LedgerServiceConnection
	users        []dbmodel.User
	accounts     []dbmodel.Account
	transactions []dbmodel.Transaction
	calls        [][]any
}

func (l *TestServiceDatabaseConnection) UserExists(ctx context.Context, username string) (bool, error) {
	l.calls = append(l.calls, []any{"UserExists", username})
	for _, user := range l.users {
		if user.Username == username {
			return true, nil
		}
	}

	return false, nil
}

func (l *TestServiceDatabaseConnection) GetUser(ctx context.Context, username string) (dbmodel.User, error) {
	l.calls = append(l.calls, []any{"GetUser", username})
	for _, user := range l.users {
		if user.Username == username {
			return user, nil
		}
	}

	return dbmodel.User{}, errors.New("Not found")
}

func (l *TestServiceDatabaseConnection) AccountExists(ctx context.Context, userId int64, accountName string) (bool, error) {
	l.calls = append(l.calls, []any{"AccountExists", userId, accountName})
	for _, account := range l.accounts {
		if account.User == userId && account.Name == accountName {
			return true, nil
		}
	}

	return false, nil
}

func (l *TestServiceDatabaseConnection) AccountIdExists(ctx context.Context, userId int64, accountId int64) (bool, error) {
	l.calls = append(l.calls, []any{"AccountIdExists", userId, accountId})
	for _, account := range l.accounts {
		if account.User == userId && account.ID == accountId {
			return true, nil
		}
	}

	return false, nil
}

func (l *TestServiceDatabaseConnection) CreateUser(ctx context.Context, user *dbmodel.User) error {
	l.calls = append(l.calls, []any{"CreateUser", user})
	newUser := *user
	newUser.ID = int64(len(l.users) + 1)
	l.users = append(l.users, newUser)

	return nil
}

func (l *TestServiceDatabaseConnection) CreateAccount(ctx context.Context, account *dbmodel.Account) error {
	l.calls = append(l.calls, []any{"CreateAccount", account})
	newAccount := *account
	newAccount.ID = int64(len(l.accounts) + 1)
	l.accounts = append(l.accounts, newAccount)

	return nil
}

func (l *TestServiceDatabaseConnection) CreateTransaction(ctx context.Context, transaction *dbmodel.Transaction) error {
	l.calls = append(l.calls, []any{"CreateTransaction", transaction})
	newTransaction := *transaction
	newTransaction.ID = int64(len(l.accounts) + 1)
	l.transactions = append(l.transactions, newTransaction)

	return nil
}

// Stub user account list to ensure the service is returning these correctly
var accounts = []dbmodel.Account{
	{
		ID:   1,
		User: 1,
		Name: "Cash",
	},
}

func (l *TestServiceDatabaseConnection) GetUserAccounts(ctx context.Context, userId int64) ([]dbmodel.Account, error) {
	l.calls = append(l.calls, []any{"GetUserAccounts", userId})
	return accounts, nil
}

// Stub user balance list to ensure the service is returning these correctly
var balances = []types.AccountBalance{
	{
		AccountID:   1,
		AccountName: "Cash",
		Balance:     "30,000",
	},
}

func (l *TestServiceDatabaseConnection) GetBalance(ctx context.Context, userId int64) ([]types.AccountBalance, error) {
	l.calls = append(l.calls, []any{"GetBalance", userId})
	return balances, nil
}

// Stub transaction list to ensure the service is returning these correctly
var transactions = []dbmodel.Transaction{
	{
		ID:            1,
		Detail:        "Some Transaction",
		Amount:        "30,000",
		CreditAccount: 1,
		DebitAccount:  2,
	},
}

func (l *TestServiceDatabaseConnection) GetUserTransactions(ctx context.Context, userId int64) ([]dbmodel.Transaction, error) {
	l.calls = append(l.calls, []any{"GetUserTransactions", userId})
	return transactions, nil
}

func TestLedgerServiceAuthentication(t *testing.T) {
	t.Setenv("JWT_TOKEN", "test123")
	service := &basicLedgerService{}

	// Attach a stub logging client
	loggingClient := stubs.TestLoggingClient{}
	ctx := context.WithValue(
		context.Background(),
		"logging-client",
		loggingClient,
	)

	// Attach a service connection for testing
	serviceConnection := TestServiceDatabaseConnection{}
	ctx = context.WithValue(
		ctx,
		"service-connection",
		&serviceConnection,
	)

	t.Run("Register", func(t *testing.T) {
		// Register an account
		service.Register(ctx, "testuser", "testpassword")

		if len(serviceConnection.users) == 0 {
			t.Fatal("New user was not inserted into the database")
		}

		// Attempt to register another account with the same username
		err := service.Register(ctx, "testuser", "testpassword")

		if err == nil || err.Error() != "That user already exists in the database" {
			t.Fatal("Duplicate username not prevented")
		}
	})

	t.Run("Authenticate", func(t *testing.T) {
		// Test authentication
		_, err := service.Authenticate(ctx, "testuser", "testpassword")

		if err != nil {
			t.Fatalf("Authentication failed: %s", err)
		}

		// Wrong password should fail
		_, err = service.Authenticate(ctx, "testuser", "testwrongpassword")

		if err == nil || err.Error() != "Invalid password" {
			t.Fatal("Wrong password not prevented.")
		}

		// Wrong username should fail
		_, err = service.Authenticate(ctx, "testnonexistentuser", "testwrongpassword")

		if err == nil || err.Error() != "User not found" {
			t.Fatal("Wrong username not prevented.")
		}
	})
}

func TestLedgerServiceUserFunctions(t *testing.T) {
	t.Setenv("JWT_TOKEN", "test123")
	service := &basicLedgerService{}

	// Attach a stub logging client
	loggingClient := stubs.TestLoggingClient{}
	ctx := context.WithValue(
		context.Background(),
		"logging-client",
		loggingClient,
	)

	// Attach a service connection for testing
	serviceConnection := TestServiceDatabaseConnection{}
	ctx = context.WithValue(
		ctx,
		"service-connection",
		&serviceConnection,
	)

	// Register an account
	err := service.Register(ctx, "testuser", "testpassword")

	if err != nil {
		t.Fatalf("Authentication failed: %s", err)
	}

	userId := serviceConnection.users[0].ID

	// Attach the logged in user ID
	ctx = context.WithValue(
		ctx,
		"user-id",
		userId,
	)

	t.Run("NewAccount", func(t *testing.T) {
		// Attempt to add a new account
		err = service.NewAccount(ctx, "Cash")

		if err != nil {
			t.Fatalf("NewAccount failed: %s", err)
		}

		if len(serviceConnection.accounts) == 0 {
			t.Fatal("New account not added.")
		}

		if serviceConnection.accounts[0].Name != "Cash" {
			t.Fatalf("Wrong account name %s added.", serviceConnection.accounts[0].Name)
		}

		if serviceConnection.accounts[0].User != userId {
			t.Fatalf("Wrong user %d assigned.", serviceConnection.accounts[0].User)
		}
	})

	t.Run("NewTransaction", func(t *testing.T) {
		// Attempt to add another account
		err = service.NewAccount(ctx, "Notes Payable")

		if err != nil {
			t.Fatalf("NewAccount failed: %s", err)
		}

		err = service.NewTransaction(ctx,
			"Bank Loan",
			strconv.FormatInt(serviceConnection.accounts[0].ID, 10),
			strconv.FormatInt(serviceConnection.accounts[1].ID, 10),
			"30,000",
		)

		if err != nil {
			t.Fatalf("NewTransaction failed: %s", err)
		}

		// Ensure the transaction was added.
		if len(serviceConnection.transactions) == 0 {
			t.Fatal("New transaction not added.")
		}

		if serviceConnection.transactions[0].CreditAccount != serviceConnection.accounts[0].ID || serviceConnection.transactions[0].DebitAccount != serviceConnection.accounts[1].ID {
			t.Fatal("New transaction accounts recorded incorrectly.")
		}
	})

	t.Run("GetAccounts", func(t *testing.T) {
		result, err := service.GetAccounts(ctx)

		if serviceConnection.calls[len(serviceConnection.calls)-1][1] != userId {
			t.Fatal("Wrong user ID passed.")
		}

		if err != nil {
			t.Fatalf("GetAccounts failed: %s", err)
		}

		if result[0] != accounts[0] {
			t.Fatal("Incorrect accounts returned.")
		}
	})

	t.Run("GetBalance", func(t *testing.T) {
		result, err := service.GetBalance(ctx)

		if serviceConnection.calls[len(serviceConnection.calls)-1][1] != userId {
			t.Fatal("Wrong user ID passed.")
		}

		if err != nil {
			t.Fatalf("GetBalance failed: %s", err)
		}

		if (*result)[0] != balances[0] {
			t.Fatal("Incorrect balances returned.")
		}
	})

	t.Run("GetTransactions", func(t *testing.T) {
		result, err := service.GetTransactions(ctx)

		if serviceConnection.calls[len(serviceConnection.calls)-1][1] != userId {
			t.Fatal("Wrong user ID passed.")
		}

		if err != nil {
			t.Fatalf("GetTransactions failed: %s", err)
		}

		if (*result)[0] != transactions[0] {
			t.Fatal("Incorrect balances returned.")
		}
	})
}
