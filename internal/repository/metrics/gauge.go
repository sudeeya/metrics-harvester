package metrics

import "fmt"

type Gauge struct {
	name  string
	value float64
}

func NewGauge(name string, value float64) *Gauge {
	return &Gauge{
		name:  name,
		value: value,
	}
}

func (g *Gauge) ChangeValue(newValue float64) {
	g.value = newValue
}

func (g Gauge) GetName() string {
	return g.name
}

func (g Gauge) GetValue() string {
	return fmt.Sprintf("%v", g.value)
}
