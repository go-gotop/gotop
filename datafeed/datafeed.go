package datafeed

import (
	"context"

	"github.com/go-gotop/gotop/stream"
)

// DataFeed 接口使用两个泛型参数，分别用于不同类型的数据流请求，比如交易和订单。
// 例如，对于Binance的实现，可以定义：
// type BinanceTradeRequest struct {...}
// type BinanceOrderRequest struct {...}
// 然后实现：DataFeed[BinanceTradeRequest, BinanceOrderRequest]
type DataFeed[TradeRequest any, OrderRequest any, Request any] interface {
    // Name 返回DataFeed的名称, 例如"BINANCE"
    Name() string

    // TradeStream 订阅交易数据
    // id: 调用方在订阅前就给定的ID，用来唯一标识该订阅。
    // request: 交易数据的订阅请求，类型为TradeRequest，在不同交易所实现中具有不同的字段。
    TradeStream(ctx context.Context, id string,request TradeRequest) error

    // OrderStream 订阅订单数据
    // 同理，id和request由调用方提供，request为OrderRequest类型，各实现可定制。
    OrderStream(ctx context.Context, id string, request OrderRequest) error

    // Streams 返回当前所有订阅的id列表
    Streams() map[string]stream.Stream[Request]

    // CloseStream 关闭单个订阅
    CloseStream(id string) error

    // Close 关闭所有订阅
    Close() error
}
