package bnexc

const (
	// 现货&杠杆
	BNEX_API_SPOT_URL = "https://api.binance.com"
	// U本位合约
	BNEX_API_FUTURES_USD_URL = "https://fapi.binance.com"
	// 币本位合约
	BNEX_API_FUTURES_COIN_URL = "https://dapi.binance.com"
)


// bnOrderACKResponse 币安下单返回响应(下单最快返回)
type bnOrderACKResponse struct {
	// 用户自定义的订单号
	ClientOrderID string `json:"clientOrderId"`
	// 系统订单号
	OrderID int64 `json:"orderId"`
	// 交易对
	Symbol string `json:"symbol"`
}

// bnTickerPriceResponse 币安ticker价格响应
type bnTickerPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}