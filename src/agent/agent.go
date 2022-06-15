package agent

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/inoth/ino-gathere/src/accumulator"
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

func (a *Agent) Run(ctx context.Context) error {
	// run, 获取配置数据
	// config.Get(采集器工作频率)
	err := a.initPlugins()
	if err != nil {
		return err
	}

	// 作为服务启动项, 阻塞程序避免结束
	next, out, err := a.startOutputs(ctx, a.ag.Outputs)
	if err != nil {
		return err
	}
	in, err := a.startInputs(next, a.ag.Inputs)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 运行 output 程序
		a.runOutputs(out)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 运行 input 程序
		a.runInputs(ctx, in)
	}()

	wg.Wait()
	fmt.Println("=========采集结束=========")
	return err
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

func (a *Agent) runInputs(ctx context.Context, in *inputUnit) {
	// 给每个采集器单独开线程执行
	var wg sync.WaitGroup
	for _, inp := range in.inputs {
		acc := accumulator.New(in.dst)
		wg.Add(1)
		go func(in input.Input) {
			defer wg.Done()
			a.gatherLoop(ctx, acc, inp)
		}(inp)
	}
	wg.Wait()
	close(in.dst)
	fmt.Println("采集器结束工作")
}

func (a *Agent) runOutputs(out *outputUnit) {
	for mrc := range out.src {
		for i, ou := range out.outputs {
			if i == len(out.outputs)-1 {
				ou.Write(mrc)
			} else {
				ou.Write(mrc.Copy())
			}
		}
	}
}

// gather runs an input's gather function periodically until the context is
// done.
func (a *Agent) gatherLoop(
	ctx context.Context,
	acc accumulator.Accumulator,
	input input.Input,
	// 采集间隔通过配置初始化到 input 模组中
) {
	defer panicRecover(input)
	// 创建一个定时器
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			err := a.gatherOnce(acc, input)
			if err != nil {
				acc.AddError(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// gatherOnce runs the input's Gather function once, logging a warning each
// interval it fails to complete before.
func (a *Agent) gatherOnce(
	acc accumulator.Accumulator,
	input input.Input,
) error {
	done := make(chan error)
	go func() {
		done <- input.GetMetrics(acc)
	}()

	// Only warn after interval seconds, even if the interval is started late.
	// Intervals can start late if the previous interval went over or due to
	// clock changes.
	// slowWarning := time.NewTicker(interval)
	// defer slowWarning.Stop()

	for {
		select {
		case err := <-done:
			return err
			// case <-slowWarning.C:
			// 	log.Printf("W! [%s] Collection took longer than expected; not complete after interval of %s",
			// 		input.LogName(), interval)
			// case <-ticker.Elapsed():
			// 	log.Printf("D! [%s] Previous collection has not completed; scheduled collection skipped",
			// 		input.LogName())
		}
	}
}

// panicRecover displays an error if an input panics.
func panicRecover(input input.Input) {
	if err := recover(); err != nil {
		trace := make([]byte, 2048)
		runtime.Stack(trace, true)
		log.Printf("E! FATAL: [] panicked: %s, Stack:\n%s", err, trace)
		log.Println("E! PLEASE REPORT THIS PANIC ON GITHUB with " +
			"stack trace, configuration, and OS information: " +
			"https://github.com/influxdata/telegraf/issues/new/choose")
	}
}
