package binance

import (
	"binance_collector/binance/binance_models"
	"binance_collector/database"
	"binance_collector/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

type Binance struct {
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

func (bin *Binance) Init(env *utils.GetEnvironment) *Binance {
	bin.url = "wss://stream.binance.com/stream"
	bin.header = nil

	var params []string

	for _, symbol := range strings.Split(env.Symbols, ",") {
		params = append(params, fmt.Sprintf("%s@depth@%sms", strings.ToLower(symbol), env.DepthMS))
		params = append(params, fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol)))
		params = append(params, fmt.Sprintf("%s@ticker", strings.ToLower(symbol)))
	}

	bin.wsRequest = binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: params,
		Id:     1,
	}

	bin.Influx = new(database.Influx).Init(env.ServerURL, env.AuthToken).Select("buysell", env.Bucket)

	return bin
}

func (bin *Binance) Start(commitMs string) error {
	if err := bin.OpenWs(); err != nil {
		log.Print(err)
	}

	if writeErr := bin.connection.WriteJSON(bin.wsRequest); writeErr != nil {
		return writeErr
	}

	bin.commit(time.ParseDuration(commitMs))

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

func (bin *Binance) OpenWs() error {
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

func (bin *Binance) choose(mixinResult *binance_models.MixinResult) *Binance {
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

func (bin *Binance) commit(duration time.Duration, err error) *Binance {
	bin.Influx.Commit()

	if err != nil {
		log.Fatal("Commit Duration Invalid!")
	}

	time.AfterFunc(time.Millisecond*duration, func() {
		bin.commit(duration, err)
	})

	return bin
}
