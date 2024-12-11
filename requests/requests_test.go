package requests

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHttpClient_DoRequest(t *testing.T) {
    // 创建测试服务器
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/success":
            // 测试成功请求
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
        case "/error":
            // 测试错误响应
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
        case "/custom-header":
            // 测试自定义请求头
            if r.Header.Get("X-Custom-Header") != "test-value" {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
        }
    }))
    defer ts.Close()

    tests := []struct {
        name       string
        method     string
        path       string
        params     interface{}
        opts       []HttpClientOption
        wantErr    bool
        statusCode int
    }{
        {
            name:       "GET success",
            method:     http.MethodGet,
            path:       "/success",
            params:     nil,
            wantErr:    false,
            statusCode: http.StatusOK,
        },
        {
            name:       "POST success",
            method:     http.MethodPost,
            path:       "/success",
            params:     map[string]string{"key": "value"},
            wantErr:    false,
            statusCode: http.StatusOK,
        },
        {
            name:       "Error response",
            method:     http.MethodGet,
            path:       "/error",
            params:     nil,
            wantErr:    true,
            statusCode: http.StatusBadRequest,
        },
        {
            name:   "Custom header",
            method: http.MethodGet,
            path:   "/custom-header",
            opts: []HttpClientOption{
                WithHeaderSetter(func(req *http.Request, _ interface{}) error {
                    req.Header.Set("X-Custom-Header", "test-value")
                    return nil
                }),
            },
            wantErr:    false,
            statusCode: http.StatusOK,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := NewHttpClient(tt.opts...)
            
            resp, err := client.DoRequest(
                context.Background(),
                tt.method,
                ts.URL+tt.path,
                tt.params,
            )

            if (err != nil) != tt.wantErr {
                t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && len(resp) == 0 {
                t.Error("DoRequest() returned empty response for success case")
            }
        })
    }
}

func TestDefaultJSONRequestEncoder(t *testing.T) {
    tests := []struct {
        name    string
        method  string
        url     string
        params  interface{}
        wantErr bool
    }{
        {
            name:    "GET without params",
            method:  http.MethodGet,
            url:     "http://example.com",
            params:  nil,
            wantErr: false,
        },
        {
            name:    "POST with params",
            method:  http.MethodPost,
            url:     "http://example.com",
            params:  map[string]string{"key": "value"},
            wantErr: false,
        },
        {
            name:    "Invalid params",
            method:  http.MethodPost,
            url:     "http://example.com",
            params:  make(chan int), // 不可序列化的类型
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := DefaultJSONRequestEncoder(tt.method, tt.url, tt.params)
            if (err != nil) != tt.wantErr {
                t.Errorf("DefaultJSONRequestEncoder() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr {
                if req.Method != tt.method {
                    t.Errorf("DefaultJSONRequestEncoder() method = %v, want %v", req.Method, tt.method)
                }

                if tt.method == http.MethodPost && req.Header.Get("Content-Type") != "application/json" {
                    t.Error("DefaultJSONRequestEncoder() Content-Type header not set for POST request")
                }
            }
        })
    }
} 