package http

import (
	"fmt"

	"github.com/inoth/ino-gathere/src/metric"
	"github.com/inoth/ino-gathere/src/output"
	"github.com/inoth/ino-gathere/src/plugins/outputs"
)

type HttpOutput struct {
}

func (HttpOutput) Init() error {
	// fmt.Println("http 做了一些初始化操作")
	return nil
}
func (HttpOutput) Connect() error { return nil }
func (HttpOutput) Close() error   { return nil }

func (HttpOutput) Write(metrics metric.Metric) error {
	// fmt.Printf("%v: %v\n%v\n", metrics.Name(), metrics.Tags(), metrics.Fields())
	fmt.Println(metrics.String())
	return nil
}

func init() {
	outputs.Add("http", func() output.Output {
		return &HttpOutput{}
	})
}
