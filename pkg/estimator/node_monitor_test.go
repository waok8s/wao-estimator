package estimator

import (
	"context"
	"errors"
	"testing"
)

func getFnPassthrough(err error) func(ctx context.Context, base *NodeStatus) error {
	return func(context.Context, *NodeStatus) error {
		return err
	}
}

func TestFakeNodeMonitor_FetchStatus(t *testing.T) {
	type fields struct {
		FetchFunc func(ctx context.Context, base *NodeStatus) error
	}
	type args struct {
		ctx  context.Context
		base *NodeStatus
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"FetchFunc=nil", fields{nil}, args{context.Background(), NewNodeStatus()}, true},
		{"ok", fields{getFnPassthrough(nil)}, args{context.Background(), NewNodeStatus()}, false},
		{"err", fields{getFnPassthrough(errors.New(""))}, args{context.Background(), NewNodeStatus()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FakeNodeMonitor{
				FetchFunc: tt.fields.FetchFunc,
			}
			err := m.FetchStatus(tt.args.ctx, tt.args.base)
			if (err != nil) != tt.wantErr {
				t.Errorf("FakeNodeMonitor.FetchStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
