package agent

import (
	"context"
	"errors"

	"github.com/inoth/ino-gathere/src/input"
	"github.com/inoth/ino-gathere/src/output"
	"github.com/inoth/ino-gathere/src/plugins/inputs"
	"github.com/inoth/ino-gathere/src/plugins/outputs"
)

// 保存采集器信息
// 输出器信息
type AgentConfig struct {
	Inputs []input.Input
	// Inputs []InputConfig
	Outputs []output.Output
	// Outputs []OutputConfig
}

// 包含采集器需要用到的配置信息
// 采集器自身, 方便从第三方写入配置
type InputConfig struct {
	Config map[string]string
	Input  input.Input
}
type OutputConfig struct {
	Config map[string]string
	OutPut input.Input
}

// 加载所有装载的输入、输出模块
func (ag *AgentConfig) Init() error {

	if len(inputs.ReadyCollectors) <= 0 {
		return errors.New("not found input")
	}
	for _, input := range inputs.ReadyCollectors {
		ag.Inputs = append(ag.Inputs, input())
		// ag.Inputs = append(ag.Inputs, InputConfig{
		// 	// 从 viper 中读取
		// 	//Config: viper.get(name),
		// 	Input: input(),
		// })
	}
	if len(outputs.ReadyOutputs) <= 0 {
		return errors.New("not found output")
	}
	for _, output := range outputs.ReadyOutputs {
		ag.Outputs = append(ag.Outputs, output())
	}
	return nil
}

// 注入到agent, 运行
func (ag *AgentConfig) ServeStart() error {
	// err := ag.InitAgent()
	// if err != nil {
	// 	return err
	// }
	return NewAgent(ag).Run(context.TODO())
}
