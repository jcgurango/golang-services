package service

import "context"

// LoggingService describes the service.
type LoggingService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	Log(ctx context.Context, s string) (err error)
}
