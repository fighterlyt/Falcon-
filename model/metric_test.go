package model

import (
	"testing"
	"time"
)

func TestAgent_Push(t *testing.T) {

	agent:=Agent{
		Server:"localhost",
		Port:"8081",
		Endpoint:"测试",
		Metric:"测试指标",
		CounterType:GAUGE,
		Step:5,
	}

	ticker:=time.NewTicker(time.Second*5)

	for i:=0;i<100;i++{
		select{
		case <-ticker.C:
			agent.Push(i)
		}
	}
}