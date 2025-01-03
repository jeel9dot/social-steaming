package handlers

import (
	"context"

	"github.com/jeel9dot/social-steam/constants"
	social_stream "github.com/jeel9dot/social-steam/protobuf/genproto/social-stream"
	"go.uber.org/zap"
)

func (s *SocialStreamGrpcHandler) PublishTrades(ctx context.Context, req *social_stream.PublisherRequest) (*social_stream.PublisherResponce, error) {
	subject := req.GetSubject()
	msg := req.GetMsg()

	if subject == "" || msg == "" {
		return &social_stream.PublisherResponce{
			Success: false,
			Message: constants.ErrMsgEmptySubjectOrMsg,
		}, nil
	}

	err := s.nc.Publish(subject, msg)
	if err != nil {
		s.logger.Error("error while publishing message", zap.Error(err))
		return &social_stream.PublisherResponce{
			Success: false,
			Message: constants.ErrMsgNotPublished,
		}, err
	}

	s.logger.Debug("message published", zap.String("subject", subject), zap.String("message", msg))
	return &social_stream.PublisherResponce{
		Success: true,
		Message: constants.SuccessMsgPublished,
	}, nil
}
