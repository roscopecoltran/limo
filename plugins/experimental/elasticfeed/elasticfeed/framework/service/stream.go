package service

import (
	"errors"
	"github.com/roscopecoltran/feedify/stream"
)

type StreamService struct {
	Message *stream.StreamMessage
}

func (s *StreamService) Name() string {
	return "stream-service"
}

func NewStream() (*StreamService, error) {
	message, err := stream.NewStreamMessage()
	if err != nil {
		return nil, errors.New("Cannot create stream service.")
	}
	return &StreamService{message}, nil
}
