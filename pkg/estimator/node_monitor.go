package estimator

import (
	"context"
	"fmt"
)

type NodeMonitor interface {
	FetchStatus(ctx context.Context, base *NodeStatus) error
}

type FakeNodeMonitor struct {
	FetchFunc func(ctx context.Context, base *NodeStatus) error
}

var _ NodeMonitor = (*FakeNodeMonitor)(nil)

func (m *FakeNodeMonitor) FetchStatus(ctx context.Context, base *NodeStatus) error {
	if m.FetchFunc == nil {
		return fmt.Errorf("FetchFunc not set (%w)", ErrNodeMonitor)
	}
	return m.FetchFunc(ctx, base)
}

type RedfishNodeMonitor struct {
	// TODO
}

var _ NodeMonitor = (*RedfishNodeMonitor)(nil)

func (m *RedfishNodeMonitor) FetchStatus(ctx context.Context, base *NodeStatus) error {
	// TODO
	return fmt.Errorf("not yet implemented (%w)", ErrNodeMonitor)
}
