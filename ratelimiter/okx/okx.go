package okx

import (
	"context"
	"errors"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/redis/go-redis/v9"
)

type OkxRateLimiterManager struct {
	redisClient *redis.Client
}

func NewOkxRateLimiterManager(redisClient *redis.Client) *OkxRateLimiterManager {
	return &OkxRateLimiterManager{
		redisClient: redisClient,
	}
}

func (m *OkxRateLimiterManager) PreCheck(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	rateLimiters := make([]ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest], 0)

	switch request.RequestType {
	case ratelimiter.RequestTypeOrder:
		rateLimiters = append(rateLimiters, NewTimesRateLimiter(m.redisClient))
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
