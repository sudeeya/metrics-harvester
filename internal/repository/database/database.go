package database

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sudeeya/metrics-harvester/internal/metric"
)

const (
	CreateMetricsTable = `
CREATE TABLE IF NOT EXISTS metrics (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	delta INTEGER,
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
	delta = EXCLUDED.delta;
`
)

type Database struct {
	*sqlx.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Database{
		DB: db,
	}, nil
}

func (db *Database) PutMetric(m metric.Metric) error {
	switch m.MType {
	case metric.Gauge:
		_, err := db.DB.ExecContext(context.Background(),
			insertGauge, m.ID, m.MType, *m.Value)
		if err != nil {
			return err
		}
	case metric.Counter:
		_, err := db.DB.ExecContext(context.Background(),
			insertCounter, m.ID, m.MType, *m.Delta)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) GetMetric(mName string) (metric.Metric, error) {
	var (
		mType string
		delta sql.NullInt64
		value sql.NullFloat64
	)
	row := db.DB.QueryRowContext(context.Background(),
		"SELECT type, delta, value FROM metrics WHERE id = $1", mName)
	if err := row.Scan(&mType, &delta, &value); err != nil {
		return metric.Metric{}, err
	}
	if delta.Valid {
		return metric.Metric{
			ID:    mName,
			MType: mType,
			Delta: &delta.Int64,
		}, nil
	}
	return metric.Metric{
		ID:    mName,
		MType: mType,
		Value: &value.Float64,
	}, nil
}

func (db *Database) GetAllMetrics() ([]metric.Metric, error) {
	allMetrics := make([]metric.Metric, 0)
	var (
		delta sql.NullInt64
		value sql.NullFloat64
	)
	rows, err := db.DB.QueryContext(context.Background(),
		"SELECT id, type, delta, value FROM metrics ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m metric.Metric
		if err := rows.Scan(&m.ID, &m.MType, &delta, &value); err != nil {
			return nil, err
		}
		switch {
		case delta.Valid:
			m.Delta = &delta.Int64
		case value.Valid:
			m.Value = &value.Float64
		}
		allMetrics = append(allMetrics, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allMetrics, nil
}

func (db *Database) Close() error {
	if err := db.DB.Close(); err != nil {
		return err
	}
	return nil
}
