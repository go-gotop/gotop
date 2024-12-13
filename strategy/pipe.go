package strategy

// Pipe[I] 表示一个计算管道节点，内部持有一个延迟计算的闭包 evaluate。
// 当最终调用 Value() 时才会执行 evaluate 函数来计算出类型为 I 的结果。
// I：该管道节点的输出类型。
type Pipe[I any] struct {
    // evaluate 是延迟执行的计算函数，在最终调用 Value() 时才运行。
    // 返回 (I, error):
    // - I：计算结果值
    // - error：若执行中出错，返回错误并终止后续计算。
    evaluate func() (I, error)
}

// Value() 触发 Pipe 中的计算过程，返回最终结果 I 和可能的 error。
func (p Pipe[I]) Value() (I, error) {
    return p.evaluate()
}

// Then 在当前管道节点的基础上添加下一个步骤函数。
// I: 输入类型，O: 输出类型
// p: 当前管道节点
// f: 转换函数，将类型 I 转换为 O
// 返回新的管道节点 Pipe[O]
func Then[I, O any](p Pipe[I], f func(I) (O, error)) Pipe[O] {
    return Pipe[O]{
        evaluate: func() (O, error) {
            // 先执行上游节点
            iv, err := p.evaluate()
            if err != nil {
                var zero O
                return zero, err
            }
            // 再执行当前步骤函数 f
            return f(iv)
        },
    }
}

// Start[I] 用于构建管道的起点，以一个已知的初始值 val 开始。
// 返回 Pipe[I]，此时 Pipeline 的输出是确定的常量值。
func Start[I any](val I) Pipe[I] {
    return Pipe[I]{
        evaluate: func() (I, error) {
            // 初始节点直接返回固定值，不会出错
            return val, nil
        },
    }
}