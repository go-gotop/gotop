package binance

import "time"

// LimitRule 定义了限流规则
type LimitRule struct {
	// 时间窗口大小
	Window time.Duration
	// 窗口内允许的最大请求数	
	Threshold int
}

// BinanceRateLimitConfig 定义了Binance的限流配置
type BinanceRateLimitConfig struct {
    Rules map[string]LimitRule // key前缀 -> 规则
}

// DefaultBinanceConfig 返回Binance特定的限流配置（实际可从配置文件、环境变量中加载）
func DefaultBinanceConfig() BinanceRateLimitConfig {
    return BinanceRateLimitConfig{
        Rules: map[string]LimitRule{
            BinanceFuturesOrder10sKey: {Window: 10 * time.Second, Threshold: 300},
            BinanceFuturesOrder1mKey:  {Window: time.Minute, Threshold: 1200},
        },
    }
}

