package binance

import (
	"github.com/go-gotop/gotop/stream"
)

type BinanceStream struct {
	stream.Stream
}

func NewBinanceStream() *BinanceStream {
	return &BinanceStream{}
}

func (s *BinanceStream) Connect() error {
	return nil
}

func (s *BinanceStream) Disconnect() error {
	return nil
}

func (s *BinanceStream) AddMessageHandler(h stream.MessageHandler) {
}

func (s *BinanceStream) AddErrorHandler(h stream.ErrorHandler) {
}

func (s *BinanceStream) AddCloseHandler(h stream.CloseHandler) {
}
