package binance

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/go-gotop/gotop/types"
	"github.com/redis/go-redis/v9"
)

// ============================ OrderRateLimiter 下单限流器 ============================
type OrderRateLimiter struct {
	timesAlgorithm *TimesAlgorithm
}

func NewOrderRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[BinanceRateLimiterRequest] {
	return &OrderRateLimiter{
		timesAlgorithm: NewTimesAlgorithm(redisClient),
	}
}

func (r *OrderRateLimiter) Check(ctx context.Context, request BinanceRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
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

func (r *OrderRateLimiter) extractRuleKey(request BinanceRateLimiterRequest) string {
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
	return fmt.Sprintf("%s:%s:%s", "binance", marketType, string(request.RequestType))
}

func (r *OrderRateLimiter) extractRedisKey(request BinanceRateLimiterRequest) string {
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
	return fmt.Sprintf("%s:%s:%s:%s", "binance", marketType, string(request.RequestType), request.AccountID)
}

// ============================ IPRateLimiter 权重限流器(权重粒度只到IP) ============================
type IPRateLimiter struct {
	weightAlgorithm *WeightAlgorithm
}

func NewIPRateLimiter(redisClient *redis.Client) ratelimiter.RateLimiter[BinanceRateLimiterRequest] {
	return &IPRateLimiter{
		weightAlgorithm: NewWeightAlgorithm(redisClient),
	}
}

func (r *IPRateLimiter) Check(ctx context.Context, request BinanceRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
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

func (r *IPRateLimiter) extractRuleKey(request BinanceRateLimiterRequest) string {

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

func (r *IPRateLimiter) extractRedisKey(request BinanceRateLimiterRequest) string {
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

func (r *IPRateLimiter) extractWeightKey(request BinanceRateLimiterRequest) string {
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

func getRules(key string) ([]ratelimiter.RateLimitRule, error) {
	defaultRules := DefaultBinanceConfig()
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
