package estimator

import (
	"testing"
)

func TestNodeStatusGetCPUUsage(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		{"10.0", 10.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewNodeStatus()
			NodeStatusSetCPUUsage(s, tt.want)
			got, err := NodeStatusGetCPUUsage(s)
			if err != nil {
				t.Errorf("NodeStatusGetCPUUsage() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("NodeStatusGetCPUUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStatusGetAmbientTemp(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		{"10.0", 10.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewNodeStatus()
			NodeStatusSetAmbientTemp(s, tt.want)
			got, err := NodeStatusGetAmbientTemp(s)
			if err != nil {
				t.Errorf("NodeStatusGetAmbientTemp() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("NodeStatusGetAmbientTemp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeStatusGetStaticPressureDiff(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		{"10.0", 10.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewNodeStatus()
			NodeStatusSetStaticPressureDiff(s, tt.want)
			got, err := NodeStatusGetStaticPressureDiff(s)
			if err != nil {
				t.Errorf("NodeStatusGetStaticPressureDiff() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("NodeStatusGetStaticPressureDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}
