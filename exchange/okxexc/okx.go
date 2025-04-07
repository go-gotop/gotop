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

// okx 账户余额
type okxBalanceResponse struct {
	Code string `json:"code"`
	Data []struct {
		// 美金层面有效保证金，适用于现货模式/跨币种保证金模式/组合保证金模式
		AdjEq string `json:"adjEq"`
		// 账户美金层面潜在借币占用保证金，仅适用于现货模式/跨币种保证金模式/组合保证金模式。在其他账户模式下为""
		BorrowFroz string `json:"borrowFroz"`
		// 各币种资产详细信息
		Details []struct {
			// 可用余额
			AvailBal string `json:"availBal"`
			// 可用保证金，适用于现货和合约模式/跨币种保证金模式/组合保证金模式
			AvailEq string `json:"availEq"`
			// 币种美金层面潜在借币占用保证金，仅适用于现货模式/跨币种保证金模式/组合保证金模式。在其他账户模式下为""
			BorrowFroz string `json:"borrowFroz"`
			// 币种余额
			CashBal string `json:"cashBal"`
			// 币种
			Ccy string `json:"ccy"`
			// 币种全仓负债额，适用于现货模式/跨币种保证金模式/组合保证金模式
			CrossLiab string `json:"crossLiab"`
			// true：质押币，false：非质押币，适用于跨币种保证金模式
			CollateralEnabled bool `json:"collateralEnabled"`
			// 美金层面币种折算权益，适用于现货模式(开通了借币功能)/跨币种保证金模式/组合保证金模式
			DisEq string `json:"disEq"`
			// 币种总权益
			Eq string `json:"eq"`
			// 币种权益美金价值
			EqUsd string `json:"eqUsd"`
			// 合约智能跟单权益，默认为0，仅适用于跟单人
			SmtSyncEq string `json:"smtSyncEq"`
			// 现货智能跟单权益，默认为0，仅适用于跟单人
			SpotCopyTradingEq string `json:"spotCopyTradingEq"`
			// 抄底宝、逃顶宝功能的币种冻结金额
			FixedBal string `json:"fixedBal"`
			// 币种占用金额
			FrozenBal string `json:"frozenBal"`
			// 币种维度全仓占用保证金，适用于现货和合约模式且有全仓仓位时
			Imr string `json:"imr"`
			// 计息，应扣未扣利息，值为正数，如 9.01，适用于现货模式/跨币种保证金模式/组合保证金模式
			Interest string `json:"interest"`
			// 币种逐仓仓位权益，适用于现货和合约模式/跨币种保证金模式/组合保证金模式
			IsoEq string `json:"isoEq"`
			// 币种逐仓负债额，适用于跨币种保证金模式/组合保证金模式
			IsoLiab string `json:"isoLiab"`
			// 逐仓未实现盈亏，适用于现货和合约模式/跨币种保证金模式/组合保证金模式
			IsoUpl string `json:"isoUpl"`
			// 币种负债额，值为正数，如 "21625.64"，适用于现货模式/跨币种保证金模式/组合保证金模式
			Liab string `json:"liab"`
			// 币种最大可借，适用于现货模式/跨币种保证金模式/组合保证金模式的全仓
			MaxLoan string `json:"maxLoan"`
			// 币种全仓保证金率，衡量账户内某项资产风险的指标，适用于现货和合约模式且有全仓仓位时
			MgnRatio string `json:"mgnRatio"`
			// 币种维度全仓维持保证金，适用于现货和合约模式且有全仓仓位时
			Mmr string `json:"mmr"`
			// 币种杠杆倍数，适用于现货和合约模式
			NotionalLever string `json:"notionalLever"`
			// 挂单冻结数量，适用于现货模式/现货和合约模式/跨币种保证金模式
			OrdFrozen string `json:"ordFrozen"`
			// 体验金余额
			RewardBal string `json:"rewardBal"`
			// 现货对冲占用数量，适用于组合保证金模式
			SpotInUseAmt string `json:"spotInUseAmt"`
			// 系统计算得到的最大可能现货占用数量，适用于组合保证金模式
			MaxSpotInUse string `json:"maxSpotInUse"`
			// 现货逐仓余额，仅适用于现货带单/跟单，适用于现货模式/现货和合约模式
			SpotIsoBal string `json:"spotIsoBal"`
			// 策略权益
			StgyEq string `json:"stgyEq"`
			// 当前负债币种触发系统自动换币的风险，0、1、2、3、4、5其中之一，数字越大代表您的负债币种触发自动换币概率越高，适用于现货模式/跨币种保证金模式/组合保证金模式
			Twap string `json:"twap"`
			// 币种余额信息的更新时间，Unix时间戳的毫秒数格式，如 1597026383085
			Utime string `json:"uTime"`
			// 未实现盈亏，适用于现货和合约模式/跨币种保证金模式/组合保证金模式
			Upl string `json:"upl"`
			// 由于仓位未实现亏损导致的负债，适用于跨币种保证金模式/组合保证金模式
			UplLiab string `json:"uplLiab"`
			// 现货余额，单位为币种，比如 BTC
			SpotBal string `json:"spotBal"`
			// 现货开仓成本价，单位 USD
			OpenAvgPx string `json:"openAvgPx"`
			// 现货累计成本价，单位 USD
			AccAvgPx string `json:"accAvgPx"`
			// 现货未实现收益，单位 USD
			SpotUpl string `json:"spotUpl"`
			// 现货未实现收益率
			SpotUplRatio string `json:"spotUplRatio"`
			// 现货累计收益，单位 USD
			TotalPnl string `json:"totalPnl"`
			// 现货累计收益率
			TotalPnlRatio string `json:"totalPnlRatio"`
		} `json:"details"`
		// 美金层面占用保证金，适用于现货模式/跨币种保证金模式/组合保证金模式
		Imr string `json:"imr"`
		// 美金层面逐仓仓位权益，适用于现货和合约模式/跨币种保证金模式/组合保证金模式
		IsoEq string `json:"isoEq"`
		// 美金层面保证金率，适用于现货模式/跨币种保证金模式/组合保证金模式
		MgnRatio string `json:"mgnRatio"`
		// 美金层面维持保证金，适用于现货模式/跨币种保证金模式/组合保证金模式
		Mmr string `json:"mmr"`
		// 以美金价值为单位的持仓数量，即仓位美金价值，适用于现货模式/跨币种保证金模式/组合保证金模式
		NotionalUsd string `json:"notionalUsd"`
		// 借币金额（美元价值），适用于现货模式/跨币种保证金模式/组合保证金模式
		NotionalUsdForBorrow string `json:"notionalUsdForBorrow"`
		// 永续合约持仓美元价值，适用于跨币种保证金模式/组合保证金模式
		NotionalUsdForSwap string `json:"notionalUsdForSwap"`
		// 交割合约持仓美元价值，适用于跨币种保证金模式/组合保证金模式
		NotionalUsdForFutures string `json:"notionalUsdForFutures"`
		// 期权持仓美元价值，适用于现货模式/跨币种保证金模式/组合保证金模式
		NotionalUsdForOption string `json:"notionalUsdForOption"`
		// 美金层面全仓挂单占用保证金，仅适用于现货模式/跨币种保证金模式/组合保证金模式
		OrdFroz string `json:"ordFroz"`
		// 美金层面权益
		TotalEq string `json:"totalEq"`
		// 账户信息的更新时间，Unix时间戳的毫秒数格式，如 1597026383085
		Utime string `json:"uTime"`
		// 账户层面全仓未实现盈亏（美元单位），适用于跨币种保证金模式/组合保证金模式
		Upl string `json:"upl"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// okx k线响应
type okxKlineResponse struct {
	Code string     `json:"code"`
	Data [][]string `json:"data"`
	Msg  string     `json:"msg"`
}
