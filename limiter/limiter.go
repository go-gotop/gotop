package limiter

import (
	"context"
	"time"
)

// RequestType 表示请求类型，例如 HTTP 或 WebSocket
type RequestType string

const (
	RequestTypeHTTP RequestType = "http"
	RequestTypeWS   RequestType = "ws"
)

// RateLimitDecision 表示对一次请求的限流判断结果
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
	// Record 在请求成功或完成后，更新限流状态
	Record(ctx context.Context, request T) error
}

// KeyExtractor 接口用于从请求类型 T 中提取关键信息（键）用于限流算法的判定。
// K 是键的类型（如字符串、元组或自定义可比较类型）。
// 对于交易所请求，可以将 (Exchange, RequestType, IP, UserID, Endpoint) 等字段拼接成键或键组。
type KeyExtractor[T any, K comparable] interface {
	ExtractKeys(request T) []K
}

// RateLimitAlgorithm 是底层的限流算法接口，不关心业务上下文，
// 只针对键类型 K 进行限流判断和记录。
// 例如可以对 K 调用内部数据结构（计数器、令牌桶）以判断是否允许请求。
type RateLimitAlgorithm[K comparable] interface {
	Check(key K) (RateLimitDecision, error)
	Record(key K) error
}

// GenericRateLimiter 是一个通用限流器实现的抽象，
// 它通过 KeyExtractor 从请求中提取键列表，
// 然后对每个键调用底层的 RateLimitAlgorithm 进行判定与记录。
type GenericRateLimiter[T any, K comparable] struct {
	Extractor KeyExtractor[T, K]
	Algorithm RateLimitAlgorithm[K]
}

// Check 实现 RateLimiter[T] 接口
func (g *GenericRateLimiter[T, K]) Check(ctx context.Context, request T) (RateLimitDecision, error) {
	// 这里不写具体逻辑，只示意：
	// 1. Extract keys
	// 2. 对每个 key 调用 Algorithm.Check 并合并结果（如有任一不允许则整体不允许）
	return RateLimitDecision{Allowed: true}, nil
}

// Record 实现 RateLimiter[T] 接口
func (g *GenericRateLimiter[T, K]) Record(ctx context.Context, request T) error {
	// 同理：对 ExtractKeys 后的每个 key 调用 Algorithm.Record
	return nil
}

// RateLimiterProvider 用于根据请求上下文 T 动态返回需要检查的限流器列表。
// 不同的交易所、请求类型、用户、IP等维度可能需要叠加多个限流规则（多个 RateLimiter）。
type RateLimiterProvider[T any] interface {
	GetRateLimiters(request T) ([]RateLimiter[T], error)
}

// RateLimitManager 是对上层业务的统一抽象接口，
// 外部系统只需调用它即可完成对请求的限流检查和记录，而内部的逻辑由 RateLimiterProvider、RateLimiter 等组成。
type RateLimitManager[T any] interface {
	PreCheck(ctx context.Context, request T) (RateLimitDecision, error)
	PostRecord(ctx context.Context, request T) error
}

// GenericRateLimitManager 是一个通用的管理器实现示意（无具体逻辑）。
// 它使用一个 RateLimiterProvider 来获得对应的限流器列表，然后在 PreCheck 和 PostRecord 中统一调用。
type GenericRateLimitManager[T any] struct {
	Provider RateLimiterProvider[T]
}

// PreCheck 实现 RateLimitManager[T] 接口
func (m *GenericRateLimitManager[T]) PreCheck(ctx context.Context, request T) (RateLimitDecision, error) {
	// 1. 从 Provider 获取需要检查的限流器列表
	// 2. 对每个限流器调用 Check，并合并结果
	return RateLimitDecision{Allowed: true}, nil
}

// PostRecord 实现 RateLimitManager[T] 接口
func (m *GenericRateLimitManager[T]) PostRecord(ctx context.Context, request T) error {
	// 1. 从 Provider 获取对应限流器列表
	// 2. 对每个限流器调用 Record 进行状态更新
	return nil
}