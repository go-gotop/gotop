package binance

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// BinanceRequest Binance Stream的配置参数
type BinanceRequest struct {
	// 交易对
	Symbol string
	// WebSocket请求的URL
	URL string
	// API Key（如果需要）
	APIKey string
	// Secret Key（如果需要）
	SecretKey string
}

// NewBinanceStream 创建新的BinanceStream实例
func NewBinanceStream(id string) *BinanceStream {
	return &BinanceStream{
		id:      id,
		closeCh: make(chan struct{}),
		doneCh:  make(chan struct{}),
	}
}

// BinanceStream 实现了 stream.Stream 接口
type BinanceStream struct {
	id          string
	conn        *websocket.Conn
	request     *BinanceRequest
	closeCh     chan struct{}
	doneCh      chan struct{}
	closeOnce   sync.Once
	doneOnce    sync.Once
	connectTime time.Time
}

// 实现 Stream[BinanceRequest] 接口
func (b *BinanceStream) ID() string {
	return b.id
}

// Connect 实现 Stream 接口，使用 StreamConfig 作为参数类型
func (b *BinanceStream) Connect(request *BinanceRequest) error {
	b.request = request
	conn, _, err := websocket.DefaultDialer.Dial(request.URL, nil)
	if err != nil {
		return err
	}
	b.conn = conn
	b.connectTime = time.Now()
	return nil
}

func (b *BinanceStream) Disconnect() error {
	// 实现断开连接逻辑
	return nil
}