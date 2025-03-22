package binance

const (
	// ============ 下单限制规则Key ================================
	// BinanceFuturesCreateOrder10sKey 合约每10s不超过300次下单请求
	BinanceFuturesCreateOrder10sKey = "binance:futures:createorder:10s"
	// BinanceFuturesCreateOrder1mKey 合约每分钟不超过1200次下单请求
	BinanceFuturesCreateOrder1mKey = "binance:futures:createorder:1m"
	// BinanceSpotCreateOrder10sKey 现货每10s不超过100次下单请求
	BinanceSpotCreateOrder10sKey = "binance:spot:createorder:10s"

	// ============ 权重限制规则Key ================================
	// BinanceSpotRequest1mKey 现货每分钟不超过6000权重请求
	BinanceSpotRequest1mKey = "binance:spot:request:1m"
	// BinanceFuturesRequest1mKey 合约每分钟不超过2400权重请求
	BinanceFuturesRequest1mKey = "binance:futures:request:1m"

	// ============ 权重Key ================================
	BinanceFuturesCreateOrderWeightKey = "binance:futures:createorder:weight"
	BinanceSpotCreateOrderWeightKey    = "binance:spot:createorder:weight"
	BinanceSpotRequestWeightKey        = "binance:spot:request:weight"
	BinanceFuturesRequestWeightKey     = "binance:futures:request:weight"
)
