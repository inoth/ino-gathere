package agent

import (
	"context"

	"github.com/inoth/ino-gathere/src/input"
	"github.com/inoth/ino-gathere/src/metric"
	"github.com/inoth/ino-gathere/src/output"
)

// 资源统合管理器
// 注册采集器
// 注册输出器
// 设定采集频率, 设定同一时刻采集开启最大阈值
type inputUnit struct {
	dst chan<- metric.Metric
	// 输入采集器列表
	inputs []input.Input
}

// 用作输出
type outputUnit struct {
	src <-chan metric.Metric
	// 输出发送器列表
	outputs []output.Output
}

type Agent struct {
	ag *AgentConfig
}

func NewAgent(ag *AgentConfig) *Agent {
	return &Agent{
		ag: ag,
	}
}

func (a *Agent) Run() {
	// run, 获取配置数据
	// config.Get(采集器工作频率)

	// 作为服务启动项, 阻塞程序避免结束
	select {}
}

func (a *Agent) initPlugins() error {
	for _, in := range a.ag.Inputs {
		err := in.Init()
		if err != nil {
			return err
		}
	}
	for _, out := range a.ag.Outputs {
		err := out.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Agent) startInputs(dst chan<- metric.Metric, inputs []input.Input) (*inputUnit, error) {
	unit := &inputUnit{
		dst: dst,
	}
	for _, input := range inputs {
		// 初始化累加器, 添加组建中
		unit.inputs = append(unit.inputs, input)
	}
	return unit, nil
}

func (a *Agent) startOutputs(ctx context.Context, outputs []output.Output) (chan<- metric.Metric, *outputUnit, error) {
	src := make(chan metric.Metric, 100)
	unit := &outputUnit{src: src}
	for _, output := range outputs {
		// 初始化连接
		unit.outputs = append(unit.outputs, output)
	}
	return src, unit, nil
}
