package strategy

import (
	"context"

	"github.com/go-gotop/gotop/types"
)

// 提供策略的基本信息：ID 与 状态
type StrategyInfo interface {
	// ID 返回策略ID
    ID() string
    // Status 返回策略状态
    Status() types.StrategyStatus
	// SetStatus 设置策略状态
	SetStatus(status types.StrategyStatus)
}

// Notifiable 可选的通知接口：策略若需接收外部通知，可实现此接口
type Notifiable[Notification any] interface {
    // Notify 接收新的通知数据，可能影响策略内部状态。
    // 返回error以便调用方感知处理通知时出现的问题。
    Notify(ctx context.Context, notification Notification) error
}

// EventProcessor 必要的事件处理接口：策略通过Next处理传入的事件并返回结果
type EventProcessor[Event any, Result any] interface {
    // Next 处理输入事件，并返回结果和error。
    // 当error不为nil表示处理事件时出现异常，调用方可进行相应的错误处理或退出。
    Next(ctx context.Context, event Event) (Result, error)
}

// Strategy 最终的Strategy接口集成上述特性。
// 如果某些策略不需要通知功能，可定义没有Notifiable的版本。
type Strategy[Event any, Notification any, Result any] interface {
    StrategyInfo
    Notifiable[Notification]
    EventProcessor[Event, Result]
}
