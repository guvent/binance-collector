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

type Binance struct {
	url        string
	header     http.Header
	connection *websocket.Conn

	wsRequest binance_models.WsRequest
	timeoutMs int

	DepthAsks []binance_models.DepthValue
	DepthBids []binance_models.DepthValue

	Tickers   []binance_models.TickerResultData
	AggTrades []binance_models.AggTradeResultData

	Symbol  string
	Limit   int
	DepthMs int

	depthCollected bool
}

func (bin *Binance) Init(symbol string, limit int, depthMs int, ticker, aggTrade bool) *Binance {
	bin.url = "wss://stream.binance.com/stream"
	bin.header = nil
	bin.depthCollected = false

	bin.Symbol = symbol
	bin.Limit = limit
	bin.DepthMs = depthMs

	var params []string

	if depthMs >= 100 {
		params = append(params, fmt.Sprintf("%s@depth@%dms", strings.ToLower(symbol), depthMs))
	}

	if aggTrade {
		params = append(params, fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol)))
	}
	if ticker {
		params = append(params, fmt.Sprintf("%s@ticker", strings.ToLower(symbol)))
	}

	bin.wsRequest = binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: params,
		Id:     1,
	}

	bin.timeoutMs = 0

	return bin
}

func (bin *Binance) SetTimeOut(timeoutMs int) *Binance {
	bin.timeoutMs = timeoutMs

	return bin
}

func (bin *Binance) Start(durationSec int) error {
	if err := bin.OpenWs(); err != nil {
		log.Print(err)
	}

	if writeErr := bin.connection.WriteJSON(bin.wsRequest); writeErr != nil {
		return writeErr
	}

	go func() {
		for {
			var mixinResult *binance_models.MixinResult

			if bin.timeoutMs > 0 {
				if setTimeoutErr := bin.connection.SetReadDeadline(time.Now().Add(
					time.Duration(bin.timeoutMs) * time.Millisecond),
				); setTimeoutErr != nil {
					log.Print(setTimeoutErr)
				}
			}

			if readErr := bin.connection.ReadJSON(&mixinResult); readErr != nil {
				log.Print(readErr)
			} else {
				go bin.choose(mixinResult)
			}

			mixinResult = nil
		}
	}()

	if bin.DepthMs >= 100 {
		if collectErr := bin.collectDepth(bin.Symbol, bin.Limit); collectErr != nil {
			return collectErr
		}
	}

	time.Sleep(time.Duration(durationSec) * time.Second)

	return nil
}

func (bin *Binance) OpenWs() error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if conn, _, err := websocket.DefaultDialer.Dial(bin.url, bin.header); err != nil {
		log.Print("dial:", err)
	} else {
		//defer func(conn *websocket.Connection) {
		//	if closeErr := conn.Close(); closeErr != nil {
		//		log.Print(closeErr)
		//	}
		//}(conn)

		bin.connection = conn
	}

	return nil
}

func (bin *Binance) CloseWs() *Binance {
	if bin.connection != nil {
		if err := bin.connection.Close(); err != nil {
			log.Print(err)
			return nil
		}
		bin.connection = nil
	}

	return bin
}

func (bin *Binance) collectDepth(symbol string, limit int) error {
	time.Sleep(time.Duration(bin.DepthMs) * time.Millisecond)

	if depth, err := new(Market).Init().GetDepth(symbol, limit); err != nil {
		return err
	} else {
		// TODO: don't worked this...
		//
		//for i, depthAsk := range bin.DepthAsks {
		//	if depth.LastUpdateId < depthAsk.FinalUpdateIDInEvent {
		//		bin.DepthAsks = bin.DepthAsks[i:] // +1
		//		break
		//	}
		//}
		//for i, depthBid := range bin.DepthBids {
		//	if depth.LastUpdateId < depthBid.FinalUpdateIDInEvent {
		//		bin.DepthBids = bin.DepthBids[i:] // +1
		//		break
		//	}
		//}
		//

		// TODO: Put one wrong item! But don't worry, sometimes do it.
		var newDepthAsks []binance_models.DepthValue
		var newDepthBids []binance_models.DepthValue

		for _, depthAsk := range bin.DepthAsks {
			if depth.LastUpdateId+1 <= depthAsk.FirstUpdateIDInEvent {
				newDepthAsks = append(newDepthAsks, depthAsk)
			}
		}

		for _, depthBid := range bin.DepthBids {
			if depth.LastUpdateId+1 <= depthBid.FirstUpdateIDInEvent {
				newDepthBids = append(newDepthBids, depthBid)
			}
		}

		bin.DepthAsks = newDepthAsks
		bin.DepthBids = newDepthBids
		// -- Put one wrong item! But don't worry, sometimes do it.

		for _, ask := range depth.Asks {
			bin.DepthAsks = append([]binance_models.DepthValue{
				{
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
					FirstUpdateIDInEvent: depth.LastUpdateId,
					FinalUpdateIDInEvent: depth.LastUpdateId,
					IsStream:             false,
				},
			}, bin.DepthAsks...)
		}
		for _, bid := range depth.Bids {
			bin.DepthBids = append([]binance_models.DepthValue{
				{
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
					FirstUpdateIDInEvent: depth.LastUpdateId,
					FinalUpdateIDInEvent: depth.LastUpdateId,
					IsStream:             false,
				},
			}, bin.DepthBids...)
		}

		bin.depthCollected = true
	}

	return nil
}

func (bin *Binance) choose(mixinResult *binance_models.MixinResult) *Binance {
	switch {

	case strings.Contains(mixinResult.Stream, "depth"):
		bin.depth(*mixinResult.Data)
		//bin.depth(binance_models.DepthResultData{
		//	EventTime:            mixinResult.Data["E"].(float64),
		//	FirstUpdateIDInEvent: mixinResult.Data["U"].(float64),
		//	AsksToBeUpdated:      mixinResult.Data["a"].([][]string),
		//	BidsToBeUpdated:      mixinResult.Data["b"].([][]string),
		//	EventType:            mixinResult.Data["e"].(string),
		//	Symbol:               mixinResult.Data["s"].(string),
		//	FinalUpdateIDInEvent: mixinResult.Data["u"].(float64),
		//})

	case strings.Contains(mixinResult.Stream, "ticker"):
		bin.ticker(binance_models.TickerResultData{
			BestAskQuantity:             (*mixinResult.Data)["A"].(string),
			BestBidQuantity:             (*mixinResult.Data)["B"].(string),
			StatisticsCloseTime:         (*mixinResult.Data)["C"].(float64),
			EventTime:                   (*mixinResult.Data)["E"].(float64),
			FirstTradeID:                (*mixinResult.Data)["F"].(float64),
			LastTradeID:                 (*mixinResult.Data)["L"].(float64),
			StatisticsOpenTime:          (*mixinResult.Data)["O"].(float64),
			PriceChangePercent:          (*mixinResult.Data)["P"].(string),
			LastQuantity:                (*mixinResult.Data)["Q"].(string),
			BestAskPrice:                (*mixinResult.Data)["a"].(string),
			BestBidPrice:                (*mixinResult.Data)["b"].(string),
			LastPrice:                   (*mixinResult.Data)["c"].(string),
			EventType:                   (*mixinResult.Data)["e"].(string),
			HighPrice:                   (*mixinResult.Data)["h"].(string),
			LowPrice:                    (*mixinResult.Data)["l"].(string),
			TotalNumberOfTrades:         (*mixinResult.Data)["n"].(float64),
			OpenPrice:                   (*mixinResult.Data)["o"].(string),
			PriceChange:                 (*mixinResult.Data)["p"].(string),
			TotalTradedQuoteAssetVolume: (*mixinResult.Data)["q"].(string),
			Symbol:                      (*mixinResult.Data)["s"].(string),
			TotalTradedBaseAssetVolume:  (*mixinResult.Data)["v"].(string),
			WeightedAveragePrice:        (*mixinResult.Data)["w"].(string),
			FirstTradeBefore24Hr:        (*mixinResult.Data)["x"].(string),
		})

	case strings.Contains(mixinResult.Stream, "aggTrade"):
		bin.aggTrade(binance_models.AggTradeResultData{
			EventTime:           (*mixinResult.Data)["E"].(float64),
			Ignore:              (*mixinResult.Data)["M"].(bool),
			TradeTime:           (*mixinResult.Data)["T"].(float64),
			AggregateTradeID:    (*mixinResult.Data)["a"].(float64),
			EventType:           (*mixinResult.Data)["e"].(string),
			FirstTradeID:        (*mixinResult.Data)["f"].(float64),
			LastTradeID:         (*mixinResult.Data)["l"].(float64),
			IsBuyerMarketMarker: (*mixinResult.Data)["m"].(bool),
			Price:               (*mixinResult.Data)["p"].(string),
			Quantity:            (*mixinResult.Data)["q"].(string),
			Symbol:              (*mixinResult.Data)["s"].(string),
		})

	}

	return bin
}

func (bin *Binance) depth(message map[string]interface{}) *Binance {
	log.Printf("Depth Fetched: %s", time.Now().UTC())

	for _, ask := range message["a"].([]interface{}) {
		priceLevel, _ := strconv.ParseFloat(ask.([]interface{})[0].(string), 64)
		quantity, _ := strconv.ParseFloat(ask.([]interface{})[1].(string), 64)

		bin.DepthAsks = append(bin.DepthAsks, binance_models.DepthValue{
			PriceLevel:           priceLevel,
			Quantity:             quantity,
			FirstUpdateIDInEvent: message["U"].(float64),
			FinalUpdateIDInEvent: message["u"].(float64),
			IsStream:             true,
		})

		if bin.depthCollected {
			bin.DepthAsks = bin.DepthAsks[1:]
		}
	}

	for _, bid := range message["b"].([]interface{}) {
		priceLevel, _ := strconv.ParseFloat(bid.([]interface{})[0].(string), 64)
		quantity, _ := strconv.ParseFloat(bid.([]interface{})[1].(string), 64)

		bin.DepthBids = append(bin.DepthBids, binance_models.DepthValue{
			PriceLevel:           priceLevel,
			Quantity:             quantity,
			FirstUpdateIDInEvent: message["U"].(float64),
			FinalUpdateIDInEvent: message["u"].(float64),
			IsStream:             true,
		})

		if bin.depthCollected {
			bin.DepthBids = bin.DepthBids[1:]
		}
	}

	log.Printf("Depth Appended: %s", time.Now().UTC())

	return bin
}

func (bin *Binance) ticker(message binance_models.TickerResultData) *Binance {
	log.Printf("Ticker Fetched: %s", time.Now().UTC())

	bin.Tickers = append(bin.Tickers, message)

	log.Printf("Ticker Appended: %s", time.Now().UTC())

	return bin
}

func (bin *Binance) aggTrade(message binance_models.AggTradeResultData) *Binance {
	log.Printf("AggTrade Fetched: %s", time.Now().UTC())

	bin.AggTrades = append(bin.AggTrades, message)

	log.Printf("AggTrade Appended: %s", time.Now().UTC())

	return bin
}
