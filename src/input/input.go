package input

import "github.com/inoth/ino-gathere/src/accumulator"

type Input interface {
	Init() error
	GetMetrics(acc accumulator.Accumulator) error
}
