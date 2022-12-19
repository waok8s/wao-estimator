package estimator

import (
	"context"
	"fmt"
)

type NodeMonitor interface {
	FetchStatus(ctx context.Context) (NodeStatus, error)
}

type FakeNodeMonitor struct {
	FetchFunc func(ctx context.Context) (NodeStatus, error)
}

var _ NodeMonitor = (*FakeNodeMonitor)(nil)

func (m *FakeNodeMonitor) FetchStatus(ctx context.Context) (NodeStatus, error) {
	if m.FetchFunc == nil {
		return NodeStatus{}, fmt.Errorf("FetchFunc not set (%w)", ErrNodeMonitor)
	}
	return m.FetchFunc(ctx)
}

type RedfishNodeMonitor struct {
	// TODO
}

var _ NodeMonitor = (*RedfishNodeMonitor)(nil)

func (m *RedfishNodeMonitor) FetchStatus(ctx context.Context) (NodeStatus, error) {
	// TODO
	return NodeStatus{}, fmt.Errorf("not yet implemented (%w)", ErrNodeMonitor)
}
