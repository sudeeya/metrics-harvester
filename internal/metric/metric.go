// Package metric provides methods for working with metrics.
package metric

import "fmt"

// Types of metrics.
const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Metric contains metric parameters.
type Metric struct {
	// ID identifies the metric.
	ID string `json:"id"`

	// MType defines the type of metric.
	MType string `json:"type"`

	// Delta stores a pointer if the type is Counter. Otherwise it is nil.
	Delta *int64 `json:"delta,omitempty"`

	// Value stores a pointer if the type is Gauge. Otherwise it is nil.
	Value *float64 `json:"value,omitempty"`
}

// Update changes the value depending on the type of metric.
func (m *Metric) Update(update any) {
	switch m.MType {
	case Gauge:
		newValue := update.(float64)
		m.Value = &newValue
	case Counter:
		delta := update.(int64)
		*m.Delta += delta
	}
}

// GetValue returns the value in string format depending on the type of metric.
func (m Metric) GetValue() string {
	var value string
	switch m.MType {
	case Gauge:
		value = fmt.Sprintf("%v", *m.Value)
	case Counter:
		value = fmt.Sprintf("%v", *m.Delta)
	}
	return value
}
