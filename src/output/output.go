package output

import "github.com/inoth/ino-gathere/src/metric"

type Output interface {
	Init() error
	Connect() error
	Close() error
	Write(metrics metric.Metric) error
}
