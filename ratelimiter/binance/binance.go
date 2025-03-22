package binance

import (
	"context"
	"errors"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/redis/go-redis/v9"
)

// ============================ RateLimiterManager ============================

type BinanceRateLimiterManager struct {
	redisClient *redis.Client
}

func NewBinanceRateLimiterManager(
	redisClient *redis.Client,
) *BinanceRateLimiterManager {
	return &BinanceRateLimiterManager{
		redisClient: redisClient,
	}
}

// GetRedisClient 返回Redis客户端，用于测试
func (m *BinanceRateLimiterManager) GetRedisClient() *redis.Client {
	return m.redisClient
}

func (m *BinanceRateLimiterManager) PreCheck(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	rateLimiters := make([]ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest], 0)

	switch request.RequestType {
	case ratelimiter.RequestTypeOrder:
		rateLimiters = append(rateLimiters, NewOrderRateLimiter(m.redisClient))
		rateLimiters = append(rateLimiters, NewIPRateLimiter(m.redisClient))
	case ratelimiter.RequestTypeNormal:
		rateLimiters = append(rateLimiters, NewIPRateLimiter(m.redisClient))
	default:
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  "unsupported request type",
		}, errors.New("unsupported request type")
	}

	for _, rateLimiter := range rateLimiters {
		decision, err := rateLimiter.Check(ctx, request)
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
