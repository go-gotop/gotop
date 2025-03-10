package stream

import (
	"context"

	"github.com/go-gotop/gotop/types"
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
	Type() types.StreamType
}
