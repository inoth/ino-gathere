package metric

import "fmt"

type MetricValue struct {
	Endpoint  string                 `json:"endpoint"`
	Metric    string                 `json:"metric"`
	Value     interface{}            `json:"value"`
	Step      int64                  `json:"step"`
	Type      string                 `json:"type"`
	Tags      map[string]interface{} `json:"tags"`
	Timestamp int64                  `json:"timestamp"`
}

func (this *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}

func (this *MetricValue) Clone() *MetricValue {
	mv := MetricValue{
		Metric: this.Metric,
		Value:  this.Value,
		Type:   this.Type,
	}
	tagsMap := map[string]interface{}{}
	for k, v := range this.Tags {
		tagsMap[k] = v
	}
	mv.Tags = tagsMap
	return &mv
}
