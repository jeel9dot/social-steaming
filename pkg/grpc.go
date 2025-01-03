package grpc_server

import (
	"net"

	"github.com/jeel9dot/trading-pub-sub/handlers"
	"github.com/jeel9dot/trading-pub-sub/nats"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr   string
	logger *zap.Logger
	nc     *nats.NatsClient
}

// NewGRPCServer creates new grpc server
func NewGRPCServer(addr string, logger *zap.Logger, nc *nats.NatsClient) *gRPCServer {
	return &gRPCServer{addr: addr, logger: logger, nc: nc}
}

// Run starts grpc server
func (s *gRPCServer) Run() error {
	// Start net server for grpc
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	// register our grpc services
	err = handlers.NewGrpcSocialStreamHandler(grpcServer, s.logger, s.nc)
	if err != nil {
		return err
	}

	s.logger.Info("Starting gRPC server", zap.String("addr", s.addr))
	return grpcServer.Serve(lis)
}
