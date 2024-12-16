package binance

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"github.com/go-gotop/gotop/requests"
)

// BinanceAdapter 实现 ExchangeAdapter，用于构建符合 Binance API 要求的 HTTP 请求。
// 假设用户传入的 URL 是全路径（例如："https://api.binance.com/api/v3/order"），在 BuildRequest 中会自动拼接成完整URL。
// 对于需要鉴权的请求（即 req.Auth 存在且 SecretKey 不为空），会注入签名相关的参数及头信息。
type BinanceAdapter struct {

}

func NewBinanceAdapter() *BinanceAdapter {
	return &BinanceAdapter{

	}
}

// BuildRequest 根据输入参数构建一个完整的 PreparedRequest。
// 内部完成：URL拼接、查询参数处理、Headers构建、签名注入、请求体序列化等所有定制逻辑。
func (b *BinanceAdapter) BuildRequest(req *requests.Request) (*requests.PreparedRequest, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	if req.URL == "" {
		return nil, errors.New("request URL is empty")
	}

	// 判断是否需要签名
	needsSignature := false
	if req.Auth != nil && req.Auth.SecretKey != "" {
		needsSignature = true
	}

	// 初始化输出结果
	prepared := &requests.PreparedRequest{
		Method:  req.Method,
		Headers: http.Header{},
	}

	// 构建完整URL
	fullURL := req.URL
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// 根据参数处理query和body
	params := req.Params
	if params == nil {
		params = make(map[string]interface{})
	}

	var queryParams url.Values
	bodyBytes := []byte{}

	switch strings.ToUpper(req.Method) {
	case http.MethodGet, http.MethodDelete:
		// 对于GET/DELETE请求，将所有业务参数放到URL的query中
		queryParams = convertParamsToQuery(params)
		if needsSignature {
			// Binance需要添加timestamp
			timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			queryParams.Set("timestamp", timestamp)
			// 对参数进行签名
			queryStr := queryParams.Encode()
			signature := sign(queryStr, req.Auth.SecretKey)
			queryParams.Set("signature", signature)
		}
		u.RawQuery = queryParams.Encode()

	case http.MethodPost, http.MethodPut:
		// 对于POST/PUT请求，Binance在SIGNED类型的请求中通常使用query参数与签名，而不是json body。
		// 此处以与官方REST API使用习惯相符的方式实现：
		// 若needsSignature为true，则使用application/x-www-form-urlencoded，将所有参数转为query参数并进行签名。
		// 否则使用JSON序列化参数至body中。
		if needsSignature {
			queryParams = convertParamsToQuery(params)
			timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
			queryParams.Set("timestamp", timestamp)

			queryStr := queryParams.Encode()
			signature := sign(queryStr, req.Auth.SecretKey)
			queryParams.Set("signature", signature)

			prepared.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
			bodyBytes = []byte(queryParams.Encode())
		} else {
			prepared.Headers.Set("Content-Type", "application/json")
			bodyBytes, err = json.Marshal(params)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", req.Method)
	}

	// 设置 Headers
	if needsSignature && req.Auth != nil && req.Auth.APIKey != "" {
		prepared.Headers.Set("X-MBX-APIKEY", req.Auth.APIKey)
	}
	// 一些公共Headers
	prepared.Headers.Set("Accept", "application/json")
	prepared.Headers.Set("User-Agent", "gotop/1.0")

	prepared.URL = u.String()
	prepared.Body = bodyBytes

	// 请求方法检查
	if prepared.Method == "" {
		return nil, errors.New("HTTP method is required")
	}

	return prepared, nil
}

// convertParamsToQuery 将map转换为url.Values并对key进行排序（有利于调试和签名一致性）。
func convertParamsToQuery(params map[string]interface{}) url.Values {
	query := url.Values{}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := params[k]
		query.Set(k, fmt.Sprintf("%v", v))
	}
	return query
}

// sign 使用HMAC-SHA256对字符串进行签名
func sign(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}