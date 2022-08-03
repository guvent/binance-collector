package binance

import (
	"binance_collector/binance/binance_models"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type BinanceSingle struct {
	Url        string
	Header     http.Header
	Connection *websocket.Conn

	DepthAsks []binance_models.DepthValue
	DepthBids []binance_models.DepthValue

	Request binance_models.WsRequest
}

func (bs *BinanceSingle) Init(symbol, depth string, limit int) *BinanceSingle {
	bs.Url = "wss://stream.binance.com/stream"
	bs.Header = nil

	bs.Request = binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: []string{
			fmt.Sprintf("%s@depth@%s", strings.ToLower(symbol), depth),
			fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol)),
			fmt.Sprintf("%s@ticker", strings.ToLower(symbol)),
		},
		Id: 1,
	}

	return bs.CollectDepth(symbol, limit)
}

func (bs *BinanceSingle) CollectDepth(symbol string, limit int) *BinanceSingle {
	if depth, err := new(Market).Init().GetDepth(symbol, limit); err != nil {
		log.Print(err)
	} else {
		for _, ask := range depth.Asks {
			bs.DepthAsks = append(bs.DepthAsks, binance_models.DepthValue{
				PriceLevel: (func(vl string) float64 {
					if v, e := strconv.ParseFloat(vl, 64); e != nil {
						return 0.0
					} else {
						return v
					}
				})(ask[0]),
				Quantity: (func(vl string) float64 {
					if v, e := strconv.ParseFloat(vl, 64); e != nil {
						return 0.0
					} else {
						return v
					}
				})(ask[1]),
			})
		}
		for _, bid := range depth.Bids {
			bs.DepthBids = append(bs.DepthBids, binance_models.DepthValue{
				PriceLevel: (func(vl string) float64 {
					if v, e := strconv.ParseFloat(vl, 64); e != nil {
						return 0.0
					} else {
						return v
					}
				})(bid[0]),
				Quantity: (func(vl string) float64 {
					if v, e := strconv.ParseFloat(vl, 64); e != nil {
						return 0.0
					} else {
						return v
					}
				})(bid[1]),
			})
		}
	}

	return bs
}

func (bs *BinanceSingle) Open() error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if conn, _, err := websocket.DefaultDialer.Dial(bs.Url, bs.Header); err != nil {
		log.Print("dial:", err)
	} else {
		//defer func(conn *websocket.Connection) {
		//	if closeErr := conn.Close(); closeErr != nil {
		//		log.Print(closeErr)
		//	}
		//}(conn)

		bs.Connection = conn
	}

	return nil
}

func (bs *BinanceSingle) Any(size int) error {
	if err := bs.Open(); err != nil {
		log.Print(err)
	}

	if writeErr := bs.Connection.WriteJSON(bs.Request); writeErr != nil {
		return writeErr
	}

	for {
		size--

		if size <= 0 {
			return nil
		}

		var mixinResult *binance_models.MixinResult

		if readErr := bs.Connection.ReadJSON(&mixinResult); readErr != nil {
			log.Print(readErr)
		} else {
			bs.Switch(mixinResult)
			//log.Print(mixinResult.Stream, " - ", time.Now().UTC())
		}
	}
}

func (bs *BinanceSingle) Switch(mixinResult *binance_models.MixinResult) *BinanceSingle {
	switch {

	case strings.Contains(mixinResult.Stream, "depth"):
		bs.Depth(*mixinResult.Data)
		//bs.Depth(binance_models.DepthResultData{
		//	EventTime:            mixinResult.Data["E"].(float64),
		//	FirstUpdateIDInEvent: mixinResult.Data["U"].(float64),
		//	AsksToBeUpdated:      mixinResult.Data["a"].([][]string),
		//	BidsToBeUpdated:      mixinResult.Data["b"].([][]string),
		//	EventType:            mixinResult.Data["e"].(string),
		//	Symbol:               mixinResult.Data["s"].(string),
		//	FinalUpdateIDInEvent: mixinResult.Data["u"].(float64),
		//})

	case strings.Contains(mixinResult.Stream, "ticker"):
		bs.Ticker(*mixinResult.Data)
		//bs.Ticker(binance_models.TickerResultData{
		//	BestAskQuantity:             mixinResult.Data["A"].(string),
		//	BestBidQuantity:             mixinResult.Data["B"].(string),
		//	StatisticsCloseTime:         mixinResult.Data["C"].(float64),
		//	EventTime:                   mixinResult.Data["E"].(float64),
		//	FirstTradeID:                mixinResult.Data["F"].(float64),
		//	LastTradeID:                 mixinResult.Data["L"].(float64),
		//	StatisticsOpenTime:          mixinResult.Data["O"].(float64),
		//	PriceChangePercent:          mixinResult.Data["P"].(string),
		//	LastQuantity:                mixinResult.Data["Q"].(string),
		//	BestAskPrice:                mixinResult.Data["a"].(string),
		//	BestBidPrice:                mixinResult.Data["b"].(string),
		//	LastPrice:                   mixinResult.Data["c"].(string),
		//	EventType:                   mixinResult.Data["e"].(string),
		//	HighPrice:                   mixinResult.Data["h"].(string),
		//	LowPrice:                    mixinResult.Data["l"].(string),
		//	TotalNumberOfTrades:         mixinResult.Data["n"].(float64),
		//	OpenPrice:                   mixinResult.Data["o"].(string),
		//	PriceChange:                 mixinResult.Data["p"].(string),
		//	TotalTradedQuoteAssetVolume: mixinResult.Data["q"].(string),
		//	Symbol:                      mixinResult.Data["s"].(string),
		//	TotalTradedBaseAssetVolume:  mixinResult.Data["v"].(string),
		//	WeightedAveragePrice:        mixinResult.Data["w"].(string),
		//	FirstTradeBefore24Hr:        mixinResult.Data["x"].(string),
		//})

	case strings.Contains(mixinResult.Stream, "aggTrade"):
		bs.AggTrade(*mixinResult.Data)
		//bs.AggTrade(binance_models.AggTradeResultData{
		//	EventTime:           mixinResult.Data["E"].(float64),
		//	Ignore:              mixinResult.Data["M"].(bool),
		//	TradeTime:           mixinResult.Data["T"].(float64),
		//	AggregateTradeID:    mixinResult.Data["a"].(float64),
		//	EventType:           mixinResult.Data["e"].(string),
		//	FirstTradeID:        mixinResult.Data["f"].(float64),
		//	LastTradeID:         mixinResult.Data["l"].(float64),
		//	IsBuyerMarketMarker: mixinResult.Data["m"].(bool),
		//	Price:               mixinResult.Data["p"].(string),
		//	Quantity:            mixinResult.Data["q"].(string),
		//	Symbol:              mixinResult.Data["s"].(string),
		//})

	}

	return bs
}

func (bs *BinanceSingle) Depth(message map[string]interface{}) *BinanceSingle {
	log.Printf("Fetched: %s", time.Now().UTC())

	for _, ask := range message["a"].([]interface{}) {
		priceLevel, _ := strconv.ParseFloat(ask.([]interface{})[0].(string), 64)
		quantity, _ := strconv.ParseFloat(ask.([]interface{})[1].(string), 64)

		bs.DepthAsks = append(bs.DepthAsks, binance_models.DepthValue{
			PriceLevel: priceLevel,
			Quantity:   quantity,
		})
		bs.DepthAsks = bs.DepthAsks[1:]
	}
	for _, bid := range message["b"].([]interface{}) {
		priceLevel, _ := strconv.ParseFloat(bid.([]interface{})[0].(string), 64)
		quantity, _ := strconv.ParseFloat(bid.([]interface{})[1].(string), 64)

		bs.DepthBids = append(bs.DepthBids, binance_models.DepthValue{
			PriceLevel: priceLevel,
			Quantity:   quantity,
		})
		bs.DepthBids = bs.DepthBids[1:]
	}

	log.Printf("Appended: %s", time.Now().UTC())

	return bs
}

func (bs *BinanceSingle) Ticker(message map[string]interface{}) *BinanceSingle {
	log.Printf("Ticker: %s", time.Now().UTC())
	return bs
}

func (bs *BinanceSingle) AggTrade(message map[string]interface{}) *BinanceSingle {
	log.Printf("AggTrade: %s", time.Now().UTC())
	return bs
}

func (bs *BinanceSingle) CloseWS() *BinanceSingle {
	if bs.Connection != nil {
		if err := bs.Connection.Close(); err != nil {
			log.Print(err)
			return nil
		}
		bs.Connection = nil
	}

	return bs
}

func (bs *BinanceSingle) GetAsks() []binance_models.DepthValue {
	return bs.DepthAsks
}

func (bs *BinanceSingle) GetBids() []binance_models.DepthValue {
	return bs.DepthBids
}
