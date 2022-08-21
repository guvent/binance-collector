package binance_models

type KLineResponse struct {
	Value []string
}

type KLineResponseValue struct {
	Value []KLineResponseData
}

type KLineResponseData struct {
	OpenTime                 int64
	Open                     string
	High                     string
	Low                      string
	Close                    string
	Volume                   string
	CloseTime                int64
	QuoteAssetVolume         string
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
	Ignore                   string
}

type KlineResultData struct {
	EventTime float64   `json:"E"`
	EventType string    `json:"e"`
	Symbol    string    `json:"s"`
	Data      KlineData `json:"k"`
}

type KlineData struct {
	StartTime          float64 `json:"t"`
	CloseTime          float64 `json:"T"`
	Symbol             string  `json:"s"`
	Interval           string  `json:"i"`
	FirstNAVUpdateTime float64 `json:"f"`
	LastNAVUpdateTime  float64 `json:"L"`
	OpenNAVPrice       string  `json:"o"`
	CloseNAVPrice      string  `json:"c"`
	HighestNAVPrice    string  `json:"h"`
	LowestNAVPrice     string  `json:"l"`
	RealLeverage       string  `json:"v"`
	NumberOfNAVUpdate  float64 `json:"n"`
	Ignore1            bool    `json:"x"`
	Ignore2            string  `json:"q"`
	Ignore3            string  `json:"V"`
	Ignore4            string  `json:"Q"`
	Ignore5            string  `json:"B"`
}
