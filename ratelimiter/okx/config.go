package okx

import (
	"time"

	"github.com/go-gotop/gotop/ratelimiter"
)

// OkxRateLimitConfig 定义了Okx的限流配置
type OkxRateLimitConfig struct {
	Rules map[string]ratelimiter.RateLimitRule // key前缀 -> 规则
}

// DefaultOkxConfig 返回Okx特定的限流配置（实际可从配置文件、环境变量中加载）
func DefaultOkxConfig() OkxRateLimitConfig {
	return OkxRateLimitConfig{
		Rules: map[string]ratelimiter.RateLimitRule{
			OkxCreateOrder2sKey: {Window: 2 * time.Second, Threshold: 60},
		},
	}
}
