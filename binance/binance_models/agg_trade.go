package binance_models

type AggTradeResult struct {
	Data struct {
		EventTime           int64  `json:"E"`
		Ignore              bool   `json:"M"`
		TradeTime           int64  `json:"T"`
		AggregateTradeID    int    `json:"a"`
		EventType           string `json:"e"`
		FirstTradeID        int    `json:"f"`
		LastTradeID         int    `json:"l"`
		IsBuyerMarketMarker bool   `json:"m"`
		Price               string `json:"p"`
		Quantity            string `json:"q"`
		Symbol              string `json:"s"`
	} `json:"data"`
	Stream string `json:"stream"`
}
