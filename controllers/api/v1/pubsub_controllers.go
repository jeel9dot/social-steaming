package v1

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/jeel9dot/social-steam/nats"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type PublishRequest struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

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

func (ctrl *PubSubController) subscribeToSubjects(subjects []string, conn *websocket.Conn, closeChan chan struct{}) error {
	for _, subject := range subjects {
		subscription, err := ctrl.nc.Subscribe(subject, func(msg *nc.Msg) {
			select {
			case <-closeChan:
				// Stop sending messages after disconnect
				return
			default:
				// Send message to WebSocket client
				ctrl.logger.Debug("Sending message to WebSocket client", zap.String("subject", subject), zap.String("message", string(msg.Data)))
				if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
					ctrl.logger.Warn("Error while sending message to WebSocket client", zap.Error(err), zap.String("subject", subject))
					// Trigger disconnection cleanup
					close(closeChan)
				}
			}
		})
		if err != nil {
			ctrl.logger.Warn("Subscription failed", zap.String("subject", subject), zap.Error(err))
			return err
		}

		// Ensure unsubscription on exit
		defer func(sub string) {
			if err := subscription.Unsubscribe(); err != nil {
				ctrl.logger.Warn("Error during unsubscription", zap.String("subject", sub), zap.Error(err))
			} else {
				ctrl.logger.Debug("Successfully unsubscribed", zap.String("subject", sub))
			}
		}(subject)
	}

	return nil
}

func (ctrl *PubSubController) Subscribe() fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		subjects := conn.Query("subjects")
		if subjects == "" {
			ctrl.logger.Error("Subscription attempt without subjects")
			_ = conn.Close()
			return
		}

		// Parse the subjects into a slice
		subjectList := strings.Split(subjects, ",")
		ctrl.logger.Info("WebSocket client connected", zap.Strings("subjects", subjectList))

		defer func() {
			// Close WebSocket connection on exit
			if err := conn.Close(); err != nil {
				ctrl.logger.Warn("Error while closing WebSocket connection", zap.Error(err))
			} else {
				ctrl.logger.Info("WebSocket connection closed", zap.Strings("subjects", subjectList))
			}
		}()

		// Create a channel to handle WebSocket disconnection
		closeChan := make(chan struct{})

		// Subscribe to the subjects
		if err := ctrl.subscribeToSubjects(subjectList, conn, closeChan); err != nil {
			ctrl.logger.Error("Subscription process failed", zap.Error(err))
			return
		}

		// Keep the connection open until the client disconnects
		ctrl.logger.Debug("Waiting for WebSocket client to disconnect", zap.Strings("subjects", subjectList))
		<-closeChan
		ctrl.logger.Info("WebSocket client disconnected", zap.Strings("subjects", subjectList))
	})
}
