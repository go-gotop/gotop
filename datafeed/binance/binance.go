package binance

import (
	"strings"
	"context"
	"time"
	"sync"
	"fmt"

	"github.com/go-gotop/gotop/types"
	"github.com/go-gotop/gotop/stream"
	binanceStream"github.com/go-gotop/gotop/stream/binance"
)

const (
	spotHTTPURL = "https://api.binance.com"
	spotWSURL   = "wss://stream.binance.com:9443/ws"

	futuresHTTPURL = "https://fapi.binance.com"
	futuresWSURL   = "wss://fstream.binance.com/ws"
)

// NewBinanceDataFeed 创建一个新的BinanceDataFeed
func NewBinanceDataFeed(opts ...Option) *BinanceDataFeed {
	o := applyOptions(opts...)
	return &BinanceDataFeed{
		opts: o,
		streams: make(map[string]stream.Stream[binanceStream.BinanceRequest]),
	}
}

// BinanceTradeRequest 是Binance的交易数据订阅请求
type BinanceTradeRequest struct {
	// Symbol 交易对，例如"BTCUSDT"
	Symbol string
	// Market 市场类型，例如"SPOT"或"FUTURES"
	Market types.MarketType
	// Handler 数据处理函数
	Handler func(data []byte)
	// ErrorHandler 错误处理函数
	ErrorHandler func(err error)
}

// BinanceOrderRequest 是Binance的订单数据订阅请求
type BinanceOrderRequest struct {
}

// listenKey 监听键
type listenKey struct {
    Key        string
    ExpireTime time.Time
}

// accountInfo 账户信息
type accountInfo struct {
	AccountID   string
    APIKey      string
    SecretKey   string
    // 不同类型账户的listenKey
    ListenKeys  map[types.MarketType]*listenKey
}

// BinanceDataFeed 是Binance的数据订阅器
type BinanceDataFeed struct {
	mu sync.Mutex
	// opts 配置选项
	opts *options
	// listenKeys 监听键
	listenKeys map[string]listenKey
	// streams 数据流
	streams map[string]stream.Stream[binanceStream.BinanceRequest]
}

// Name 返回DataFeed的名称, BINANCE
func (b *BinanceDataFeed) Name() string {
	return types.BinanceExchange
}

// TradeStream 订阅交易数据
// id: 调用方在订阅前就给定的ID，用来唯一标识该订阅。
// request: 交易数据的订阅请求，类型为BinanceTradeRequest。
func (b *BinanceDataFeed) TradeStream(ctx context.Context, id string, request BinanceTradeRequest) error {
	var url string
	switch request.Market {
	case types.MarketTypeSpot:
		url = fmt.Sprintf("/%s/%s@trade", spotWSURL, strings.ToLower(request.Symbol))
	case types.MarketTypeFuturesUSDMargined:
		url = fmt.Sprintf("/%s/%s@aggTrade", futuresWSURL, strings.ToLower(request.Symbol))
	default:
		return fmt.Errorf("invalid market type: %v", request.Market)
	}

	s := binanceStream.NewBinanceStream(id, stream.StreamTypeTrade)

	if err := s.Connect(ctx, binanceStream.BinanceRequest{
		URL: url,
		Handler: request.Handler,
		ErrorHandler: request.ErrorHandler,
	}); err != nil {
		return err
	}

	b.mu.Lock()
	b.streams[id] = s
	b.mu.Unlock()

	return nil
}

// OrderStream 订阅订单数据
// id: 调用方在订阅前就给定的ID，用来唯一标识该订阅。
// request: 订单数据的订阅请求，类型为BinanceOrderRequest。
func (b *BinanceDataFeed) OrderStream(ctx context.Context, id string, request BinanceOrderRequest) error {
	return nil
}

// Streams 返回当前所有订阅的id列表
func (b *BinanceDataFeed) Streams() map[string]stream.Stream[binanceStream.BinanceRequest] {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.streams
}

// CloseStream 关闭单个订阅
func (b *BinanceDataFeed) CloseStream(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if stream, ok := b.streams[id]; ok {
		if err := stream.Disconnect(); err != nil {
			return err
		}
		delete(b.streams, id)
	}
	return nil
}

// Close 关闭所有订阅
func (b *BinanceDataFeed) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, stream := range b.streams {
		if err := stream.Disconnect(); err != nil {
			return err
		}
	}
	return nil
}
