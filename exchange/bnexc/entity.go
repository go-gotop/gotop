package bnexc

// bnOrderACKResponse BNEX Order Ack Response (下单最快返回)
type bnOrderACKResponse struct {
	// 用户自定义的订单号
	ClientOrderID string `json:"clientOrderId"`
	// 系统订单号
	OrderID int64 `json:"orderId"`
	// 交易对
	Symbol string `json:"symbol"`
}

// bnTickerPriceResponse BNEX Ticker Price Response
type bnTickerPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
