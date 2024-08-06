package metric

import "fmt"

type Counter struct {
	name  string
	value int64
}

func NewCounter(name string, value int64) *Counter {
	return &Counter{
		name:  name,
		value: value,
	}
}

func (c *Counter) IncreaseValue(additive int64) {
	c.value += additive
}

func (c Counter) GetName() string {
	return c.name
}

func (c Counter) GetValue() string {
	return fmt.Sprintf("%v", c.value)
}
