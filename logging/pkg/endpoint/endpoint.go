package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	service "logging/pkg/service"
)

// LogRequest collects the request parameters for the Log method.
type LogRequest struct {
	Service   string `json:"service"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// LogResponse collects the response parameters for the Log method.
type LogResponse struct {
	Logged bool  `json:"logged"`
	Err    error `json:"err"`
}

// MakeLogEndpoint returns an endpoint that invokes Log on the service.
func MakeLogEndpoint(s service.LoggingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LogRequest)
		logged, err := s.Log(ctx, req.Service, req.Timestamp, req.Message)
		return LogResponse{
			Err:    err,
			Logged: logged,
		}, nil
	}
}

// Failed implements Failer.
func (r LogResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Log implements Service. Primarily useful in a client.
func (e Endpoints) Log(ctx context.Context, service string, timestamp int64, message string) (logged bool, err error) {
	request := LogRequest{
		Message:   message,
		Service:   service,
		Timestamp: timestamp,
	}
	response, err := e.LogEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(LogResponse).Logged, response.(LogResponse).Err
}
