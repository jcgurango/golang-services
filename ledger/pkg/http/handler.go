package http

import (
	"context"
	"encoding/json"
	"errors"
	endpoint "ledger/pkg/endpoint"
	"net/http"

	http1 "github.com/go-kit/kit/transport/http"
)

// makeRegisterHandler creates the handler logic
func makeRegisterHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/register", http1.NewServer(endpoints.RegisterEndpoint, decodeRegisterRequest, encodeRegisterResponse, options...))
}

// decodeRegisterRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeRegisterResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeRegisterResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeAuthenticateHandler creates the handler logic
func makeAuthenticateHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/authenticate", http1.NewServer(endpoints.AuthenticateEndpoint, decodeAuthenticateRequest, encodeAuthenticateResponse, options...))
}

// decodeAuthenticateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeAuthenticateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.AuthenticateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeAuthenticateResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeAuthenticateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeNewAccountHandler creates the handler logic
func makeNewAccountHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/new-account", http1.NewServer(endpoints.NewAccountEndpoint, decodeNewAccountRequest, encodeNewAccountResponse, options...))
}

// decodeNewAccountRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeNewAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.NewAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeNewAccountResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeNewAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeNewTransactionHandler creates the handler logic
func makeNewTransactionHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/new-transaction", http1.NewServer(endpoints.NewTransactionEndpoint, decodeNewTransactionRequest, encodeNewTransactionResponse, options...))
}

// decodeNewTransactionRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeNewTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.NewTransactionRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeNewTransactionResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeNewTransactionResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetBalanceHandler creates the handler logic
func makeGetBalanceHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/get-balance", http1.NewServer(endpoints.GetBalanceEndpoint, decodeGetBalanceRequest, encodeGetBalanceResponse, options...))
}

// decodeGetBalanceRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetBalanceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.GetBalanceRequest{}
	return req, nil
}

// encodeGetBalanceResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetBalanceResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// makeGetTransactionsHandler creates the handler logic
func makeGetTransactionsHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/get-transactions", http1.NewServer(endpoints.GetTransactionsEndpoint, decodeGetTransactionsRequest, encodeGetTransactionsResponse, options...))
}

// decodeGetTransactionsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.GetTransactionsRequest{}
	return req, nil
}

// encodeGetTransactionsResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetTransactionsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}

// makeGetAccountsHandler creates the handler logic
func makeGetAccountsHandler(m *http.ServeMux, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Handle("/get-accounts", http1.NewServer(endpoints.GetAccountsEndpoint, decodeGetAccountsRequest, encodeGetAccountsResponse, options...))
}

// decodeGetAccountsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeGetAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.GetAccountsRequest{}
	return req, nil
}

// encodeGetAccountsResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeGetAccountsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
