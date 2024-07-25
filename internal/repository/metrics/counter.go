package metrics

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
