package telegram

import "github.com/go-gotop/gotop/notifier"

type TelegramNotifier struct {
	
}

func (t *TelegramNotifier) Notify(msg notifier.Message) error {
	return nil
}