package datafeed

import (
	"context"
)

// DataFeed is a stream of data.
type DataFeed[T any] interface {
	// Name 返回DataFeed的名称, 例如"BINANCE"
	Name() string
	// ListStream 返回所有订阅id
	ListStream() []string
	// TradeStream 订阅交易数据
	TradeStream(ctx context.Context, request T) error
	// OrderStream 订阅订单数据
	OrderStream(ctx context.Context, request T) error
	// CloseStream 关闭单个订阅
	CloseStream(id string) error
	// Close 关闭所有订阅
	Close() error
}
