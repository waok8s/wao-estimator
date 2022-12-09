package estimator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type NodeStatus struct {
	Timestamp time.Time

	// CPUSockets is the number of sockets the node has.
	// e.g. 2
	CPUSockets int
	// CPUCores is the number of cores per CPU.
	// e.g. 4
	CPUCores int
	// CPUUsages is the list of CPU usages per core.
	// e.g. [[10.0, 10.0, 10.0, 10.0], [20.0, 20.0, 20.0, 20.0]]
	CPUUsages [][]float64
	// CPUTemps is the list of CPU temperatures in Celsius per core.
	// e.g. [[30.0, 30.0, 30.0, 30.0], [50.0, 50.0, 50.0, 50.0]]
	CPUTemps [][]float64

	// AmbientSensors is the number of ambient temperature sensors the node has.
	// e.g. 4
	AmbientSensors int
	// AmbientTemps is the list of temperatures in Celsius.
	// e.g. [20.0, 20.0, 20.0, 20.0]
	AmbientTemps []float64
}

func (s *NodeStatus) AverageCPUUsage() (float64, error) {
	var c int
	var sum float64
	for _, temps := range s.CPUUsages {
		for _, e := range temps {
			sum += e
			c++
		}
	}
	if c == 0 {
		return sum, ErrNodeStatus
	}
	return sum / float64(c), nil
}

func (s *NodeStatus) AverageCPUTemp() (float64, error) {
	var c int
	var sum float64
	for _, temps := range s.CPUTemps {
		for _, e := range temps {
			sum += e
			c++
		}
	}
	if c == 0 {
		return sum, ErrNodeStatus
	}
	return sum / float64(c), nil
}

func (s *NodeStatus) AverageAmbientTemp() (float64, error) {
	var c int
	var sum float64
	for _, e := range s.AmbientTemps {
		sum += e
		c++
	}
	if c == 0 {
		return sum, ErrNodeStatus
	}
	return sum / float64(c), nil
}

// deepcopy copies structs.
// Both dst and src must be pointers.
func deepcopy(dst any, src any) error {
	if dst == nil || src == nil {
		return errors.New("both dst and src cannot be nil")
	}
	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := json.Unmarshal(bytes, dst); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

type Node struct {
	Name string

	mu     sync.Mutex
	stopCh chan struct{}

	Monitor    NodeMonitor
	nmInterval time.Duration
	status     *NodeStatus

	PCPredictor PowerConsumptionPredictor
}

func NewNode(name string, nm NodeMonitor, nodeStatusRefreshInterval time.Duration, pcp PowerConsumptionPredictor) *Node {
	n := Node{
		Name:        name,
		Monitor:     nm,
		PCPredictor: pcp,
		nmInterval:  nodeStatusRefreshInterval,
	}
	return &n
}

func (n *Node) start() {
	lg.Info().Msgf("Node.start() Name=%v", n.Name)
	go func() {
		for {
			select {
			case <-n.stopCh:
				return
			case <-time.After(n.nmInterval):
				timeout := n.nmInterval / 2
				ctx, cncl := context.WithTimeout(context.Background(), timeout)
				status, err := n.Monitor.FetchStatus(ctx)
				cncl()
				if err != nil {
					lg.Error().Msgf("could not fetch NodeStatus: %v", err)
				}
				n.mu.Lock()
				n.status = status
				n.mu.Unlock()
			}
		}
	}()
}

func (n *Node) stop() {
	lg.Info().Msgf("Node.stop() Name=%v", n.Name)
	close(n.stopCh)
}

func (n *Node) GetStatus() *NodeStatus {
	n.mu.Lock()
	var s NodeStatus
	deepcopy(&s, n.status)
	return &s
}

type Nodes struct {
	count int32
	m     sync.Map
}

func (m *Nodes) Get(k string) (*Node, bool) {
	v, ok := m.m.Load(k)
	if !ok {
		return nil, false
	}
	return v.(*Node), true
}

func (m *Nodes) Add(k string, v *Node) bool {
	_, ok := m.m.Load(k)
	if ok {
		return false
	}
	v.start()
	m.m.Store(k, v)
	atomic.AddInt32(&m.count, 1)
	return true
}

func (m *Nodes) Delete(k string) {
	n, ok := m.Get(k)
	if ok {
		n.stop()
	}
	m.m.Delete(k)
}

func (m *Nodes) Range(f func(k string, v *Node) bool) {
	m.m.Range(func(kk any, vv any) bool { return f(kk.(string), vv.(*Node)) })
}

func (m *Nodes) Len() int {
	// FIXME: this is O(N)
	var i int
	m.Range(func(_ string, _ *Node) bool {
		i++
		return true
	})
	return i
}
