package requests

import (
    "net/http"
    "testing"
    "time"
)

func TestWithTimeout(t *testing.T) {
    client := &httpClientImpl{}
    timeout := 5 * time.Second

    WithTimeout(timeout)(client)

    if client.client.Timeout != timeout {
        t.Errorf("WithTimeout() = %v, want %v", client.client.Timeout, timeout)
    }
}

func TestWithClient(t *testing.T) {
    customClient := &http.Client{
        Timeout: 30 * time.Second,
    }
    client := &httpClientImpl{}

    WithClient(customClient)(client)

    if client.client != customClient {
        t.Error("WithClient() did not set the custom client correctly")
    }
}

func TestWithRequestSigner(t *testing.T) {
    client := &httpClientImpl{}
    signer := func(req *http.Request, params interface{}) error {
        req.Header.Set("Authorization", "Bearer test-token")
        return nil
    }

    WithRequestSigner(signer)(client)

    if client.signer == nil {
        t.Error("WithRequestSigner() did not set the signer")
        return
    }

    // 测试签名功能
    req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
    err := client.signer(req, nil)
    if err != nil {
        t.Errorf("Signer returned unexpected error: %v", err)
    }
    if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
        t.Errorf("Signer did not set header correctly, got %v", got)
    }
}

func TestWithErrorHandler(t *testing.T) {
    client := &httpClientImpl{}
    handler := func(statusCode int, body []byte) error {
        return nil
    }

    WithErrorHandler(handler)(client)

    if client.errorHandler == nil {
        t.Error("WithErrorHandler() did not set the error handler")
    }
}

func TestWithHeaderSetter(t *testing.T) {
    client := &httpClientImpl{}
    setter := func(req *http.Request, params interface{}) error {
        req.Header.Set("X-Test", "test-value")
        return nil
    }

    WithHeaderSetter(setter)(client)

    if client.headerSetter == nil {
        t.Error("WithHeaderSetter() did not set the header setter")
        return
    }

    // 测试header设置功能
    req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
    err := client.headerSetter(req, nil)
    if err != nil {
        t.Errorf("HeaderSetter returned unexpected error: %v", err)
    }
    if got := req.Header.Get("X-Test"); got != "test-value" {
        t.Errorf("HeaderSetter did not set header correctly, got %v", got)
    }
}

func TestWithRequestEncoder(t *testing.T) {
    client := &httpClientImpl{}
    encoder := func(method, urlStr string, params interface{}) (*http.Request, error) {
        return http.NewRequest(method, urlStr, nil)
    }

    WithRequestEncoder(encoder)(client)

    if client.requestEncoder == nil {
        t.Error("WithRequestEncoder() did not set the request encoder")
        return
    }

    // 测试编码器功能
    req, err := client.requestEncoder("GET", "http://example.com", nil)
    if err != nil {
        t.Errorf("RequestEncoder returned unexpected error: %v", err)
    }
    if req.Method != "GET" {
        t.Errorf("RequestEncoder did not set method correctly, got %v", req.Method)
    }
}

func TestOptionsChaining(t *testing.T) {
    client := &httpClientImpl{}
    timeout := 5 * time.Second
    customClient := &http.Client{Timeout: timeout}

    // 测试多个选项的链式调用
    opts := []HttpClientOption{
        WithClient(customClient),
        WithHeaderSetter(func(req *http.Request, _ interface{}) error {
            req.Header.Set("X-Test", "test-value")
            return nil
        }),
        WithRequestSigner(func(req *http.Request, _ interface{}) error {
            req.Header.Set("Authorization", "Bearer test-token")
            return nil
        }),
    }

    for _, opt := range opts {
        opt(client)
    }

    // 验证所有选项是否都正确应用
    if client.client != customClient {
        t.Error("Client was not set correctly in chained options")
    }
    if client.headerSetter == nil {
        t.Error("HeaderSetter was not set correctly in chained options")
    }
    if client.signer == nil {
        t.Error("RequestSigner was not set correctly in chained options")
    }

    // 测试设置的功能是否正常
    req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
    
    err := client.headerSetter(req, nil)
    if err != nil || req.Header.Get("X-Test") != "test-value" {
        t.Error("HeaderSetter not working in chained options")
    }

    err = client.signer(req, nil)
    if err != nil || req.Header.Get("Authorization") != "Bearer test-token" {
        t.Error("RequestSigner not working in chained options")
    }
} 