package hooks

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type fireFunc func(entry *logrus.Entry, hook *NsqHook) error

type NsqHook struct {
	client    *nsq.Producer
	topic     string
	levels    []logrus.Level
	ctx       context.Context
	ctxCancel context.CancelFunc
	fireFunc  fireFunc
}

type message struct {
	Timestamp string `json:"@timestamp"`
	File      string `json:"file"`
	Func      string `json:"func"`
	Message   string `json:"message"`
	Data      logrus.Fields
	Level     string `json:"level"`
}

func NewNsqHook(client *nsq.Producer, topic string, level logrus.Level) (*NsqHook, error) {
	return newHookFuncAndFireFunc(client, topic, level, syncFireFunc)
}

func NewAsyncNsqHook(client *nsq.Producer, topic string, level logrus.Level) (*NsqHook, error) {
	return newHookFuncAndFireFunc(client, topic, level, asyncFireFunc)
}

func newHookFuncAndFireFunc(client *nsq.Producer, topic string, level logrus.Level, fireFunc fireFunc) (*NsqHook, error) {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &NsqHook{
		client:    client,
		topic:     topic,
		levels:    levels,
		ctx:       ctx,
		ctxCancel: cancel,
		fireFunc:  fireFunc,
	}, nil
}

func (hook *NsqHook) Fire(entry *logrus.Entry) error {
	return hook.fireFunc(entry, hook)
}

func asyncFireFunc(entry *logrus.Entry, hook *NsqHook) error {
	client := hook.client
	msg := createMessage(entry, hook)
	txt, _ := json.Marshal(msg)
	return client.Publish(hook.topic, txt)
}

func syncFireFunc(entry *logrus.Entry, hook *NsqHook) error {
	client := hook.client
	msg := createMessage(entry, hook)
	txt, _ := json.Marshal(msg)
	return client.PublishAsync(hook.topic, txt, nil)
}

func createMessage(entry *logrus.Entry, hook *NsqHook) *message {
	level := entry.Level.String()

	if e, ok := entry.Data[logrus.ErrorKey]; ok && e != nil {
		if err, ok := e.(error); ok {
			entry.Data[logrus.ErrorKey] = err.Error()
		}
	}

	var file string
	var function string
	if entry.HasCaller() {
		file = entry.Caller.File
		function = entry.Caller.Function
	}

	return &message{
		entry.Time.UTC().Format(time.RFC3339Nano),
		file,
		function,
		entry.Message,
		entry.Data,
		strings.ToUpper(level),
	}
}

func (hook *NsqHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *NsqHook) Cancel() {
	hook.ctxCancel()
}
