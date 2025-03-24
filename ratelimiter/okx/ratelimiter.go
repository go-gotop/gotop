package okx

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/redis/go-redis/v9"
)

// ============================ TimesRateLimiter 次数限流器 ============================
type TimesRateLimiter struct {
	timesAlgorithm *TimesAlgorithm
}

func NewTimesRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest] {
	return &TimesRateLimiter{
		timesAlgorithm: NewTimesAlgorithm(redisClient),
	}
}

func (r *TimesRateLimiter) Check(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	rules, err := getRules(r.extractRuleKey(request))
	if err != nil {
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  err.Error(),
		}, err
	}

	redisKey := r.extractRedisKey(request)
	decision, err := r.timesAlgorithm.Check(redisKey, rules)
	if err != nil {
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  err.Error(),
		}, err
	}
	return decision, nil
}

func (r *TimesRateLimiter) extractRuleKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	switch request.RequestType {
	case ratelimiter.RequestTypeOrder:
		return fmt.Sprintf("%s:%s", "okx", request.RequestType)
	default:
		return ""
	}
}

func (r *TimesRateLimiter) extractRedisKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	switch request.RequestType {
	case ratelimiter.RequestTypeOrder:
		return fmt.Sprintf("%s:%s:%s:%s", "okx", request.MarketType.String(), request.RequestType, request.AccountID)
	default:
		return ""
	}
}

func getRules(key string) ([]ratelimiter.RateLimitRule, error) {
	defaultRules := DefaultOkxConfig()
	rules := []ratelimiter.RateLimitRule{}

	// 检查是否开启调试模式
	debugMode := os.Getenv("DEBUG") != ""

	// 输出调试信息
	if debugMode {
		fmt.Printf("尝试获取规则，输入键: %s\n", key)
		fmt.Printf("可用规则: %v\n", defaultRules.Rules)
	}

	// 1. 首先尝试精确匹配
	if rule, exists := defaultRules.Rules[key]; exists {
		if debugMode {
			fmt.Printf("找到精确匹配规则: %s -> %+v\n", key, rule)
		}
		rules = append(rules, rule)
		return rules, nil
	}

	// 2. 提取基础键
	baseKey := key
	// 去掉账户ID部分（如果有）
	parts := strings.Split(key, ":")
	if len(parts) > 3 {
		baseKey = strings.Join(parts[:3], ":")
		if debugMode {
			fmt.Printf("提取基础键: %s\n", baseKey)
		}
	}

	// 3. 针对特定规则类型匹配
	for k, v := range defaultRules.Rules {
		// 规则键以基础键开头
		if strings.HasPrefix(k, baseKey) {
			if debugMode {
				fmt.Printf("找到匹配规则: %s -> %+v\n", k, v)
			}
			rules = append(rules, v)
		}
	}

	if debugMode {
		fmt.Printf("总共找到 %d 个规则\n", len(rules))
	}

	return rules, nil
}
