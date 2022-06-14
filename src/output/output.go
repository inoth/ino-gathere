package output

type Output interface {
	Init() error
	Connect() error
	Close() error
	Write() error
}

type Creator func() Output

// 准备好的采集器
var ReadyOutputs = map[string]Output{}

// 添加准备工作的采集器
func Add(name string, creator Creator) {
	ReadyOutputs[name] = creator()
}
