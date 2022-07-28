package grpc

import (
	"context"
	endpoint "logging/pkg/endpoint"
	pb "logging/pkg/grpc/pb"

	grpc "github.com/go-kit/kit/transport/grpc"
)

// makeLogHandler creates the handler logic
func makeLogHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.LogEndpoint, decodeLogRequest, encodeLogResponse, options...)
}

// decodeLogResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Log request.
func decodeLogRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LogRequest)

	return endpoint.LogRequest{
		Service:   req.Service,
		Timestamp: req.Timestamp,
		Message:   req.Message,
	}, nil
}

// encodeLogResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeLogResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.LogResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	return &pb.LogReply{
		Logged: resp.Logged,
	}, nil
}
func (g *grpcServer) Log(ctx context.Context, req *pb.LogRequest) (*pb.LogReply, error) {
	_, rep, err := g.log.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LogReply), nil
}
