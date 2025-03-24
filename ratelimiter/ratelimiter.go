package ratelimiter

import (
	"context"
	"time"

	"github.com/go-gotop/gotop/types"
)

// RequestType 表示请求类型，例如 HTTP 或 WebSocket
type RequestType string

const (
	// RequestTypeOrder 下订单请求
	RequestTypeOrder RequestType = "createorder"
	// RequestTypeNormal 普通请求
	RequestTypeNormal RequestType = "request"
)

// ExchangeRateLimiterRequest 表示交易所限流器请求的上下文
type ExchangeRateLimiterRequest struct {
	// 交易所类型
	Exchange string
	// IP
	IP string
	// 用户ID
	AccountID string
	// 市场类型
	MarketType types.MarketType
	// 请求类型
	RequestType RequestType
}

// RateLimitRule 表示限流规则
// 通用限流规则结构体，用于描述限流规则的窗口大小和阈值
type RateLimitRule struct {
	// 窗口大小
	Window time.Duration
	// 阈值
	Threshold int
}

// RateLimitDecision 表示对一次请求的限流判断结果
// 如果允许，则 RetryAfter 为 0，Reason 为空
// 如果不允许，则 RetryAfter 为需要等待的时间，Reason 为拒绝原因
type RateLimitDecision struct {
	Allowed    bool
	RetryAfter time.Duration
	Reason     string
}

// RateLimiter 泛型上下文 T 是一次请求的上下文类型，
// 不限制其具体结构，由上层业务自行定义。
// 对于交易所，可以是一个包含交易所类型、请求类型、IP、UserID、Endpoint等信息的结构体。
type RateLimiter[T any] interface {
	// Check 根据请求上下文决定是否允许通过
	Check(ctx context.Context, request T) (RateLimitDecision, error)
}

// KeyExtractor 接口用于从请求类型 T 中提取关键信息（键）用于限流算法的判定。
// K 是键的类型（如字符串、元组或自定义可比较类型）。
// 对于交易所请求，可以将 (Exchange, RequestType, IP, UserID, Endpoint) 等字段拼接成键或键组。
type KeyExtractor[T any, K comparable] interface {
	ExtractKeys(request T) []K
}

// RateLimitManager 是对上层业务的统一抽象接口，
// 外部系统只需调用它即可完成对请求的限流检查和记录，而内部的逻辑由 RateLimiterProvider、RateLimiter 等组成。
type RateLimitManager[T any] interface {
	// PreCheck 在请求发送前进行检查，返回是否允许发送
	PreCheck(ctx context.Context, request T) (RateLimitDecision, error)
}
