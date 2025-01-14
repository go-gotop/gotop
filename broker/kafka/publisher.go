package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/go-gotop/gotop/broker"
	"github.com/go-gotop/gotop/tracing"
	"github.com/go-kratos/kratos/v2/log"
	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"go.opentelemetry.io/otel/attribute"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type kafkaPublisher struct {
	sync.RWMutex

	options   *broker.Options
	retries   int
	connected bool

	saslMechanism sasl.Mechanism
	tlsConfig     *tls.Config
	// 不同主题对应不同写入器
	writer         *Writer
	writerConfig   WriterConfig
	producerTracer *tracing.Tracer

	logger *log.Helper
}

func NewPublisher(logger *log.Helper, opts ...broker.Option) (broker.Publisher, error) {
	po := broker.NewOptions(opts...)

	kp := &kafkaPublisher{
		options: &po,
		writerConfig: WriterConfig{
			Brokers:      []string{},
			Balancer:     &kafkaGo.LeastBytes{},
			Logger:       nil,
			ErrorLogger:  &ErrorLogger{logger: logger},
			BatchTimeout: 10 * time.Millisecond,
			Async:        true, // 异步发送，不等待异常响应
		},
		logger: logger,
	}

	// 应用初始化配置
	kp.applyOptions(po.Context)

	// 应用Writer配置
	kp.applyWriterConfig(po.Context)

	// 创建写入器
	kp.writer = NewWriter()

	return kp, nil
}

// Publish 发布消息到指定的Kafka主题
//
// 参数:
//   - topic: 目标主题名称，消息将被发送到该主题
//   - key: 消息的键，用于分区路由，可以为nil。相同key的消息会被路由到相同分区
//   - value: 消息的实际内容
//   - headers: 消息头信息，包含自定义的键值对元数据。用于传递额外的消息属性，如追踪信息等
//
// 返回值:
//   - error: 发送失败时返回错误信息，成功时返回nil
//
// 示例:
//
//	err := publisher.Publish(
//	    "my-topic",
//	    []byte("user-123"),    // 用户ID作为key
//	    []byte("hello world"), // 消息内容
//	    map[string]string{     // 自定义头信息
//	        "version": "1.0",
//	        "type": "greeting"
//	    },
//	)
func (p *kafkaPublisher) Publish(ctx context.Context, message *broker.Message, opts ...broker.Option) error {
	// 消息构建
	msg := kafkaGo.Message{
		Topic: message.Topic,
		Value: message.Value,
	}

	if message.Headers != nil {
		for k, v := range message.Headers {
			msg.Headers = append(msg.Headers, kafkaGo.Header{Key: k, Value: []byte(v)})
		}
	}

	if message.Key != "" {
		msg.Key = []byte(message.Key)
	}

	// 获取写入器
	var cached bool
	p.Lock()
	writer, ok := p.writer.Writers[message.Topic]
	if !ok {
		writer = p.writer.CreateProducer(p.writerConfig, p.saslMechanism, p.tlsConfig)
		p.writer.Writers[message.Topic] = writer
	} else {
		cached = true
	}
	p.Unlock()

	var err error

	// 开始生产者追踪
	span := p.startProducerSpan(ctx, &msg)
	defer p.finishProducerSpan(span, int32(msg.Partition), 0, err)

	// 写入消息
	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		p.logger.Errorf("WriteMessages error: %s", err.Error())
		err = p.handleWriterError(ctx, err, cached, msg, writer)
	}
	return err
}

// Close 关闭Kafka发布者
func (p *kafkaPublisher) Close() error {
	p.RLock()
	if !p.connected {
		p.RUnlock()
		return nil
	}
	p.RUnlock()

	p.Lock()
	defer p.Unlock()

	p.writer.Close()

	p.connected = false
	return nil
}

func (p *kafkaPublisher) applyOptions(ctx context.Context) {
	if value, ok := ctx.Value(addrsKey{}).([]string); ok {
		for _, addr := range value {
			if len(addr) > 0 {
				p.writerConfig.Brokers = append(p.writerConfig.Brokers, addr)
			}
		}
	}

	if value, ok := ctx.Value(tracingKey{}).([]tracing.Option); ok {
		p.producerTracer = tracing.NewTracer(trace.SpanKindProducer, "kafka-producer", value...)
	}

	if value, ok := ctx.Value(retriesKey{}).(int); ok {
		p.retries = value
	}

	if value, ok := ctx.Value(saslMechanismKey{}).(sasl.Mechanism); ok {
		p.saslMechanism = value
	}

	if value, ok := ctx.Value(tlsConfigKey{}).(*tls.Config); ok {
		p.tlsConfig = value
	}

}

func (p *kafkaPublisher) applyWriterConfig(ctx context.Context) {
	if value, ok := ctx.Value(batchSizeKey{}).(int); ok {
		p.writerConfig.BatchSize = value
	}

	if value, ok := ctx.Value(batchTimeoutKey{}).(time.Duration); ok {
		p.writerConfig.BatchTimeout = value
	}

	if value, ok := ctx.Value(batchBytesKey{}).(int64); ok {
		p.writerConfig.BatchBytes = value
	}

	if value, ok := ctx.Value(asyncKey{}).(bool); ok {
		p.writerConfig.Async = value
	}

	if value, ok := ctx.Value(maxAttemptsKey{}).(int); ok {
		p.writerConfig.MaxAttempts = value
	}

	if value, ok := ctx.Value(readTimeoutKey{}).(time.Duration); ok {
		p.writerConfig.ReadTimeout = value
	}

	if value, ok := ctx.Value(writeTimeoutKey{}).(time.Duration); ok {
		p.writerConfig.WriteTimeout = value
	}

	if value, ok := ctx.Value(allowPublishAutoTopicCreationKey{}).(bool); ok {
		p.writerConfig.AllowAutoTopicCreation = value
	}

	if value, ok := ctx.Value(balancerKey{}).(*balancerValue); ok {
		switch value.Name {
		default:
		case LeastBytesBalancer:
			p.writerConfig.Balancer = &kafkaGo.LeastBytes{}
		case RoundRobinBalancer:
			p.writerConfig.Balancer = &kafkaGo.RoundRobin{}
		case HashBalancer:
			p.writerConfig.Balancer = &kafkaGo.Hash{
				Hasher: value.Hasher,
			}
		case ReferenceHashBalancer:
			p.writerConfig.Balancer = &kafkaGo.ReferenceHash{
				Hasher: value.Hasher,
			}
		case Crc32Balancer:
			p.writerConfig.Balancer = &kafkaGo.CRC32Balancer{
				Consistent: value.Consistent,
			}
		case Murmur2Balancer:
			p.writerConfig.Balancer = &kafkaGo.Murmur2Balancer{
				Consistent: value.Consistent,
			}
		}
	}
}

// 处理写入错误
func (p *kafkaPublisher) handleWriterError(ctx context.Context, err error, cached bool, msg kafkaGo.Message, writer *kafkaGo.Writer) error {
	switch cached {
	case false:
		// 非缓存写入器, 重试
		var kerr kafkaGo.Error
		if errors.As(err, &kerr) {
			if kerr.Temporary() && !kerr.Timeout() {
				time.Sleep(200 * time.Millisecond)
				err = writer.WriteMessages(ctx, msg)
			}
		}
	case true:
		// 缓存写入器, 重试
		p.Lock()
		if err = writer.Close(); err != nil {
			p.Unlock()
			break
		}
		delete(p.writer.Writers, msg.Topic)
		p.Unlock()

		writer = p.writer.CreateProducer(p.writerConfig, p.saslMechanism, p.tlsConfig)
		for i := 0; i < p.retries; i++ {
			if err = writer.WriteMessages(ctx, msg); err == nil {
				p.Lock()
				p.writer.Writers[msg.Topic] = writer
				p.Unlock()
				break
			}
		}
	}
	return err
}

// 开始生产者追踪
func (p *kafkaPublisher) startProducerSpan(ctx context.Context, msg *kafkaGo.Message) trace.Span {
	if p.producerTracer == nil {
		return nil
	}

	carrier := NewMessageCarrier(msg)

	attrs := []attribute.KeyValue{
		semConv.MessagingSystemKey.String("kafka"),
		semConv.MessagingDestinationKindTopic,
		semConv.MessagingDestinationKey.String(msg.Topic),
	}

	var span trace.Span
	_, span = p.producerTracer.Start(ctx, carrier, attrs...)

	return span
}

// 结束生产者追踪
func (p *kafkaPublisher) finishProducerSpan(span trace.Span, partition int32, offset int64, err error) {
	if p.producerTracer == nil {
		return
	}

	attrs := []attribute.KeyValue{
		semConv.MessagingMessageIDKey.String(strconv.FormatInt(offset, 10)),
		semConv.MessagingKafkaPartitionKey.Int64(int64(partition)),
	}

	p.producerTracer.End(context.Background(), span, err, attrs...)
}
