package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/jeel9dot/trading-pub-sub/nats"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type PublishRequest struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// WebSocket Upgrader
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

type PubSubController struct {
	logger *zap.Logger
	nc     *nats.NatsClient
}

func NewPubSubController(logger *zap.Logger, nc *nats.NatsClient) *PubSubController {
	return &PubSubController{logger: logger, nc: nc}
}

func (ctrl *PubSubController) Publish(c *fiber.Ctx) error {
	req := new(PublishRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	ctrl.logger.Info("publishing message", zap.String("subject", req.Subject), zap.String("message", req.Message))
	err := ctrl.nc.Publish(req.Subject, req.Message)
	if err != nil {
		ctrl.logger.Error("error while publishing message", zap.Error(err))
		return err
	}
	ctrl.logger.Info("published message", zap.String("subject", req.Subject), zap.String("message", req.Message))

	return c.JSON(fiber.Map{
		"message": "published message successfully",
	})
}

// Subscribe to a subject
func (ctrl *PubSubController) Subscribe() fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		subject := conn.Params("subject")
		ctrl.logger.Info("WebSocket client connected", zap.String("subject", subject))
		defer func() {
			// Close WebSocket connection on exit
			err := conn.Close()
			if err != nil {
				ctrl.logger.Error("Error while closing WebSocket connection", zap.Error(err))
			}
			ctrl.logger.Info("WebSocket client disconnected", zap.String("subject", subject))
		}()

		// Create a channel to handle WebSocket disconnection
		closeChan := make(chan struct{})

		// Subscribe to the subject
		subscription, err := ctrl.nc.Subscribe(subject, func(msg *nc.Msg) {
			select {
			case <-closeChan:
				// Stop sending messages after disconnect
				return
			default:

				// Send message to WebSocket client
				ctrl.logger.Info("Received message", zap.String("subject", subject), zap.String("message", string(msg.Data)))
				err := conn.WriteMessage(websocket.TextMessage, msg.Data)
				if err != nil {
					ctrl.logger.Error("Error while sending message to WebSocket client", zap.Error(err))
					// Trigger disconnection cleanup
					close(closeChan)
				}
			}
		})
		if err != nil {
			ctrl.logger.Error("Error while subscribing to subject", zap.Error(err))
			return
		}
		defer subscription.Unsubscribe() // Unsubscribe when client disconnects

		// Keep the connection open until the client disconnects
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Trigger cleanup on WebSocket disconnection
				close(closeChan)
				break
			}
		}
	})
}
