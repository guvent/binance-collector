package binance_models

type MixinResult struct {
	Stream string                  `json:"stream"`
	Data   *map[string]interface{} `json:"data"`
	Result interface{}             `json:"result"`
	Id     float64                 `json:"id"`
}

type MixinResultData struct {
	EventTime int64  `json:"E"`
	EventType string `json:"e"`
	Symbol    string `json:"s"`

	FirstUpdateIDInEvent int64      `json:"U"`
	AsksToBeUpdated      [][]string `json:"a"` // 0: Price level to be updated, 1: Quantity
	BidsToBeUpdated      [][]string `json:"b"` // 0: Price level to be updated, 1: Quantity
	FinalUpdateIDInEvent int64      `json:"u"`

	BestAskQuantity     string `json:"A"`
	BestBidQuantity     string `json:"B"`
	StatisticsCloseTime int64  `json:"C"`
	//FirstTradeID                int    `json:"F"`
	//LastTradeID                 int    `json:"L"`
	StatisticsOpenTime int64  `json:"O"`
	PriceChangePercent string `json:"P"`
	LastQuantity       string `json:"Q"`
	//BestAskPrice                string `json:"a"`
	//BestBidPrice                string `json:"b"`
	LastPrice                   string `json:"c"`
	HighPrice                   string `json:"h"`
	LowPrice                    string `json:"l"`
	TotalNumberOfTrades         int    `json:"n"`
	OpenPrice                   string `json:"o"`
	PriceChange                 string `json:"p"`
	TotalTradedQuoteAssetVolume string `json:"q"`
	TotalTradedBaseAssetVolume  string `json:"v"`
	WeightedAveragePrice        string `json:"w"`
	FirstTradeBefore24Hr        string `json:"x"`

	Ignore    bool  `json:"M"`
	TradeTime int64 `json:"T"`
	//AggregateTradeID    int    `json:"a"`
	FirstTradeID int `json:"f"`
	//LastTradeID         int    `json:"l"`
	IsBuyerMarketMarker bool `json:"m"`
	//Price               string `json:"p"`
	//Quantity            string `json:"q"`
}
