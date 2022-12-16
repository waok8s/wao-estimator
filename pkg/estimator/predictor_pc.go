package estimator

import (
	"context"
	"fmt"
)

type PowerConsumptionPredictor interface {
	Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error)
}

type FakePCPredictor struct {
	PredictFunc func(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error)
}

var _ PowerConsumptionPredictor = (*FakePCPredictor)(nil)

func (p *FakePCPredictor) Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error) {
	if p.PredictFunc == nil {
		return 0.0, fmt.Errorf("PredictFunc not set (%w)", ErrPCPredictor)
	}
	return p.PredictFunc(ctx, requestCPUMilli, status)
}

// PredictPCFnDummy returns ( mCPU/100 + avg(CPUUsage) + avg(AmbientTemp) )
func PredictPCFnDummy(_ context.Context, mcpu int, status NodeStatus) (float64, error) {
	aat, err := status.AverageAmbientTemp()
	if err != nil {
		return 0.0, err
	}
	acu, err := status.AverageCPUUsage()
	if err != nil {
		return 0.0, err
	}
	w := float64(mcpu)/100 + aat + acu
	return w, nil
}

type MLServerPCPredictor struct {
	// TODO
}

var _ PowerConsumptionPredictor = (*MLServerPCPredictor)(nil)

func (p *MLServerPCPredictor) Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error) {
	// TODO
	return 0.0, fmt.Errorf("not yet implemented (%w)", ErrPCPredictor)
}
