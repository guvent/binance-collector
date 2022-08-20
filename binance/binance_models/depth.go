package binance_models

type DepthResultData struct {
	EventTime            float64       `json:"E"`
	FirstUpdateIDInEvent float64       `json:"U"`
	AsksToBeUpdated      []interface{} `json:"a"` // 0: Price level to be updated, 1: Quantity
	BidsToBeUpdated      []interface{} `json:"b"` // 0: Price level to be updated, 1: Quantity
	EventType            string        `json:"e"`
	Symbol               string        `json:"s"`
	FinalUpdateIDInEvent float64       `json:"u"`
}

type DepthResponse struct {
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
	LastUpdateId float64    `json:"lastUpdateId"`
}

type DepthValue struct {
	PriceLevel           float64
	Quantity             float64
	EventTime            float64
	FirstUpdateIDInEvent float64
	FinalUpdateIDInEvent float64
	IsStream             bool
}
