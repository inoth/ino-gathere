package collector

import "github.com/inoth/ino-gathere/src/metric"

// 用于检查注册再字典中的采集器信息, 负责调用暴露接口
// 采集器中可以记录出现错误次数, 超过阈值转换为脱机状体啊, 避免浪费调用资源

type ICollector interface {
	// Init() error
	GetMetrics() ([]metric.MetricValue, error)
}

type Creator func() ICollector

// 准备好的采集器
var readyCollectors = map[string]ICollector{}

// 需要工作的采集器
var workCollectors = map[string]ICollector{}

// 添加准备工作的采集器
func Add(name string, creator Creator) {
	readyCollectors[name] = creator()
}
