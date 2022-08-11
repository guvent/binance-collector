package utils

import (
	"log"
	"os"
	"strings"
)

type GetEnvironment struct {
	ServerURL string
	AuthToken string
	Bucket    string
	Org       string
	CommitMS  string
	DepthMS   string
	Symbols   string
}

func (env *GetEnvironment) Init() *GetEnvironment {

	env.ServerURL = os.Getenv("INFLUX_URL")
	env.AuthToken = os.Getenv("INFLUX_TOKEN")
	env.Bucket = os.Getenv("INFLUX_BUCKET")
	env.CommitMS = os.Getenv("INFLUX_COMMIT_MS")
	env.Symbols = os.Getenv("SYMBOLS")

	env.Org = "buysell"
	env.DepthMS = "100"

	return env.load()
}

func (env *GetEnvironment) load() *GetEnvironment {
	env.ServerURL = strings.ReplaceAll(env.ServerURL, "'", "")
	env.AuthToken = strings.ReplaceAll(env.AuthToken, "'", "")
	env.Bucket = strings.ReplaceAll(env.Bucket, "'", "")
	env.CommitMS = strings.ReplaceAll(env.CommitMS, "'", "")
	env.Symbols = strings.ReplaceAll(env.Symbols, "'", "")

	log.Printf("* INFLUX_URL : %s", env.ServerURL)
	log.Printf("* INFLUX_TOKEN : %s", env.AuthToken)
	log.Printf("* INFLUX_BUCKET : %s", env.Bucket)
	log.Printf("* INFLUX_COMMIT_MS : %s", env.CommitMS)
	log.Printf("* SYMBOLS: %v", env.Bucket)

	return env.check()
}

func (env *GetEnvironment) check() *GetEnvironment {
	if strings.Trim(env.ServerURL, " ") == "" {
		log.Fatal("Require Influx URL!")
	}
	if strings.Trim(env.AuthToken, " ") == "" {
		log.Fatal("Require Influx Auth Token!")
	}
	if strings.Trim(env.Bucket, " ") == "" {
		log.Fatal("Require Influx Bucket Name!")
	}
	if strings.Trim(env.Org, " ") == "" {
		log.Fatal("Require Influx Organization Name!")
	}
	if strings.Trim(env.CommitMS, " ") == "" {
		log.Fatal("Require Influx Commit Ms!")
	}
	if strings.Trim(env.DepthMS, " ") == "" {
		log.Fatal("Require Depth Value (ms)!")
	}

	return env
}
