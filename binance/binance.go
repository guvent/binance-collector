package binance

import (
	"binance_collector/binance/binance_models"
	"log"
)

type Binance struct {
	Market *Market
	Stream *Stream

	Asks []binance_models.DepthValue
	Bids []binance_models.DepthValue

	Symbol string
}

func (b *Binance) Init(symbol string) *Binance {
	b.Symbol = symbol

	b.Stream = new(Stream).Init()
	b.Market = new(Market).Init()

	return b
}

func (b *Binance) CollectDepth() *Binance {
	if depth, err := b.Market.GetDepth(b.Symbol, 5000); err != nil {
		log.Print(err)
	} else {
		for _, ask := range depth.Asks {
			b.Asks = append(b.Asks, binance_models.DepthValue{
				PriceLevel: ask[0],
				Quantity:   ask[1],
			})
		}
		for _, bid := range depth.Bids {
			b.Bids = append(b.Bids, binance_models.DepthValue{
				PriceLevel: bid[0],
				Quantity:   bid[1],
			})
		}
	}

	return b
}

func (b *Binance) FetchDepth(max int) *Binance {
	if err := b.Stream.Open(); err != nil {
		log.Print(err)
	}

	var loop = 0
	if err := b.Stream.Depth(
		[]string{b.Symbol}, "1000ms",
		func(message binance_models.DepthResult) {
			loop++

			if loop >= max {
				b.Stream.Close()
			}

			for _, ask := range message.Data.AsksToBeUpdated {
				b.Asks = append(b.Asks, binance_models.DepthValue{
					PriceLevel: ask[0],
					Quantity:   ask[1],
				})
			}
			for _, bid := range message.Data.BidsToBeUpdated {
				b.Bids = append(b.Bids, binance_models.DepthValue{
					PriceLevel: bid[0],
					Quantity:   bid[1],
				})
			}
		},
	); err != nil {
		log.Print(err)
	}

	return b
}

func (b *Binance) GetAsks() []binance_models.DepthValue {
	return b.Asks
}

func (b *Binance) GetBids() []binance_models.DepthValue {
	return b.Bids
}
