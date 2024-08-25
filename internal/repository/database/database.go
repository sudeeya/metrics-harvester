package database

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sudeeya/metrics-harvester/internal/metric"
)

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

type Database struct {
	*sqlx.DB
}

func NewDatabase(dsn string) *Database {
	db := sqlx.MustOpen("pgx", dsn)
	return &Database{
		DB: db,
	}
}

func (db *Database) PutMetric(m metric.Metric) error {
	switch m.MType {
	case metric.Gauge:
		_, err := db.ExecContext(context.TODO(),
			insertGauge, m.ID, m.MType, *m.Value)
		if err != nil {
			return err
		}
	case metric.Counter:
		_, err := db.ExecContext(context.TODO(),
			insertCounter, m.ID, m.MType, *m.Delta)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) PutBatch(metrics []metric.Metric) error {
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	stmtGauge, err := tx.PrepareContext(context.TODO(), insertGauge)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()
	stmtCounter, err := tx.PrepareContext(context.TODO(), insertCounter)
	if err != nil {
		return err
	}
	defer stmtCounter.Close()
	for _, m := range metrics {
		switch m.MType {
		case metric.Gauge:
			_, err := stmtGauge.ExecContext(context.TODO(), m.ID, m.MType, *m.Value)
			if err != nil {
				return tx.Rollback()
			}
		case metric.Counter:
			_, err := stmtCounter.ExecContext(context.TODO(), m.ID, m.MType, *m.Delta)
			if err != nil {
				return tx.Rollback()
			}
		}
	}
	return tx.Commit()
}

func (db *Database) GetMetric(mName string) (metric.Metric, error) {
	var dbm DBMetric
	if err := db.GetContext(context.TODO(), &dbm,
		"SELECT id, type, delta, value FROM metrics WHERE id = $1", mName); err != nil {
		return metric.Metric{}, err
	}
	return dbm.ToMetric(), nil
}

func (db *Database) GetAllMetrics() ([]metric.Metric, error) {
	var dbMetrics []DBMetric
	if err := db.SelectContext(context.TODO(), &dbMetrics,
		"SELECT id, type, delta, value FROM metrics ORDER BY id"); err != nil {
		return nil, err
	}
	allMetrics := make([]metric.Metric, len(dbMetrics))
	for i, dbm := range dbMetrics {
		allMetrics[i] = dbm.ToMetric()
	}
	return allMetrics, nil
}
