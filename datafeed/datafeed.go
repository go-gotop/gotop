package datafeed

import (
	"context"

	"github.com/go-gotop/gotop/types"
)

// TradeHandler 处理交易数据
type TradeHandler func(trade types.TradeEvent)

// ErrorHandler 处理错误
type ErrorHandler func(err error)

// TradeStream is a stream of trades.
type TradeStream interface {
	// Stream 开始流式处理
	Stream(ctx context.Context, handler TradeHandler, errHandler ErrorHandler) error
	// Close 关闭流
	Close() error
}
