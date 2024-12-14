package kafka

import (
	"crypto/tls"
	"hash"
	"time"

	"github.com/go-gotop/gotop/broker"
	"github.com/go-gotop/gotop/tracing"
	"github.com/segmentio/kafka-go/sasl"
)

const (
	DefaultAddr = "127.0.0.1:9092"

	LeastBytesBalancer    BalancerName = "LeastBytes"
	RoundRobinBalancer    BalancerName = "RoundRobin"
	HashBalancer          BalancerName = "Hash"
	ReferenceHashBalancer BalancerName = "ReferenceHash"
	Crc32Balancer         BalancerName = "CRC32Balancer"
	Murmur2Balancer       BalancerName = "Murmur2Balancer"
)

// context key 定义
type (
	// 初始化配置
	addrsKey         struct{}
	tracingKey       struct{}
	retriesKey       struct{}
	saslMechanismKey struct{}
	tlsConfigKey     struct{}

	// common
	maxAttemptsKey struct{}

	// writer配置
	balancerKey     struct{}
	batchSizeKey    struct{}
	batchTimeoutKey struct{}
	batchBytesKey   struct{}
	asyncKey        struct{}

	readTimeoutKey                   struct{}
	writeTimeoutKey                  struct{}
	allowPublishAutoTopicCreationKey struct{}

	// reader配置
	autoAckKey                struct{}
	queueKey                  struct{}
	queueCapacityKey          struct{}
	minBytesKey               struct{}
	maxBytesKey               struct{}
	maxWaitKey                struct{}
	readLagIntervalKey        struct{}
	heartbeatIntervalKey      struct{}
	commitIntervalKey         struct{}
	partitionWatchIntervalKey struct{}
	watchPartitionChangesKey  struct{}
	sessionTimeoutKey         struct{}
	rebalanceTimeoutKey       struct{}
	retentionTimeKey          struct{}
	startOffsetKey            struct{}
	dialerConfigKey           struct{}
	dialerTimeoutKey          struct{}
)

type BalancerName string

type balancerValue struct {
	Name       BalancerName
	Consistent bool
	Hasher     hash.Hash32
}

type BrokerOption func(*brokerOptions)

type brokerOptions struct {
	Publisher  *kafkaPublisher
	Subscriber *kafkaSubscriber
}

// WithPublisher 设置Kafka发布者实例
func WithPublisher(pub *kafkaPublisher) BrokerOption {
	return func(bo *brokerOptions) {
		bo.Publisher = pub
	}
}

// WithSubscriber 设置Kafka订阅者实例
func WithSubscriber(sub *kafkaSubscriber) BrokerOption {
	return func(bo *brokerOptions) {
		bo.Subscriber = sub
	}
}

// //////////////////////////////////////////////////////////

// WithAddrs 设置Kafka集群的broker地址列表
func WithAddrs(addrs ...string) broker.Option {
	return broker.OptionsContextWithValue(addrsKey{}, addrs)
}

// WithTracings 设置追踪选项
func WithTracings(opts ...tracing.Option) broker.Option {
	return broker.OptionsContextWithValue(tracingKey{}, opts)
}

// WithRetries 设置重试次数，默认不重试
func WithRetries(retries int) broker.Option {
	return broker.OptionsContextWithValue(retriesKey{}, retries)
}

// WithSASLMechanism 设置SASL机制
func WithSASLMechanism(mechanism sasl.Mechanism) broker.Option {
	return broker.OptionsContextWithValue(saslMechanismKey{}, mechanism)
}

// WithTLSConfig 设置TLS配置
func WithTLSConfig(config *tls.Config) broker.Option {
	return broker.OptionsContextWithValue(tlsConfigKey{}, config)
}

// WithLeastBytesBalancer LeastBytes负载均衡器
func WithLeastBytesBalancer() broker.Option {
	return broker.OptionsContextWithValue(balancerKey{},
		&balancerValue{
			Name:   LeastBytesBalancer,
			Hasher: nil,
		},
	)
}

// WithRoundRobinBalancer RoundRobin负载均衡器，默认均衡器。
func WithRoundRobinBalancer() broker.Option {
	return broker.OptionsContextWithValue(balancerKey{},
		&balancerValue{
			Name: RoundRobinBalancer,
		},
	)
}

// WithHashBalancer Hash负载均衡器
func WithHashBalancer(hasher hash.Hash32) broker.Option {
	return broker.OptionsContextWithValue(balancerKey{},
		&balancerValue{
			Name:   HashBalancer,
			Hasher: hasher,
		},
	)
}

// WithReferenceHashBalancer ReferenceHash负载均衡器
func WithReferenceHashBalancer(hasher hash.Hash32) broker.Option {
	return broker.OptionsContextWithValue(balancerKey{},
		&balancerValue{
			Name:   ReferenceHashBalancer,
			Hasher: hasher,
		},
	)
}

// WithCrc32Balancer CRC32负载均衡器
func WithCrc32Balancer(consistent bool) broker.Option {
	return broker.OptionsContextWithValue(balancerKey{},
		&balancerValue{
			Name:       Crc32Balancer,
			Consistent: consistent,
		},
	)
}

// WithMurmur2Balancer Murmur2负载均衡器
func WithMurmur2Balancer(consistent bool) broker.Option {
	return broker.OptionsContextWithValue(balancerKey{},
		&balancerValue{
			Name:       Murmur2Balancer,
			Consistent: consistent,
		},
	)
}

// 发送到分区前缓冲的消息数量限制
//
// 默认的目标批次大小是100条消息
func WithBatchSize(size int) broker.Option {
	return broker.OptionsContextWithValue(batchSizeKey{}, size)
}

// 未满批次的消息发送到Kafka的时间间隔
//
// 默认至少每秒刷新一次
func WithBatchTimeout(timeout time.Duration) broker.Option {
	return broker.OptionsContextWithValue(batchTimeoutKey{}, timeout)
}

// 发送到分区前请求的最大字节数限制
//
// 默认使用Kafka的默认值1048576字节
func WithBatchBytes(bytes int64) broker.Option {
	return broker.OptionsContextWithValue(batchBytesKey{}, bytes)
}

// 设置为true时，WriteMessages方法将永不阻塞
// 这也意味着错误会被忽略，因为调用者无法接收返回值
// 仅在不关心消息是否成功写入Kafka时使用此选项
func WithAsync(async bool) broker.Option {
	return broker.OptionsContextWithValue(asyncKey{}, async)
}

// 消息投递的最大重试次数
//
// 默认最多重试10次
func WithMaxAttempts(attempts int) broker.Option {
	return broker.OptionsContextWithValue(maxAttemptsKey{}, attempts)
}

// Writer执行读操作的超时时间
//
// 默认10秒
func WithReadTimeout(timeout time.Duration) broker.Option {
	return broker.OptionsContextWithValue(readTimeoutKey{}, timeout)
}

// Writer执行写操作的超时时间
//
// 默认10秒
func WithWriteTimeout(timeout time.Duration) broker.Option {
	return broker.OptionsContextWithValue(writeTimeoutKey{}, timeout)
}

// 设置为true时，允许在发布消息时自动创建主题
//
// 默认不允许自动创建主题
func WithAllowPublishAutoTopicCreation(allow bool) broker.Option {
	return broker.OptionsContextWithValue(allowPublishAutoTopicCreationKey{}, allow)
}

////////////////////////////////////////////////////////////

type SubscriberOption func(*subscriberOptions)

type subscriberOptions struct {
}
