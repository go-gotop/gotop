package okxexc

import (
	"strings"

	"github.com/go-gotop/gotop/types"
)

const (
	OKX_API_BASE_URL = "https://www.okx.com"
)

func toOkxSide(side types.SideType) string {
	return strings.ToLower(side.String())
}

func toOkxOrderType(orderType types.OrderType) string {
	return strings.ToLower(orderType.String())
}

func toOkxPositionSide(positionSide types.PositionSide) string {
	return strings.ToLower(positionSide.String())
}

func toOkxPosMode(posMode types.PosMode) string {
	return strings.ToLower(posMode.String())
}

// okx创建订单响应
type okxOrderResponse struct {
	Code string `json:"code"`
	Data []struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
		Tag     string `json:"tag"`
		Ts      string `json:"ts"`
	} `json:"data"`
	InTime  string `json:"inTime"`
	Msg     string `json:"msg"`
	OutTime string `json:"outTime"`
}

// okx深度响应
type okxDepthResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
		Ts   string     `json:"ts"`
	} `json:"data"`
}

// okx ticker 响应
type okxTickerResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		// 产品类型
		InstType string `json:"instType"`
		// 产品ID
		InstId string `json:"instId"`
		// 最新成交价
		Last string `json:"last"`
		// 最新成交数量
		LastSz string `json:"lastSz"`
		// 卖一价
		AskPx string `json:"askPx"`
		// 卖一量
		AskSz string `json:"askSz"`
		// 买一价
		BidPx string `json:"bidPx"`
		// 买一量
		BidSz string `json:"bidSz"`
		// 24小时开盘价
		Open24h string `json:"open24h"`
		// 24小时最高价
		High24h string `json:"high24h"`
		// 24小时最低价
		Low24h string `json:"low24h"`
		// 24小时成交量（计价货币）
		VolCcy24h string `json:"volCcy24h"`
		// 24小时成交量（交易货币）
		Vol24h string `json:"vol24h"`
		// 时间戳
		Ts string `json:"ts"`
		// 24小时开盘价（UTC0）
		SodUtc0 string `json:"sodUtc0"`
		// 24小时开盘价（UTC8）
		SodUtc8 string `json:"sodUtc8"`
	} `json:"data"`
}
