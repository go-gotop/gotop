package okx

import (
	"time"

	"github.com/go-gotop/gotop/ratelimiter"
)

type OkxRateLimitConfig struct {
	TimesRules map[string]ratelimiter.RateLimitRule // key前缀 -> 规则
}

func DefaultOkxConfig() OkxRateLimitConfig {
	return OkxRateLimitConfig{
		TimesRules: map[string]ratelimiter.RateLimitRule{
			OkxCreateOrder2sKey: {Window: 2 * time.Second, Threshold: 60},
			OkxCancelOrder2sKey: {Window: 2 * time.Second, Threshold: 60},
		},
	}
}
