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

func (bs *Stream) Connect(symbols []string, callback func(interface{})) error {
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

		for {
			var message interface{}

			if readErr := c.ReadJSON(&message); readErr != nil {
				return readErr
			}

			callback(message)
		}
	}

	return nil
}
