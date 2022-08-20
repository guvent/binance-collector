package binance

import (
	"binance_collector/binance/binance_models"
	"binance_collector/database"
	"binance_collector/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BinanceBoth struct {
	url        string
	header     http.Header
	connection *websocket.Conn

	DepthAsks map[string][]binance_models.DepthValue
	DepthBids map[string][]binance_models.DepthValue

	Tickers   map[string][]binance_models.TickerResultData
	AggTrades map[string][]binance_models.AggTradeResultData

	wsRequest binance_models.WsRequest

	Influx *database.Influx
}

func (bin *BinanceBoth) Init(env *utils.GetEnvironment) *BinanceBoth {
	bin.url = "wss://stream.binance.com/stream"
	bin.header = nil

	var params []string

	bin.DepthAsks = map[string][]binance_models.DepthValue{}
	bin.DepthBids = map[string][]binance_models.DepthValue{}
	bin.Tickers = map[string][]binance_models.TickerResultData{}
	bin.AggTrades = map[string][]binance_models.AggTradeResultData{}

	for _, symbol := range strings.Split(env.Symbols, ",") {

		bin.DepthAsks[strings.ToUpper(symbol)] = []binance_models.DepthValue{{}}
		bin.DepthBids[strings.ToUpper(symbol)] = []binance_models.DepthValue{{}}
		bin.Tickers[strings.ToUpper(symbol)] = []binance_models.TickerResultData{{}}
		bin.AggTrades[strings.ToUpper(symbol)] = []binance_models.AggTradeResultData{{}}

		params = append(params, fmt.Sprintf("%s@depth@%sms", strings.ToLower(symbol), env.DepthMS))
		params = append(params, fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol)))
		params = append(params, fmt.Sprintf("%s@ticker", strings.ToLower(symbol)))
	}

	bin.wsRequest = binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: params,
		Id:     1,
	}

	bin.Influx = new(database.Influx).Init(env.ServerURL, env.AuthToken)
	bin.Influx = bin.Influx.WriteInit(env.Org, env.Bucket)
	bin.Influx = bin.Influx.QueryInit(env.Org)

	return bin
}

func (bin *BinanceBoth) QueryTest() {
	bin.Influx.ReadData()
}

func (bin *BinanceBoth) Start(commitMs string) error {
	if err := bin.OpenWs(); err != nil {
		log.Print(err)
	}

	if writeErr := bin.connection.WriteJSON(bin.wsRequest); writeErr != nil {
		return writeErr
	}

	bin.commit(time.ParseDuration(commitMs))

	bin.summary()

	for {
		var mixinResult *binance_models.MixinResult

		if readErr := bin.connection.ReadJSON(&mixinResult); readErr != nil {
			log.Print(readErr)
		} else {
			go bin.choose(mixinResult)
		}

		mixinResult = nil
	}
}

func (bin *BinanceBoth) OpenWs() error {
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

func (bin *BinanceBoth) CloseWs() *BinanceBoth {
	if bin.connection != nil {
		if err := bin.connection.Close(); err != nil {
			log.Print(err)
			return nil
		}
		bin.connection = nil
	}

	return bin
}

func (bin *BinanceBoth) choose(mixinResult *binance_models.MixinResult) *BinanceBoth {
	mutex.Lock()

	switch {

	case strings.Contains(mixinResult.Stream, "depth"):
		go bin.Influx.WriteDepth(*mixinResult.Data)

		(func(message binance_models.DepthResultData) {
			for _, ask := range message.AsksToBeUpdated {
				priceLevel, _ := strconv.ParseFloat(ask.([]interface{})[0].(string), 64)
				quantity, _ := strconv.ParseFloat(ask.([]interface{})[1].(string), 64)

				bin.DepthAsks[message.Symbol] = append(bin.DepthAsks[message.Symbol], binance_models.DepthValue{
					PriceLevel:           priceLevel,
					Quantity:             quantity,
					EventTime:            message.EventTime,
					FirstUpdateIDInEvent: message.FirstUpdateIDInEvent,
					FinalUpdateIDInEvent: message.FinalUpdateIDInEvent,
					IsStream:             true,
				})
			}

			for _, bid := range message.BidsToBeUpdated {
				priceLevel, _ := strconv.ParseFloat(bid.([]interface{})[0].(string), 64)
				quantity, _ := strconv.ParseFloat(bid.([]interface{})[1].(string), 64)

				bin.DepthBids[message.Symbol] = append(bin.DepthBids[message.Symbol], binance_models.DepthValue{
					PriceLevel:           priceLevel,
					Quantity:             quantity,
					EventTime:            message.EventTime,
					FirstUpdateIDInEvent: message.FirstUpdateIDInEvent,
					FinalUpdateIDInEvent: message.FinalUpdateIDInEvent,
					IsStream:             true,
				})
			}

		})(binance_models.DepthResultData{
			EventTime:            (*mixinResult.Data)["E"].(float64),
			FirstUpdateIDInEvent: (*mixinResult.Data)["U"].(float64),
			AsksToBeUpdated:      (*mixinResult.Data)["a"].([]interface{}),
			BidsToBeUpdated:      (*mixinResult.Data)["b"].([]interface{}),
			EventType:            (*mixinResult.Data)["e"].(string),
			Symbol:               (*mixinResult.Data)["s"].(string),
			FinalUpdateIDInEvent: (*mixinResult.Data)["u"].(float64),
		})

	case strings.Contains(mixinResult.Stream, "ticker"):
		go bin.Influx.WriteTicker(*mixinResult.Data)

		(func(message binance_models.TickerResultData) {
			bin.Tickers[message.Symbol] = append(bin.Tickers[message.Symbol], message)
		})(binance_models.TickerResultData{
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
		go bin.Influx.WriteAggTrade(*mixinResult.Data)

		(func(message binance_models.AggTradeResultData) {
			bin.AggTrades[message.Symbol] = append(bin.AggTrades[message.Symbol], message)
		})(binance_models.AggTradeResultData{
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

	mutex.Unlock()

	return bin
}

func (bin *BinanceBoth) commit(duration time.Duration, err error) *BinanceBoth {
	bin.Influx.Commit()

	if err != nil {
		log.Fatal("Commit Duration Invalid!")
	}

	time.AfterFunc(time.Millisecond*duration, func() {
		bin.commit(duration, err)
	})

	return bin
}

func (bin *BinanceBoth) summary() *BinanceBoth {
	mutex.Lock()

	askSummary := ""

	for symbol, depthValues := range bin.DepthAsks {
		askSummary = fmt.Sprintf("%s | %s: %d", askSummary, symbol, len(depthValues))
	}

	log.Printf("%s |", askSummary)

	mutex.Unlock()

	time.AfterFunc(10*time.Second, func() {
		bin.summary()
	})

	return bin
}
