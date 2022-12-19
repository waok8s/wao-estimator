package estimator

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
)

type Estimator struct {
	Nodes *Nodes
	init  sync.Once
}

func (e *Estimator) initOnce() {
	e.init.Do(func() {
		if e.Nodes == nil {
			e.Nodes = &Nodes{}
		}
	})
}

// EstimatePowerConsumption is a thread-safe function that
// estimates power consumption with the given parameters.
func (e *Estimator) EstimatePowerConsumption(ctx context.Context, cpuMilli, numWorkloads int) ([]float64, error) {
	e.initOnce()

	// init wattMatrix[node][workload]
	wattMatrix := make([][]float64, e.Nodes.Len())
	for i := range wattMatrix {
		wattMatrix[i] = make([]float64, numWorkloads)
	}

	// prediction
	for i := 0; i < numWorkloads; i++ {
		j := 0
		e.Nodes.Range(func(_ string, node *Node) bool {
			watt, err := node.Predict(ctx, cpuMilli*(i+1), node.GetStatus())
			if err != nil {
				watt = math.MaxFloat64
			}
			wattMatrix[i][j] = watt
			j++
			return true
		})
	}

	// search
	minCosts, err := ComputeLeastCostsFn(e.Nodes.Len(), numWorkloads, wattMatrix)
	if err != nil {
		return nil, err
	}

	lg.Debug().Msgf("minCosts=%v", minCosts)

	return minCosts, nil
}

func (e *Estimator) stop() {
	e.initOnce()

	lg.Info().Msgf("Estimator.stop()")
	e.Nodes.Range(func(k string, _ *Node) bool {
		e.Nodes.Delete(k)
		return true
	})
}

type Estimators struct {
	c int32
	m sync.Map
}

func (m *Estimators) Get(k string) (*Estimator, bool) {
	v, ok := m.m.Load(k)
	if !ok {
		return nil, ok
	}
	return v.(*Estimator), true
}

func (m *Estimators) Add(k string, v *Estimator) bool {
	if v == nil {
		return false
	}
	_, ok := m.m.Load(k)
	if ok {
		return false
	}
	m.m.Store(k, v)
	atomic.AddInt32(&(m.c), 1)
	return true
}

func (m *Estimators) Delete(k string) {
	v, ok := m.Get(k)
	if !ok {
		return
	}
	v.stop()
	atomic.AddInt32(&(m.c), -1)
	m.m.Delete(k)
}

func (m *Estimators) Range(f func(k string, v *Estimator) bool) {
	m.m.Range(func(kk any, vv any) bool { return f(kk.(string), vv.(*Estimator)) })
}

func (m *Estimators) Len() int { return int(m.c) }
