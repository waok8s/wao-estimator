package estimator

import (
	"context"
	"fmt"
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

	if e.Nodes.Len() == 0 {
		return nil, fmt.Errorf("no nodes available (%w)", ErrEstimatorNoNodesAvailable)
	}

	// init wattMatrix[node][workload]
	lg.Debug().Msgf("init wattMatrix[%d][%d]", e.Nodes.Len(), numWorkloads+1)
	wattMatrix := make([][]float64, e.Nodes.Len())
	for i := range wattMatrix {
		wattMatrix[i] = make([]float64, numWorkloads+1)
	}

	// prediction
	wg := sync.WaitGroup{}
	i := 0
	e.Nodes.Range(func(nodeName string, node *Node) bool {
		nodeIdx := i
		wg.Add(1)
		// NOTE: no need to sync, the goroutines below only write different slice elements
		go func() {
			defer wg.Done()
			for j := 0; j < numWorkloads+1; j++ {
				watt, err := node.Predict(ctx, cpuMilli*(j), node.GetStatus())
				if err != nil {
					lg.Warn().Msgf("node.Predict() for name=%s got error at wattMatrix[%d][%d] err=%v", node.Name, nodeIdx, j, err)
					watt = math.Inf(1)
				}
				lg.Debug().Msgf("call node.Predict() for name=%s wattMatrix[%d][%d] watt=%f", node.Name, nodeIdx, j, watt)
				wattMatrix[nodeIdx][j] = watt
			}
		}()
		i++
		return true
	})
	wg.Wait()
	lg.Debug().Msgf("wattMatrix=%v", wattMatrix)

	patchWattMatrix(wattMatrix)

	// search
	wattDiffs, err := toDiff(wattMatrix)
	if err != nil {
		return nil, err
	}
	lg.Debug().Msgf("wattDiffs=%v", wattDiffs)
	minCosts, err := ComputeLeastCostsFn(e.Nodes.Len(), numWorkloads, wattDiffs)
	if err != nil {
		return nil, err
	}
	lg.Debug().Msgf("minCosts=%v", minCosts)

	return minCosts, nil
}

func patchWattMatrix(wattMatrix [][]float64) {
	////////////////
	// Do some replacements to ensure toDiff() returns [+Inf ...] for error rows.
	////////////////
	// 1. detect errors: any row includes +Inf
	// e.g. wattMatrix=[[1 2 3] [1 +Inf 3]] then errs={1: struct{}}
	errs := map[int]struct{}{}
	for i, row := range wattMatrix {
		for _, elem := range row {
			if elem == math.Inf(1) {
				errs[i] = struct{}{}
			}
		}
	}
	// 2. set errors: set error rows [0 +Inf +Inf ...]
	for i := range wattMatrix {
		if _, isErr := errs[i]; isErr {
			for j := range wattMatrix[i] {
				if j == 0 {
					wattMatrix[i][j] = 0
				} else {
					wattMatrix[i][j] = math.Inf(1)
				}
			}
		}
	}
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
