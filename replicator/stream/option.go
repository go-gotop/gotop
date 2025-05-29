package stream

// StreamReplicatorOption 流复制器选项
type StreamReplicatorOption func(*streamReplicatorManager)

// WithDataHandler 设置数据处理函数
func WithDataHandler(handler func(data []byte)) StreamReplicatorOption {
	return func(m *streamReplicatorManager) {
		m.dataHandler = handler
	}
}

// WithErrorHandler 设置错误处理函数
func WithErrorHandler(handler func(streamID string, err error)) StreamReplicatorOption {
	return func(m *streamReplicatorManager) {
		m.errorHandler = handler
	}
}
