package okx

import (
	"context"
	"strings"

	"github.com/go-gotop/gotop/ratelimiter"
	extractkey "github.com/go-gotop/gotop/ratelimiter/okx/extractKey"
	"github.com/redis/go-redis/v9"
)

type GeneralRateLimiter struct {
	extractTimesRuleKey  *extractkey.TimesRuleKey
	extractTimesRedisKey *extractkey.TimesRedisKey
	timesAlgorithm       *ratelimiter.TimesAlgorithm
}

func NewGeneralRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest] {
	return &GeneralRateLimiter{
		extractTimesRuleKey:  &extractkey.TimesRuleKey{},
		extractTimesRedisKey: &extractkey.TimesRedisKey{},
		timesAlgorithm:       ratelimiter.NewTimesAlgorithm(redisClient),
	}
}

func (r *GeneralRateLimiter) Check(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	// 1. 提取次数规则键
	timesRuleKey := r.extractTimesRuleKey.ExtractKeys(request)
	// 2. 提取次数redis key
	timesRedisKey := r.extractTimesRedisKey.ExtractKeys(request)
	// 3. 提取次数规则
	matchTimesRules := []ratelimiter.RateLimitRule{}
	if timesRuleKey != "" {
		timesRules, err := getRules(timesRuleKey, DefaultOkxConfig().TimesRules)
		if err == nil && len(timesRules) > 0 {
			matchTimesRules = timesRules
		}
	}

	if timesRedisKey != "" && len(matchTimesRules) > 0 {
		decision, err := r.timesAlgorithm.Check(timesRedisKey, matchTimesRules)
		if err != nil {
			return ratelimiter.RateLimitDecision{
				Allowed: false,
				Reason:  err.Error(),
			}, err
		}
		if !decision.Allowed {
			return decision, nil
		}
	}
	return ratelimiter.RateLimitDecision{
		Allowed: true,
	}, nil
}

func getRules(key string, defaultRules map[string]ratelimiter.RateLimitRule) ([]ratelimiter.RateLimitRule, error) {
	rules := []ratelimiter.RateLimitRule{}

	// 1. 首先尝试精确匹配
	if rule, exists := defaultRules[key]; exists {
		rules = append(rules, rule)
		return rules, nil
	}

	// 2. 提取基础键
	baseKey := key
	// 去掉账户ID部分（如果有）
	parts := strings.Split(key, ":")
	if len(parts) > 3 {
		baseKey = strings.Join(parts[:3], ":")
	}

	// 3. 针对特定规则类型匹配
	for k, v := range defaultRules {
		// 规则键以基础键开头
		if strings.HasPrefix(k, baseKey) {
			rules = append(rules, v)
		}
	}

	return rules, nil
}
