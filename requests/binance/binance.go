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
	"net/url"
	"sort"

	"github.com/go-gotop/gotop/requests"
)

// BinanceAdapter 封装了构建 Binance API 请求的逻辑。
// - Binance官方文档：https://binance-docs.github.io/apidocs/spot/en
// 通常流程：
//   1. 将 params 编码为查询字符串，或根据 method 决定放在请求体中。
//   2. 在请求参数中加入 timestamp, recvWindow 等字段 (如果需要)。
//   3. 使用 HMAC SHA256 对参数签名，将签名附加到查询参数中。
//   4. 设置适当的 Headers，包括API Key。
//   5. 根据实际情况对 GET 使用 query string，对 POST 等使用 body。
type BinanceAdapter struct {
	BaseURL  string
	APIKey   string
	SecretKey string
	// 可选的 recvWindow 设置
	RecvWindow int64
}

// BuildRequest 为 Binance 构建请求
func (b *BinanceAdapter) BuildRequest(method, endpoint string, params map[string]interface{}) (*requests.PreparedRequest, error) {
	if b.APIKey == "" || b.SecretKey == "" {
		return nil, errors.New("binance adapter requires apiKey and secretKey")
	}

	fullURL := b.BaseURL
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	fullURL += endpoint

	// Binance通常要求timestamp和signature
	if params == nil {
		params = make(map[string]interface{})
	}
	if b.RecvWindow > 0 {
		params["recvWindow"] = b.RecvWindow
	}
	params["timestamp"] = time.Now().UnixNano() / int64(time.Millisecond)

	// 构建query string
	queryStr, err := buildQueryString(params)
	if err != nil {
		return nil, err
	}

	// 对params进行签名
	signature := hmacSignSHA256(queryStr, b.SecretKey)
	queryStr += "&signature=" + signature

	var req *requests.PreparedRequest
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodDelete:
		// 对于GET,DELETE请求，将query置于URL上
		u := fullURL + "?" + queryStr
		req = &requests.PreparedRequest{
			Method:  strings.ToUpper(method),
			URL:     u,
			Headers: http.Header{},
			Body:    nil,
		}
	case http.MethodPost, http.MethodPut:
		// 对于POST,PUT请求，也可考虑使用 query string 或者放 Body 中
		// Binance大部分交易接口要求 form 表单格式 (Content-Type: application/x-www-form-urlencoded)
		// 这里使用query string作为body传递。
		req = &requests.PreparedRequest{
			Method:  strings.ToUpper(method),
			URL:     fullURL,
			Headers: http.Header{},
			Body:    []byte(queryStr),
		}
		req.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	// 设置头部
	req.Headers.Set("X-MBX-APIKEY", b.APIKey)

	return req, nil
}

// 辅助函数：构建查询字符串（key按字典序排序）
func buildQueryString(params map[string]interface{}) (string, error) {
	if len(params) == 0 {
		return "", nil
	}
	kv := make([]string, 0, len(params))
	for k, v := range params {
		// 将value序列化为string
		valStr := fmt.Sprintf("%v", v)
		kv = append(kv, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(valStr)))
	}
	sort.Strings(kv)
	return strings.Join(kv, "&"), nil
}

// 辅助函数：HMAC SHA256签名 (Binance)
func hmacSignSHA256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}