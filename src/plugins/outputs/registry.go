package outputs

import (
	"fmt"

	"github.com/inoth/ino-gathere/src/output"
)

type Creator func() output.Output

// 准备好的采集器
var ReadyOutputs = make(map[string]Creator)

// 添加准备工作的采集器
func Add(name string, creator Creator) {
	fmt.Printf("加载: %v 输出器\n", name)
	ReadyOutputs[name] = creator
}
