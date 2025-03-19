package okxexc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/requests"
	okxreq "github.com/go-gotop/gotop/requests/okx"
	"github.com/go-gotop/gotop/types"
	"github.com/shopspring/decimal"
)

var _ exchange.OrderManager = &OkxOrderManager{}

type OkxOrderManager struct {
	client requests.RequestClient
}

func NewOkxOrderManager() *OkxOrderManager {
	adapter := okxreq.NewOKXAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &OkxOrderManager{
		client: client,
	}
}

// CreateOrder 创建订单
func (o *OkxOrderManager) CreateOrder(ctx context.Context, req *exchange.CreateOrderRequest) (*exchange.CreateOrderResponse, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/trade/order"

	params, err := o.toOrderParams(req)
	if err != nil {
		return nil, err
	}

	resp, err := o.client.DoRequest(&requests.Request{
		Method: http.MethodPost,
		URL:    apiUrl,
		Params: params,
		Auth: &requests.AuthInfo{
			APIKey:     req.APIKey,
			SecretKey:  req.SecretKey,
			Passphrase: req.Passphrase,
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var respData okxOrderResponse
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	if respData.Code != "0" || len(respData.Data) == 0 || respData.Data[0].SCode != "0" {
		msg := respData.Msg
		code := respData.Code
		if len(respData.Data) > 0 {
			msg = respData.Data[0].SMsg
			code = respData.Data[0].SCode
		}
		return nil, fmt.Errorf("operation failed, code: %s, message: %s", code, msg)
	}

	return &exchange.CreateOrderResponse{
		OrderID:       respData.Data[0].OrdId,
		ClientOrderID: respData.Data[0].ClOrdId,
		Symbol:        respData.Data[0].SCode,
	}, nil
}

// CancelOrder 取消订单
func (o *OkxOrderManager) CancelOrder(ctx context.Context, req *exchange.CancelOrderRequest) (*exchange.CancelOrderResponse, error) {
	return nil, errors.New("not implemented")
}

// GetOrder 获取订单
func (o *OkxOrderManager) GetOrder(ctx context.Context, req *exchange.GetOrderRequest) (*exchange.GetOrderResponse, error) {
	return nil, errors.New("not implemented")
}

// getTickerPrice 获取ticker价格
func (o *OkxOrderManager) getTickerPrice(symbol types.Symbol) (decimal.Decimal, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/market/ticker"

	params := map[string]any{
		"instId": symbol.OriginalSymbol,
	}

	resp, err := o.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Params: params,
	})
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var tickerResponse okxTickerResponse
	err = json.Unmarshal(body, &tickerResponse)
	if err != nil {
		return decimal.Zero, err
	}

	if len(tickerResponse.Data) == 0 {
		return decimal.Zero, fmt.Errorf("ticker response data is empty")
	}

	return decimal.NewFromString(tickerResponse.Data[0].Last)
}

// 下单请求参数
func (o *OkxOrderManager) toOrderParams(req *exchange.CreateOrderRequest) (map[string]any, error) {
	params := map[string]any{
		"instId":  req.Symbol.OriginalSymbol,
		"clOrdId": req.ClientOrderID,
		"side":    toOkxSide(req.Side),
		"ordType": toOkxOrderType(req.OrderType),
	}

	if req.MarketType == types.MarketTypeFuturesUSDMargined ||
		req.MarketType == types.MarketTypePerpetualUSDMargined ||
		req.MarketType == types.MarketTypeFuturesCoinMargined ||
		req.MarketType == types.MarketTypePerpetualCoinMargined {
		// 合约规则
		// 1. size 单位是张
		opType := "open"
		if req.Side == types.SideTypeSell && req.PositionSide == types.PositionSideLong ||
			req.Side == types.SideTypeBuy && req.PositionSide == types.PositionSideShort {
			opType = "close"
		}
		// 币转张
		size, err := o.convertContractCoin("1", req.Symbol, fmt.Sprintf("%v", req.Size), opType)
		if err != nil {
			return nil, err
		}
		params["sz"] = size.String()
		// 默认全仓
		params["tdMode"] = toOkxPosMode(types.PosModeCross)
		// 持仓方向
		params["posSide"] = toOkxPositionSide(req.PositionSide)
	} else if req.MarketType == types.MarketTypeSpot {
		// 现货规则
		// 1. 以限价单成交，size 单位指交易货币
		// 2. 以市价单成交，size 单位由tgtCcy指定
		params["tgtCcy"] = "base_ccy"
		// 现货默认现金模式
		params["tdMode"] = "cash"
		params["sz"] = fmt.Sprintf("%v", req.Size)

	} else if req.MarketType == types.MarketTypeMargin {
		// 杠杆规则：
		// 1. 以限价单成交，size 单位指交易货币
		// 2. 以市价单买入，size 单位指计价货币
		// 3. 以市价单卖出，size 单位指交易货币
		// 默认全仓
		params["tdMode"] = toOkxPosMode(types.PosModeCross)
		// 保证金固定为USDT
		params["ccy"] = "USDT"
		if req.Side == types.SideTypeBuy && req.OrderType == types.OrderTypeMarket {
			tickerPrice, err := o.getTickerPrice(req.Symbol)
			if err != nil {
				return nil, err
			}
			size := req.Size.Mul(tickerPrice)
			params["sz"] = size.String()
		} else {
			params["sz"] = fmt.Sprintf("%v", req.Size)
		}
	}

	// 限价单
	if req.OrderType == types.OrderTypeLimit && !req.Price.IsZero() {
		params["px"] = fmt.Sprintf("%v", req.Price)
	}

	return params, nil
}

// typ：1-币转张 2-张转币; symbol: 交易对; sz：数量; opTyp: open（舍位），close（四舍五入）
func (o *OkxOrderManager) convertContractCoin(typ string, symbol types.Symbol, sz string, opTyp string) (decimal.Decimal, error) {
	if symbol.Type != types.MarketTypeFuturesUSDMargined && symbol.Type != types.MarketTypePerpetualUSDMargined && symbol.Type != types.MarketTypeFuturesCoinMargined && symbol.Type != types.MarketTypePerpetualCoinMargined {
		return decimal.Zero, fmt.Errorf("invalid symbol type: %v", symbol.Type)
	}

	if symbol.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("invalid ctVal: %v", symbol.CtVal)
	}

	if opTyp == "" {
		opTyp = "open"
	}

	size, err := decimal.NewFromString(sz)
	if err != nil {
		return decimal.Zero, err
	}

	switch typ {
	case "1":
		if symbol.Type == types.MarketTypeFuturesUSDMargined || symbol.Type == types.MarketTypePerpetualUSDMargined {
			// U本位，合约面值以币计算
			size = size.Div(symbol.CtVal)
			size, err = o.sizePrecision(size, symbol, opTyp)
			if err != nil {
				return decimal.Zero, err
			}
		} else {
			// 币本位，合约面值以usdt计算
			price, err := o.getTickerPrice(symbol)
			if err != nil {
				return decimal.Zero, err
			}
			totalUsdt := size.Mul(price)
			size, err = o.sizePrecision(totalUsdt.Div(symbol.CtVal), symbol, opTyp)
			if err != nil {
				return decimal.Zero, err
			}
		}

		return size, nil
	case "2":
		if symbol.Type == types.MarketTypeFuturesUSDMargined || symbol.Type == types.MarketTypePerpetualUSDMargined {
			// U本位，合约面值以币计算
			size = size.Mul(symbol.CtVal)
			return size, nil
		} else {
			// 币本位，合约面值以usdt计算
			price, err := o.getTickerPrice(symbol)
			if err != nil {
				return decimal.Zero, err
			}
			totalUsdt := size.Mul(symbol.CtVal)
			size = totalUsdt.Div(price)
			return size, nil
		}

	}
	return decimal.Zero, fmt.Errorf("invalid type: %v", typ)
}

// size 精度处理
// TODO: 没有做stepSize的限制
func (o *OkxOrderManager) sizePrecision(size decimal.Decimal, symbol types.Symbol, opType string) (decimal.Decimal, error) {
	if symbol.MaxSize.Equal(decimal.Zero) {
		return decimal.Zero, errors.New("invalid symbol max size")
	}

	if symbol.MinSize.Equal(decimal.Zero) {
		return decimal.Zero, errors.New("invalid symbol min size")
	}

	orderQuantity := size
	if opType == "open" {
		// 向下取整到指定精度
		orderQuantity = orderQuantity.Truncate(symbol.SizePrecision)
	} else {
		// 四舍五入到指定精度
		orderQuantity = orderQuantity.Round(symbol.SizePrecision)
	}

	// 2. 限制最大值
	if orderQuantity.GreaterThan(symbol.MaxSize) {
		orderQuantity = symbol.MaxSize
	}

	// 3. 限制最小值
	if orderQuantity.LessThan(symbol.MinSize) {
		orderQuantity = symbol.MinSize
	}
	return orderQuantity, nil
}
