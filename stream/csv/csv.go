package file

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/shopspring/decimal"
	"github.com/go-gotop/gotop/types"
)

// CSVFileRequest 请求参数
type CSVFileRequest struct {
	// Dir 数据目录
	Dir string
	// Start 开始时间
	Start int64
	// End 结束时间
	End int64
	// Handler 处理函数，用于处理每一个TradeEvent
	Handler func(trade types.TradeEvent)
	// ErrorHandler 错误处理函数
	ErrorHandler func(err error)
	// CloseHandler 关闭处理函数
	CloseHandler func()
}

// CSVFile CSV文件逐笔数据流
type CSVFile struct {
	id       string
	ctx      context.Context
	cancel   context.CancelFunc
	request  *CSVFileRequest
	eventID  uint64  // 新增，用于自增事件ID
}

// NewCSVFile 创建CSV数据流
func NewCSVFile() *CSVFile {
	return &CSVFile{
		ctx: context.Background(),
	}
}

// ID 返回该Stream的唯一ID
func (f *CSVFile) ID() string {
	return f.id
}

// Connect 流式读取CSV文件
func (f *CSVFile) Connect(ctx context.Context, id string, request *CSVFileRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if request.Handler == nil {
		return fmt.Errorf("request.Handler cannot be nil")
	}

	f.id = id
	f.ctx, f.cancel = context.WithCancel(ctx)
	f.request = request

	files, err := readCSVFileNames(request.Dir, request.Start, request.End)
	if err != nil {
		return fmt.Errorf("read file names error: %w", err)
	}

	if len(files) == 0 {
		// 没有符合范围的文件
		return fmt.Errorf("no CSV files found in the given time range")
	}

	go f.readCSVFiles(files)
	return nil
}

// Disconnect 关闭CSV数据流
func (f *CSVFile) Disconnect() error {
	if f.cancel != nil {
		f.cancel()
	}
	if f.request != nil && f.request.CloseHandler != nil {
		f.request.CloseHandler()
	}
	return nil
}

// readCSVFiles 读取CSV文件
func (f *CSVFile) readCSVFiles(files []string) {
    eventChan := make(chan types.TradeEvent, 10)
    errorChan := make(chan error, 1)
    var wg sync.WaitGroup

    // 启动事件处理协程
    wg.Add(1)
    go func() {
        defer wg.Done()
        for event := range eventChan {
            select {
            case <-f.ctx.Done():
                return
            default:
                if f.request.Handler != nil {
                    f.request.Handler(event)
                }
            }
        }
    }()

    // 启动文件处理协程
    go func() {
        defer close(eventChan)
        defer close(errorChan)

        for _, file := range files {
            select {
            case <-f.ctx.Done():
                return
            default:
                if err := f.processFile(file, eventChan, f.request.Start, f.request.End); err != nil {
                    errorChan <- fmt.Errorf("process file %s error: %w", file, err)
                    return
                }
            }
        }

        // 正常完成时，不往errorChan写入任何错误，这样errorChan会在结束时被关闭
    }()

    // 等待错误或处理完成
    var finalErr error
    select {
    case finalErr = <-errorChan:
        // 如果有错误，此处finalErr会非空
        // 如果errorChan被关闭但未发送错误，则finalErr为nil表示正常完成
    case <-f.ctx.Done():
        // 上下文取消
    }

    // 等待所有事件处理完成
    wg.Wait()

    // 根据结果进行相应的处理
    if f.ctx.Err() != nil {
        // 上下文已取消，不调用CloseHandler，也不额外处理
        return
    }

    if finalErr != nil {
        // 有错误
        if f.request.ErrorHandler != nil {
            f.request.ErrorHandler(finalErr)
        }
    } else {
        // 无错误且上下文未取消，正常完成
        if f.request.CloseHandler != nil {
            f.request.CloseHandler()
        }
	}
}

func (f *CSVFile) processFile(filePath string, eventChan chan<- types.TradeEvent, start int64, end int64) error {
	data, err := readCSVFile(filePath)
	if err != nil {
		return fmt.Errorf("read file %s error: %w", filePath, err)
	}

	for _, v := range data {
		select {
		case <-f.ctx.Done():
			return f.ctx.Err()
		default:
			if !isInTimeRange(v.TradedAt, start, end) {
				continue
			}

			tradeEvent, err := convertToTradeEvent(v)
			if err != nil {
				return fmt.Errorf("convert tick error: %w", err)
			}

			// 为每个事件分配自增ID
			// types.TradeEvent 有一个 uint64 类型的 ID 字段
			tradeEvent.ID = atomic.AddUint64(&f.eventID, 1) - 1

			eventChan <- tradeEvent
		}
	}
	return nil
}

type tradeData struct {
	TradeID  uint64
	Size     string
	Price    string
	Side     string
	Symbol   string
	Quote    string
	TradedAt int64
}

func readCSVFile(f string) ([]*tradeData, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)

	headers, err := r.Read()
	if err != nil {
		if err == io.EOF {
			// 空文件
			return []*tradeData{}, nil
		}
		return nil, fmt.Errorf("read header error: %w", err)
	}

	// 基于期望的列名进行检查（可选）
	expectedHeaders := []string{"trade_id", "size", "price", "side", "quote", "traded_at"}
	if !validateHeaders(headers, expectedHeaders) {
		return nil, fmt.Errorf("invalid or missing headers, got: %v, expected at least: %v", headers, expectedHeaders)
	}

	rows := make([]*tradeData, 0, 3000)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read record error: %w", err)
		}

		row, err := toTradeData(headers, record)
		if err != nil {
			return nil, fmt.Errorf("convert record error: %w", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func readCSVFileNames(path string, start, end int64) ([]string, error) {
	var fileNames []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查是否是 CSV 文件
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".csv") {
			return nil
		}

		// 匹配文件名格式
		re := regexp.MustCompile(`^(\d+)\.csv$`)
		match := re.FindStringSubmatch(info.Name())
		if len(match) != 2 {
			return nil
		}

		// 从文件名解析时间戳
		timestamp, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil
		}

		// 检查时间戳是否在范围内
		if (start == 0 || timestamp >= start) && (end == 0 || timestamp <= end) {
			fileNames = append(fileNames, filePath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk directory error: %w", err)
	}

	// 按文件名排序
	sort.Strings(fileNames)
	return fileNames, nil
}

func toTradeData(headers []string, record []string) (*tradeData, error) {
	row := &tradeData{}
	for i, value := range record {
		switch headers[i] {
		case "trade_id":
			tradeID, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse trade_id error: %w", err)
			}
			row.TradeID = tradeID
		case "size":
			row.Size = value
		case "price":
			row.Price = value
		case "side":
			row.Side = value
		case "quote":
			row.Quote = value
		case "traded_at":
			tradedAt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse traded_at error: %w", err)
			}
			row.TradedAt = tradedAt
		}
	}
	return row, nil
}

func isInTimeRange(timestamp, start, end int64) bool {
	return (start == 0 || timestamp >= start) && (end == 0 || timestamp <= end)
}

// convertToTradeEvent 将 tradeData 转换为 types.TradeEvent
func convertToTradeEvent(data *tradeData) (types.TradeEvent, error) {
	price, err := decimal.NewFromString(data.Price)
	if err != nil {
		return types.TradeEvent{}, fmt.Errorf("parse price error: %w", err)
	}
	if price.IsNegative() {
		return types.TradeEvent{}, fmt.Errorf("invalid price: %s", price.String())
	}

	size, err := decimal.NewFromString(data.Size)
	if err != nil {
		return types.TradeEvent{}, fmt.Errorf("parse size error: %w", err)
	}
	if size.IsNegative() {
		return types.TradeEvent{}, fmt.Errorf("invalid size: %s", size.String())
	}

	return types.TradeEvent{
		// 将ID赋值由调用者统一完成
		Timestamp: data.TradedAt,
		Price:     price,
		Size:      size,
		Side:      parseSide(data.Side),
	}, nil
}

// parseSide 解析交易方向
func parseSide(side string) types.SideType {
	switch strings.ToUpper(side) {
	case "SELL":
		return types.SideTypeSell
	case "BUY":
		return types.SideTypeBuy
	default:
		// 如果有其他需要处理的逻辑，请添加
		return types.SideTypeBuy
	}
}

// validateHeaders 校验CSV表头是否包含必需字段
func validateHeaders(gotHeaders, required []string) bool {
	headerMap := make(map[string]bool)
	for _, h := range gotHeaders {
		headerMap[strings.ToLower(h)] = true
	}
	for _, req := range required {
		if !headerMap[strings.ToLower(req)] {
			return false
		}
	}
	return true
}
