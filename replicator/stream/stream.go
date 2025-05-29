package stream

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-gotop/gotop/replicator"
)

type streamReplicatorManager struct {
	config            replicator.ConnectionReplicatorConfig
	connectionFactory replicator.ConnectionFactory
	connections       map[string]replicator.Connection

	// 状态管理
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
	isRunning bool
	startTime time.Time

	// 连接队列
	pendingQueue       chan pendingConnection
	lastConnectionID   string
	lastConnectionTime time.Time

	// 事件处理
	dataHandler  func(data []byte)
	errorHandler func(streamID string, err error)
}

type pendingConnection struct {
	id     string
	config replicator.ConnectionConfig
}

type ConnectionEvent struct {
	StreamID  string
	EventType string
	Timestamp time.Time
	Error     error
}

// NewStreamReplicatorManager 创建流复制器管理器
func NewStreamReplicatorManager(
	config replicator.ConnectionReplicatorConfig,
	connectionFactory replicator.ConnectionFactory,
	options ...StreamReplicatorOption,
) *streamReplicatorManager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &streamReplicatorManager{
		config:             config,
		connectionFactory:  connectionFactory,
		connections:        make(map[string]replicator.Connection),
		ctx:                ctx,
		cancel:             cancel,
		pendingQueue:       make(chan pendingConnection, config.GetConnectionCount()*10),
		lastConnectionTime: time.Now().Add(-config.GetReplicaDelay()),
	}

	// 应用选项
	for _, opt := range options {
		opt(m)
	}

	return m
}

// Connect 建立连接
func (m *streamReplicatorManager) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("stream replicator is already running")
	}

	// 创建所有连接
	for i := 0; i < m.config.GetConnectionCount(); i++ {
		connectionID := fmt.Sprintf("%s-%d", m.config.GetConnectionID(), i)

		connConfig, err := m.config.GenerateConnectionConfig(connectionID)
		if err != nil {
			return fmt.Errorf("failed to generate connection config: %w", err)
		}

		conn := m.connectionFactory.CreateConnection(connConfig)

		// 设置处理器
		if m.dataHandler != nil {
			conn.SetDataHandler(m.dataHandler)
		}
		conn.SetErrorHandler(m.createErrorHandler(connectionID))

		m.connections[connectionID] = conn

		// 加入连接队列
		select {
		case m.pendingQueue <- pendingConnection{
			id:     connectionID,
			config: connConfig,
		}:
		default:
			return fmt.Errorf("connection queue is full")
		}
	}

	m.isRunning = true
	m.startTime = time.Now()

	go m.processConnectionQueue()

	return nil
}

// Disconnect 断开连接
func (m *streamReplicatorManager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	// 停止连接
	if m.cancel != nil {
		m.cancel()
	}

	var lastErr error
	for _, conn := range m.connections {
		if conn.IsConnected() {
			if err := conn.Disconnect(); err != nil {
				lastErr = err
			}
		}
	}

	m.isRunning = false
	return lastErr
}

// IsConnected 检查是否连接
func (m *streamReplicatorManager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, conn := range m.connections {
		if conn.IsConnected() {
			return true
		}
	}

	return false
}

// GetConnectedCount 获取当前连接数
func (m *streamReplicatorManager) GetConnectedCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, conn := range m.connections {
		if conn.IsConnected() {
			count++
		}
	}

	return count
}

func (m *streamReplicatorManager) GetConnectionIDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.connections))
	for id := range m.connections {
		ids = append(ids, id)
	}

	return ids
}

// createErrorHandler 创建错误处理函数
func (m *streamReplicatorManager) createErrorHandler(connectionID string) func(err error) {
	return func(err error) {
		if m.errorHandler != nil {
			m.errorHandler(connectionID, err)
		}

		// 处理重连
		m.handlerConnectionError(connectionID, err)
	}
}

// processConnectionQueue 处理连接队列
func (m *streamReplicatorManager) processConnectionQueue() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case pending := <-m.pendingQueue:
			m.processConnection(pending)
		}
	}
}

// processConnection 处理单个连接
func (m *streamReplicatorManager) processConnection(pending pendingConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查连接间隔
	now := time.Now()
	if waitDuration := m.config.GetReplicaDelay() - now.Sub(m.lastConnectionTime); waitDuration > 0 {
		// 需要等待，重新加入队列
		go func() {
			time.Sleep(waitDuration)
			select {
			case m.pendingQueue <- pending:
			case <-m.ctx.Done():
			}
		}()
		return
	}

	// 尝试连接
	if conn, exists := m.connections[pending.id]; exists {
		if err := conn.Connect(m.ctx); err != nil {
			// 连接失败，重新加入队列
			go func() {
				time.Sleep(time.Second * 5) // 等待5秒后重试
				select {
				case m.pendingQueue <- pending:
				case <-m.ctx.Done():
				}
			}()
			return
		} else {
			// 连接成功，更新状态
			m.lastConnectionTime = now
			m.lastConnectionID = pending.id
		}
	}
}

// handlerConnectionError 处理连接错误
func (m *streamReplicatorManager) handlerConnectionError(connectionID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.connections[connectionID]; exists {
		// 断开当前连接
		_ = conn.Disconnect()
		// 删除连接
		delete(m.connections, connectionID)

		// 生成新配置并重新连接
		if config, err := m.config.GenerateConnectionConfig(connectionID); err == nil {
			newConn := m.connectionFactory.CreateConnection(config)
			newConn.SetDataHandler(m.dataHandler)
			newConn.SetErrorHandler(m.createErrorHandler(connectionID))

			m.connections[connectionID] = newConn

			// 加入连接队列
			go func() {
				select {
				case m.pendingQueue <- pendingConnection{
					id:     connectionID,
					config: config,
				}:
				case <-m.ctx.Done():
				}
			}()
		}
	}
}
