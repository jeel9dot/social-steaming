package cli

import (
	"github.com/jeel9dot/trading-pub-sub/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	apiCmd := GetAPICommandDef(cfg, logger)
	rootCmd := &cobra.Command{Use: "social-trading-pub-sub"}
	rootCmd.AddCommand(&apiCmd)
	return rootCmd.Execute()
}