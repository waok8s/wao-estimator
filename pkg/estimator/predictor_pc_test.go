package estimator

import (
	"context"
	"testing"
)

func newNodeStatus(cpuUsage, ambientTemp float64) *NodeStatus {
	v := NewNodeStatus()
	NodeStatusSetCPUUsage(v, cpuUsage)
	NodeStatusSetAmbientTemp(v, ambientTemp)
	return v
}

func TestFakePCPredictor_Predict(t *testing.T) {
	type fields struct {
		PredictFunc func(ctx context.Context, requestCPUMilli int, status *NodeStatus) (watt float64, err error)
	}
	type args struct {
		ctx             context.Context
		requestCPUMilli int
		status          *NodeStatus
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantWatt float64
		wantErr  bool
	}{
		{"PredictFunc=nil", fields{nil}, args{context.Background(), 2000, newNodeStatus(30.0, 20.0)}, 0.0, true},
		{"70", fields{PredictPCFnDummy}, args{context.Background(), 2000, newNodeStatus(30.0, 20.0)}, 70.0, false},
		{"80", fields{PredictPCFnDummy}, args{context.Background(), 2000, newNodeStatus(35.0, 25.0)}, 80.0, false},
		{"err", fields{PredictPCFnDummy}, args{context.Background(), 2000, NewNodeStatus()}, 0.0, true},
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
