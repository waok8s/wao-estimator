package estimator

import (
	"context"
	"math"
	"sync"
)

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
	for i := range wattMatrix {
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

func (e *Estimator) stop() {
	lg.Info().Msgf("Estimator.stop()")
	e.Nodes.Range(func(k string, _ *Node) bool {
		e.Nodes.Delete(k)
		return true
	})
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

func (m *Estimators) Add(k string, v *Estimator) bool {
	_, ok := m.m.Load(k)
	if ok {
		return false
	}
	m.m.Store(k, v)
	return true
}

func (m *Estimators) Delete(k string) {
	n, ok := m.Get(k)
	if ok {
		n.stop()
	}
	m.m.Delete(k)
}
