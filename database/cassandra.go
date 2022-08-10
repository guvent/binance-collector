package database

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type Cassandra struct {
	Cluster *gocql.ClusterConfig
}

func (cassandra *Cassandra) Init() *Cassandra {
	cassandra.Cluster = gocql.NewCluster("10.50.42.59")
	cassandra.Cluster.Keyspace = "example"
	cassandra.Cluster.Consistency = gocql.Quorum

	return cassandra
}

func (cassandra *Cassandra) WriteData() *Cassandra {
	if session, sessionErr := cassandra.Cluster.CreateSession(); sessionErr != nil {
		log.Print(sessionErr)
	} else {
		mapVal := map[string]interface{}{
			"id":       gocql.TimeUUID().String(),
			"price":    1.12345678912345678,
			"quantity": 3.0000000000001,
			"time":     time.Now().UTC(),
			"asks": []map[string]float64{
				{
					"price":    1.12345678912345678,
					"quantity": 3.0000000000001,
				},
				{
					"price":    133.12345678912345678,
					"quantity": 4.0000000000001,
				},
			},
			"bids": []map[string]float64{
				{
					"price":    2.12345678912345678,
					"quantity": 5.0000000000001,
				},
				{
					"price":    3.12345678912345678,
					"quantity": 7.0000000000001,
				},
			},
		}
		mapJson, _ := json.Marshal(mapVal)
		if queryErr := session.Query(
			"INSERT INTO btcusdt JSON ?", string(mapJson),
		).Exec(); queryErr != nil {
			log.Print(queryErr)
		} else {
			log.Print("Success...")
		}
	}

	return cassandra
}
