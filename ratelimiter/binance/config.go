package binance

import (
	"time"

	"github.com/go-gotop/gotop/ratelimiter"
)

// BinanceRateLimitConfig 定义了Binance的限流配置
type BinanceRateLimitConfig struct {
	Rules  map[string]ratelimiter.RateLimitRule // key前缀 -> 规则
	Weight map[string]int                       // key前缀 -> 权重
}

// DefaultBinanceConfig 返回Binance特定的限流配置（实际可从配置文件、环境变量中加载）
func DefaultBinanceConfig() BinanceRateLimitConfig {
	return BinanceRateLimitConfig{
		Rules: map[string]ratelimiter.RateLimitRule{
			BinanceFuturesCreateOrder10sKey: {Window: 10 * time.Second, Threshold: 300},
			BinanceFuturesCreateOrder1mKey:  {Window: time.Minute, Threshold: 1200},
			BinanceSpotCreateOrder10sKey:    {Window: 10 * time.Second, Threshold: 100},
			BinanceSpotRequest1mKey:         {Window: time.Minute, Threshold: 6000},
			BinanceFuturesRequest1mKey:      {Window: time.Minute, Threshold: 2400},
		},
		Weight: map[string]int{
			BinanceFuturesCreateOrderWeightKey: 0,
			BinanceSpotCreateOrderWeightKey:    1,
			BinanceSpotRequestWeightKey:        1,
			BinanceFuturesRequestWeightKey:     1,
		},
	}
}
