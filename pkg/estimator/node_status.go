package estimator

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

type NodeStatusKey string

type NodeStatus struct {
	timestamp time.Time
	data      sync.Map
}

func NewNodeStatus() *NodeStatus { return &NodeStatus{timestamp: time.Now()} }

func (s *NodeStatus) Timestamp() time.Time { return s.timestamp }

func (s *NodeStatus) Get(k NodeStatusKey) (string, bool) {
	v, ok := s.data.Load(k)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (s *NodeStatus) Set(k NodeStatusKey, v string) { s.data.Store(k, v) }

func (s *NodeStatus) Delete(k NodeStatusKey) { s.data.Delete(k) }

func (s *NodeStatus) Range(f func(k NodeStatusKey, v string) bool) {
	s.data.Range(func(kk any, vv any) bool { return f(kk.(NodeStatusKey), vv.(string)) })
}

const NodeStatusCPUUsage NodeStatusKey = "cpuUsage"

func NodeStatusSetCPUUsage(s *NodeStatus, cpuUsage float64) {
	s.Set(NodeStatusCPUUsage, strconv.FormatFloat(cpuUsage, 'f', -1, 64))
}

func NodeStatusGetCPUUsage(s *NodeStatus) (float64, error) {
	v, ok := s.Get(NodeStatusCPUUsage)
	if !ok {
		return -math.MaxFloat64, fmt.Errorf("key=%+v not found", NodeStatusCPUUsage)
	}
	vv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return -math.MaxFloat64, err
	}
	return vv, nil
}

const NodeStatusAmbientTemp NodeStatusKey = "ambientTemp"

func NodeStatusSetAmbientTemp(s *NodeStatus, ambientTemp float64) {
	s.Set(NodeStatusAmbientTemp, strconv.FormatFloat(ambientTemp, 'f', -1, 64))
}

func NodeStatusGetAmbientTemp(s *NodeStatus) (float64, error) {
	v, ok := s.Get(NodeStatusAmbientTemp)
	if !ok {
		return -math.MaxFloat64, fmt.Errorf("key=%+v not found", NodeStatusAmbientTemp)
	}
	vv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return -math.MaxFloat64, err
	}
	return vv, nil
}

const NodeStatusStaticPressureDiff NodeStatusKey = "staticPressureDiff"

func NodeStatusSetStaticPressureDiff(s *NodeStatus, spd float64) {
	s.Set(NodeStatusStaticPressureDiff, strconv.FormatFloat(spd, 'f', -1, 64))
}

func NodeStatusGetStaticPressureDiff(s *NodeStatus) (float64, error) {
	v, ok := s.Get(NodeStatusStaticPressureDiff)
	if !ok {
		return -math.MaxFloat64, fmt.Errorf("key=%+v not found", NodeStatusStaticPressureDiff)
	}
	vv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return -math.MaxFloat64, err
	}
	return vv, nil
}

const NodeStatusLogicalProcessors NodeStatusKey = "logicalProcessors"

func NodeStatusSetLogicalProcessors(s *NodeStatus, num int) {
	s.Set(NodeStatusLogicalProcessors, strconv.Itoa(num))
}

func NodeStatusGetLogicalProcessors(s *NodeStatus) (int, error) {
	v, ok := s.Get(NodeStatusLogicalProcessors)
	if !ok {
		return 0, fmt.Errorf("key=%+v not found", NodeStatusLogicalProcessors)
	}
	vv, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return vv, nil
}
