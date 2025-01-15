package kafka

import (
	"github.com/go-gotop/gotop/broker"
)

type kafkaBroker struct {
	publisher  broker.Publisher
	subscriber broker.Subscriber
}

func NewBroker(options ...broker.Option) (broker.Broker, error) {

	kp, err := NewPublisher(options...)
	if err != nil {
		return nil, err
	}

	ks, err := NewSubscriber(options...)
	if err != nil {
		return nil, err
	}

	return &kafkaBroker{
		publisher:  kp,
		subscriber: ks,
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
