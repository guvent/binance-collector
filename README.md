## Binance Collector

### Go Version 1.16.15

---

Run:

```
export GO111MODULE="on"

export INFLUX_URL='http://<influx_host>:<influx_port>'
export INFLUX_TOKEN='<influx_token>'
export INFLUX_BUCKET='<your_bucket_name>'

export SYMBOLS='BTCUSDT,ETHUSDT,ETCUSDT,.....'

go mod download

go run main
```

Build:

```
export GO111MODULE="on"

export INFLUX_URL='http://<influx_host>:<influx_port>'
export INFLUX_TOKEN='<influx_token>'
export INFLUX_BUCKET='<your_bucket_name>'

export SYMBOLS='BTCUSDT,ETHUSDT,ETCUSDT,.....'

go mod download

go build -o main

chomd +x main

./main
```
