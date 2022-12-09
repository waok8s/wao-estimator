package estimator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidRequest      = errors.New("ErrInvalidRequest")
	ErrEstimator           = errors.New("ErrEstimator")
	ErrEstimatorNotFound   = errors.New("ErrEstimatorNotFound")
	ErrNodeMonitor         = errors.New("ErrNodeMonitor")
	ErrNodeMonitorNotFound = errors.New("ErrNodeMonitorNotFound")
	ErrNodeStatus          = errors.New("ErrNodeStatus")
	ErrNodeStatusNotFound  = errors.New("ErrNodeStatusNotFound")
	ErrPCPredictor         = errors.New("ErrPCPredictor")
	ErrPCPredictorNotFound = errors.New("ErrPCPredictorNotFound")
)

type PowerConsumptionPredictor interface {
	Predict(ctx context.Context, requestCPUMilli int, status *NodeStatus) (watt float64, err error)
}

type FakePCPredictor struct {
	PredictFunc func(ctx context.Context, requestCPUMilli int, status *NodeStatus) (watt float64, err error)
}

var _ PowerConsumptionPredictor = (*FakePCPredictor)(nil)

func (p *FakePCPredictor) Predict(ctx context.Context, requestCPUMilli int, status *NodeStatus) (watt float64, err error) {
	return p.PredictFunc(ctx, requestCPUMilli, status)
}

type MLServerPCPredictor struct {
	// TODO
}

var _ PowerConsumptionPredictor = (*MLServerPCPredictor)(nil)

func (p *MLServerPCPredictor) Predict(ctx context.Context, requestCPUMilli int, status *NodeStatus) (watt float64, err error) {
	// TODO
	return 0.0, fmt.Errorf("not yet implemented (%w)", ErrPCPredictor)
}

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

type NodeMonitor interface {
	FetchStatus(ctx context.Context) (*NodeStatus, error)
}

type FakeNodeMonitor struct {
	GetFunc func(ctx context.Context) (*NodeStatus, error)
}

var _ NodeMonitor = (*FakeNodeMonitor)(nil)

func (m *FakeNodeMonitor) FetchStatus(ctx context.Context) (*NodeStatus, error) {
	return m.GetFunc(ctx)
}

type RedfishNodeMonitor struct {
	// TODO
}

var _ NodeMonitor = (*RedfishNodeMonitor)(nil)

func (m *RedfishNodeMonitor) FetchStatus(ctx context.Context) (*NodeStatus, error) {
	// TODO
	return nil, fmt.Errorf("not yet implemented (%w)", ErrNodeMonitor)
}

func NewNode(name string, nm NodeMonitor, pcp PowerConsumptionPredictor, nodeStatusRefreshInterval time.Duration) *Node {
	n := Node{
		Name:        name,
		Monitor:     nm,
		PCPredictor: pcp,
		nmInterval:  nodeStatusRefreshInterval,
	}
	return &n
}

func (n *Node) start() {
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
					// TODO: print logs
				}
				n.mu.Lock()
				n.status = status
				n.mu.Unlock()
			}
		}
	}()
}

func (n *Node) stop() {
	close(n.stopCh)
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

type Estimator struct {
	Nodes *Nodes
}

func NewEstimator(nodes *Nodes) *Estimator {
	return &Estimator{Nodes: nodes}
}

// EstimatePowerConsumption is a thread-safe function that
// estimates power consumption with the given parameters.
func (e *Estimator) EstimatePowerConsumption(ctx context.Context, cpuMilli, numWorkloads int) ([]float64, error) {

	// init wattMatrix[node][workload]
	wattMatrix := make([][]float64, e.Nodes.Len())
	for i, _ := range wattMatrix {
		wattMatrix[i] = make([]float64, numWorkloads)
	}

	// prediction
	for i := 0; i < numWorkloads; i++ {
		j := 0
		e.Nodes.Range(func(_ string, node *Node) bool {
			watt, err := node.PCPredictor.Predict(ctx, cpuMilli*(i+1), node.GetStatus())
			if err != nil {
				watt = math.MaxFloat64
			}
			wattMatrix[i][j] = watt
			j++
			return true
		})
	}

	// search
	// TODO: exhaustive search wattMatrix[][] -> wattIncrease[]
	wattIncrease := make([]float64, numWorkloads)
	return wattIncrease, nil
}

type Estimators struct {
	m sync.Map
}

func (m *Estimators) Get(k string) (*Estimator, bool) {
	v, ok := m.m.Load(k)
	if !ok {
		return nil, ok
	}
	x, ok := v.(*Estimator)
	return x, ok
}

func (m *Estimators) Set(k string, v *Estimator) {
	m.m.Store(k, v)
}
