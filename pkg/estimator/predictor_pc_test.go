package estimator

import (
	"context"
	"testing"
)

func TestFakePCPredictor_Predict(t *testing.T) {
	type fields struct {
		PredictFunc func(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error)
	}
	type args struct {
		ctx             context.Context
		requestCPUMilli int
		status          NodeStatus
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantWatt float64
		wantErr  bool
	}{
		{"70", fields{PredictPCFnDummy}, args{context.Background(), 2000, NodeStatus{
			CPUUsages:    [][]float64{{30.0}},
			AmbientTemps: []float64{20.0},
		}}, 70.0, false},
		{"80", fields{PredictPCFnDummy}, args{context.Background(), 2000, NodeStatus{
			CPUUsages:    [][]float64{{35.0}},
			AmbientTemps: []float64{25.0},
		}}, 80.0, false},
		{"err", fields{PredictPCFnDummy}, args{context.Background(), 2000, NodeStatus{}}, 0.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FakePCPredictor{
				PredictFunc: tt.fields.PredictFunc,
			}
			gotWatt, err := p.Predict(tt.args.ctx, tt.args.requestCPUMilli, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("FakePCPredictor.Predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWatt != tt.wantWatt {
				t.Errorf("FakePCPredictor.Predict() = %v, want %v", gotWatt, tt.wantWatt)
			}
		})
	}
}
