package analytics

import (
	"fmt"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/settings"
)

type influxDb struct {
	sync.Mutex
	agent    client.Client
	replicas []client.Client
	batch    client.BatchPoints
	database string
	interval time.Duration
	ticker   *time.Ticker
	enabled  bool
	logger   ILogger
}

func NewInfluxDbProvider(configuration IConfiguration, logger ILogger) IMeasurementsProvider {
	influxConfiguration := configuration.GetInfluxConfiguration()

	agent, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxConfiguration.GetAddress(),
		Username: influxConfiguration.GetUsername(),
		Password: influxConfiguration.GetPassword(),
	})

	if err != nil {
		return nil
	}

	replicas := make([]client.Client, 0)
	for _, address := range influxConfiguration.GetReplicas() {
		replicaAgent, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     address,
			Username: influxConfiguration.GetUsername(),
			Password: influxConfiguration.GetPassword(),
		})

		if err != nil {
			return nil
		}

		replicas = append(replicas, replicaAgent)
	}

	interval := 5 * time.Second
	provider := &influxDb{
		agent:    agent,
		replicas: replicas,
		database: influxConfiguration.GetDatabase(),
		interval: interval,
		ticker:   time.NewTicker(interval),
		enabled:  influxConfiguration.IsEnabled(),
		logger:   logger,
	}

	provider.resetBatch()

	go func() {
		for {
			select {
			case <-provider.ticker.C:
				go provider.flush(provider.batch)
				provider.resetBatch()
			}
		}
	}()

	return provider
}

func (influxDb *influxDb) SubmitMeasurementAsync(key string, tags Tags, fields Fields) {
	if !influxDb.enabled {
		return
	}

	go func() {
		influxDb.SubmitMeasurement(key, tags, fields)
	}()
}

func (influxDb *influxDb) SubmitMeasurement(key string, tags Tags, fields Fields) {
	if !influxDb.enabled {
		return
	}

	defer func() {
		if reason := recover(); reason != nil {
			influxDb.logger.Panic(fmt.Sprintf("IFX: %s", reason))
		}
	}()

	point, err := client.NewPoint(key, tags, fields, time.Now())
	if err != nil {
		influxDb.logger.Error(fmt.Sprintf("IFX/PNT: %s", err))
		return
	}

	influxDb.Lock()
	defer influxDb.Unlock()
	influxDb.batch.AddPoint(point)
	if len(influxDb.batch.Points()) >= 7500 {
		go influxDb.flush(influxDb.batch)
		influxDb.resetBatch()
	}
}

func (influxDb *influxDb) resetBatch() {
	var err error
	if influxDb.batch, err = client.NewBatchPoints(client.BatchPointsConfig{Database: influxDb.database}); err != nil {
		influxDb.logger.Error(fmt.Sprintf("IFX: %s", err))
	}
}

func (influxDb *influxDb) flush(batch client.BatchPoints) {
	if len(batch.Points()) < 1 {
		return
	}

	if err := influxDb.agent.Write(batch); err != nil {
		influxDb.logger.Error(fmt.Sprintf("IFX: %s", err))
	}

	for _, replica := range influxDb.replicas {
		if err := replica.Write(batch); err != nil {
			influxDb.logger.Error(fmt.Sprintf("IFX: REPLICA %s", err))
		}
	}
}
