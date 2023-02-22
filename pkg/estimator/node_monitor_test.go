package estimator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func getFnCopy(s *NodeStatus, err error) func(ctx context.Context, base *NodeStatus) (*NodeStatus, error) {
	return func(context.Context, *NodeStatus) (*NodeStatus, error) {
		return s, err
	}
}

func TestFakeNodeMonitor_FetchStatus(t *testing.T) {
	var testNodeStatus = NewNodeStatus()

	type fields struct {
		FetchFunc func(ctx context.Context, base *NodeStatus) (*NodeStatus, error)
	}
	type args struct {
		ctx  context.Context
		base *NodeStatus
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *NodeStatus
		wantErr bool
	}{
		{"FetchFunc=nil", fields{nil}, args{context.Background(), nil}, nil, true},
		{"ok", fields{getFnCopy(testNodeStatus, nil)}, args{context.Background(), testNodeStatus}, testNodeStatus, false},
		{"err", fields{getFnCopy(nil, errors.New(""))}, args{context.Background(), testNodeStatus}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FakeNodeMonitor{
				FetchFunc: tt.fields.FetchFunc,
			}
			got, err := m.FetchStatus(tt.args.ctx, tt.args.base)
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
