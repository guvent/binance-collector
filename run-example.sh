#!/bin/bash

export INFLUX_URL=http://127.0.0.1:8086
export INFLUX_TOKEN='0AzvZOwOwe1Nhzg6nw6aedtFL7QTnhyoBXjXwR9Jw7WXmTNhNyqwnQekl75NAo__nYYNkjWE0SKuBDqdV-7rCw=='

export INFLUX_BUCKET=buysell
export SYMBOLS=BTCUSDT,ETHUSDT,ETCUSDT,XRPUSDT,LTCUSDT,SOLUSDT,XMRUSDT,BCHUSDT,ATOMUSDT,LINKUSDT,DOGEUSDT,AVAXUSDT

export INFLUX_COMMIT_MS=500ms

export GO111MODULE="on"
export TZ=Europe/Istanbul

#######

git fetch

git pull origin master

#######

if  ! [ -e main ]
then
   rm main
fi

go mod download

go build -o main -tags musl

chmod 755 main

if  ! [ -d logs ]
then
   mkdir logs
fi

while ! [ -e .crashed ]
do
   a=$(date +%F-%H-%M-%S)
   ./main >> logs/out-$a.log
done

echo 'terminated...'
