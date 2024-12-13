package notifier

import (
    "time"
)

// Message 表示一条需要发送的通知消息。
type Message struct {
	// From 发送者
	From      string
	// To 接收者
	To        []string
	// Subject 主题
	Subject   string
	// Body 内容
	Body      string
	// Metadata 元数据
	Metadata  map[string]string
	// Priority 优先级
	Priority  int
	// Timestamp 时间戳
	Timestamp time.Time
}

// Notifier 接口代表一个通用的通知发送器，实现者应根据 Message 的内容将其发送到对应的渠道。
type Notifier interface {
	// Notify 发送通知消息
	Notify(msg Message) error
}

// NotifierFactory 接口定义了创建 Notifier 实例的工厂方法。
// 不同类型的 Notifier（如SMS、Telegram、Email）可通过此工厂根据配置创建实例。
type NotifierFactory interface {
	// CreateNotifier 创建 Notifier 实例
	CreateNotifier(cfg interface{}) (Notifier, error)
}

// NotifierRouter 接口用于对消息进行路由，决定使用哪些 Notifier 来发送消息。
// 可在实现时根据消息类型、优先级或元数据选择合适的 Notifier。
type NotifierRouter interface {
	// RegisterNotifier 注册 Notifier
	RegisterNotifier(name string, notifier Notifier)
	// GetNotifier 获取 Notifier
	GetNotifier(name string) (Notifier, bool)
	// Route 路由消息
	Route(msg Message) ([]Notifier, error)
}

// 通用配置结构，可根据实际需要扩展。
// 不同类型的 Notifier 可在 Notifiers 字段中定义自己的配置结构。
type GlobalConfig struct {
	// Notifiers 通知器配置
	Notifiers map[string]interface{} `yaml:"notifiers"`
}

// 可选的中间件接口，用来对 Notifier 包装。例如添加重试、日志、熔断等逻辑。
type NotifierMiddleware interface {
	// Wrap 包装 Notifier
	Wrap(notifier Notifier) Notifier
}

// 可选的观察者接口，用于在发送前后进行日志、监控或其他钩子处理。
type NotifierObserver interface {
	// BeforeSend 发送前钩子
	BeforeSend(msg Message)
	// AfterSend 发送后钩子
	AfterSend(msg Message, err error)
}
