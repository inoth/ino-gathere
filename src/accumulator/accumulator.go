package accumulator

import (
	"time"

	"github.com/inoth/ino-gathere/src/metric"
)

type Accumulator interface {
	AddFields(measurement string, fields map[string]interface{}, tags map[string]string, tm time.Time)
}

type accumulator struct {
	metrics chan<- metric.Metric
}

func New(metrics chan<- metric.Metric) Accumulator {
	acc := accumulator{
		metrics: metrics,
	}
	return &acc
}

func (acc *accumulator) AddFields(measurement string, fields map[string]interface{}, tags map[string]string, tm time.Time) {
	m := metric.New(measurement, tags, fields, tm)
	acc.metrics <- m
}
