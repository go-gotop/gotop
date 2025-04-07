package okx

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/redis/go-redis/v9"
)

// TimesAlgorithm 基于时间窗口的次数限流算法
type TimesAlgorithm struct {
	// redis客户端
	redisClient *redis.Client
}

func NewTimesAlgorithm(
	redisClient *redis.Client,
) *TimesAlgorithm {
	return &TimesAlgorithm{
		redisClient: redisClient,
	}
}

func (r *TimesAlgorithm) Check(
	key string,
	rules []ratelimiter.RateLimitRule,
) (ratelimiter.RateLimitDecision, error) {
	ctx := context.Background()
	now := time.Now().Unix()

	// 检查是否开启调试模式
	debugMode := os.Getenv("DEBUG") != ""

	if debugMode {
		fmt.Printf("检查键: %s，规则数: %d\n", key, len(rules))
		for i, rule := range rules {
			fmt.Printf("规则 #%d: 窗口=%v, 阈值=%d\n", i+1, rule.Window, rule.Threshold)
		}
	}

	// 如果没有规则，则默认允许
	if len(rules) == 0 {
		if debugMode {
			fmt.Println("没有找到匹配的规则，默认允许请求")
		}
		return ratelimiter.RateLimitDecision{
			Allowed: true,
			Reason:  "no rules found",
		}, nil
	}

	// 使用Lua脚本保证原子性
	luaScript := `
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	
	-- 对每个规则分别检查
	for i=3, #ARGV, 2 do
		local window = tonumber(ARGV[i])
		local threshold = tonumber(ARGV[i+1])
		
		-- 计算窗口的起始时间
		local minTime = now - window
		
		-- 计算窗口内的请求数
		local count = redis.call('ZCOUNT', key, minTime, '+inf')
		
		-- 检查是否超过阈值
		if count >= threshold then
			return "false:Too many requests: limit " .. threshold .. " per " .. window .. " seconds reached"
		end
	end
	
	-- 允许请求：添加记录
	-- 添加当前时间戳作为记录
	redis.call('ZADD', key, now, now .. ":" .. math.random(9999999))
	
	-- 设置过期时间（所有规则窗口的最大值的2倍）
	local maxWindow = 86400  -- 默认1天（秒）
	for i=3, #ARGV, 2 do
		maxWindow = math.max(maxWindow, tonumber(ARGV[i]))
	end
	redis.call('EXPIRE', key, maxWindow * 2)
	
	-- 清理过期数据（使用最小窗口时间）
	local minWindow = 86400
	for i=3, #ARGV, 2 do
		minWindow = math.min(minWindow, tonumber(ARGV[i]))
	end
	local oldestTime = now - minWindow
	redis.call('ZREMRANGEBYSCORE', key, 0, oldestTime)
	
	return "true:"
	`

	// 准备参数
	keys := []string{key}
	args := []interface{}{now, debugMode}

	// 添加所有规则参数
	for i, rule := range rules {
		args = append(args, int64(rule.Window.Seconds()), rule.Threshold)

		if debugMode {
			fmt.Printf("添加规则参数 #%d: 窗口=%d秒, 阈值=%d\n",
				i+1, int64(rule.Window.Seconds()), rule.Threshold)
		}
	}

	// 执行Lua脚本
	result, err := r.redisClient.Eval(ctx, luaScript, keys, args...).Result()
	if err != nil {
		if debugMode {
			fmt.Printf("Redis错误: %v\n", err)
		}
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  fmt.Sprintf("redis error: %v", err),
		}, err
	}

	// 解析结果
	if resultStr, ok := result.(string); ok {
		parts := strings.SplitN(resultStr, ":", 2)
		allowed := parts[0] == "true"
		reason := ""
		if len(parts) > 1 {
			reason = parts[1]
		}

		if debugMode {
			fmt.Printf("决策结果: allowed=%v, reason=%s\n", allowed, reason)
		}

		return ratelimiter.RateLimitDecision{
			Allowed: allowed,
			Reason:  reason,
		}, nil
	}

	// 返回默认结果
	if debugMode {
		fmt.Println("无法解析结果，返回默认允许决策")
	}

	return ratelimiter.RateLimitDecision{
		Allowed: true,
	}, nil
}
