package strategy

// Strategy 策略接口
type Strategy[T any, S any] interface {
	// Next 处理事件, 返回处理结果
	Next(event T) S
	// Name 策略名称
	Name() string
}
