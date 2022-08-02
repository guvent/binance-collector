package binance

import (
	"binance_collector/binance/binance_models"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Stream struct {
	Url        string
	Header     http.Header
	Connection *websocket.Conn
}

func (bs *Stream) Init() *Stream {
	bs.Url = "wss://stream.binance.com/stream"
	bs.Header = nil

	return bs
}

func (bs *Stream) Open() error {
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

func (bs *Stream) Any(symbols []string, callback func(result interface{})) error {
	if bs.Connection == nil {
		return errors.New("connection closed")
	}

	req := binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: symbols,
		Id:     1,
	}

	if writeErr := bs.Connection.WriteJSON(req); writeErr != nil {
		return writeErr
	}

	var first *binance_models.FirstInfo

	for {
		var message interface{}

		if first == nil {
			if readErr := bs.Connection.ReadJSON(&first); readErr != nil {
				return readErr
			}
		}

		if readErr := bs.Connection.ReadJSON(&message); readErr != nil {
			return readErr
		}

		callback(message)
	}
}

func (bs *Stream) Ticker(symbols []string, callback func(result binance_models.TickerResult)) error {
	if bs.Connection == nil {
		return errors.New("connection closed")
	}

	req := binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: (func(symbols []string) []string {
			var strings []string

			for _, symbol := range symbols {
				strings = append(strings, string(symbol)+"@ticker")
			}

			return strings
		})(symbols),
		Id: 1,
	}

	if writeErr := bs.Connection.WriteJSON(req); writeErr != nil {
		return writeErr
	}

	var first *binance_models.FirstInfo

	for {
		if bs.Connection == nil {
			break
		}

		var message binance_models.TickerResult

		if first == nil {
			if readErr := bs.Connection.ReadJSON(&first); readErr != nil {
				return readErr
			}
		}

		if readErr := bs.Connection.ReadJSON(&message); readErr != nil {
			return readErr
		}

		callback(message)
	}

	return nil
}

func (bs *Stream) AggTrade(symbols []string, callback func(result binance_models.AggTradeResult)) error {
	if bs.Connection == nil {
		return errors.New("connection closed")
	}

	req := binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: (func(symbols []string) []string {
			var strings []string

			for _, symbol := range symbols {
				strings = append(strings, string(symbol)+"@aggTrade")
			}

			return strings
		})(symbols),
		Id: 1,
	}

	if writeErr := bs.Connection.WriteJSON(req); writeErr != nil {
		return writeErr
	}

	var first *binance_models.FirstInfo

	for {
		if bs.Connection == nil {
			break
		}

		var message binance_models.AggTradeResult

		if first == nil {
			if readErr := bs.Connection.ReadJSON(&first); readErr != nil {
				return readErr
			}
		}

		if readErr := bs.Connection.ReadJSON(&message); readErr != nil {
			return readErr
		}

		callback(message)
	}

	return nil
}

func (bs *Stream) Depth(symbols []string, retention string, callback func(result binance_models.DepthResult)) error {
	if bs.Connection == nil {
		return errors.New("connection closed")
	}

	req := binance_models.WsRequest{
		Method: "SUBSCRIBE",
		Params: (func(symbols []string) []string {
			var strings []string

			for _, symbol := range symbols {
				strings = append(strings, fmt.Sprintf("%s@depth@%s", symbol, retention))
			}

			return strings
		})(symbols),
		Id: 1,
	}

	if writeErr := bs.Connection.WriteJSON(req); writeErr != nil {
		return writeErr
	}

	var first *binance_models.FirstInfo

	for {
		if bs.Connection == nil {
			break
		}

		var message binance_models.DepthResult

		if first == nil {
			if readErr := bs.Connection.ReadJSON(&first); readErr != nil {
				return readErr
			}
		}

		if readErr := bs.Connection.ReadJSON(&message); readErr != nil {
			return readErr
		}

		callback(message)
	}

	return nil
}

func (bs *Stream) Close() *Stream {
	if bs.Connection != nil {
		if err := bs.Connection.Close(); err != nil {
			log.Print(err)
			return nil
		}
		bs.Connection = nil
	}

	return bs
}
