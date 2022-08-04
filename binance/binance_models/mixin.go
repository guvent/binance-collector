package binance_models

type MixinResult struct {
	Stream string                  `json:"stream"`
	Data   *map[string]interface{} `json:"data"`
	Result interface{}             `json:"result"`
	Id     float64                 `json:"id"`
}
