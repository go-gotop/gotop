package file

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/go-gotop/gotop/types"
	"github.com/go-gotop/gotop/datafeed"
)

// CSVFile CSV文件逐笔数据流
type CSVFile struct {
	// dir 数据目录
	dir string
	// start 开始时间
	start int64
	// end 结束时间
	end int64
	ctx    context.Context
	cancel context.CancelFunc
	handler datafeed.TradeHandler
	errorHandler datafeed.ErrorHandler
}

// NewCSVFile 创建CSV数据流
func NewCSVFile(dir string, start, end int64) *CSVFile {
	return &CSVFile{
		dir:    dir,
		start:  start,
		end:    end,
	}
}

// Stream 流式读取CSV文件
func (f *CSVFile) Stream(ctx context.Context, tradeHandler datafeed.TradeHandler, errorHandler datafeed.ErrorHandler) error {
	f.handler = tradeHandler
	f.errorHandler = errorHandler
	f.ctx, f.cancel = context.WithCancel(ctx)
	files, err := readCSVFileNames(f.dir, f.start, f.end)
	if err != nil {
		return fmt.Errorf("read file names error: %w", err)
	}

	go f.readCSVFiles(files)
	return nil
}

// Close 关闭CSV数据流
func (f *CSVFile) Close() error {
	f.cancel()
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
				f.handler(event)
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
				if err := f.processFile(file, eventChan, f.start, f.end); err != nil {
					errorChan <- fmt.Errorf("process file %s error: %w", file, err)
					return
				}
			}
		}
	}()

	// 等待错误或处理完成
	select {
	case err := <-errorChan:
		if f.errorHandler != nil {
			f.errorHandler(err)
		}
	case <-f.ctx.Done():
		return
	}

	// 等待所有事件处理完成
	wg.Wait()
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

            
            tradeEvent, err := convertToTradeEvent(v);
            if err != nil {
                return fmt.Errorf("convert tick error: %w", err)
            }
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

	// 读取 csv 文件中的表头
	headers, err := r.Read()
	if err != nil {
		return nil, err
	}
	headers = append(headers, "ignore")

	rows := make([]*tradeData, 0, 3000)
	for {
		record, err := r.Read()
		if err != nil && record == nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		row, err := toTradeData(headers, record)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func readCSVFileNamesBackup(path string, start int64, end int64) ([]string, error) {
	// 确保路径以斜杠结尾
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// 获取目录下所有文件
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	// 过滤CSV文件并存储文件名和时间戳
	var fileNames []string
	var fileTimestamps []int64
	re := regexp.MustCompile(`^(\d+)\.csv$`)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			match := re.FindStringSubmatch(file.Name())
			if len(match) != 2 {
				// 文件名不是时间戳类型，跳过
				continue
			}
			// 从文件名解析时间戳
			timestamp, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				return nil, err
			}

			// 检查时间戳是否在指定范围内
			if (start == 0 || timestamp >= normalizeTimestamp(start)) && (end == 0 || timestamp <= normalizeTimestamp(end)) {
				fileTimestamps = append(fileTimestamps, timestamp)
			}
		}
	}

	// 按照文件名排序
	sort.Slice(fileTimestamps, func(i, j int) bool {
		return fileTimestamps[i] < fileTimestamps[j]
	})
	for _, ts := range fileTimestamps {
		fileNames = append(fileNames, strconv.FormatInt(ts, 10)+".csv")
	}
	return fileNames, nil
}

func readCSVFileNames(path string, start, end int64) ([]string, error) {
    var fileNames []string

    // 遍历目录
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
        if timestamp >= start && timestamp <= end {
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

func normalizeTimestamp(timestamp int64) int64 {
	return timestamp - timestamp%(3600*1000)
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

// isInTimeRange 检查时间戳是否在指定范围内
func isInTimeRange(timestamp, start, end int64) bool {
    inRange := (start == 0 || timestamp >= start) && (end == 0 || timestamp <= end)
    return inRange
}

// convertToTradeEvent 将 tradeData 转换为 types.TradeEvent
func convertToTradeEvent(data *tradeData) (types.TradeEvent, error) {
    price, err := decimal.NewFromString(data.Price)
    if err != nil {
        return types.TradeEvent{}, fmt.Errorf("parse price error: %w", err)
    }
    
    size, err := decimal.NewFromString(data.Size)
    if err != nil {
        return types.TradeEvent{}, fmt.Errorf("parse size error: %w", err)
    }

    return types.TradeEvent{
        Timestamp:  data.TradedAt,
        Price:      price,
        Size:       size,
        Side:       parseSide(data.Side),
    }, nil
}

// parseSide 解析交易方向
func parseSide(side string) types.SideType {
    if strings.ToUpper(side) == "SELL" {
        return types.SideTypeSell
    }
    return types.SideTypeBuy
}
