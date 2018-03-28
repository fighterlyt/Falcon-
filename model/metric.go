package model

import (
	"strings"
	"fmt"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
	"os"
)

//CounterType 数值类型,只能是COUNTER或者GAUGE二选一
type CounterType int

const (
	GAUGE   CounterType = iota //即用户上传什么样的值，就原封不动的存储
	COUNTER                    //指标在存储和展现的时候，会被计算为speed，即（当前值 - 上次值）/ 时间间隔
)

func (c CounterType) String() string {
	switch c {
	case GAUGE:
		return "GAUGE"
	case COUNTER:
		return "COUNTER"
	default:
		return "UNKNOWN"
	}
}

func (c CounterType) MarshalJSON() ([]byte, error) {
	return []byte(c.String()), nil
}

//Tags 一组逗号分割的键值对, 对metric进一步描述和细化, 可以是空字符串. 比如idc=lg，比如service=xbox等
type Tags map[string]string

//多个tag之间用逗号分割
func (t Tags) String() string {

	result := make([]string, 0, len(t))

	for k, v := range t {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(result, ",")
}

func (t Tags) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

type Metric struct {
	Endpoint    string      `json:"endpoint"`
	Metric      string      `json:"metric"`
	Timestamp   int64       `json:"timestamp"`
	Step        int         `json:"step"`
	Value       int         `json:"value"`
	CounterType CounterType `json:"counterType"`
	Tags        Tags        `json:"tags"`
}

func NewMetric(endpoint,metric string,timestamp int64,step,value int,counterType CounterType,tags Tags) Metric{
	return Metric{
		Endpoint:endpoint,
		Metric:metric,
		Timestamp:timestamp,
		Step:step,
		Value:value,
		CounterType:counterType,
		Tags:tags,
	}
}


type Agent struct{
	Server  string
	Port string
	Endpoint  string
	Metric string
	Step int
	CounterType CounterType
	Tags Tags
}

func (c Agent) NewMetric(value int) Metric{
	return NewMetric(c.Endpoint,c.Metric,time.Now().Unix(),value,c.Step,c.CounterType,c.Tags)
}

func (c Agent) Push(value int)error{
	metric:=c.NewMetric(value)
	return c.push(metric)
}

func (c Agent) push(metric Metric) error{
	buf:=&bytes.Buffer{}

	json.NewEncoder(buf).Encode(metric)
	if resp,err:=http.Post(fmt.Sprintf("http://%s:%s/v1/push",c.Server,c.Port),"application/json",buf);err!=nil{
		return err
	}else{
		resp.Write(os.Stdout)
		return nil
	}

}