package binance_models

type DepthResult struct {
	Data   DepthResultData `json:"data"`
	Stream string          `json:"stream"`
}

type DepthResultData struct {
	EventTime            float64    `json:"E"`
	FirstUpdateIDInEvent float64    `json:"U"`
	AsksToBeUpdated      [][]string `json:"a"` // 0: Price level to be updated, 1: Quantity
	BidsToBeUpdated      [][]string `json:"b"` // 0: Price level to be updated, 1: Quantity
	EventType            string     `json:"e"`
	Symbol               string     `json:"s"`
	FinalUpdateIDInEvent float64    `json:"u"`
}

type DepthResponse struct {
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
	LastUpdateId float64    `json:"lastUpdateId"`
}

type DepthValue struct {
	PriceLevel           float64
	Quantity             float64
	FirstUpdateIDInEvent float64
	FinalUpdateIDInEvent float64
	IsStream             bool
}
