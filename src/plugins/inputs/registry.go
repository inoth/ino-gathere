package inputs

import (
	"fmt"
	"sync"

	"github.com/inoth/ino-gathere/src/input"
)

type Creator func() input.Input

// 准备好的采集器
var ReadyCollectors = make(map[string]Creator)

// 需要工作的采集器
var workCollector = workCollectors{
	m:     sync.Mutex{},
	works: make(map[string]Creator),
}

type workCollectors struct {
	m     sync.Mutex
	works map[string]Creator
}

// 添加准备工作的采集器
func Add(name string, creator Creator) {
	fmt.Printf("加载: %v 采集器\n", name)
	ReadyCollectors[name] = creator
}

// 变更正在工作中的采集器, 也许加上生命周期, 切换时停止
func (wc *workCollectors) ChangeWorkCollectors(keys ...string) {
	wc.m.Lock()
	defer wc.m.Unlock()
	for k := range wc.works {
		delete(wc.works, k)
	}
	for k, v := range ReadyCollectors {
		wc.works[k] = v
	}
}
