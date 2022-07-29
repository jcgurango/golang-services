package stubs

import (
	"context"
	"fmt"
	logging "logging/pkg/grpc/pb"
	"time"

	log "github.com/go-kit/kit/log"

	"google.golang.org/grpc"
)

type TestLoggingClient struct {
	logging.LoggingClient
	Logger log.Logger
}

func (s TestLoggingClient) Log(ctx context.Context, in *logging.LogRequest, opts ...grpc.CallOption) (*logging.LogReply, error) {
	parsedTimestamp := time.Unix(in.Timestamp, 0)
	fmt.Printf("[%s] [%s] %s\n", in.Service, parsedTimestamp.Format(time.RFC3339), in.Message)

	return &logging.LogReply{
		Logged: true,
	}, nil
}
