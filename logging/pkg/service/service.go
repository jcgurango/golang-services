package service

import (
	"context"
	"fmt"
	"time"
)

// LoggingService describes the service.
type LoggingService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	Log(ctx context.Context, service string, timestamp int64, message string) (logged bool, err error)
}

type basicLoggingService struct{}

func (b *basicLoggingService) Log(ctx context.Context, service string, timestamp int64, message string) (logged bool, err error) {
	parsedTimestamp := time.Unix(timestamp, 0)
	fmt.Printf("[%s] [%s] %s\n", service, parsedTimestamp.Format("YYYY-MM-DD hh:mm:ss"), message)
	return true, nil
}

// NewBasicLoggingService returns a naive, stateless implementation of LoggingService.
func NewBasicLoggingService() LoggingService {
	return &basicLoggingService{}
}

// New returns a LoggingService with all of the expected middleware wired in.
func New(middleware []Middleware) LoggingService {
	var svc LoggingService = NewBasicLoggingService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
