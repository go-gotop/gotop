package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-gotop/gotop/requests"
)

// NewOKXAdapter 创建一个新的 OKXAdapter 实例。
func NewOKXAdapter() *OKXAdapter {
	return &OKXAdapter{}
}

// OKXAdapter 实现了 ExchangeAdapter 接口，用于根据 OKX 的 HTTP 请求签名流程构建请求。
type OKXAdapter struct {
	// 如果有其他通用配置项可在此处扩展。例如日志记录器、调试开关等。
}

// BuildRequest 根据 OKX 的要求构建一个完整的请求。
// 内部步骤:
// 1. 参数处理（若是GET等请求，需要将参数以 QueryString 的形式追加到 URL 后面；对于 POST/PUT/DELETE 等请求将参数序列化为 JSON 并放入 Body）。
// 2. 获取当前 UTC 时间戳，并按要求进行签名串的构建。
// 3. 使用 secretKey 对签名串进行 HmacSHA256 后 Base64 编码，生成签名。
// 4. 构建 Headers，包括鉴权所需的 OK-ACCESS-KEY, OK-ACCESS-SIGN, OK-ACCESS-TIMESTAMP, OK-ACCESS-PASSPHRASE。
// 5. 返回构建好的 PreparedRequest 对象。
func (o *OKXAdapter) BuildRequest(req *requests.Request) (*requests.PreparedRequest, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.Method == "" {
		return nil, errors.New("missing HTTP method")
	}
	if req.URL == "" {
		return nil, errors.New("missing request URL")
	}

	method := strings.ToUpper(req.Method)

	// 构建URL和Body
	var (
		finalURL   = req.URL
		bodyBytes  []byte
		err        error
		queryPairs = url.Values{}
	)

	if req.Params != nil && len(req.Params) > 0 {
		keys := make([]string, 0, len(req.Params))
		for k := range req.Params {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		switch method {
		case http.MethodGet, http.MethodDelete:
			for _, k := range keys {
				val := fmt.Sprintf("%v", req.Params[k])
				queryPairs.Add(k, val)
			}
			if strings.Contains(finalURL, "?") {
				finalURL += "&" + queryPairs.Encode()
			} else {
				finalURL += "?" + queryPairs.Encode()
			}

		case http.MethodPost, http.MethodPut, http.MethodPatch:
			jsonData := make(map[string]interface{}, len(keys))
			for _, k := range keys {
				jsonData[k] = req.Params[k]
			}
			bodyBytes, err = json.Marshal(jsonData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}

		default:
			// 对于不常见的HTTP方法，也按JSON序列化处理
			jsonData := make(map[string]interface{}, len(keys))
			for _, k := range keys {
				jsonData[k] = req.Params[k]
			}
			bodyBytes, err = json.Marshal(jsonData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
		}
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	// 如果存在 AuthInfo，则构建签名
	if req.Auth != nil {
		parsedURL, err := url.Parse(finalURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}

		pathWithQuery := parsedURL.Path
		if parsedURL.RawQuery != "" {
			pathWithQuery += "?" + parsedURL.RawQuery
		}

		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		var signData string
		if len(bodyBytes) > 0 {
			signData = timestamp + method + pathWithQuery + string(bodyBytes)
		} else {
			signData = timestamp + method + pathWithQuery
		}

		sign, err := signMessage(signData, req.Auth.SecretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to sign message: %w", err)
		}

		headers.Set("OK-ACCESS-KEY", req.Auth.APIKey)
		headers.Set("OK-ACCESS-SIGN", sign)
		headers.Set("OK-ACCESS-TIMESTAMP", timestamp)
		headers.Set("OK-ACCESS-PASSPHRASE", req.Auth.Passphrase)
	}

	return &requests.PreparedRequest{
		Method:  method,
		URL:     finalURL,
		Headers: headers,
		Body:    bodyBytes,
	}, nil
}

// signMessage 使用 secretKey 对 signData 进行 HMAC-SHA256 签名，并进行 Base64 编码
func signMessage(signData string, secretKey string) (string, error) {
	h := hmac.New(sha256.New, []byte(secretKey))
	_, err := h.Write([]byte(signData))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}