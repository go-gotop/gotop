package ratelimiter

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

// RateLimiterProvider 用于根据请求上下文 T 动态返回需要检查的限流器列表。
// 不同的交易所、请求类型、用户、IP等维度可能需要叠加多个限流规则（多个 RateLimiter）。
type RateLimiterProvider[T any] interface {
	// GetRateLimiters 获取限流器列表
	GetRateLimiters(request T) ([]RateLimiter[T], error)
}

// RateLimitManager 是对上层业务的统一抽象接口，
// 外部系统只需调用它即可完成对请求的限流检查和记录，而内部的逻辑由 RateLimiterProvider、RateLimiter 等组成。
type RateLimitManager[T any] interface {
	// PreCheck 在请求发送前进行检查，返回是否允许发送
	PreCheck(ctx context.Context, request T) (RateLimitDecision, error)

	// PostRecord 在请求完成后进行记录，更新限流状态
	PostRecord(ctx context.Context, request T) error
}

// GenericRateLimiter 是一个通用限流器实现的抽象，
// 它通过 KeyExtractor 从请求中提取键列表，
// 然后对每个键调用底层的 RateLimitAlgorithm 进行判定与记录。
type GenericRateLimiter[T any, K comparable] struct {
	Extractor KeyExtractor[T, K]
	Algorithm RateLimitAlgorithm[K]
}

// Check 实现 RateLimiter[T] 接口逻辑：
// 1. 从请求中抽取出键列表
// 2. 对列表中的每个键调用 Algorithm.Check(key)
// 3. 如果任一键不允许，则本次请求不允许通过，合并相关信息到最终的Decision中
// 4. 如果所有键都允许，则 Allowed=true
func (g *GenericRateLimiter[T, K]) Check(ctx context.Context, request T) (RateLimitDecision, error) {
	keys := g.Extractor.ExtractKeys(request)
	if len(keys) == 0 {
		// 如果连键都提不出来，那通常意味着不需要限流规则，可直接通过
		return RateLimitDecision{Allowed: true}, nil
	}

	finalDecision := RateLimitDecision{
		Allowed:    true,
		RetryAfter: 0,
		Reason:     "",
	}

	for _, key := range keys {
		decision, err := g.Algorithm.Check(key)
		if err != nil {
			// 算法调用出错，可能是数据存取错误等，实际中可加入重试、日志等逻辑
			return RateLimitDecision{}, err
		}

		if !decision.Allowed {
			// 一旦有一个键不允许，这个请求就不允许。
			finalDecision.Allowed = false
			// 合并Reason和RetryAfter（可根据需求决定合并策略）
			if decision.Reason != "" {
				finalDecision.Reason = decision.Reason
			}
			// 取较大的RetryAfter保证客户端等待足够久
			if decision.RetryAfter > finalDecision.RetryAfter {
				finalDecision.RetryAfter = decision.RetryAfter
			}
			// 如果策略是只要有一个键不允许就立刻返回，可以直接return
			// 这里假设需要合并所有键的信息（通常没必要）
			// return finalDecision, nil
		}
	}

	return finalDecision, nil
}

// Record 实现 RateLimiter[T] 接口逻辑：
// 1. 从请求中抽取出键列表
// 2. 对每个键调用 Algorithm.Record(key) 更新计数器或令牌状态
// 这里假设所有键都需要进行记录。
func (g *GenericRateLimiter[T, K]) Record(ctx context.Context, request T) error {
	keys := g.Extractor.ExtractKeys(request)
	if len(keys) == 0 {
		// 没有键需要记录，不做任何事情。
		return nil
	}

	for _, key := range keys {
		if err := g.Algorithm.Record(key); err != nil {
			// 如有错误，根据需求决定是否忽略还是直接返回
			return err
		}
	}
	return nil
}
