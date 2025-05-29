package replicator

import (
	"context"
	"time"
)

// ConnectionConfig 连接配置
type ConnectionConfig struct {
	ID     string
	URL    string
	Type   string
	Params map[string]interface{}
}

// ConnectionReplicator 连接复制器
type ConnectionReplicator interface {
	// Connect 建立连接
	Connect(ctx context.Context) error

	// Disconnect 断开连接
	Disconnect() error

	// GetConnectedCount 获取当前连接数
	GetConnectedCount() int

	// GetConnectionIDs 获取所有连接ID
	GetConnectionIDs() []string
}

// ConnectionConfig 连接配置接口
type ConnectionReplicatorConfig interface {
	// GetConnectionID 获取连接ID
	GetConnectionID() string

	// GetConnectionCount 获取连接数量
	GetConnectionCount() int

	// GetReplicaDelay 获取副本延迟
	GetReplicaDelay() time.Duration

	// GenerateConnectionConfig 生成连接配置
	GenerateConnectionConfig(connectionID string) (ConnectionConfig, error)

	// Validate 验证配置
	Validate() error
}

// Connection 连接接口
type Connection interface {
	// Connect 建立连接
	Connect(ctx context.Context) error

	// Disconnect 断开连接
	Disconnect() error

	// GetConnectionID 获取连接ID
	GetConnectionID() string

	// ID 获取连接ID
	ID() string

	// Send 发送数据
	Send(data []byte) error

	// SetDataHandler 设置数据处理函数
	SetDataHandler(handler func(data []byte))

	// SetErrorHandler 设置错误处理函数
	SetErrorHandler(handler func(err error))

	// IsConnected 检查是否连接
	IsConnected() bool
}

// ConnectionFactory 连接工厂接口
type ConnectionFactory interface {
	// CreateConnection 创建连接
	CreateConnection(config ConnectionConfig) Connection

	// GetSupportedTypes 获取支持的连接类型
	GetSupportedTypes() []string
}
