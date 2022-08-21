package binance

import (
	"binance_collector/binance/binance_models"
	"binance_collector/calculator"
	"binance_collector/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type BinanceCalc struct {
	url        string
	header     http.Header
	connection *websocket.Conn

	wsRequest binance_models.WsRequest

	calculate *calculator.Calculate
}

func (bin *BinanceCalc) Init(env *utils.GetEnvironment) *BinanceCalc {
	bin.url = "wss://stream.binance.com/stream"
	bin.header = nil

	var params []string
	var symbols = strings.Split(env.Symbols, ",")

	for _, symbol := range symbols {
		//params = append(params, fmt.Sprintf("%s@depth@%sms", strings.ToLower(symbol), env.DepthMS))
		//params = append(params, fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol)))
		//params = append(params, fmt.Sprintf("%s@ticker", strings.ToLower(symbol)))
		params = append(params, fmt.Sprintf("%s@kline_4h", strings.ToLower(symbol)))
	}

	bin.wsRequest = binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: params,
		Id:     1,
	}

	bin.calculate = new(calculator.Calculate).Init(symbols)

	return bin
}

func (bin *BinanceCalc) Start() error {
	if err := bin.OpenWs(); err != nil {
		log.Print(err)
	}

	if writeErr := bin.connection.WriteJSON(bin.wsRequest); writeErr != nil {
		return writeErr
	}

	for {
		var mixinResult *binance_models.MixinResult

		if readErr := bin.connection.ReadJSON(&mixinResult); readErr != nil {
			log.Print(readErr)
		} else {
			go bin.collect(mixinResult)
		}

		mixinResult = nil
	}
}

func (bin *BinanceCalc) OpenWs() error {
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

func (bin *BinanceCalc) CloseWs() *BinanceCalc {
	if bin.connection != nil {
		if err := bin.connection.Close(); err != nil {
			log.Print(err)
			return nil
		}
		bin.connection = nil
	}

	return bin
}

func (bin *BinanceCalc) collect(mixinResult *binance_models.MixinResult) *BinanceCalc {
	mutex.Lock()

	switch {

	case strings.Contains(mixinResult.Stream, "depth"):
		(func(message binance_models.DepthResultData) {
			for _, ask := range message.AsksToBeUpdated {
				priceLevel, _ := strconv.ParseFloat(ask.([]interface{})[0].(string), 64)
				quantity, _ := strconv.ParseFloat(ask.([]interface{})[1].(string), 64)

				bin.calculate.Ask(binance_models.DepthValue{
					PriceLevel:           priceLevel,
					Quantity:             quantity,
					EventTime:            message.EventTime,
					FirstUpdateIDInEvent: message.FirstUpdateIDInEvent,
					FinalUpdateIDInEvent: message.FinalUpdateIDInEvent,
					IsStream:             true,
				}, message.Symbol)
			}

			for _, bid := range message.BidsToBeUpdated {
				priceLevel, _ := strconv.ParseFloat(bid.([]interface{})[0].(string), 64)
				quantity, _ := strconv.ParseFloat(bid.([]interface{})[1].(string), 64)

				bin.calculate.Bid(binance_models.DepthValue{
					PriceLevel:           priceLevel,
					Quantity:             quantity,
					EventTime:            message.EventTime,
					FirstUpdateIDInEvent: message.FirstUpdateIDInEvent,
					FinalUpdateIDInEvent: message.FinalUpdateIDInEvent,
					IsStream:             true,
				}, message.Symbol)
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
		bin.calculate.Tick(binance_models.TickerResultData{
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
		bin.calculate.AggTrade(binance_models.AggTradeResultData{
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

	case strings.Contains(mixinResult.Stream, "kline"):
		kLineData := (*mixinResult.Data)["k"].(map[string]interface{})

		bin.calculate.KLine(binance_models.KlineData{
			StartTime:          (kLineData)["t"].(float64),
			CloseTime:          (kLineData)["T"].(float64),
			Symbol:             (kLineData)["s"].(string),
			Interval:           (kLineData)["i"].(string),
			FirstNAVUpdateTime: (kLineData)["f"].(float64),
			LastNAVUpdateTime:  (kLineData)["L"].(float64),
			OpenNAVPrice:       (kLineData)["o"].(string),
			CloseNAVPrice:      (kLineData)["c"].(string),
			HighestNAVPrice:    (kLineData)["h"].(string),
			LowestNAVPrice:     (kLineData)["l"].(string),
			RealLeverage:       (kLineData)["v"].(string),
			NumberOfNAVUpdate:  (kLineData)["n"].(float64),
			Ignore1:            (kLineData)["x"].(bool),
			Ignore2:            (kLineData)["q"].(string),
			Ignore3:            (kLineData)["V"].(string),
			Ignore4:            (kLineData)["Q"].(string),
			Ignore5:            (kLineData)["B"].(string),
		})
	}

	mutex.Unlock()

	return bin
}
