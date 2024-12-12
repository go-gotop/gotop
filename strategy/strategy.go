package strategy

import (
	"context"

	"github.com/go-gotop/gotop/types"
)

// Strategy 策略接口
type Strategy[Event any, NotifyData any, Result any] interface {
	// SetStatus 设置策略状态
	SetStatus(status types.StrategyStatus)
	// Status 返回策略状态
	Status() types.StrategyStatus
	// Notify 通知策略信号变化
	Notify(ctx context.Context, notifyData NotifyData)
	// Next 处理事件, 返回处理结果
	Next(ctx context.Context, event Event) Result
	// ID 策略ID
	ID() string
}
