package estimator

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

var testNS1 = NodeStatus{
	Timestamp:      time.Now(),
	CPUSockets:     2,
	CPUCores:       4,
	CPUUsages:      [][]float64{{10.0, 10.0, 10.0, 10.0}, {10.0, 10.0, 10.0, 10.0}},
	CPUTemps:       [][]float64{{30.0, 30.0, 30.0, 30.0}, {30.0, 30.0, 30.0, 30.0}},
	AmbientSensors: 2,
	AmbientTemps:   []float64{20.0, 20.0},
}

func getFnCopy(s NodeStatus, err error) func(ctx context.Context) (NodeStatus, error) {
	return func(_ context.Context) (NodeStatus, error) {
		return s, err
	}
}

func TestFakeNodeMonitor_FetchStatus(t *testing.T) {
	type fields struct {
		FetchFunc func(ctx context.Context) (NodeStatus, error)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    NodeStatus
		wantErr bool
	}{
		{"FetchFunc=nil", fields{nil}, args{context.Background()}, NodeStatus{}, true},
		{"ok", fields{getFnCopy(testNS1, nil)}, args{context.Background()}, testNS1, false},
		{"err", fields{getFnCopy(NodeStatus{}, errors.New(""))}, args{context.Background()}, NodeStatus{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FakeNodeMonitor{
				FetchFunc: tt.fields.FetchFunc,
			}
			got, err := m.FetchStatus(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FakeNodeMonitor.FetchStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FakeNodeMonitor.FetchStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
