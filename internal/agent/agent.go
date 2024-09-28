package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sudeeya/metrics-harvester/internal/metric"
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
		sigChan      = make(chan os.Signal, 1)
	)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var symmetricKey []byte
	if a.cfg.CryptoKey != "" {
		symmetricKey = generateSymmetricKey()
		publicKey := extractPublicKey(a.cfg.CryptoKey)
		a.shareSymmetricKey(symmetricKey, publicKey)
	}

	go func() {
		for range pollTicker.C {
			a.logger.Info("Updating metric values")
			go func() {
				metrics.Update()
			}()
			go func() {
				if err := metrics.UpdatePSUtil(); err != nil {
					a.logger.Error(err.Error())
				}
			}()
		}
	}()
	go func() {
		for range reportTicker.C {
			a.logger.Info("Sending all metrics")
			a.SendMetrics(metrics, symmetricKey)
		}
	}()
	go func() {
		<-sigChan
		a.logger.Info("Agent is shutting down")
		a.Shutdown()
	}()
	select {}
}

func (a *Agent) shareSymmetricKey(symmetricKey []byte, publicKey *rsa.PublicKey) {
	for _, backoff := range a.backoffSchedule {
		if err := a.tryShare(symmetricKey, publicKey); err != nil {
			a.logger.Error(err.Error())
			time.Sleep(backoff)
			continue
		}
		return
	}
}

func (a *Agent) tryShare(symmetricKey []byte, publicKey *rsa.PublicKey) error {
	request := a.client.R()
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return err
	}
	if a.cfg.Key != "" {
		h := hmac.New(sha256.New, []byte(a.cfg.Key))
		if _, err := h.Write(encryptedKey); err != nil {
			return err
		}
		request.SetHeader("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	}
	response, err := request.
		SetBody(encryptedKey).
		Post("/key/")
	if err != nil {
		return err
	}
	defer response.RawResponse.Body.Close()
	return nil
}

func (a *Agent) SendMetrics(metrics *Metrics, symmetricKey []byte) {
	mSlice := metrics.List()
	for _, backoff := range a.backoffSchedule {
		if err := a.trySend(mSlice, symmetricKey); err != nil {
			a.logger.Error(err.Error())
			time.Sleep(backoff)
			continue
		}
		return
	}
}

func (a *Agent) trySend(mSlice []metric.Metric, symmetricKey []byte) error {
	var buf bytes.Buffer
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
	if a.cfg.CryptoKey != "" {
		block, err := aes.NewCipher(symmetricKey)
		if err != nil {
			return err
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}
		encryptedBody := gcm.Seal(nonce, nonce, body, nil)
		body = encryptedBody
	}
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

func (a *Agent) Shutdown() {
	os.Exit(0)
}
