package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-gotop/gotop/broker"
	"github.com/go-gotop/gotop/tracing"
	"github.com/google/uuid"
	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"go.opentelemetry.io/otel/attribute"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type subscriber struct {
	sync.RWMutex

	topic   string
	options broker.Options
	handler broker.Handler
	reader  *kafkaGo.Reader
	stopCh  chan struct{}
}

type kafkaSubscriber struct {
	sync.RWMutex
	// 订阅者map, key 为 topic
	suberMap       map[string]*subscriber
	retries        int
	consumerTracer *tracing.Tracer
	readerConfig   kafkaGo.ReaderConfig
	saslMechanism  sasl.Mechanism
	tlsConfig      *tls.Config
	options        *broker.Options
}

func NewSubscriber(opts ...broker.Option) (broker.Subscriber, error) {
	po := broker.NewOptions(opts...)

	ks := &kafkaSubscriber{
		suberMap: make(map[string]*subscriber),
		readerConfig: kafkaGo.ReaderConfig{
			WatchPartitionChanges: true,
			MaxWait:               500 * time.Millisecond,
		},
		options: &po,
	}

	ks.applyOptions(po.Context)
	ks.applyReaderConfig(po.Context)

	if ks.tlsConfig != nil {
		if ks.readerConfig.Dialer == nil {
			ks.readerConfig.Dialer = &kafkaGo.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
				TLS:       ks.tlsConfig,
			}
		} else {
			ks.readerConfig.Dialer.TLS = ks.tlsConfig
		}
	}

	return ks, nil
}

func (s *kafkaSubscriber) Subscribe(ctx context.Context, topics []string, handler broker.Handler, opts ...broker.Option) error {
	po := broker.NewOptions(opts...)

	for _, topic := range topics {
		if err := s.subscribe(ctx, topic, handler, po); err != nil {
			return err
		}
	}
	return nil
}

func (s *kafkaSubscriber) Unsubscribe(ctx context.Context, topics []string) error {
	for _, topic := range topics {
		s.Lock()
		sub, ok := s.suberMap[topic]
		if ok {
			sub.reader.Close()
			close(sub.stopCh)
			delete(s.suberMap, topic)
		}
		s.Unlock()
	}
	return nil
}

func (s *kafkaSubscriber) subscribe(ctx context.Context, topic string, handler broker.Handler, po broker.Options) error {
	autoAck := true
	queue := uuid.New().String()

	if value, ok := ctx.Value(autoAckKey{}).(bool); ok {
		autoAck = value
	}

	if value, ok := ctx.Value(queueKey{}).(string); ok {
		queue = value
	}

	readerConfig := s.readerConfig
	readerConfig.Topic = topic

	readerConfig.GroupID = queue

	sub := &subscriber{
		topic:   topic,
		options: po,
		handler: handler,
		reader:  kafkaGo.NewReader(readerConfig),
		stopCh:  make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-sub.stopCh:
				return
			case <-po.Context.Done():
				return
			default:
				msg, err := sub.reader.FetchMessage(po.Context)
				if err != nil {
					if errors.Is(err, io.EOF) {
						continue
					} else {
						continue
					}
				}

				ctx, span := s.startConsumerSpan(po.Context, &msg)

				m := &broker.Message{
					Topic: msg.Topic,
					Key:   string(msg.Key),
					Value: msg.Value,
				}

				if err = sub.handler(ctx, m); err != nil {
					s.finishConsumerSpan(span, err)
					continue
				}

				if autoAck {
					err = sub.reader.CommitMessages(ctx, msg)
					if err != nil {
						s.finishConsumerSpan(span, err)
						continue
					}
				}

				s.finishConsumerSpan(span, err)
			}
		}
	}()

	s.Lock()
	s.suberMap[topic] = sub
	s.Unlock()

	return nil
}

func (s *kafkaSubscriber) Close() error {
	s.Lock()
	for _, sub := range s.suberMap {
		sub.reader.Close()
		close(sub.stopCh)
	}
	s.suberMap = make(map[string]*subscriber)
	s.Unlock()

	return nil
}

func (s *kafkaSubscriber) applyOptions(ctx context.Context) {
	if value, ok := ctx.Value(addrsKey{}).([]string); ok {
		for _, addr := range value {
			if len(addr) > 0 {
				s.readerConfig.Brokers = append(s.readerConfig.Brokers, addr)
			}
		}
	}

	if value, ok := ctx.Value(tracingKey{}).([]tracing.Option); ok {
		s.consumerTracer = tracing.NewTracer(trace.SpanKindConsumer, "kafka-consumer", value...)
	}

	if value, ok := ctx.Value(retriesKey{}).(int); ok {
		s.retries = value
	}

	if value, ok := ctx.Value(saslMechanismKey{}).(sasl.Mechanism); ok {
		s.saslMechanism = value
	}

	if value, ok := ctx.Value(tlsConfigKey{}).(*tls.Config); ok {
		s.tlsConfig = value
	}

}

func (s *kafkaSubscriber) applyReaderConfig(ctx context.Context) {
	if s.readerConfig.Dialer == nil {
		s.readerConfig.Dialer = kafkaGo.DefaultDialer
	}

	if value, ok := ctx.Value(queueCapacityKey{}).(int); ok {
		s.readerConfig.QueueCapacity = value
	}

	if value, ok := ctx.Value(minBytesKey{}).(int); ok {
		s.readerConfig.MinBytes = value
	}

	if value, ok := ctx.Value(maxBytesKey{}).(int); ok {
		s.readerConfig.MaxBytes = value
	}

	if value, ok := ctx.Value(maxWaitKey{}).(time.Duration); ok {
		s.readerConfig.MaxWait = value
	}

	if value, ok := ctx.Value(readLagIntervalKey{}).(time.Duration); ok {
		s.readerConfig.ReadLagInterval = value
	}

	if value, ok := ctx.Value(heartbeatIntervalKey{}).(time.Duration); ok {
		s.readerConfig.HeartbeatInterval = value
	}

	if value, ok := ctx.Value(commitIntervalKey{}).(time.Duration); ok {
		s.readerConfig.CommitInterval = value
	}

	if value, ok := ctx.Value(partitionWatchIntervalKey{}).(time.Duration); ok {
		s.readerConfig.PartitionWatchInterval = value
	}

	if value, ok := ctx.Value(watchPartitionChangesKey{}).(bool); ok {
		s.readerConfig.WatchPartitionChanges = value
	}

	if value, ok := ctx.Value(sessionTimeoutKey{}).(time.Duration); ok {
		s.readerConfig.SessionTimeout = value
	}

	if value, ok := ctx.Value(rebalanceTimeoutKey{}).(time.Duration); ok {
		s.readerConfig.RebalanceTimeout = value
	}

	if value, ok := ctx.Value(retentionTimeKey{}).(time.Duration); ok {
		s.readerConfig.RetentionTime = value
	}

	if value, ok := ctx.Value(startOffsetKey{}).(int64); ok {
		s.readerConfig.StartOffset = value
	}

	if value, ok := ctx.Value(maxAttemptsKey{}).(int); ok {
		s.readerConfig.MaxAttempts = value
	}

	if value, ok := ctx.Value(saslMechanismKey{}).(sasl.Mechanism); ok {
		s.saslMechanism = value

		if s.readerConfig.Dialer == nil {
			dialer := &kafkaGo.Dialer{
				Timeout:       10 * time.Second,
				DualStack:     true,
				SASLMechanism: s.saslMechanism,
			}
			s.readerConfig.Dialer = dialer
		} else {
			s.readerConfig.Dialer.SASLMechanism = s.saslMechanism
		}
	}

	if value, ok := ctx.Value(dialerTimeoutKey{}).(time.Duration); ok {
		if s.readerConfig.Dialer != nil {
			s.readerConfig.Dialer.Timeout = value
		}
	}

	if value, ok := ctx.Value(dialerConfigKey{}).(*kafkaGo.Dialer); ok {
		s.readerConfig.Dialer = value
	}
}

func (s *kafkaSubscriber) startConsumerSpan(ctx context.Context, msg *kafkaGo.Message) (context.Context, trace.Span) {
	if s.consumerTracer == nil {
		return ctx, nil
	}

	carrier := NewMessageCarrier(msg)

	attrs := []attribute.KeyValue{
		semConv.MessagingSystemKey.String("kafka"),
		semConv.MessagingDestinationKindTopic,
		semConv.MessagingDestinationKey.String(msg.Topic),
		semConv.MessagingOperationReceive,
		semConv.MessagingMessageIDKey.String(strconv.FormatInt(msg.Offset, 10)),
		semConv.MessagingKafkaPartitionKey.Int64(int64(msg.Partition)),
	}

	var span trace.Span
	ctx, span = s.consumerTracer.Start(ctx, carrier, attrs...)

	return ctx, span
}

func (s *kafkaSubscriber) finishConsumerSpan(span trace.Span, err error) {
	if s.consumerTracer == nil {
		return
	}

	s.consumerTracer.End(context.Background(), span, err)
}
