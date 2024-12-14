package kafka

import (
	"crypto/tls"
	"time"

	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
)

type WriterConfig struct {
	// Kafka集群的broker地址列表
	Brokers []string

	// 用于在分区之间分配消息的均衡器
	//
	// 默认使用轮询分配方式
	Balancer kafkaGo.Balancer

	// 消息投递的最大重试次数
	//
	// 默认最多重试10次
	MaxAttempts int

	// 发送到分区前缓冲的消息数量限制
	//
	// 默认的目标批次大小是100条消息
	BatchSize int

	// 发送到分区前请求的最大字节数限制
	//
	// 默认使用Kafka的默认值1048576字节
	BatchBytes int64

	// 未满批次的消息发送到Kafka的时间间隔
	//
	// 默认至少每秒刷新一次
	BatchTimeout time.Duration

	// Writer执行读操作的超时时间
	//
	// 默认10秒
	ReadTimeout time.Duration

	// Writer执行写操作的超时时间
	//
	// 默认10秒
	WriteTimeout time.Duration

	// 生产请求需要等待的分区副本确认数
	// 默认值为-1，表示等待所有副本确认
	// 大于0的值表示需要多少个副本确认才算成功
	//
	// 当前kafka-go版本(v0.3)不支持设置为0
	// 如果需要该功能，需要升级到v0.4版本
	RequiredAcks kafkaGo.RequiredAcks

	// 设置为true时，WriteMessages方法将永不阻塞
	// 这也意味着错误会被忽略，因为调用者无法接收返回值
	// 仅在不关心消息是否成功写入Kafka时使用此选项
	Async bool

	// 如果不为nil，用于报告Writer内部变化的日志记录器
	Logger kafkaGo.Logger

	// 用于报告错误的日志记录器
	// 如果为nil，Writer将使用Logger代替
	ErrorLogger kafkaGo.Logger

	// 允许Writer在主题不存在时自动创建主题
	AllowAutoTopicCreation bool
}

type Writer struct {
	Writer                  *kafkaGo.Writer
	Writers                 map[string]*kafkaGo.Writer
}

func NewWriter() *Writer {
	return &Writer{
		Writers: make(map[string]*kafkaGo.Writer),
	}
}

func (w *Writer) Close() {
	if w.Writer != nil {
		_ = w.Writer.Close()
	}
	for _, writer := range w.Writers {
		_ = writer.Close()
	}
	w.Writer = nil
	w.Writers = nil
}

// CreateProducer 创建一个kafka-go Writer实例
func (w *Writer) CreateProducer(writerConfig WriterConfig, saslMechanism sasl.Mechanism, tlsConfig *tls.Config) *kafkaGo.Writer {
	sharedTransport := &kafkaGo.Transport{
		SASL: saslMechanism,
		TLS:  tlsConfig,
	}

	writer := &kafkaGo.Writer{
		Transport: sharedTransport,

		Addr:                   kafkaGo.TCP(writerConfig.Brokers...),
		Balancer:               writerConfig.Balancer,
		MaxAttempts:            writerConfig.MaxAttempts,
		BatchSize:              writerConfig.BatchSize,
		BatchBytes:             writerConfig.BatchBytes,
		BatchTimeout:           writerConfig.BatchTimeout,
		ReadTimeout:            writerConfig.ReadTimeout,
		WriteTimeout:           writerConfig.WriteTimeout,
		RequiredAcks:           writerConfig.RequiredAcks,
		Async:                  writerConfig.Async,
		Logger:                 writerConfig.Logger,
		ErrorLogger:            writerConfig.ErrorLogger,
		AllowAutoTopicCreation: writerConfig.AllowAutoTopicCreation,
	}

	return writer
}
