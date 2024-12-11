package stream

// MessageHandler 是消息处理函数
type MessageHandler func(msg []byte)

// ErrorHandler 是错误处理函数
type ErrorHandler func(err error)

// CloseHandler 是连接关闭处理函数
type CloseHandler func()

// Stream 表示一个连接到特定交易所的单一数据流的WebSocket客户端接口。
// 不包含对心跳或ping/pong的固定接口要求。每个实现类可在Connect中自行启动心跳线程、
// 定期发送ping，或根据该交易所要求处理pong等细节。
type Stream interface {
    // Connect 建立WebSocket连接并开始接收该数据流的数据。
    // 一旦连接成功，如果交易所会立即推送数据，则会通过MessageHandler传递出去。
    Connect() error

    // Disconnect 关闭连接并停止接收数据。
    Disconnect() error

    // 回调函数应在Connect之前设置好，以确保在连接成功后收到的首条消息能被正确处理。
    AddMessageHandler(h MessageHandler)
    AddErrorHandler(h ErrorHandler)
    AddCloseHandler(h CloseHandler)
}