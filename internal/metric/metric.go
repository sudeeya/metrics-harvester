package metric

type Metric interface {
	GetName() string
	GetValue() string
}
