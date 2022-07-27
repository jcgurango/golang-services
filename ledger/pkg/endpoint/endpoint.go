package endpoint

import (
	"context"
	service "ledger/pkg/service"
	"ledger/pkg/service/types"

	endpoint "github.com/go-kit/kit/endpoint"
	"jcgurango.com/ledger/dbmodel"
)

// RegisterRequest collects the request parameters for the Register method.
type RegisterRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// RegisterResponse collects the response parameters for the Register method.
type RegisterResponse struct {
	Err error `json:"err"`
}

// MakeRegisterEndpoint returns an endpoint that invokes Register on the service.
func MakeRegisterEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		err := s.Register(ctx, req.User, req.Pass)
		return RegisterResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r RegisterResponse) Failed() error {
	return r.Err
}

// AuthenticateRequest collects the request parameters for the Authenticate method.
type AuthenticateRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// AuthenticateResponse collects the response parameters for the Authenticate method.
type AuthenticateResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

// MakeAuthenticateEndpoint returns an endpoint that invokes Authenticate on the service.
func MakeAuthenticateEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		token, err := s.Authenticate(ctx, req.User, req.Pass)
		return AuthenticateResponse{
			Err:   err,
			Token: token,
		}, nil
	}
}

// Failed implements Failer.
func (r AuthenticateResponse) Failed() error {
	return r.Err
}

// GetBalanceRequest collects the request parameters for the GetBalance method.
type GetBalanceRequest struct{}

// GetBalanceResponse collects the response parameters for the GetBalance method.
type GetBalanceResponse struct {
	Balances *[]types.AccountBalance `json:"balances"`
	Err      error                   `json:"err"`
}

// MakeGetBalanceEndpoint returns an endpoint that invokes GetBalance on the service.
func MakeGetBalanceEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		balances, err := s.GetBalance(ctx)
		return GetBalanceResponse{
			Balances: balances,
			Err:      err,
		}, nil
	}
}

// Failed implements Failer.
func (r GetBalanceResponse) Failed() error {
	return r.Err
}

// GetTransactionsRequest collects the request parameters for the GetTransactions method.
type GetTransactionsRequest struct{}

// GetTransactionsResponse collects the response parameters for the GetTransactions method.
type GetTransactionsResponse struct {
	Transactions *[]dbmodel.Transaction `json:"transactions"`
	Err          error                  `json:"err"`
}

// MakeGetTransactionsEndpoint returns an endpoint that invokes GetTransactions on the service.
func MakeGetTransactionsEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		transactions, err := s.GetTransactions(ctx)
		return GetTransactionsResponse{
			Err:          err,
			Transactions: transactions,
		}, nil
	}
}

// Failed implements Failer.
func (r GetTransactionsResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Register implements Service. Primarily useful in a client.
func (e Endpoints) Register(ctx context.Context, user string, pass string) (err error) {
	request := RegisterRequest{
		Pass: pass,
		User: user,
	}
	response, err := e.RegisterEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(RegisterResponse).Err
}

// Authenticate implements Service. Primarily useful in a client.
func (e Endpoints) Authenticate(ctx context.Context, user string, pass string) (token string, err error) {
	request := AuthenticateRequest{
		Pass: pass,
		User: user,
	}
	response, err := e.AuthenticateEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(AuthenticateResponse).Token, response.(AuthenticateResponse).Err
}

// GetBalance implements Service. Primarily useful in a client.
func (e Endpoints) GetBalance(ctx context.Context) (balances *[]types.AccountBalance, err error) {
	request := GetBalanceRequest{}
	response, err := e.GetBalanceEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetBalanceResponse).Balances, response.(GetBalanceResponse).Err
}

// GetTransactions implements Service. Primarily useful in a client.
func (e Endpoints) GetTransactions(ctx context.Context) (transactions *[]dbmodel.Transaction, err error) {
	request := GetTransactionsRequest{}
	response, err := e.GetTransactionsEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetTransactionsResponse).Transactions, response.(GetTransactionsResponse).Err
}

// NewTransactionRequest collects the request parameters for the NewTransaction method.
type NewTransactionRequest struct {
	Detail        string `json:"detail"`
	CreditAccount string `json:"credit_account"`
	DebitAccount  string `json:"debit_account"`
	Amount        string `json:"amount"`
}

// NewTransactionResponse collects the response parameters for the NewTransaction method.
type NewTransactionResponse struct {
	Err error `json:"err"`
}

// MakeNewTransactionEndpoint returns an endpoint that invokes NewTransaction on the service.
func MakeNewTransactionEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NewTransactionRequest)
		err := s.NewTransaction(ctx, req.Detail, req.CreditAccount, req.DebitAccount, req.Amount)
		return NewTransactionResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r NewTransactionResponse) Failed() error {
	return r.Err
}

// NewTransaction implements Service. Primarily useful in a client.
func (e Endpoints) NewTransaction(ctx context.Context, detail string, creditAccount string, debitAccount string, amount string) (err error) {
	request := NewTransactionRequest{
		Amount:        amount,
		CreditAccount: creditAccount,
		DebitAccount:  debitAccount,
		Detail:        detail,
	}
	response, err := e.NewTransactionEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(NewTransactionResponse).Err
}

// GetAccountsRequest collects the request parameters for the GetAccounts method.
type GetAccountsRequest struct{}

// GetAccountsResponse collects the response parameters for the GetAccounts method.
type GetAccountsResponse struct {
	Accounts []dbmodel.Account `json:"accounts"`
	Err      error             `json:"err"`
}

// MakeGetAccountsEndpoint returns an endpoint that invokes GetAccounts on the service.
func MakeGetAccountsEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		accounts, err := s.GetAccounts(ctx)
		return GetAccountsResponse{
			Accounts: accounts,
			Err:      err,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAccountsResponse) Failed() error {
	return r.Err
}

// GetAccounts implements Service. Primarily useful in a client.
func (e Endpoints) GetAccounts(ctx context.Context) (accounts []dbmodel.Account, err error) {
	request := GetAccountsRequest{}
	response, err := e.GetAccountsEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetAccountsResponse).Accounts, response.(GetAccountsResponse).Err
}

// NewAccountRequest collects the request parameters for the NewAccount method.
type NewAccountRequest struct {
	Name string `json:"name"`
}

// NewAccountResponse collects the response parameters for the NewAccount method.
type NewAccountResponse struct {
	Err error `json:"err"`
}

// MakeNewAccountEndpoint returns an endpoint that invokes NewAccount on the service.
func MakeNewAccountEndpoint(s service.LedgerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NewAccountRequest)
		err := s.NewAccount(ctx, req.Name)
		return NewAccountResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r NewAccountResponse) Failed() error {
	return r.Err
}

// NewAccount implements Service. Primarily useful in a client.
func (e Endpoints) NewAccount(ctx context.Context, name string) (err error) {
	request := NewAccountRequest{Name: name}
	response, err := e.NewAccountEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(NewAccountResponse).Err
}
