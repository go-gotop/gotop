package okx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"

	"github.com/go-gotop/gotop/types"
)

type OkxRequest struct {
	// WebSocket请求的URL
	URL string
	// Logger 可选的日志记录器，用于调试
	Logger *slog.Logger
	// ConnectedHandler 连接成功处理函数(okx订阅频道是通过连接建立之后，发送消息进行订阅的，所以需要通过这个函数进行订阅)
	ConnectedHandler func(conn *websocket.Conn)
	// Handler 数据处理函数
	Handler func(data []byte)
	// ErrorHandler 错误处理函数
	ErrorHandler func(err error)
}

// dialFunc定义用于方便在测试时mock连接的逻辑
type dialFunc func(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)

// OkxStream 是Okx Stream的核心结构体
type OkxStream struct {
	mu     sync.Mutex
	id     string
	st     types.StreamType
	cfg    OkxRequest
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

// NewOkxStream 创建一个新的OkxStream
func NewOkxStream(id string, st types.StreamType, opts ...Option) *OkxStream {
	b := &OkxStream{
		id:                id,
		st:                st,
		dialer:            defaultDialer,
		pingInterval:      10 * time.Second,
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

// Connect 连接到Okx Stream
func (b *OkxStream) Connect(ctx context.Context, cfg OkxRequest) error {
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
func (b *OkxStream) ID() string {
	return b.id
}

// Type 返回当前连接的类型
func (b *OkxStream) Type() types.StreamType {
	return b.st
}

// Disconnect 断开当前连接
func (b *OkxStream) Disconnect() error {
	b.mu.Lock()
	// 取消当前上下文
	if b.cancel != nil {
		b.cancel()
	}
	// 关闭当前连接
	if b.conn != nil {
		_ = b.conn.Close()
	}
	b.mu.Unlock()

	// 等待所有goroutine结束
	b.wg.Wait()
	b.log("Disconnected")
	return nil
}

// connect 函数用于建立WebSocket连接
// 主要功能:
// 1. 使用dialer建立WebSocket连接
// 2. 设置连接的读取超时时间
// 3. 设置pong消息的处理函数
func (b *OkxStream) connect() error {
	// 使用dialer建立WebSocket连接
	c, _, err := b.dialer(b.cfg.URL, nil)
	if err != nil {
		// 如果连接失败,调用错误处理函数并返回错误
		b.handleErr(err)
		return err
	}

	// 保存连接实例
	b.conn = c
	// 设置读取超时时间为pongWait
	b.conn.SetReadDeadline(time.Now().Add(b.pongWait))

	// 调用连接成功处理函数
	if b.cfg.ConnectedHandler != nil {
		b.cfg.ConnectedHandler(b.conn)
	}

	return nil
}

func (b *OkxStream) startGoroutines() {
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

func (b *OkxStream) readLoop() {
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

			// 检查是否是 pong 消息
			if string(msg) == "pong" {
				// 如果是 pong 消息，更新读取超时时间
				conn.SetReadDeadline(time.Now().Add(b.pongWait))
				continue // 不需要传递给 Handler
			}

			j, err := simplejson.NewJson(msg)
			if err != nil {
				b.handleErr(err)
				continue
			}

			// 获取event
			event := j.Get("event").MustString()
			if event == "error" {
				b.handleErr(errors.New(j.Get("msg").MustString()))
				continue
			}

			if event != "" {
				continue
			}

			if b.cfg.Handler != nil {
				b.cfg.Handler(msg)
			}
		}
	}
}

func (b *OkxStream) pingLoop() {
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
			if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				b.handleErr(err)
				// 异步重连
				go b.attemptReconnect()
				return
			}
		}
	}
}

func (b *OkxStream) autoReconnectLoop() {
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

func (b *OkxStream) attemptReconnect() {
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

	b.log("Reconnecting...")

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

func (b *OkxStream) handleErr(err error) {
	if b.cfg.ErrorHandler != nil {
		b.cfg.ErrorHandler(err)
	}
	b.log("Error: " + err.Error())
}

func (b *OkxStream) log(msg string) {
	if b.cfg.Logger != nil {
		b.cfg.Logger.Info(msg)
	}
}
