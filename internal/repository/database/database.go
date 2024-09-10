package database

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/sudeeya/metrics-harvester/internal/metric"
	"github.com/sudeeya/metrics-harvester/internal/repository"
)

const limitInSeconds = 10

const (
	CreateMetricsTable = `
CREATE TABLE IF NOT EXISTS metrics (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	delta BIGINT,
	value DOUBLE PRECISION
);
`
	insertGauge = `
INSERT INTO metrics (id, type, value)
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET
	value = EXCLUDED.value;
`
	insertCounter = `
INSERT INTO metrics (id, type, delta)
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET
	delta = metrics.delta + EXCLUDED.delta;
`
)

var _ repository.Repository = (*Database)(nil)

type Database struct {
	*sqlx.DB
}

func NewDatabase(dsn string) *Database {
	db := sqlx.MustOpen("pgx", dsn)
	return &Database{
		DB: db,
	}
}

func (db *Database) PutMetric(ctx context.Context, m metric.Metric) error {
	ctx, cancel := context.WithTimeout(ctx, limitInSeconds*time.Second)
	defer cancel()
	switch m.MType {
	case metric.Gauge:
		_, err := db.ExecContext(ctx, insertGauge, m.ID, m.MType, *m.Value)
		if err != nil {
			return err
		}
	case metric.Counter:
		_, err := db.ExecContext(ctx, insertCounter, m.ID, m.MType, *m.Delta)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) PutBatch(ctx context.Context, metrics []metric.Metric) error {
	ctx, cancel := context.WithTimeout(ctx, limitInSeconds*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmtGauge, err := tx.PrepareContext(ctx, insertGauge)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()
	stmtCounter, err := tx.PrepareContext(ctx, insertCounter)
	if err != nil {
		return err
	}
	defer stmtCounter.Close()
	for _, m := range metrics {
		switch m.MType {
		case metric.Gauge:
			_, err := stmtGauge.ExecContext(ctx, m.ID, m.MType, *m.Value)
			if err != nil {
				return tx.Rollback()
			}
		case metric.Counter:
			_, err := stmtCounter.ExecContext(ctx, m.ID, m.MType, *m.Delta)
			if err != nil {
				return tx.Rollback()
			}
		}
	}
	return tx.Commit()
}

func (db *Database) GetMetric(ctx context.Context, mName string) (metric.Metric, error) {
	var dbm DBMetric
	ctx, cancel := context.WithTimeout(ctx, limitInSeconds*time.Second)
	defer cancel()
	if err := db.GetContext(ctx, &dbm,
		"SELECT id, type, delta, value FROM metrics WHERE id = $1", mName); err != nil {
		return metric.Metric{}, err
	}
	return dbm.ToMetric(), nil
}

func (db *Database) GetAllMetrics(ctx context.Context) ([]metric.Metric, error) {
	var dbMetrics []DBMetric
	ctx, cancel := context.WithTimeout(ctx, limitInSeconds*time.Second)
	defer cancel()
	if err := db.SelectContext(ctx, &dbMetrics,
		"SELECT id, type, delta, value FROM metrics ORDER BY id"); err != nil {
		return nil, err
	}
	allMetrics := make([]metric.Metric, len(dbMetrics))
	for i, dbm := range dbMetrics {
		allMetrics[i] = dbm.ToMetric()
	}
	return allMetrics, nil
}
