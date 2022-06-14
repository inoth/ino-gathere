package agent

import (
	"errors"

	"github.com/inoth/ino-gathere/src/input"
	"github.com/inoth/ino-gathere/src/output"
)

// 保存采集器信息
// 输出器信息
type AgentConfig struct {
	Inputs  []input.Input
	Outputs []output.Output
}

// 加载所有装载的输入、输出模块
func (ag *AgentConfig) Init() error {
	if len(input.ReadyCollectors) <= 0 {
		return errors.New("not found input")
	}
	for _, input := range input.ReadyCollectors {
		ag.Inputs = append(ag.Inputs, input)
	}
	if len(output.ReadyOutputs) <= 0 {
		return errors.New("not found output")
	}
	for _, output := range output.ReadyOutputs {
		ag.Outputs = append(ag.Outputs, output)
	}
	return nil
}

// 注入到agent, 运行
func (ag *AgentConfig) RunServer() {
	NewAgent(ag).Run()
	select {}
}
