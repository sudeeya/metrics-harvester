package database

import (
	"database/sql"

	"github.com/sudeeya/metrics-harvester/internal/metric"
)

// DBMetric is an auxiliary structure into which the database response is written.
type DBMetric struct {
	ID    string          `db:"id"`
	MType string          `db:"type"`
	Delta sql.NullInt64   `db:"delta"`
	Value sql.NullFloat64 `db:"value"`
}

// ToMetric converts DBMetric to metric.Metric.
func (dbm DBMetric) ToMetric() metric.Metric {
	var m metric.Metric
	m.ID = dbm.ID
	m.MType = dbm.MType
	if dbm.Delta.Valid {
		m.Delta = &dbm.Delta.Int64
	}
	if dbm.Value.Valid {
		m.Value = &dbm.Value.Float64
	}
	return m
}
