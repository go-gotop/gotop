package binance

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/go-gotop/gotop/stream"
)

// BinanceRequest Binance Stream的配置参数
type BinanceRequest struct {
	// WebSocket请求的URL
	URL string
	// Logger 可选的日志记录器，用于调试
	Logger *slog.Logger
	// Handler 数据处理函数
	Handler func(data []byte)
	// ErrorHandler 错误处理函数
	ErrorHandler func(err error)
}

// dialFunc定义用于方便在测试时mock连接的逻辑
type dialFunc func(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)

// BinanceStream 是Binance Stream的核心结构体
type BinanceStream struct {
	mu     sync.Mutex
	id     string
	st stream.StreamType
	cfg    BinanceRequest
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc

	// 用于重连的dial函数，可在测试中mock
	dialer dialFunc

	// 心跳与超时配置
	pingInterval time.Duration
	pongWait     time.Duration
	writeWait    time.Duration

	// 用于在24小时后重连的ticker和控制重连的通道
	reconnectInterval time.Duration
	reconnectTicker   *time.Ticker

	// 用于等待后台goroutine的结束
	wg sync.WaitGroup

	// 用于标识后台goroutine的完成，以便Disconnect时等待
	doneCh chan struct{}

	// 用于记录数据序列的占位变量，未来可用于无缝重连的数据去重处理
	lastSequence int64

	// 是否正在尝试重连
	reconnecting bool
}

// NewBinanceStream 创建一个新的BinanceStream
func NewBinanceStream(id string, st stream.StreamType, opts ...Option) *BinanceStream {
	b := &BinanceStream{
		id:     id,
		st:     st,
		dialer:            defaultDialer,
		pingInterval:      30 * time.Second,
		pongWait:          60 * time.Second,
		writeWait:         5 * time.Second,
		reconnectInterval: 23 * time.Hour, // 24小时自动重连周期
	}
	applyOptions(b, opts...)
	return b
}

// defaultDialer 默认的dial函数
func defaultDialer(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error) {
	var d websocket.Dialer
	return d.Dial(urlStr, requestHeader)
}

// Connect 连接到Binance Stream
func (b *BinanceStream) Connect(ctx context.Context, cfg BinanceRequest) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.cfg = cfg
	b.ctx, b.cancel = context.WithCancel(ctx)
	b.reconnecting = false

	if err := b.connect(); err != nil {
		return err
	}

	// 启动后台goroutine
	b.startGoroutines()

	b.log("Connected successfully")
	return nil
}

// ID 返回当前连接的ID
func (b *BinanceStream) ID() string {
	return b.id
}

// Type 返回当前连接的类型
func (b *BinanceStream) Type() stream.StreamType {
	return b.st
}

// Disconnect 断开当前连接
func (b *BinanceStream) Disconnect() error {
	b.mu.Lock()
	// 取消当前上下文
	if b.cancel != nil {
		b.cancel()
	}
	// 关闭当前连接
	if b.conn != nil {
		_ = b.conn.Close()
		b.conn = nil
	}
	b.mu.Unlock()

	// 等待所有goroutine结束
	b.wg.Wait()
	b.log("Disconnected")
	return nil
}

func (b *BinanceStream) connect() error {
	c, _, err := b.dialer(b.cfg.URL, nil)
	if err != nil {
		b.handleErr(err)
		return err
	}

	b.conn = c
	b.conn.SetReadDeadline(time.Now().Add(b.pongWait))
	b.conn.SetPongHandler(func(string) error {
		b.conn.SetReadDeadline(time.Now().Add(b.pongWait))
		return nil
	})
	return nil
}

// startGoroutines 启动后台goroutine
func (b *BinanceStream) startGoroutines() {
	// 开始readLoop
	b.wg.Add(1)
	go b.readLoop()

	// 开始pingLoop
	b.wg.Add(1)
	go b.pingLoop()

	// 开始autoReconnectLoop
	b.wg.Add(1)
	go b.autoReconnectLoop()
}

// readLoop 读取数据
func (b *BinanceStream) readLoop() {
	defer b.wg.Done()

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			b.mu.Lock()
			conn := b.conn
			b.mu.Unlock()

			if conn == nil {
				// 如果没有连接，稍等等待重连
				time.Sleep(500 * time.Millisecond)
				continue
			}

			_, msg, err := conn.ReadMessage()
			if err != nil {
				// 通知错误
				b.handleErr(err)

				// 异步尝试重连，readLoop立即返回，防止死锁
				go b.attemptReconnect()

				return // 结束当前readLoop
			}

			if b.cfg.Handler != nil {
				b.cfg.Handler(msg)
			}
		}
	}
}

func (b *BinanceStream) pingLoop() {
	defer b.wg.Done()

	ticker := time.NewTicker(b.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			b.mu.Lock()
			conn := b.conn
			b.mu.Unlock()

			if conn == nil {
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(b.writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				b.handleErr(err)
				// 异步重连
				go b.attemptReconnect()
				return
			}
		}
	}
}

// autoReconnectLoop 自动重连循环
func (b *BinanceStream) autoReconnectLoop() {
	defer b.wg.Done()

	ticker := time.NewTicker(b.reconnectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			b.log("Time-based reconnect triggered")
			go b.attemptReconnect()
			return
		}
	}
}

func (b *BinanceStream) attemptReconnect() {
	b.mu.Lock()
	if b.reconnecting {
		b.mu.Unlock()
		return
	}
	b.reconnecting = true
	b.mu.Unlock()

	b.log("Attempting reconnect...")

	// 停止当前上下文，等待goroutine全部退出
	b.mu.Lock()
	if b.cancel != nil {
		b.cancel()
	}
	conn := b.conn
	b.conn = nil
	b.mu.Unlock()

	if conn != nil {
		conn.Close()
	}

	// 等待所有当前goroutine结束
	b.wg.Wait()

	// 创建新上下文
	ctx, cancel := context.WithCancel(context.Background())

	b.mu.Lock()
	b.ctx = ctx
	b.cancel = cancel
	b.reconnecting = false
	b.mu.Unlock()

	// 简单重试机制
	for i := 0; i < 5; i++ { // 尝试重连5次
		if err := b.connect(); err == nil {
			b.log("Reconnected successfully")
			b.startGoroutines()
			return
		}
		time.Sleep(5 * time.Second)
	}

	b.log("Failed to reconnect after 5 attempts")
	// 如果5次都失败，可以选择再次尝试或者直接返回上层处理
}

func (b *BinanceStream) handleErr(err error) {
	if b.cfg.ErrorHandler != nil {
		b.cfg.ErrorHandler(err)
	}
	b.log("Error: " + err.Error())
}

func (b *BinanceStream) log(msg string) {
	if b.cfg.Logger != nil {
		b.cfg.Logger.Info(msg)
	}
}
