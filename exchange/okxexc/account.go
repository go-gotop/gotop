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
	"github.com/shopspring/decimal"
)

var _ exchange.AccountManager = &OkxAccountManager{}

type OkxAccountManager struct {
	client requests.RequestClient
}

func NewOkxAccountManager() *OkxAccountManager {
	adapter := okxreq.NewOKXAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &OkxAccountManager{
		client: client,
	}
}

func (m *OkxAccountManager) GetBalances(ctx context.Context, authInfo exchange.AuthInfo) (*exchange.GetBalancesResponse, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/account/balance"

	resp, err := m.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Auth: &requests.AuthInfo{
			APIKey:     authInfo.APIKey,
			SecretKey:  authInfo.SecretKey,
			Passphrase: authInfo.Passphrase,
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
		return nil, fmt.Errorf("get balances failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var balanceResp okxBalanceResponse
	err = json.Unmarshal(body, &balanceResp)
	if err != nil {
		return nil, err
	}

	if balanceResp.Code != "0" {
		return nil, fmt.Errorf("operation failed, code: %s, message: %s", balanceResp.Code, balanceResp.Msg)
	}

	result := &exchange.GetBalancesResponse{
		Balances: make([]exchange.Balance, 0),
	}

	// 处理返回的数据
	if len(balanceResp.Data) > 0 {
		for _, data := range balanceResp.Data {
			for _, detail := range data.Details {
				availBal := decimal.Zero
				if detail.AvailBal != "" {
					availBal, err = decimal.NewFromString(detail.AvailBal)
					if err != nil {
						availBal = decimal.Zero
					}
				}

				frozenBal := decimal.Zero
				if detail.FrozenBal != "" {
					frozenBal, err = decimal.NewFromString(detail.FrozenBal)
					if err != nil {
						frozenBal = decimal.Zero
					}
				}

				// 只返回有余额的资产
				if !availBal.IsZero() || !frozenBal.IsZero() {
					result.Balances = append(result.Balances, exchange.Balance{
						Asset:     detail.Ccy,
						Available: availBal,
						Locked:    frozenBal,
					})
				}
			}
		}
	}

	return result, nil
}

func (m *OkxAccountManager) GetBalance(ctx context.Context, authInfo exchange.AuthInfo, asset string) (*exchange.GetBalanceResponse, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/account/balance"

	if asset == "" {
		return nil, errors.New("asset is required")
	}

	params := map[string]any{
		"ccy": asset,
	}

	resp, err := m.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Params: params,
		Auth: &requests.AuthInfo{
			APIKey:     authInfo.APIKey,
			SecretKey:  authInfo.SecretKey,
			Passphrase: authInfo.Passphrase,
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
		return nil, fmt.Errorf("get balance failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var balanceResp okxBalanceResponse
	err = json.Unmarshal(body, &balanceResp)
	if err != nil {
		return nil, err
	}

	if balanceResp.Code != "0" {
		return nil, fmt.Errorf("operation failed, code: %s, message: %s", balanceResp.Code, balanceResp.Msg)
	}

	if len(balanceResp.Data) == 0 || len(balanceResp.Data[0].Details) == 0 {
		return nil, fmt.Errorf("asset %s not found", asset)
	}

	// 查找指定资产
	for _, data := range balanceResp.Data {
		for _, detail := range data.Details {
			if detail.Ccy == asset {
				availBal := decimal.Zero
				if detail.AvailBal != "" {
					availBal, err = decimal.NewFromString(detail.AvailBal)
					if err != nil {
						availBal = decimal.Zero
					}
				}

				frozenBal := decimal.Zero
				if detail.FrozenBal != "" {
					frozenBal, err = decimal.NewFromString(detail.FrozenBal)
					if err != nil {
						frozenBal = decimal.Zero
					}
				}

				return &exchange.GetBalanceResponse{
					Balance: exchange.Balance{
						Asset:     detail.Ccy,
						Available: availBal,
						Locked:    frozenBal,
					},
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("asset %s not found", asset)
}
