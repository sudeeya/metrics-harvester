package metric

import "fmt"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

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
