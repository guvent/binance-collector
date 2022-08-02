package binance_models

type DepthResult struct {
	Data struct {
		EventTime            int64      `json:"E"`
		FirstUpdateIDInEvent int64      `json:"U"`
		AsksToBeUpdated      [][]string `json:"a"` // 0: Price level to be updated, 1: Quantity
		BidsToBeUpdated      [][]string `json:"b"` // 0: Price level to be updated, 1: Quantity
		EventType            string     `json:"e"`
		Symbol               string     `json:"s"`
		FinalUpdateIDInEvent int64      `json:"u"`
	} `json:"data"`
	Stream string `json:"stream"`
}

type DepthResponse struct {
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
	LastUpdateId int64      `json:"lastUpdateId"`
}

type DepthValue struct {
	PriceLevel string
	Quantity   string
}
