package database

import (
	"binance_collector/utils"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Influx struct {
	Client   influxdb2.Client
	WriteApi api.WriteAPI
	QueryApi api.QueryAPI
}

func (influx *Influx) Init(serverURL, authToken string) *Influx {
	influx.Client = influxdb2.NewClientWithOptions(
		serverURL, authToken,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)

	if ping, err := influx.Client.Ping(context.Background()); err != nil {
		log.Print(err)
	} else {
		log.Printf("Influx Connected : %v", ping)
	}

	return influx
}

func (influx *Influx) QueryInit(org string) *Influx {

	influx.QueryApi = influx.Client.QueryAPI(org)

	return influx
}

func (influx *Influx) ReadData() *Influx {
	/*
		from(bucket: "buysell")
		  |> range(start: -1m)
		  |> filter(fn: (r) => r["_measurement"] == "ATOMUSDT@depth")
		  |> filter(fn: (r) => r["rotation"] == "bid")

		  |> filter( fn: (r) => r["_field"] == "price" or r["_field"] == "quantity" )

		  |> group(columns: ["_value", "_start", "_stop", "_field", "rotation", "symbol"])
		  |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
		  |> map(fn: (r) => ({ r with _value: r.price * r.quantity }))
		  |> sort(columns: ["_value"], desc: true)
		  |> limit(n:10)
	*/

	queryBuild := new(utils.InfluxQueryBuilder)
	queryBuild = queryBuild.From("buysell")
	queryBuild = queryBuild.Range("-1m", "")
	queryBuild = queryBuild.Filter(utils.KeyValueItem{Key: "_measurement", Value: "ATOMUSDT@depth"})
	queryBuild = queryBuild.Filter(utils.KeyValueItem{Key: "rotation", Value: "bid"})

	queryBuild = queryBuild.Filters([]utils.KeyValueItem{
		{Key: "_field", Value: "price"},
		{Key: "_field", Value: "quantity"},
	}, "or")

	queryBuild = queryBuild.Group([]string{"_value", "_start", "_stop", "_field", "rotation", "symbol"})
	queryBuild = queryBuild.Pivot("_time", "_field", "_value")
	queryBuild = queryBuild.Map("{ r with _value: r.price * r.quantity }")
	queryBuild = queryBuild.Sort([]string{"_value"}, true)
	queryBuild = queryBuild.Limit(10)

	if err := queryBuild.Execute(influx.QueryApi, func(result *api.QueryTableResult) {
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		fmt.Printf("value: %v\n", result.Record().Values())
	}); err != nil {
		log.Fatal(err)
	}

	return influx
}

func (influx *Influx) WriteInit(org, bucket string) *Influx {

	influx.WriteApi = influx.Client.WriteAPI(org, bucket)

	return influx
}

func (influx *Influx) WriteData(measurement string) *Influx {

	for i := 0; i < 100; i++ {
		// create point
		p := influxdb2.NewPoint(
			measurement,
			map[string]string{
				"id":       fmt.Sprintf("rack_%v", i%10),
				"vendor":   "AWS",
				"hostname": fmt.Sprintf("host_%v", i%100),
			},
			map[string]interface{}{
				"temperature": rand.Float64() * 80.0,
				"disk_free":   rand.Float64() * 1000.0,
				"disk_total":  (i/10 + 1) * 1000000,
				"mem_total":   (i/100 + 1) * 10000000,
				"mem_free":    rand.Uint64(),
			},
			time.Now())
		// write asynchronously
		influx.WriteApi.WritePoint(p)
	}

	return influx
}

func (influx *Influx) WriteDepth(message map[string]interface{}) *Influx {

	for _, ask := range message["a"].([]interface{}) {
		influx.WriteApi.WritePoint(
			influxdb2.NewPoint(
				fmt.Sprintf("%s@depth", message["s"].(string)),
				map[string]string{
					"rotation":              "ask",
					"event_type":            message["e"].(string),
					"symbol":                message["s"].(string),
					"event_time":            time.UnixMilli(int64(message["E"].(float64))).Format(time.RFC3339Nano),
					"final_update_id_event": fmt.Sprintf("%f", message["u"].(float64)),
					"first_update_id_event": fmt.Sprintf("%f", message["u"].(float64)),
				},
				map[string]interface{}{
					"price": (func(vl string) float64 {
						if v, e := strconv.ParseFloat(vl, 64); e != nil {
							return 0.0
						} else {
							return v
						}
					})(ask.([]interface{})[0].(string)),
					"quantity": (func(vl string) float64 {
						if v, e := strconv.ParseFloat(vl, 64); e != nil {
							return 0.0
						} else {
							return v
						}
					})(ask.([]interface{})[1].(string)),
				},
				time.Now(),
			),
		)
	}

	for _, bid := range message["b"].([]interface{}) {
		influx.WriteApi.WritePoint(
			influxdb2.NewPoint(
				fmt.Sprintf("%s@depth", message["s"].(string)),
				map[string]string{
					"rotation":              "bid",
					"event_type":            message["e"].(string),
					"symbol":                message["s"].(string),
					"event_time":            time.UnixMilli(int64(message["E"].(float64))).Format(time.RFC3339Nano),
					"final_update_id_event": fmt.Sprintf("%f", message["u"].(float64)),
					"first_update_id_event": fmt.Sprintf("%f", message["u"].(float64)),
				},
				map[string]interface{}{
					"price": (func(vl string) float64 {
						if v, e := strconv.ParseFloat(vl, 64); e != nil {
							return 0.0
						} else {
							return v
						}
					})(bid.([]interface{})[0].(string)),
					"quantity": (func(vl string) float64 {
						if v, e := strconv.ParseFloat(vl, 64); e != nil {
							return 0.0
						} else {
							return v
						}
					})(bid.([]interface{})[1].(string)),
				},
				time.Now(),
			),
		)
	}

	return influx
}

func (influx *Influx) WriteTicker(message map[string]interface{}) *Influx {

	influx.WriteApi.WritePoint(
		influxdb2.NewPoint(
			fmt.Sprintf("%s@ticker", message["s"].(string)),
			map[string]string{
				"event_type": message["e"].(string),
				"symbol":     message["s"].(string),
				"event_time": time.UnixMilli(
					int64(message["E"].(float64)),
				).Format(time.RFC3339Nano),
				"now": time.Now().UTC().Format(time.RFC3339Nano),
			},
			map[string]interface{}{
				"event_time":                      time.UnixMilli(int64(message["E"].(float64))).UTC(),
				"statics_open_time":               time.UnixMilli(int64(message["O"].(float64))).UTC(),
				"statics_close_time":              time.UnixMilli(int64(message["C"].(float64))).UTC(),
				"first_trade_id":                  message["F"].(float64),
				"last_trade_id":                   message["L"].(float64),
				"total_number_of_trades":          message["n"].(float64),
				"best_ask_quantity":               message["A"].(string),
				"best_bid_quantity":               message["B"].(string),
				"price_change_percent":            message["P"].(string),
				"last_quantity":                   message["Q"].(string),
				"best_ask_price":                  message["a"].(string),
				"best_bid_price":                  message["b"].(string),
				"last_price":                      message["c"].(string),
				"high_price":                      message["h"].(string),
				"low_price":                       message["l"].(string),
				"open_price":                      message["o"].(string),
				"price_change":                    message["p"].(string),
				"total_traded_quote_asset_volume": message["q"].(string),
				"total_traded_base_asset_volume":  message["v"].(string),
				"weighted_average_price":          message["w"].(string),
				"first_trade_before_24hr":         message["x"].(string),
			},
			time.Now(),
		),
	)

	return influx
}

func (influx *Influx) WriteAggTrade(message map[string]interface{}) *Influx {
	influx.WriteApi.WritePoint(
		influxdb2.NewPoint(
			fmt.Sprintf("%s@aggTrade", message["s"].(string)),
			map[string]string{
				"event_type": message["e"].(string),
				"symbol":     message["s"].(string),
				"event_time": time.UnixMilli(
					int64(message["E"].(float64)),
				).Format(time.RFC3339Nano),
			},
			map[string]interface{}{
				"event_time":             time.UnixMilli(int64(message["E"].(float64))).UTC(),
				"first_trade_id":         message["f"].(float64),
				"last_trade_id":          message["l"].(float64),
				"is_buyer_market_marker": message["m"].(bool),
				"ignore":                 message["M"].(bool),
				"price":                  message["p"].(string),
				"quantity":               message["q"].(string),
				"trade_time":             message["T"].(float64),
				"aggregate_trade_time":   message["a"].(float64),
			},
			time.Now(),
		),
	)

	return influx
}

func (influx *Influx) Commit() *Influx {
	influx.WriteApi.Flush()

	return influx
}
