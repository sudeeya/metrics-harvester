package metrics

type Metric interface {
	GetName() string
	GetValue() string
}
