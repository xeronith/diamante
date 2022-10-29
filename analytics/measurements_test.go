package analytics_test

import (
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/xeronith/diamante/contracts/analytics"
)

func Test_SubmitMeasurement(test *testing.T) {

	// line := "mini-games,type=score game-id=1570887509866549i,game-name=\"FLAPPY_BIRD\",origin=\"q\",score=5i,user=2173373u 1570887583106751044"

	agent, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "",
		Password: "",
	})

	if err != nil {
		test.Fatal(err)
	}

	point, err := client.NewPoint("mini_games", analytics.Tags{"type": "score"}, analytics.Fields{
		"user":      int64(2173373),
		"score":     int64(5),
		"origin":    "q",
		"game_id":   int64(1570887509866549),
		"game_name": "FLAPPY_BIRD",
	}, time.Now())

	if err != nil {
		test.Fatal(err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: "test"})

	if err != nil {
		test.Fatal(err)
	}

	bp.AddPoint(point)

	if err := agent.Write(bp); err != nil {
		test.Fatal(err)
	}

}
