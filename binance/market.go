package binance

import (
	"binance_collector/binance/binance_models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Market struct {
	Url string
}

func (mr *Market) Init() *Market {
	mr.Url = "https://api1.binance.com/api/v3"
	return mr
}

func (mr *Market) GetTicker(period string) ([]binance_models.TickerResponse, error) {
	var responseBody []binance_models.TickerResponse

	url := fmt.Sprintf("%s/ticker/%s", mr.Url, period)

	if response, responseErr := new(http.Client).Get(url); responseErr != nil {
		return nil, responseErr
	} else {
		if body, bodyErr := ioutil.ReadAll(response.Body); bodyErr != nil {
			return nil, bodyErr
		} else {
			if jsonErr := json.Unmarshal(body, &responseBody); jsonErr != nil {
				return nil, jsonErr
			}
		}
	}

	return responseBody, nil
}

func (mr *Market) GetDepth(symbol string, limit int) (*binance_models.DepthResponse, error) {
	var responseBody *binance_models.DepthResponse

	url := fmt.Sprintf("%s/depth?symbol=%s&limit=%d", mr.Url, strings.ToUpper(symbol), limit)

	if response, responseErr := new(http.Client).Get(url); responseErr != nil {
		return nil, responseErr
	} else {
		if body, bodyErr := ioutil.ReadAll(response.Body); bodyErr != nil {
			return nil, bodyErr
		} else {
			if jsonErr := json.Unmarshal(body, &responseBody); jsonErr != nil {
				return nil, jsonErr
			}
		}
	}

	return responseBody, nil
}
func (mr *Market) GetKLines(symbol string, limit int, interval string) (*binance_models.DepthResponse, error) {
	var responseBody *binance_models.DepthResponse

	url := fmt.Sprintf(
		"%s/uiKlines?limit=%s&symbol=%s&interval=%s",
		mr.Url, limit, strings.ToUpper(symbol), interval,
	)

	if response, responseErr := new(http.Client).Get(url); responseErr != nil {
		return nil, responseErr
	} else {
		if body, bodyErr := ioutil.ReadAll(response.Body); bodyErr != nil {
			return nil, bodyErr
		} else {
			if jsonErr := json.Unmarshal(body, &responseBody); jsonErr != nil {
				return nil, jsonErr
			}
		}
	}

	return responseBody, nil
}
