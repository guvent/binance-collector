package binance

import (
	"binance_collector/binance/binance_models"
	"binance_collector/database"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

type BinanceDatabase struct {
	url        string
	header     http.Header
	connection *websocket.Conn

	wsRequest binance_models.WsRequest

	DepthAsks []binance_models.DepthValue
	DepthBids []binance_models.DepthValue

	Tickers   []binance_models.TickerResultData
	AggTrades []binance_models.AggTradeResultData

	Influx *database.Influx
}

func (bin *BinanceDatabase) Init(
	symbols []string, depthMs int,
	serverURL, authToken, bucket string,
) *BinanceDatabase {
	bin.url = "wss://stream.binance.com/stream"
	bin.header = nil

	var params []string

	for _, symbol := range symbols {
		params = append(params, fmt.Sprintf("%s@depth@%dms", strings.ToLower(symbol), depthMs))
		params = append(params, fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol)))
		params = append(params, fmt.Sprintf("%s@ticker", strings.ToLower(symbol)))
	}

	bin.wsRequest = binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: params,
		Id:     1,
	}

	bin.Influx = new(database.Influx).Init(serverURL, authToken).Select("buysell", bucket)

	return bin
}

func (bin *BinanceDatabase) Start() error {
	if err := bin.OpenWs(); err != nil {
		log.Print(err)
	}

	if writeErr := bin.connection.WriteJSON(bin.wsRequest); writeErr != nil {
		return writeErr
	}

	bin.commit()

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

func (bin *BinanceDatabase) OpenWs() error {
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

func (bin *BinanceDatabase) CloseWs() *BinanceDatabase {
	if bin.connection != nil {
		if err := bin.connection.Close(); err != nil {
			log.Print(err)
			return nil
		}
		bin.connection = nil
	}

	return bin
}

func (bin *BinanceDatabase) choose(mixinResult *binance_models.MixinResult) *BinanceDatabase {
	switch {

	case strings.Contains(mixinResult.Stream, "depth"):
		bin.Influx.WriteDepth(*mixinResult.Data)
	case strings.Contains(mixinResult.Stream, "ticker"):
		bin.Influx.WriteTicker(*mixinResult.Data)
	case strings.Contains(mixinResult.Stream, "aggTrade"):
		bin.Influx.WriteAggTrade(*mixinResult.Data)

	}

	return bin
}

func (bin *BinanceDatabase) commit() *BinanceDatabase {
	bin.Influx.Commit()

	time.AfterFunc(time.Millisecond*600, func() {
		bin.commit()
	})

	log.Print("committed....")

	return bin
}
