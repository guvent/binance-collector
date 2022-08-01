package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Market struct {
	Url string
}

func (m *Market) Init() *Market {
	m.Url = "https://api1.binance.com/api/v3"

	return m
}

func (m *Market) GetTicker(period string) (interface{}, error) {
	var responseBody interface{}

	client := new(http.Client)
	url := fmt.Sprintf("%s/ticker/%s", m.Url, period)

	if response, responseErr := client.Get(url); responseErr != nil {
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

func (m *Market) GetDepth(symbol string) (interface{}, error) {
	var responseBody interface{}

	client := new(http.Client)
	url := fmt.Sprintf("%s/depth?symbol=%s&limit=5000", m.Url, symbol)

	if response, responseErr := client.Get(url); responseErr != nil {
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