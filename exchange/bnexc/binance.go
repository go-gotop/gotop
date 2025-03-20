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

// bnCapitalRecoveryResponse 币安资金归集响应
type bnCapitalRecoveryResponse struct {
	// 资产名
	Coin string `json:"coin"`
	// 是否可以充值
	DepositAllEnable bool `json:"depositAllEnable"`
	// 可用余额
	Free string `json:"free"`
	// 冻结余额
	Freeze string `json:"freeze"`
	// 可申购余额
	Ipoable string `json:"ipoable"`
	// 可申购余额
	Ipoing string `json:"ipoing"`
	// 是否是合法货币
	IsLegalMoney bool `json:"isLegalMoney"`
	// 锁定余额
	Locked string `json:"locked"`
	// 资产名
	Name string `json:"name"`
	// 网络列表
	NetworkList []struct {
		// 地址正则
		AddressRegex string `json:"addressRegex"`
		// 资产名
		Coin string `json:"coin"`
		// 充值描述
		DepositDesc string `json:"depositDesc"`
		// 是否可以充值
		DepositEnable bool `json:"depositEnable"`
		// 是否是默认网络
		IsDefault bool `json:"isDefault"`
		// 备注正则
		MemoRegex string `json:"memoRegex"`
		// 最小确认数
		MinConfirm int `json:"minConfirm"`
		// 资产名
		Name string `json:"name"`
		// 网络
		Network string `json:"network"`
		// 特殊提示
		SpecialTips string `json:"specialTips"`
		// 解锁需要的确认数
		UnLockConfirm int `json:"unLockConfirm"`
		// 提现描述
		WithdrawDesc string `json:"withdrawDesc"`
		// 是否可以提现
		WithdrawEnable bool `json:"withdrawEnable"`
		// 提现手续费
		WithdrawFee string `json:"withdrawFee"`
		// 提现最小数量
		WithdrawMin string `json:"withdrawMin"`
		// 提现最大数量
		WithdrawMax string `json:"withdrawMax"`
		// 内部转账最小提现数
		WithdrawInternalMin string `json:"withdrawInternalMin"`
		// 是否需要memo
		SameAddress bool `json:"sameAddress"`
		// 预计到达时间
		EstimatedArrivalTime int `json:"estimatedArrivalTime"`
		// 是否繁忙
		Busy bool `json:"busy"`
		// 合约地址URL
		ContractAddressUrl string `json:"contractAddressUrl"`
		// 合约地址
		ContractAddress string `json:"contractAddress"`
		// 面额
		Denomination int `json:"denomination"`
	} `json:"networkList"`
}
