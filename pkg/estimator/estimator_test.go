package estimator

import (
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
			{add, "n1", NewEstimator(nil)},
		}, 1},
		{"0,1,1", &Estimators{}, []action{
			{add, "n1", NewEstimator(nil)},
			{add, "n1", NewEstimator(nil)},
		}, 1},
		{"0,1,2", &Estimators{}, []action{
			{add, "n1", NewEstimator(nil)},
			{add, "n2", NewEstimator(nil)},
		}, 2},
		{"0,1,0", &Estimators{}, []action{
			{add, "n1", NewEstimator(nil)},
			{del, "n1", nil},
		}, 0},
		{"0,1,1", &Estimators{}, []action{
			{add, "n1", NewEstimator(nil)},
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
