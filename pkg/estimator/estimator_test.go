package estimator

import (
	"math"
	"testing"
)

func TestEstimators_Len(t *testing.T) {
	type op int
	const (
		nop = iota
		add
		del
	)
	type action struct {
		Op   op
		Name string
		Est  *Estimator
	}
	tests := []struct {
		name   string
		ests   *Estimators
		action []action
		want   int
	}{
		{"0", &Estimators{}, []action{}, 0},
		{"0,0", &Estimators{}, []action{
			{add, "n1", nil},
		}, 0},
		{"0,1", &Estimators{}, []action{
			{add, "n1", &Estimator{Nodes: nil}},
		}, 1},
		{"0,1,1", &Estimators{}, []action{
			{add, "n1", &Estimator{Nodes: nil}},
			{add, "n1", &Estimator{Nodes: nil}},
		}, 1},
		{"0,1,2", &Estimators{}, []action{
			{add, "n1", &Estimator{Nodes: nil}},
			{add, "n2", &Estimator{Nodes: nil}},
		}, 2},
		{"0,1,0", &Estimators{}, []action{
			{add, "n1", &Estimator{Nodes: nil}},
			{del, "n1", nil},
		}, 0},
		{"0,1,1", &Estimators{}, []action{
			{add, "n1", &Estimator{Nodes: nil}},
			{del, "n2", nil},
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, act := range tt.action {
				switch act.Op {
				case add:
					tt.ests.Add(act.Name, act.Est)
				case del:
					tt.ests.Delete(act.Name)
				}
			}
			if got := tt.ests.Len(); got != tt.want {
				t.Errorf("Estimators.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_patchWattMatrix(t *testing.T) {
	type args struct {
		wattMatrix [][]float64
	}
	tests := []struct {
		name           string
		args           args
		wantWattMatrix [][]float64
	}{
		{"no error", args{wattMatrix: testCostsReq1}, testCostsReq1},
		{"an error", args{wattMatrix: [][]float64{
			{1, 2, 3},
			{1, 2, math.Inf(1)},
			{1, 2, 3},
		}}, [][]float64{
			{1, 2, 3},
			{0, math.Inf(1), math.Inf(1)},
			{1, 2, 3},
		}},
		{"all error", args{wattMatrix: [][]float64{
			{math.Inf(1), 2, 3},
			{1, math.Inf(1), 3},
			{1, 2, math.Inf(1)},
		}}, [][]float64{
			{0, math.Inf(1), math.Inf(1)},
			{0, math.Inf(1), math.Inf(1)},
			{0, math.Inf(1), math.Inf(1)},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patchWattMatrix(tt.args.wattMatrix)
			for i := range tt.args.wattMatrix {
				for j := range tt.args.wattMatrix[i] {
					if tt.args.wattMatrix[i][j] != tt.wantWattMatrix[i][j] {
						t.Errorf("want=%v but got=%v", tt.wantWattMatrix, tt.args.wattMatrix)
					}
				}
			}
		})
	}
}
