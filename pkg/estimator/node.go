package estimator

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	Name string

	mu     sync.Mutex
	stopCh chan struct{}

	monitors   []NodeMonitor
	nmInterval time.Duration
	status     *NodeStatus

	pcPredictor PowerConsumptionPredictor
}

var _ NodeMonitor = (*Node)(nil)
var _ PowerConsumptionPredictor = (*Node)(nil)

func (n *Node) FetchStatus(ctx context.Context, base *NodeStatus) error {
	if base == nil {
		base = NewNodeStatus()
	}
	for i, nm := range n.monitors {
		err := nm.FetchStatus(ctx, base)
		if err != nil {
			lg.Warn().Msgf("FetchStatus failed NodeMonitor[%d] err=%v", i, err)
		}
	}
	return nil
}

func (n *Node) Predict(ctx context.Context, requestCPUMilli int, status *NodeStatus) (watt float64, err error) {
	if n.pcPredictor == nil {
		return 0.0, ErrPCPredictorNotFound
	}
	return n.pcPredictor.Predict(ctx, requestCPUMilli, status)
}

func NewNode(name string, nms []NodeMonitor, nodeStatusRefreshInterval time.Duration, pcp PowerConsumptionPredictor) *Node {
	n := Node{
		Name:        name,
		stopCh:      make(chan struct{}),
		monitors:    nms,
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
		status := NewNodeStatus()
		_ = n.FetchStatus(ctx, status) // this does not return errors
		cncl()
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

// GetStatus returns current status of the Node,
// if Node.status is nil then returns a new empty NodeStatus.
func (n *Node) GetStatus() *NodeStatus {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.status == nil {
		return NewNodeStatus()
	}
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
