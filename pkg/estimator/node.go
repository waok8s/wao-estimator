package estimator

import (
	"context"
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

type Node struct {
	Name string

	mu     sync.Mutex
	stopCh chan struct{}

	monitor    NodeMonitor
	nmInterval time.Duration
	status     NodeStatus

	pcPredictor PowerConsumptionPredictor
}

var _ NodeMonitor = (*Node)(nil)
var _ PowerConsumptionPredictor = (*Node)(nil)

func (n *Node) FetchStatus(ctx context.Context) (NodeStatus, error) {
	if n.monitor == nil {
		return NodeStatus{}, ErrNodeMonitorNotFound
	}
	return n.monitor.FetchStatus(ctx)
}

func (n *Node) Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error) {
	if n.pcPredictor == nil {
		return 0.0, ErrPCPredictorNotFound
	}
	return n.pcPredictor.Predict(ctx, requestCPUMilli, status)
}

func NewNode(name string, nm NodeMonitor, nodeStatusRefreshInterval time.Duration, pcp PowerConsumptionPredictor) *Node {
	n := Node{
		Name:        name,
		stopCh:      make(chan struct{}),
		monitor:     nm,
		pcPredictor: pcp,
		nmInterval:  nodeStatusRefreshInterval,
	}
	lg.Info().Msgf("NewNode() n=%+v", &n)
	return &n
}

func (n *Node) start() {
	lg.Info().Msgf("Node.start() Name=%v", n.Name)

	updateStatus := func() {
		timeout := n.nmInterval / 2
		ctx, cncl := context.WithTimeout(context.Background(), timeout)
		status, err := n.FetchStatus(ctx)
		cncl()
		if err != nil {
			lg.Error().Msgf("could not fetch NodeStatus: %v", err)
			return
		}
		n.mu.Lock()
		n.status = status
		n.mu.Unlock()
	}

	updateStatus() // first time exec

	go func() {
		for {
			select {
			case <-n.stopCh:
				return
			case <-time.After(n.nmInterval):
				updateStatus()
			}
		}
	}()
}

func (n *Node) stop() {
	lg.Info().Msgf("Node.stop() Name=%v", n.Name)
	close(n.stopCh)
}

func (n *Node) GetStatus() NodeStatus {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.status
}

type Nodes struct {
	c int32
	m sync.Map
}

func (m *Nodes) Get(k string) (*Node, bool) {
	v, ok := m.m.Load(k)
	if !ok {
		return nil, false
	}
	return v.(*Node), true
}

func (m *Nodes) Add(k string, v *Node) bool {
	if v == nil {
		return false
	}
	_, ok := m.m.Load(k)
	if ok {
		return false
	}
	m.m.Store(k, v)
	atomic.AddInt32(&(m.c), 1)
	v.start()
	return true
}

func (m *Nodes) Delete(k string) {
	v, ok := m.Get(k)
	if !ok {
		return
	}
	v.stop()
	atomic.AddInt32(&(m.c), -1)
	m.m.Delete(k)
}

func (m *Nodes) Range(f func(k string, v *Node) bool) {
	m.m.Range(func(kk any, vv any) bool { return f(kk.(string), vv.(*Node)) })
}

func (m *Nodes) Len() int { return int(m.c) }
