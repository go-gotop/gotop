# README

## 简介

本限流器框架旨在满足多元化和可扩展的限流需求。针对不同交易所、不同请求类型（HTTP/WS）、不同用户或IP维度，以及特定请求（如某交易所的下单接口），该框架通过一组抽象的接口和通用的逻辑流程，为业务方提供统一的限流调用入口。

设计目标：

- **灵活性**：能够适应多种限流场景（全局限流、用户级限流、IP级限流、Endpoint级限流等）。
- **可扩展性**：通过清晰的接口抽象，让新增或修改限流策略变得轻松。
- **统一接口**：提供统一的 `RateLimitManager` 接口以简化上层业务调用。无论底层逻辑多复杂，上层只需调用 `PreCheck` 和 `PostRecord` 即可。

## 接口定义

以下是核心接口与数据结构的定义（在实际项目中，你可以将这些接口定义在相应的 `.go` 源文件中）：

```go
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
```

## 调用流程

```

       +------------------------+
       |      RateLimitManager |
       |      (PreCheck)       |
       +-----------+-----------+
                   |
                   v
       +------------------------+
       |   RateLimiterProvider  |
       |   GetRateLimiters(...) |
       +-----------+-----------+
                   |
                   v (list of RateLimiter[T])
       +------------------------+
       |     RateLimiter[T]     |
       |     (multiple)         |
       +-----------+-----------+
          /           |          \
         v            v           v
  +----------+   +----------+   +----------+
  | RateLimiter| | RateLimiter| | RateLimiter|
  | (Check)    | | (Check)    | | (Check)    |
  +----------+   +----------+   +----------+
     |                           ^
     | if any not allowed         |
     | stops and returns          |
     | allowed=false              |
     +----------------------------+

(If allowed)
       +------------------------+
       |   Actual Request       |
       |   sent to Exchange     |
       +-----------+-----------+
                   |
                  (request finishes)
                   |
                   v
       +------------------------+
       |  RateLimitManager      |
       |   (PostRecord)         |
       +-----------+-----------+
                   |
                   v
       +------------------------+
       |  RateLimiterProvider  |
       |  GetRateLimiters(...) |
       +-----------+-----------+
                   |
                   v (list of RateLimiter[T])
             +------------+
             | RateLimiter|
             | (Record)   |
             +-----+------+
                   |
                  ...
