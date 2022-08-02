package binance_models

type TickerResult struct {
	Data struct {
		BestAskQuantity             string `json:"A"`
		BestBidQuantity             string `json:"B"`
		StatisticsCloseTime         int64  `json:"C"`
		EventTime                   int64  `json:"E"`
		FirstTradeID                int    `json:"F"`
		LastTradeID                 int    `json:"L"`
		StatisticsOpenTime          int64  `json:"O"`
		PriceChangePercent          string `json:"P"`
		LastQuantity                string `json:"Q"`
		BestAskPrice                string `json:"a"`
		BestBidPrice                string `json:"b"`
		LastPrice                   string `json:"c"`
		EventType                   string `json:"e"`
		HighPrice                   string `json:"h"`
		LowPrice                    string `json:"l"`
		TotalNumberOfTrades         int    `json:"n"`
		OpenPrice                   string `json:"o"`
		PriceChange                 string `json:"p"`
		TotalTradedQuoteAssetVolume string `json:"q"`
		Symbol                      string `json:"s"`
		TotalTradedBaseAssetVolume  string `json:"v"`
		WeightedAveragePrice        string `json:"w"`
		FirstTradeBefore24Hr        string `json:"x"`
	} `json:"data"`
	Stream string `json:"stream"`
}

type TickerResponse struct {
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	CloseTime          int64  `json:"closeTime"`
	Count              int    `json:"count"`
	FirstId            int    `json:"firstId"`
	HighPrice          string `json:"highPrice"`
	LastId             int    `json:"lastId"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	LowPrice           string `json:"lowPrice"`
	OpenPrice          string `json:"openPrice"`
	OpenTime           int64  `json:"openTime"`
	PrevClosePrice     string `json:"prevClosePrice"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	QuoteVolume        string `json:"quoteVolume"`
	Symbol             string `json:"symbol"`
	Volume             string `json:"volume"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
}