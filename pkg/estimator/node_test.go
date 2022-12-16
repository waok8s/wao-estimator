package estimator

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestNodeStatus_AverageCPUUsage(t *testing.T) {
	type fields struct {
		CPUSockets     int
		CPUCores       int
		CPUCoreUsages  [][]float64
		CPUCoreTemps   [][]float64
		AmbientSensors int
		AmbientTemps   []float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{"err_nil", fields{0, 0, nil, nil, 0, nil}, 0.0, true},
		{"err_0cpu0core", fields{0, 0, [][]float64{}, nil, 0, nil}, 0.0, true},
		{"1cpu8core", fields{0, 0, [][]float64{{20.0, 20.0, 30.0, 30.0, 40.0, 40.0, 50.0, 50.0}}, nil, 0, nil}, 35.0, false},
		{"4cpu8core", fields{0, 0, [][]float64{{20.0, 20.0}, {30.0, 30.0}, {40.0, 40.0}, {50.0, 50.0}}, nil, 0, nil}, 35.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NodeStatus{
				CPUSockets:     tt.fields.CPUSockets,
				CPUCores:       tt.fields.CPUCores,
				CPUUsages:      tt.fields.CPUCoreUsages,
				CPUTemps:       tt.fields.CPUCoreTemps,
				AmbientSensors: tt.fields.AmbientSensors,
				AmbientTemps:   tt.fields.AmbientTemps,
			}
			got, err := s.AverageCPUUsage()
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStatus.AverageCPUUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NodeStatus.AverageCPUUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStatus_AverageCPUTemp(t *testing.T) {
	type fields struct {
		CPUSockets     int
		CPUCores       int
		CPUCoreUsages  [][]float64
		CPUCoreTemps   [][]float64
		AmbientSensors int
		AmbientTemps   []float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{"err_nil", fields{0, 0, nil, nil, 0, nil}, 0.0, true},
		{"err_0cpu0core", fields{0, 0, nil, [][]float64{}, 0, nil}, 0.0, true},
		{"1cpu8core", fields{0, 0, nil, [][]float64{{20.0, 20.0, 30.0, 30.0, 40.0, 40.0, 50.0, 50.0}}, 0, nil}, 35.0, false},
		{"4cpu8core", fields{0, 0, nil, [][]float64{{20.0, 20.0}, {30.0, 30.0}, {40.0, 40.0}, {50.0, 50.0}}, 0, nil}, 35.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NodeStatus{
				CPUSockets:     tt.fields.CPUSockets,
				CPUCores:       tt.fields.CPUCores,
				CPUUsages:      tt.fields.CPUCoreUsages,
				CPUTemps:       tt.fields.CPUCoreTemps,
				AmbientSensors: tt.fields.AmbientSensors,
				AmbientTemps:   tt.fields.AmbientTemps,
			}
			got, err := s.AverageCPUTemp()
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStatus.AverageCPUTemp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NodeStatus.AverageCPUTemp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStatus_AverageAmbientTemp(t *testing.T) {
	type fields struct {
		CPUSockets     int
		CPUCores       int
		CPUCoreUsages  [][]float64
		CPUCoreTemps   [][]float64
		AmbientSensors int
		AmbientTemps   []float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{"err_nil", fields{0, 0, nil, nil, 0, nil}, 0.0, true},
		{"err_0sensor", fields{0, 0, nil, nil, 0, []float64{}}, 0.0, true},
		{"1sensor", fields{0, 0, nil, nil, 0, []float64{35.0}}, 35.0, false},
		{"4sensor", fields{0, 0, nil, nil, 0, []float64{20.0, 30.0, 40.0, 50.0}}, 35.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NodeStatus{
				CPUSockets:     tt.fields.CPUSockets,
				CPUCores:       tt.fields.CPUCores,
				CPUUsages:      tt.fields.CPUCoreUsages,
				CPUTemps:       tt.fields.CPUCoreTemps,
				AmbientSensors: tt.fields.AmbientSensors,
				AmbientTemps:   tt.fields.AmbientTemps,
			}
			got, err := s.AverageAmbientTemp()
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeStatus.AverageAmbientTemp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NodeStatus.AverageAmbientTemp() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	tests := []struct {
		name    string
		node    *Node
		want    NodeStatus
		wantErr bool
	}{
		{"nm==nil", &Node{
			Name:    "n1",
			monitor: nil,
		}, NodeStatus{}, true},
		{"nm!=nil", &Node{
			Name:    "n1",
			monitor: &FakeNodeMonitor{GetFunc: func(context.Context) (NodeStatus, error) { return testNS1, nil }},
		}, testNS1, false},
		{"failed", &Node{
			Name:    "n1",
			monitor: &FakeNodeMonitor{GetFunc: func(context.Context) (NodeStatus, error) { return NodeStatus{}, ErrNodeMonitor }},
		}, NodeStatus{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.node.FetchStatus(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.FetchStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.FetchStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_Predict(t *testing.T) {
	type args struct {
		requestCPUMilli int
		status          NodeStatus
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
		}, args{2000, testNS1}, 0.0, true},
		{"pcp!=nil", &Node{
			Name:        "n1",
			pcPredictor: &FakePCPredictor{PredictFunc: PredictPCFnDummy},
		}, args{2000, testNS1}, 50.0, false},
		{"failed", &Node{
			Name:        "n1",
			pcPredictor: &FakePCPredictor{PredictFunc: func(context.Context, int, NodeStatus) (float64, error) { return 0.0, ErrPCPredictor }},
		}, args{2000, testNS1}, 0.0, true},
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
