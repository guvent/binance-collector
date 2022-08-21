package calculator

import (
	"binance_collector/binance/binance_models"
	"fmt"
	"log"
	"strings"
	"time"
)

type Buff struct {
	Price    float64
	Quantity float64
	Total    float64
}

type Calculate struct {
	Asks Buff
	Bids Buff

	AskLog string
	BidLog string

	AskCount float64

	DepthAsks map[string][]binance_models.DepthValue
	DepthBids map[string][]binance_models.DepthValue

	Tickers   map[string][]binance_models.TickerResultData
	AggTrades map[string][]binance_models.AggTradeResultData

	KLines map[string][]binance_models.KlineData
}

func (calculate *Calculate) Init(symbols []string) *Calculate {
	calculate.Asks = Buff{}
	calculate.Bids = Buff{}

	calculate.AskCount = 0

	calculate.AskLog = ""
	calculate.BidLog = ""

	calculate.DepthAsks = map[string][]binance_models.DepthValue{}
	calculate.DepthBids = map[string][]binance_models.DepthValue{}
	calculate.Tickers = map[string][]binance_models.TickerResultData{}
	calculate.AggTrades = map[string][]binance_models.AggTradeResultData{}
	calculate.KLines = map[string][]binance_models.KlineData{}

	for _, symbol := range symbols {
		calculate.DepthAsks[strings.ToUpper(symbol)] = []binance_models.DepthValue{{}}
		calculate.DepthBids[strings.ToUpper(symbol)] = []binance_models.DepthValue{{}}
		calculate.Tickers[strings.ToUpper(symbol)] = []binance_models.TickerResultData{{}}
		calculate.AggTrades[strings.ToUpper(symbol)] = []binance_models.AggTradeResultData{{}}
		calculate.KLines[strings.ToUpper(symbol)] = []binance_models.KlineData{{}}
	}

	time.AfterFunc(1*time.Second, func() {
		calculate.Summary()
	})

	return calculate
}

func (calculate *Calculate) Pick(rotation string, value binance_models.DepthValue) {
	rot := fmt.Sprintf("BALANCED (%f)", calculate.AskCount)

	if calculate.AskCount > 0 {
		rot = fmt.Sprintf("TOO ASK (%f)", calculate.AskCount)
	}
	if calculate.AskCount < 0 {
		rot = fmt.Sprintf("TOO BID (%f)", -1*calculate.AskCount)
	}

	calculate.AskLog = fmt.Sprintf("*** ASK: %f + %f = %f  ??  (%s) %f - %f  ||  %s",
		calculate.Asks.Price, calculate.Asks.Quantity, calculate.Asks.Total,
		rotation, value.PriceLevel, value.Quantity, rot,
	)

	calculate.BidLog = fmt.Sprintf("*** BID: %f + %f = %f  ??  (%s) %f - %f  ||  %s",
		calculate.Bids.Price, calculate.Bids.Quantity, calculate.Bids.Total,
		rotation, value.PriceLevel, value.Quantity, rot,
	)
}

func (calculate *Calculate) Ask(value binance_models.DepthValue, symbol string) {
	calculate.Asks = Buff{
		Price:    calculate.Asks.Price + value.PriceLevel,
		Quantity: calculate.Asks.Quantity + value.Quantity,
		Total:    calculate.Asks.Total + value.PriceLevel*value.Quantity,
	}

	calculate.AskCount++
	calculate.Pick(">>> ASK", value)

	calculate.DepthAsks[symbol] = append(calculate.DepthAsks[symbol], value)
}

func (calculate *Calculate) Bid(value binance_models.DepthValue, symbol string) {
	calculate.Bids = Buff{
		Price:    calculate.Bids.Price + value.PriceLevel,
		Quantity: calculate.Bids.Quantity + value.Quantity,
		Total:    calculate.Bids.Total + value.PriceLevel*value.Quantity,
	}

	calculate.AskCount--
	calculate.Pick("<<< BID", value)

	calculate.DepthBids[symbol] = append(calculate.DepthBids[symbol], value)
}

func (calculate *Calculate) Tick(value binance_models.TickerResultData) {
	calculate.Tickers[value.Symbol] = append(calculate.Tickers[value.Symbol], value)
}

func (calculate *Calculate) AggTrade(value binance_models.AggTradeResultData) {
	calculate.AggTrades[value.Symbol] = append(calculate.AggTrades[value.Symbol], value)
}

func (calculate *Calculate) KLine(value binance_models.KlineData) {
	calculate.KLines[value.Symbol] = append(calculate.KLines[value.Symbol], value)
}

func (calculate *Calculate) Summary() {
	log.Println(calculate.AskLog)
	log.Println(calculate.BidLog)

	time.AfterFunc(1*time.Second, func() {
		calculate.Summary()
	})
}
