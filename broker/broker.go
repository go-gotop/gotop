package broker

import "context"

// Message 表示传递的消息载体。
// Value为消息内容主体，Key为可选的用于路由/分区的键，Topic表示逻辑上的主题/频道。
// Metadata则是实现方可选填充的元数据，用于区别实现的特定属性（例如：Kafka的Offset、Partition、Redis的Stream ID等）。
type Message struct {
	Key      []byte
	Value    []byte
	Topic    string
	Headers  map[string]string
	Metadata map[string]interface{} // 用于存放特定实现的额外信息，如偏移量、分区信息等
}

// MessageHandler 用于上层业务处理消息的回调接口。
// 不同的底层实现可以在调用此方法前后进行Offset提交、Ack确认等动作，但上层并不需要感知。
type MessageHandler interface {
	HandleMessage(ctx context.Context, msg *Message) error
}

// Publisher 定义通用的消息发布者接口。
// 不关心底层是Kafka Topic、Redis Channel、RabbitMQ Exchange，统一为Publish操作。
type Publisher interface {
	// Publish发布消息到指定的主题（或通道）。
	// 返回错误用于告知发布失败，具体重试策略由上层或底层实现负责。
	Publish(ctx context.Context, topic string, key []byte, value []byte, headers map[string]string, opts ...Option) error

	// Close用于释放发布者相关资源。
	Close() error
}

// Subscriber 定义通用的消息订阅者接口。
// 上层通过为Subscriber注册一个MessageHandler来处理指定主题的消息。
// 底层实现例如Kafka的Consumer、Redis的Subscriber等负责调用handler。
type Subscriber interface {
	// Subscribe接收一个或多个主题，并将收到的消息交给handler处理。
	// 订阅可以是阻塞式，也可以异步在内部协程中进行消息消费。
	Subscribe(ctx context.Context, topics []string, handler MessageHandler, opts ...Option) error

	// Close用于释放订阅者相关资源，断开与底层系统的连接。
	Close() error
}

// Broker 定义统一的Broker客户端接口。
// 通过此接口可获得Publisher与Subscriber实例。
// 此接口也可包括如健康检查(HealthCheck)或全局Close等方法。
// 不论是Kafka、Redis还是其他系统的实现，都应遵从此接口。
type Broker interface {
	// Publisher返回一个Publisher实例。
	Publisher() (Publisher, error)

	// Subscriber返回一个Subscriber实例。
	Subscriber() (Subscriber, error)

	// Close关闭整个BrokerClient，释放底层连接和资源。
	Close() error
}
