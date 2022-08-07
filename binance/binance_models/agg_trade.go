package binance_models

type AggTradeResultData struct {
	EventTime           float64 `json:"E"`
	Ignore              bool    `json:"M"`
	TradeTime           float64 `json:"T"`
	AggregateTradeID    float64 `json:"a"`
	EventType           string  `json:"e"`
	FirstTradeID        float64 `json:"f"`
	LastTradeID         float64 `json:"l"`
	IsBuyerMarketMarker bool    `json:"m"`
	Price               string  `json:"p"`
	Quantity            string  `json:"q"`
	Symbol              string  `json:"s"`
}
