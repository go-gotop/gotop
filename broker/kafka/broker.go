package kafka

import (
	"github.com/go-gotop/gotop/broker"
)

type kafkaBroker struct {
	publisher  *kafkaPublisher
	subscriber *kafkaSubscriber
}

func NewBroker(options ...BrokerOption) (broker.Broker, error) {
	bo := &brokerOptions{}

	for _, opt := range options {
		opt(bo)
	}

	return &kafkaBroker{
		publisher:  bo.Publisher,
		subscriber: bo.Subscriber,
	}, nil
}

func (k *kafkaBroker) Publisher() (broker.Publisher, error) {
	if k.publisher == nil {
		return nil, broker.ErrPublisherNotConfigured
	}
	return k.publisher, nil
}

func (k *kafkaBroker) Subscriber() (broker.Subscriber, error) {
	if k.subscriber == nil {
		return nil, broker.ErrSubscriberNotConfigured
	}
	return k.subscriber, nil
}

func (k *kafkaBroker) Close() error {
	var err error
	if k.publisher != nil {
		if e := k.publisher.Close(); e != nil {
			err = e
		}
	}
	if k.subscriber != nil {
		if e := k.subscriber.Close(); e != nil {
			err = e
		}
	}
	return err
}
