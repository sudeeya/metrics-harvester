package agent

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
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
		backoffSchedule[i] = time.Duration(value) * time.Second
	}
	return backoffSchedule
}

func (a *Agent) Run() {
	a.logger.Info("Agent is running")
	var (
		metrics      = NewMetrics()
		pollTicker   = time.NewTicker(time.Duration(a.cfg.PollInterval) * time.Second)
		reportTicker = time.NewTicker(time.Duration(a.cfg.ReportInterval) * time.Second)
		rwMutex      sync.RWMutex
	)
	go func() {
		for range pollTicker.C {
			a.logger.Info("Updating metric values")
			rwMutex.Lock()
			UpdateMetrics(metrics)
			rwMutex.Unlock()
		}
	}()
	go func() {
		for range reportTicker.C {
			a.logger.Info("Sending all metrics")
			rwMutex.RLock()
			a.SendMetrics(metrics)
			rwMutex.RUnlock()
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
		if err := a.trySend(mSlice); err != nil {
			a.logger.Error(err.Error())
			time.Sleep(backoff)
			continue
		}
		return
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
	body := buf.Bytes()
	request := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip")
	if a.cfg.Key != "" {
		h := hmac.New(sha256.New, []byte(a.cfg.Key))
		if _, err := h.Write(body); err != nil {
			return err
		}
		request.SetHeader("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	}
	response, err := request.
		SetBody(body).
		Post("/updates/")
	if err != nil {
		return err
	}
	defer response.RawResponse.Body.Close()
	return nil
}
