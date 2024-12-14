package kafka

import (
	"context"
	"sync"

	"github.com/go-gotop/gotop/broker"
)

type kafkaSubscriber struct {
	sync.RWMutex

	options *broker.Options
}

func NewSubscriber(opts ...broker.Option) (broker.Subscriber, error) {
	po := broker.NewOptions(opts...)

	return &kafkaSubscriber{
		options: &po,
	}, nil
}

func (s *kafkaSubscriber) Subscribe(ctx context.Context, topics []string, handler broker.MessageHandler, opts ...broker.Option) error {
	return nil
}

func (s *kafkaSubscriber) Close() error {
	return nil
}
