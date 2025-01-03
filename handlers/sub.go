package handlers

import (
	"github.com/jeel9dot/social-steam/constants"
	social_stream "github.com/jeel9dot/social-steam/protobuf/genproto/social-stream"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SocialStreamGrpcHandler) SubcribeTrades(req *social_stream.SubcribeRequest, stream social_stream.SocialSteamService_SubcribeTradesServer) error {
	if len(req.Subjects) == 0 {
		return status.Errorf(codes.InvalidArgument, constants.ErrMsgEmptySubject)
	}

	s.logger.Info("gRPC client subscribed", zap.Strings("subjects", req.Subjects))

	// Channel to handle cleanup on stream closure
	closeChan := make(chan struct{})

	// Function to subscribe to a subject and stream messages to the client
	subscribeToSubjects := func(subjects []string) error {
		for _, subject := range subjects {
			subscription, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
				select {
				case <-closeChan:
					// Stop sending messages after the client disconnects
					return
				default:
					// Send the message to the gRPC stream
					s.logger.Debug("Sending message to gRPC client", zap.String("subject", subject), zap.String("message", string(msg.Data)))
					if err := stream.Send(&social_stream.SubcribeResponce{Subject: subject, Msg: string(msg.Data)}); err != nil {
						s.logger.Warn("Error sending message to gRPC client", zap.Error(err), zap.String("subject", subject))
						// Trigger cleanup
						close(closeChan)
					}
				}
			})
			if err != nil {
				s.logger.Warn("Subscription failed", zap.String("subject", subject), zap.Error(err))
				return err
			}

			// Ensure unsubscription on exit
			defer func(sub string) {
				if err := subscription.Unsubscribe(); err != nil {
					s.logger.Warn("Error during unsubscription", zap.String("subject", sub), zap.Error(err))
				} else {
					s.logger.Debug("Successfully unsubscribed", zap.String("subject", sub))
				}
			}(subject)
		}
		return nil
	}

	// Subscribe to the requested subjects
	if err := subscribeToSubjects(req.Subjects); err != nil {
		return status.Errorf(codes.Internal, "Failed to subscribe to subjects: %v", err)
	}

	// Wait for the client to disconnect or an error to occur
	select {
	case <-stream.Context().Done():
		// Client has disconnected
		s.logger.Info("gRPC client disconnected", zap.Strings("subjects", req.Subjects))
		close(closeChan)
		return nil
	case <-closeChan:
		// Cleanup triggered due to an error
		s.logger.Warn("Cleanup triggered due to error", zap.Strings("subjects", req.Subjects))
		return status.Errorf(codes.Internal, "Subscription terminated due to an error")
	}
}
