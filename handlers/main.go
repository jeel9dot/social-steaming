package handlers

import (
	"github.com/jeel9dot/social-steam/nats"
	social_stream "github.com/jeel9dot/social-steam/protobuf/genproto/social-stream"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type SocialStreamGrpcHandler struct {
	logger *zap.Logger
	nc     *nats.NatsClient
	social_stream.UnimplementedSocialSteamServiceServer
}

func NewGrpcSocialStreamHandler(grpc *grpc.Server, logger *zap.Logger, nc *nats.NatsClient) error {
	gRPCHandler := &SocialStreamGrpcHandler{
		logger: logger,
		nc:     nc,
	}

	// register the social stream service
	social_stream.RegisterSocialSteamServiceServer(grpc, gRPCHandler)
	return nil
}
