package stream

// Stream 表示一个连接到特定交易所的单一数据流的WebSocket客户端接口。
// 不包含对心跳或ping/pong的固定接口要求。每个实现类可在Connect中自行启动心跳线程、
// 定期发送ping，或根据该交易所要求处理pong等细节, 且需要自行实现重连。
// 例如, Binance的连接是24小时之后就断开，需要自行实现重连。
type Stream[T any] interface {
    // ID 返回该Stream的唯一ID
    // 便于StreamManager管理多个Stream
    ID() string
    // Connect 建立WebSocket连接并开始接收该数据流的数据。
    // 一旦连接成功，如果交易所会立即推送数据，则会通过T类型中的Handler传递出去。
    Connect(T) error

    // Disconnect 关闭连接并停止接收数据。
    Disconnect() error
}
