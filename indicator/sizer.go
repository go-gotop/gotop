package indicator

import "github.com/shopspring/decimal"

// PositionSizer 定义了一个用于计算仓位大小的接口。实现该接口的类型应根据
// 特定的业务逻辑（例如风险控制、资金管理策略、市场条件）来计算建议持有的头寸数量。
// 
// 通常情况下，不同的实现将基于一系列参数（如当前价格、账户余额、最大风险敞口、
// 标的物流动性和交易费用等）来确定适宜的买入或卖出数量。
// 
// 实现该接口时应确保：
// 1. CalculatePositionSize 返回的值为一个 decimal.Decimal 类型，确保在处理价格、
//    数量等时不会因浮点数精度问题导致偏差。
// 2. 当无法正确计算仓位大小（例如缺乏必要数据或出现业务逻辑异常）时，应返回非 nil 的 error。
// 3. 当计算成功时，error 应为 nil，返回值应为建议的仓位大小（可为 0 或正数，取决于策略）。
type PositionSizer interface {
    // CalculatePositionSize 根据实现的策略和条件计算仓位大小。
    //
    // 返回值：
    // - decimal.Decimal：建议的仓位大小，通常为正数（表示应开仓的数量），
    //   也可能为 0（表示不增仓）或其它值（根据具体策略定义）。
    // - error：当无法正确计算时返回错误；如果计算成功则返回 nil。
    CalculatePositionSize() (decimal.Decimal, error)
}
