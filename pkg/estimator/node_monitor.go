package estimator

import (
	"context"
	"fmt"
)

type NodeMonitor interface {
	FetchStatus(ctx context.Context) (*NodeStatus, error)
}

type FakeNodeMonitor struct {
	GetFunc func(ctx context.Context) (*NodeStatus, error)
}

var _ NodeMonitor = (*FakeNodeMonitor)(nil)

func (m *FakeNodeMonitor) FetchStatus(ctx context.Context) (*NodeStatus, error) {
	return m.GetFunc(ctx)
}

type RedfishNodeMonitor struct {
	// TODO
}

var _ NodeMonitor = (*RedfishNodeMonitor)(nil)

func (m *RedfishNodeMonitor) FetchStatus(ctx context.Context) (*NodeStatus, error) {
	// TODO
	return nil, fmt.Errorf("not yet implemented (%w)", ErrNodeMonitor)
}
