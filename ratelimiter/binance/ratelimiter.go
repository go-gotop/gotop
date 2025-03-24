package binance

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-gotop/gotop/ratelimiter"
	extractkey "github.com/go-gotop/gotop/ratelimiter/binance/extractKey"
	"github.com/go-gotop/gotop/types"
	"github.com/redis/go-redis/v9"
)

// ============================ OrderRateLimiter 下单限流器 ============================
type GeneralRateLimiter struct {
	extractRuleKey   *extractkey.RuleKey
	extractWeightKey *extractkey.WeightKey
	extractRedisKey  *extractkey.RedisKey
	timesAlgorithm   *TimesAlgorithm
}

func NewGeneralRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest] {
	return &GeneralRateLimiter{
		extractRuleKey:   &extractkey.RuleKey{},
		extractWeightKey: &extractkey.WeightKey{},
		extractRedisKey:  &extractkey.RedisKey{},
		timesAlgorithm:   NewTimesAlgorithm(redisClient),
	}
}

func (r *GeneralRateLimiter) Check(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	// 提取规则键
	ruleKeys := r.extractRuleKey.ExtractKeys(request)
	weightKeys := r.extractWeightKey.ExtractKeys(request)
	redisKeys := r.extractRedisKey.ExtractKeys(request)

	// 匹配规则
	matchRules := []ratelimiter.RateLimitRule{}

	for _, key := range ruleKeys {
		timesRules, err := getRules(key, DefaultBinanceConfig().TimesRules)
		weightRules, err := getRules(key, DefaultBinanceConfig().WeightRules)
		if err != nil {
			return ratelimiter.RateLimitDecision{
				Allowed: false,
				Reason:  err.Error(),
			}, err
		}
		matchRules = append(matchRules, timesRules...)
		matchRules = append(matchRules, weightRules...)
	}

	weight
	// 获取redis key
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

// ============================ IPRateLimiter 权重限流器(权重粒度只到IP) ============================
type IPRateLimiter struct {
	weightAlgorithm *WeightAlgorithm
}

func NewIPRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest] {
	return &IPRateLimiter{
		weightAlgorithm: NewWeightAlgorithm(redisClient),
	}
}

func (r *IPRateLimiter) Check(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	rules, err := getRules(r.extractRuleKey(request))
	if err != nil {
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  err.Error(),
		}, err
	}

	redisKey := r.extractRedisKey(request)
	weightKey := r.extractWeightKey(request)
	weight := DefaultBinanceConfig().Weight[weightKey]
	decision, err := r.weightAlgorithm.Check(redisKey, weight, rules)
	if err != nil {
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  err.Error(),
		}, err
	}
	return decision, nil
}

func (r *IPRateLimiter) extractRuleKey(request ratelimiter.ExchangeRateLimiterRequest) string {

	marketType := ""
	switch request.MarketType {
	case types.MarketTypeSpot, types.MarketTypeMargin:
		marketType = "spot"
	case types.MarketTypeFuturesUSDMargined,
		types.MarketTypeFuturesCoinMargined,
		types.MarketTypePerpetualUSDMargined,
		types.MarketTypePerpetualCoinMargined:
		marketType = "futures"
	default:
		return ""
	}
	return fmt.Sprintf("binance:%s:request", marketType)
}

func (r *IPRateLimiter) extractRedisKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	ip := request.IP
	if ip == "" {
		_ip := os.Getenv("HOST_IP")
		if _ip == "" {
			_ip = "unknown"
		}
		ip = _ip
	}

	marketType := ""
	switch request.MarketType {
	case types.MarketTypeSpot, types.MarketTypeMargin:
		marketType = "spot"
	case types.MarketTypeFuturesUSDMargined,
		types.MarketTypeFuturesCoinMargined,
		types.MarketTypePerpetualUSDMargined,
		types.MarketTypePerpetualCoinMargined:
		marketType = "futures"
	default:
		return ""
	}

	return fmt.Sprintf("binance:%s:request:%s", marketType, ip)
}

func (r *IPRateLimiter) extractWeightKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	marketType := ""
	switch request.MarketType {
	case types.MarketTypeSpot, types.MarketTypeMargin:
		marketType = "spot"
	case types.MarketTypeFuturesUSDMargined,
		types.MarketTypeFuturesCoinMargined,
		types.MarketTypePerpetualUSDMargined,
		types.MarketTypePerpetualCoinMargined:
		marketType = "futures"
	default:
		return ""
	}
	return fmt.Sprintf("binance:%s:%s:weight", marketType, request.RequestType)
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
