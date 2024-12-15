package telegram

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/go-gotop/gotop/notifier"
)

// TelegramConfig 定义了创建 TelegramNotifier 所需的配置。
// 此配置结构最终应当从 GlobalConfig.Notifiers["telegram"] 中解析而来。
type TelegramConfig struct {
    BotToken string `yaml:"bot_token"`
    ChatID   string `yaml:"chat_id"`
    APIURL   string `yaml:"api_url"` // 可选，默认为 https://api.telegram.org
}

// TelegramNotifier 是对 Telegram 通知渠道的具体实现。
type TelegramNotifier struct {
    botToken string
    chatID   string
    apiURL   string
    client   *http.Client
}

// NewTelegramNotifier 根据配置创建一个 TelegramNotifier 实例。
// 通常由工厂在 CreateNotifier 中调用。
func NewTelegramNotifier(cfg TelegramConfig) *TelegramNotifier {
    apiURL := cfg.APIURL
    if apiURL == "" {
        apiURL = "https://api.telegram.org"
    }

    return &TelegramNotifier{
        botToken: cfg.BotToken,
        chatID:   cfg.ChatID,
        apiURL:   apiURL,
        client: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

// Notify 实现 Notifier 接口。根据传入的 Message 构造请求并发送到 Telegram API。
// 若 Message.Metadata 中存在 "chat_id" 字段，则使用该值优先发送。
func (t *TelegramNotifier) Notify(msg notifier.Message) error {
    chatID := t.chatID
    if metaChatID, ok := msg.Metadata["chat_id"]; ok && metaChatID != "" {
        chatID = metaChatID
    }

    if chatID == "" {
        return fmt.Errorf("no chat_id provided for Telegram message")
    }

    // 构造 Telegram sendMessage API 请求参数
    requestBody := map[string]interface{}{
        "chat_id": chatID,
        "text":    fmt.Sprintf("%s\n%s", msg.Subject, msg.Body),
    }

    bodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        return fmt.Errorf("failed to marshal request body: %w", err)
    }

    req, err := http.NewRequest("POST", fmt.Sprintf("%s/bot%s/sendMessage", t.apiURL, t.botToken), bytes.NewReader(bodyBytes))
    if err != nil {
        return fmt.Errorf("failed to create new request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := t.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request to Telegram API: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code from Telegram API: %d", resp.StatusCode)
    }

    return nil
}

// TelegramNotifierFactory 根据 TelegramConfig 创建 TelegramNotifier。
type TelegramNotifierFactory struct{}

// CreateNotifier 实现 NotifierFactory 接口。
// 入参 cfg 应该是一个可以解析为 TelegramConfig 的结构，如：
//
//  notifiers:
//    telegram:
//      bot_token: "123456789:ABC..."
//      chat_id: "987654321"
//      api_url: "https://api.telegram.org"   # 可选
//
func (f *TelegramNotifierFactory) CreateNotifier(cfg interface{}) (notifier.Notifier, error) {
    confMap, ok := cfg.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid config type for Telegram")
    }

    botToken, ok := confMap["bot_token"].(string)
    if !ok || botToken == "" {
        return nil, fmt.Errorf("telegram bot_token not provided")
    }

    chatID, _ := confMap["chat_id"].(string)
    apiURL, _ := confMap["api_url"].(string)

    telegramCfg := TelegramConfig{
        BotToken: botToken,
        ChatID:   chatID,
        APIURL:   apiURL,
    }

    return NewTelegramNotifier(telegramCfg), nil
}