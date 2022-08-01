package stubs

import (
	"context"
	logging "logging/pkg/grpc/pb"
	"time"

	log "github.com/go-kit/kit/log"

	"google.golang.org/grpc"
)

type StubLoggingClient struct {
	logging.LoggingClient
	Logger log.Logger
}

func (s StubLoggingClient) Log(ctx context.Context, in *logging.LogRequest, opts ...grpc.CallOption) (*logging.LogReply, error) {
	parsedTimestamp := time.Unix(in.Timestamp, 0)
	s.Logger.Log("service", in.Service, "timestamp", parsedTimestamp.Format(time.RFC3339), "message", in.Message)

	return &logging.LogReply{
		Logged: true,
	}, nil
}
