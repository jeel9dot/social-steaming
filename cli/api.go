package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/jeel9dot/social-steam/config"
	"github.com/jeel9dot/social-steam/nats"
	"github.com/jeel9dot/social-steam/routes"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetAPICommandDef runs app
func GetAPICommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Create fiber app
			app := fiber.New(fiber.Config{})

			// Nats Connection init
			nc, err := nats.NewNatsClient(cfg.Nats.Url)
			if err != nil {
				logger.Panic("error while connecting to nats", zap.Error(err))
			}

			// setup routes
			err = routes.Setup(app, logger, cfg, nc)
			if err != nil {
				return err
			}

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				if err := app.Listen(cfg.Port); err != nil {
					logger.Panic(err.Error())
				}
			}()

			<-interrupt
			logger.Info("gracefully shutting down...")
			nc.Close()
			if err := app.Shutdown(); err != nil {
				logger.Panic("error while shutdown server", zap.Error(err))
			}

			logger.Info("server stopped to receive new requests or connection.")
			return nil
		},
	}

	return apiCommand
}
