package service

import (
	"context"
	"errors"
	"fmt"
	"ledger/pkg/service/service_interface"
	"ledger/pkg/service/types"
	"logging/pkg/grpc/pb"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt"
	"jcgurango.com/ledger/dbmodel"
)

// LedgerService describes the service.
type LedgerService interface {
	Register(ctx context.Context, user string, pass string) (err error)
	Authenticate(ctx context.Context, user string, pass string) (token string, err error)
	NewAccount(ctx context.Context, name string) (err error)
	GetAccounts(ctx context.Context) (accounts []dbmodel.Account, err error)
	NewTransaction(ctx context.Context, detail string, creditAccount string, debitAccount string, amount string) (err error)

	GetBalance(ctx context.Context) (balances *[]types.AccountBalance, err error)
	GetTransactions(ctx context.Context) (transactions *[]dbmodel.Transaction, err error)
}

type basicLedgerService struct{}

func readLoggedInUser(ctx context.Context) int64 {
	var loggedInUserId int64 = ctx.Value("user-id").(int64)
	return loggedInUserId
}

func getLogger(ctx context.Context) pb.LoggingClient {
	var logger pb.LoggingClient = ctx.Value("logging-client").(pb.LoggingClient)
	return logger
}

func log(ctx context.Context, message string) {
	logger := getLogger(ctx)
	logger.Log(ctx, &pb.LogRequest{
		Service:   "ledger",
		Timestamp: time.Now().Unix(),
		Message:   message,
	})
}

func (b *basicLedgerService) Register(ctx context.Context, user string, pass string) (err error) {
	log(ctx, fmt.Sprintf("Register called with parameters [%s] [redacted]", user))

	if user == "" || pass == "" {
		return errors.New("Username and password are required")
	}

	// First check if a user with that username exists in the database
	exists, err := service_interface.UserExists(ctx, user)

	if err != nil {
		log(ctx, "Error encoutered in UserExists: "+err.Error())
		return errors.New("Internal error")
	}

	if exists {
		return errors.New("That user already exists in the database")
	}

	// Then create the user
	hashedPassword, err := argon2id.CreateHash(pass, argon2id.DefaultParams)

	if err != nil {
		log(ctx, "Error encoutered in CreateHash: "+err.Error())
		return errors.New("Internal error")
	}

	databaseUser := &dbmodel.User{
		Username: user,
		Password: hashedPassword,
	}

	err = service_interface.CreateUser(ctx, databaseUser)

	if err != nil {
		log(ctx, "Error encoutered in CreateUser: "+err.Error())
		return errors.New("Internal error")
	}

	return nil
}

func (b *basicLedgerService) Authenticate(ctx context.Context, user string, pass string) (token string, err error) {
	log(ctx, fmt.Sprintf("Authenticate called with parameters [%s] [redacted]", user))

	// First check if a user with that username exists in the database
	exists, err := service_interface.UserExists(ctx, user)

	if err != nil {
		log(ctx, "Error encoutered in UserExists: "+err.Error())
		return "", errors.New("Internal error")
	}

	if !exists {
		return "", errors.New("User not found")
	}

	// Retrieve the user
	dbUser, err := service_interface.GetUser(ctx, user)

	if err != nil {
		log(ctx, "Error encoutered in GetUser: "+err.Error())
		return "", errors.New("Internal error")
	}

	// Verify the password
	if dbUser.ID != 0 {
		match, _, err := argon2id.CheckHash(pass, dbUser.Password)

		if err != nil {
			log(ctx, "Error encoutered in CheckHash: "+err.Error())
			return "", errors.New("Internal error")
		}

		if match {
			// Generate an ID token for the user, note that this is just a stub for "real" authentication.
			// Best security method will always depend on actual implementation.
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": strconv.FormatInt(dbUser.ID, 10),
				"nonce":   time.Now(),
			})
			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

			if err != nil {
				return "", err
			}

			return tokenString, nil
		}

		return "", errors.New("Invalid password")
	}

	return "", errors.New("User not found")
}

func (b *basicLedgerService) NewAccount(ctx context.Context, name string) (err error) {
	userId := readLoggedInUser(ctx)
	log(ctx, fmt.Sprintf("NewAccount called by [%s] with parameters [%s] [redacted]", strconv.FormatInt(userId, 10), name))
	accountExists, err := service_interface.AccountExists(ctx, userId, name)

	if err != nil {
		log(ctx, "Error encoutered in AccountExists: "+err.Error())
		return errors.New("Internal error")
	}

	if accountExists {
		return errors.New("An account already exists with that name.")
	}

	newAccount := dbmodel.Account{
		Name: name,
		User: userId,
	}

	err = service_interface.CreateAccount(ctx, &newAccount)

	if err != nil {
		log(ctx, "Error encoutered in CreateAccount: "+err.Error())
		return errors.New("Internal error")
	}

	return nil
}

func (b *basicLedgerService) GetAccounts(ctx context.Context) (accounts []dbmodel.Account, err error) {
	userId := readLoggedInUser(ctx)
	log(ctx, fmt.Sprintf("GetAccounts called by [%s]", strconv.FormatInt(userId, 10)))
	accounts, err = service_interface.GetUserAccounts(ctx, userId)

	if err != nil {
		log(ctx, "Error encoutered in GetUserAccounts: "+err.Error())
		return nil, errors.New("Internal error")
	}

	return accounts, nil
}

func (b *basicLedgerService) NewTransaction(ctx context.Context, detail string, creditAccount string, debitAccount string, amount string) (err error) {
	userId := readLoggedInUser(ctx)
	log(ctx, fmt.Sprintf("NewAccount called by [%s] with parameters [%s] [%s] [%s] [%s]", strconv.FormatInt(userId, 10), detail, creditAccount, debitAccount, amount))
	creditAccountInt64, err := strconv.ParseInt(creditAccount, 10, 64)

	if err != nil {
		log(ctx, "Error encoutered in ParseInt: "+err.Error())
		return errors.New("Internal error")
	}

	creditAccountExists, err := service_interface.AccountIdExists(ctx, userId, creditAccountInt64)

	if err != nil {
		log(ctx, "Error encoutered in AccountIdExists: "+err.Error())
		return errors.New("Internal error")
	}

	debitAccountInt64, err := strconv.ParseInt(debitAccount, 10, 64)

	if err != nil {
		log(ctx, "Error encoutered in ParseInt: "+err.Error())
		return errors.New("Internal error")
	}

	debitAccountExists, err := service_interface.AccountIdExists(ctx, userId, debitAccountInt64)

	if err != nil {
		log(ctx, "Error encoutered in AccountIdExists: "+err.Error())
		return errors.New("Internal error")
	}

	if !creditAccountExists || !debitAccountExists {
		return errors.New("Credit or debit account does not exist")
	}

	err = service_interface.CreateTransaction(ctx, &dbmodel.Transaction{
		CreditAccount: creditAccountInt64,
		DebitAccount:  debitAccountInt64,
		Detail:        detail,
		Amount:        amount,
	})

	if err != nil {
		log(ctx, "Error encoutered in CreateTransaction: "+err.Error())
		return errors.New("Internal error")
	}

	return err
}

func (b *basicLedgerService) GetBalance(ctx context.Context) (balances *[]types.AccountBalance, err error) {
	userId := readLoggedInUser(ctx)
	log(ctx, fmt.Sprintf("GetBalance called by [%s]", strconv.FormatInt(userId, 10)))
	accountBalances, err := service_interface.GetBalance(ctx, userId)

	if err != nil {
		log(ctx, "Error encoutered in GetBalance: "+err.Error())
		return nil, errors.New("Internal error")
	}

	return &accountBalances, err
}

func (b *basicLedgerService) GetTransactions(ctx context.Context) (transactions *[]dbmodel.Transaction, err error) {
	userId := readLoggedInUser(ctx)
	log(ctx, fmt.Sprintf("GetTransactions called by [%s]", strconv.FormatInt(userId, 10)))
	userTransactions, err := service_interface.GetUserTransactions(ctx, userId)

	if err != nil {
		log(ctx, "Error encoutered in GetUserTransactions: "+err.Error())
		return nil, errors.New("Internal error")
	}

	return &userTransactions, nil
}

// NewBasicLedgerService returns a naive, stateless implementation of LedgerService.
func NewBasicLedgerService() LedgerService {
	return &basicLedgerService{}
}

// New returns a LedgerService with all of the expected middleware wired in.
func New(middleware []Middleware) LedgerService {
	var svc LedgerService = NewBasicLedgerService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
