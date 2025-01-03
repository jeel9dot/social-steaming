package cli

import (
	"github.com/jeel9dot/social-steam/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	apiCmd := GetAPICommandDef(cfg, logger)
	grpcCmd := GetGrpcCommandDef(cfg, logger)
	rootCmd := &cobra.Command{Use: "social-stream"}
	rootCmd.AddCommand(&apiCmd, &grpcCmd)
	return rootCmd.Execute()
}
