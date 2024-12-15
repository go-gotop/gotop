package stream

import (
	"context"
)

// StreamType 用于区分不同的订阅数据流类型，例如 trade、order、balance、ticker、kline、depth 等。
// 每个 StreamType 表示数据流中传递的主要信息类别。
type StreamType string

const (
	// StreamTypeTrade 交易数据流（成交信息）
	StreamTypeTrade StreamType = "TRADE"

	// StreamTypeOrder 订单数据流（用户订单更新）
	StreamTypeOrder StreamType = "ORDER"

	// StreamTypeBalance 余额数据流（用户账户余额变动）
	StreamTypeBalance StreamType = "BALANCE"

	// StreamTypeTicker 行情Ticker数据流（最新成交价、24小时涨跌等汇总信息）
	StreamTypeTicker StreamType = "TICKER"

	// StreamTypeKline K线数据流（周期价格数据，比如1m、5m、1h k线）
	StreamTypeKline StreamType = "KLINE"

	// StreamTypeDepth 市场深度数据流（订单簿更新）
	StreamTypeDepth StreamType = "DEPTH"

	// StreamTypeBookTicker 最优盘口数据流（最优买一、卖一报价更新）
	StreamTypeBookTicker StreamType = "BOOK_TICKER"

	// StreamTypeMarkPrice 标记价格数据流（期货、永续合约的标记价格）
	StreamTypeMarkPrice StreamType = "MARK_PRICE"

	// 如果有更多需要，可继续补充
)

// Stream 表示一个连接到特定交易所的单一数据流的WebSocket客户端接口。
// 不包含对心跳或ping/pong的固定接口要求。每个实现类可在Connect中自行启动心跳线程、
// 定期发送ping，或根据该交易所要求处理pong等细节, 且需要自行实现重连。
// 例如, Binance的连接是24小时之后就断开，需要自行实现重连。
type Stream[T any] interface {
    // Connect 建立连接并开始数据流接收。
    // cfg: 配置参数（任意类型），具体实现在接收cfg后根据交易所要求建立连接。
    Connect(ctx context.Context, cfg T) error

    // ID 返回该Stream的唯一ID，这个ID是在Connect时由调用方传入并存储下来的。
    ID() string

    // Disconnect 关闭连接并停止接收数据。
    Disconnect() error

    // Type 返回该Stream的类型
    Type() StreamType
}
