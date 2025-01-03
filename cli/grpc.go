package cli

import (
	"github.com/jeel9dot/social-steam/config"
	"github.com/jeel9dot/social-steam/nats"
	grpc_server "github.com/jeel9dot/social-steam/pkg"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetAPICommandDef runs app
func GetGrpcCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "grpc",
		Short: "To start grpc score service server",
		Long:  `To start grpc score service server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Nats Connection init
			nc, err := nats.NewNatsClient(cfg.Nats.Url)
			if err != nil {
				logger.Panic("error while connecting to nats", zap.Error(err))
			}
			defer func() {
				logger.Info("closing nats connection")
				nc.Close()
			}()

			// Create grpc server
			rpcServer := grpc_server.NewGRPCServer(cfg.Grpc.Port, logger, nc)
			if err := rpcServer.Run(); err != nil {
				logger.Error("error while running grpc server", zap.Error(err))
				return err
			}
			return nil
		},
	}

	return apiCommand
}
