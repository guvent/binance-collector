package binance

import (
	"binance_collector/binance/binance_models"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Stream struct {
	Url    string
	Header http.Header
}

func (bs *Stream) Init() *Stream {
	bs.Url = "wss://stream.binance.com/stream"
	bs.Header = nil

	return bs
}

func (bs *Stream) Connect(symbols []string, callback func(result interface{})) error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if c, _, err := websocket.DefaultDialer.Dial(bs.Url, bs.Header); err != nil {
		log.Fatal("dial:", err)
	} else {
		defer func(c *websocket.Conn) {
			if closeErr := c.Close(); closeErr != nil {
				log.Fatal(closeErr)
			}
		}(c)

		req := binance_models.WsRequest{
			Method: "SUBSCRIBE",
			Params: symbols,
			Id:     1,
		}

		if writeErr := c.WriteJSON(req); writeErr != nil {
			return writeErr
		}

		var first *binance_models.FirstInfo

		for {
			var message interface{}

			if first == nil {
				if readErr := c.ReadJSON(&first); readErr != nil {
					return readErr
				}
			}

			if readErr := c.ReadJSON(&message); readErr != nil {
				return readErr
			}

			callback(message)
		}
	}

	return nil
}

func (bs *Stream) Ticker(symbols []string, callback func(result binance_models.TickerResult)) error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if c, _, err := websocket.DefaultDialer.Dial(bs.Url, bs.Header); err != nil {
		log.Fatal("dial:", err)
	} else {
		defer func(c *websocket.Conn) {
			if closeErr := c.Close(); closeErr != nil {
				log.Fatal(closeErr)
			}
		}(c)

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

		if writeErr := c.WriteJSON(req); writeErr != nil {
			return writeErr
		}

		var first *binance_models.FirstInfo

		for {
			var message binance_models.TickerResult

			if first == nil {
				if readErr := c.ReadJSON(&first); readErr != nil {
					return readErr
				}
			}

			if readErr := c.ReadJSON(&message); readErr != nil {
				return readErr
			}

			callback(message)
		}
	}

	return nil
}

func (bs *Stream) AggTrade(symbols []string, callback func(result binance_models.AggTradeResult)) error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if c, _, err := websocket.DefaultDialer.Dial(bs.Url, bs.Header); err != nil {
		log.Fatal("dial:", err)
	} else {
		defer func(c *websocket.Conn) {
			if closeErr := c.Close(); closeErr != nil {
				log.Fatal(closeErr)
			}
		}(c)

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

		if writeErr := c.WriteJSON(req); writeErr != nil {
			return writeErr
		}

		var first *binance_models.FirstInfo

		for {
			var message binance_models.AggTradeResult

			if first == nil {
				if readErr := c.ReadJSON(&first); readErr != nil {
					return readErr
				}
			}

			if readErr := c.ReadJSON(&message); readErr != nil {
				return readErr
			}

			callback(message)
		}
	}

	return nil
}

func (bs *Stream) Depth(symbols []string, callback func(result binance_models.DepthResult)) error {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if c, _, err := websocket.DefaultDialer.Dial(bs.Url, bs.Header); err != nil {
		log.Fatal("dial:", err)
	} else {
		defer func(c *websocket.Conn) {
			if closeErr := c.Close(); closeErr != nil {
				log.Fatal(closeErr)
			}
		}(c)

		req := binance_models.WsRequest{
			Method: "SUBSCRIBE",
			Params: (func(symbols []string) []string {
				var strings []string

				for _, symbol := range symbols {
					strings = append(strings, string(symbol)+"@depth@1000ms")
				}

				return strings
			})(symbols),
			Id: 1,
		}

		if writeErr := c.WriteJSON(req); writeErr != nil {
			return writeErr
		}

		var first *binance_models.FirstInfo

		for {
			var message binance_models.DepthResult

			if first == nil {
				if readErr := c.ReadJSON(&first); readErr != nil {
					return readErr
				}
			}

			if readErr := c.ReadJSON(&message); readErr != nil {
				return readErr
			}

			callback(message)
		}
	}

	return nil
}
