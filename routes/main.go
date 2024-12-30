package routes

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jeel9dot/trading-pub-sub/config"
	controller "github.com/jeel9dot/trading-pub-sub/controllers/api/v1"
	"github.com/jeel9dot/trading-pub-sub/nats"

	"go.uber.org/zap"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, logger *zap.Logger, config config.AppConfig, nc *nats.NatsClient) error {
	mu.Lock()

	// app.Use(swagger.New(swagger.Config{
	// 	BasePath: "/api/v1/",
	// 	FilePath: "./assets/swagger.json",
	// 	Path:     "docs",
	// 	Title:    "Swagger API Docs",
	// }))

	router := app.Group("/api")
	v1 := router.Group("/v1")

	// For WS group
	ws := app.Group("/ws")
	ws_v1 := ws.Group("/v1")

	err := setupPubSubController(v1, ws_v1, logger, nc)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func setupPubSubController(v1 fiber.Router, ws_v1 fiber.Router, logger *zap.Logger, nc *nats.NatsClient) error {
	pubSubController := controller.NewPubSubController(logger, nc)
	v1.Post("/publish", pubSubController.Publish)

	ws_v1.Get("/subscribe/:subject", pubSubController.Subscribe())
	return nil
}
