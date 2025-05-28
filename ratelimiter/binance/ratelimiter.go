package binance

import (
	"context"
	"strings"

	"github.com/go-gotop/gotop/ratelimiter"
	extractkey "github.com/go-gotop/gotop/ratelimiter/binance/extractKey"
	"github.com/redis/go-redis/v9"
)

// ============================ GeneralRateLimiter 通用限流器 ============================
type GeneralRateLimiter struct {
	extractTimesRuleKey   *extractkey.TimesRuleKey
	extractWeightRuleKey  *extractkey.WeightRuleKey
	extractTimesRedisKey  *extractkey.TimesRedisKey
	extractWeightRedisKey *extractkey.WeightRedisKey
	timesAlgorithm        *ratelimiter.TimesAlgorithm
	weightAlgorithm       *ratelimiter.WeightAlgorithm
}

func NewGeneralRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest] {
	return &GeneralRateLimiter{
		extractTimesRuleKey:   &extractkey.TimesRuleKey{},
		extractWeightRuleKey:  &extractkey.WeightRuleKey{},
		extractTimesRedisKey:  &extractkey.TimesRedisKey{},
		extractWeightRedisKey: &extractkey.WeightRedisKey{},
		timesAlgorithm:        ratelimiter.NewTimesAlgorithm(redisClient),
		weightAlgorithm:       ratelimiter.NewWeightAlgorithm(redisClient),
	}
}

func (r *GeneralRateLimiter) Check(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	// 1. 提取次数规则键
	timesRuleKey := r.extractTimesRuleKey.ExtractKeys(request)
	// 2. 提取权重规则键
	weightRuleKey := r.extractWeightRuleKey.ExtractKeys(request)
	// 3. 提取次数redis key
	timesRedisKey := r.extractTimesRedisKey.ExtractKeys(request)
	// 4. 提取权重redis key
	weightRedisKey := r.extractWeightRedisKey.ExtractKeys(request)
	// 5. 提取权重键
	weightKey := r.extractWeightRuleKey.ExtractKeys(request)

	// 匹配规则
	matchTimesRules := []ratelimiter.RateLimitRule{}
	matchWeightRules := []ratelimiter.RateLimitRule{}

	if timesRuleKey != "" {
		timesRules, err := getRules(timesRuleKey, DefaultBinanceConfig().TimesRules)
		if err == nil && len(timesRules) > 0 {
			matchTimesRules = timesRules
		}
	}

	if weightRuleKey != "" {
		weightRules, err := getRules(weightRuleKey, DefaultBinanceConfig().WeightRules)
		if err == nil && len(weightRules) > 0 {
			matchWeightRules = weightRules
		}
	}

	weight := 0
	if weightKey != "" {
		weight = DefaultBinanceConfig().Weight[weightKey]
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

	if weightRedisKey != "" && len(matchWeightRules) > 0 {
		decision, err := r.weightAlgorithm.Check(weightRedisKey, weight, matchWeightRules)
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
