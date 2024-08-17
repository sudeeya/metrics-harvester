package agent

import (
	"compress/gzip"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sudeeya/metrics-harvester/internal/metric"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
)

type Agent struct {
	cfg             *Config
	logger          *zap.Logger
	client          *resty.Client
	backoffSchedule []time.Duration
}

func NewAgent(logger *zap.Logger, cfg *Config) *Agent {
	logger.Info("Initializing client")
	client := resty.New().SetBaseURL(cfg.Address)
	logger.Info("Initializing backoff schedule")
	backoffSchedule := initializeBackoffSchedule(logger, cfg)
	return &Agent{
		cfg:             cfg,
		logger:          logger,
		client:          client,
		backoffSchedule: backoffSchedule,
	}
}

func initializeBackoffSchedule(logger *zap.Logger, cfg *Config) []time.Duration {
	tmp := strings.Split(cfg.BackoffSchedule, ",")
	backoffSchedule := make([]time.Duration, len(tmp))
	for i, str := range tmp {
		value, err := strconv.Atoi(str)
		if err != nil {
			logger.Fatal(err.Error())
		}
		backoffSchedule[i] = time.Duration(value) * time.Millisecond
	}
	return backoffSchedule
}

func (a *Agent) Run() {
	a.logger.Info("Agent is running")
	var (
		metrics      = NewMetrics()
		pollTicker   = time.NewTicker(time.Duration(a.cfg.PollInterval) * time.Second)
		reportTicker = time.NewTicker(time.Duration(a.cfg.ReportInterval) * time.Second)
	)
	go func() {
		for range pollTicker.C {
			a.logger.Info("Updating metric values")
			UpdateMetrics(metrics)
		}
	}()
	go func() {
		for range reportTicker.C {
			a.logger.Info("Sending all metrics")
			a.SendMetrics(metrics)
		}
	}()
	select {}
}

func (a *Agent) SendMetrics(metrics *Metrics) {
	mSlice := make([]metric.Metric, len(metrics.Values))
	i := 0
	for _, m := range metrics.Values {
		mSlice[i] = *m
		i++
	}
	for _, backoff := range a.backoffSchedule {
		err := a.trySend(mSlice)
		if err == nil {
			break
		}
		a.logger.Error(err.Error())
		time.Sleep(backoff)
	}
}

func (a *Agent) trySend(mSlice []metric.Metric) error {
	var buf buffer.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return err
	}
	err = json.NewEncoder(gzipWriter).Encode(mSlice)
	if err != nil {
		return err
	}
	err = gzipWriter.Close()
	if err != nil {
		return err
	}
	response, err := a.client.R().
		SetHeader("content-type", "application/json").
		SetHeader("content-encoding", "gzip").
		SetHeader("accept-encoding", "gzip").
		SetBody(buf.Bytes()).
		Post("/updates/")
	if err != nil {
		return err
	}
	defer response.RawResponse.Body.Close()
	return nil
}
