package estimator

import (
	"context"
	"testing"
	"time"
)

func TestNodes_Len(t *testing.T) {
	type op int
	const (
		nop = iota
		add
		del
	)
	type action struct {
		Op   op
		Name string
		Node *Node
	}
	tests := []struct {
		name   string
		nodes  *Nodes
		action []action
		want   int
	}{
		{"0", &Nodes{}, []action{}, 0},
		{"0,0", &Nodes{}, []action{
			{add, "n1", nil},
		}, 0},
		{"0,1", &Nodes{}, []action{
			{add, "n1", NewNode("n1", nil, time.Second, nil)},
		}, 1},
		{"0,1,1", &Nodes{}, []action{
			{add, "n1", NewNode("n1", nil, time.Second, nil)},
			{add, "n1", NewNode("n1", nil, time.Second, nil)},
		}, 1},
		{"0,1,2", &Nodes{}, []action{
			{add, "n1", NewNode("n1", nil, time.Second, nil)},
			{add, "n2", NewNode("n2", nil, time.Second, nil)},
		}, 2},
		{"0,1,0", &Nodes{}, []action{
			{add, "n1", NewNode("n1", nil, time.Second, nil)},
			{del, "n1", nil},
		}, 0},
		{"0,1,1", &Nodes{}, []action{
			{add, "n1", NewNode("n1", nil, time.Second, nil)},
			{del, "n2", nil},
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, act := range tt.action {
				switch act.Op {
				case add:
					tt.nodes.Add(act.Name, act.Node)
				case del:
					tt.nodes.Delete(act.Name)
				}
			}
			if got := tt.nodes.Len(); got != tt.want {
				t.Errorf("Nodes.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_FetchStatus(t *testing.T) {
	testNodeStatus := NewNodeStatus()
	tests := []struct {
		name    string
		node    *Node
		want    *NodeStatus
		wantErr bool
	}{
		{"nm==nil", &Node{
			Name:     "n1",
			monitors: nil,
		}, nil, false},
		{"nm!=nil", &Node{
			Name: "n1",
			monitors: []NodeMonitor{
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return nil }},
			},
		}, testNodeStatus, false},
		{"all ok", &Node{
			Name: "n1",
			monitors: []NodeMonitor{
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return nil }},
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return nil }},
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return nil }},
			},
		}, nil, false},
		{"some NodeMonitor return error", &Node{
			Name: "n1",
			monitors: []NodeMonitor{
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return nil }},
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return ErrNodeMonitor }},
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return ErrNodeMonitor }},
				&FakeNodeMonitor{FetchFunc: func(context.Context, *NodeStatus) error { return nil }},
			},
		}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.FetchStatus(context.Background(), testNodeStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.FetchStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNode_Predict(t *testing.T) {
	type args struct {
		requestCPUMilli int
		status          *NodeStatus
	}
	tests := []struct {
		name     string
		node     *Node
		args     args
		wantWatt float64
		wantErr  bool
	}{
		{"pcp==nil", &Node{
			Name:        "n1",
			pcPredictor: nil,
		}, args{2000, newNodeStatus(20.0, 20.0, 0.0)}, 0.0, true},
		{"pcp!=nil", &Node{
			Name:        "n1",
			pcPredictor: &FakePCPredictor{PredictFunc: PredictPCFnDummy},
		}, args{2000, newNodeStatus(20.0, 10.0, 0.0)}, 50.0, false},
		{"failed", &Node{
			Name:        "n1",
			pcPredictor: &FakePCPredictor{PredictFunc: func(context.Context, int, *NodeStatus) (float64, error) { return 0.0, ErrPCPredictor }},
		}, args{2000, newNodeStatus(20.0, 20.0, 0.0)}, 0.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWatt, err := tt.node.Predict(context.Background(), tt.args.requestCPUMilli, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.Predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWatt != tt.wantWatt {
				t.Errorf("Node.Predict() = %v, want %v", gotWatt, tt.wantWatt)
			}
		})
	}
}
