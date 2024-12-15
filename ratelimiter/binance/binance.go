package binance

import (
    "time"

    "github.com/go-gotop/gotop/ratelimiter"	
)

// BinanceAlgorithm 是一个针对binance键的限流算法
type BinanceAlgorithm struct {
    // 在真实环境中这里会有一个 redisClient 等，用于存取计数
    // 这里用内存map模拟（非线程安全示例）
    counters map[string]int
    expiries map[string]time.Time
}

// NewBinanceAlgorithm 创建一个示例算法实例
func NewBinanceAlgorithm() *BinanceAlgorithm {
    return &BinanceAlgorithm{
        counters: make(map[string]int),
        expiries: make(map[string]time.Time),
    }
}

func (b *BinanceAlgorithm) Check(key string) (ratelimiter.RateLimitDecision, error) {
    // 判定规则示例：
    // global key: "binance:global:http"
    //   - 时间窗口：1秒
    //   - 上限：10次
    // api key key: "binance:api_key:xxx"
    //   - 时间窗口：1分钟
    //   - 上限：1200次

    allowed := true
    retryAfter := time.Duration(0)
    reason := ""

    now := time.Now()

    var window time.Duration
    var limit int

    if key == "binance:global:http" {
        window = time.Second
        limit = 10
    } else if len(key) > len("binance:api_key:") && key[:len("binance:api_key:")] == "binance:api_key:" {
        window = time.Minute
        limit = 1200
    } else {
        // 不认识的key不处理，直接允许
        return ratelimiter.RateLimitDecision{Allowed: true}, nil
    }

    // 清理过期
    if exp, ok := b.expiries[key]; ok && now.After(exp) {
        delete(b.counters, key)
        delete(b.expiries, key)
    }

    count := b.counters[key]

    // 判断是否超标
    if count >= limit {
        allowed = false
        retryAfter = window // 简单示例，让调用方等一个窗口再试
        reason = "rate limit exceeded"
    }

    return ratelimiter.RateLimitDecision{
        Allowed:    allowed,
        RetryAfter: retryAfter,
        Reason:     reason,
    }, nil
}

func (b *BinanceAlgorithm) Record(key string) error {
    return nil
}
