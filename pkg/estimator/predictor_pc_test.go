package estimator

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestFakePCPredictor_Predict(t *testing.T) {
	type fields struct {
		PredictFunc func(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error)
	}
	type args struct {
		ctx             context.Context
		requestCPUMilli int
		status          NodeStatus
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantWatt float64
		wantErr  bool
	}{
		{"PredictFunc=nil", fields{nil}, args{context.Background(), 2000, NodeStatus{
			CPUUsages:    [][]float64{{30.0}},
			AmbientTemps: []float64{20.0},
		}}, 0.0, true},
		{"70", fields{PredictPCFnDummy}, args{context.Background(), 2000, NodeStatus{
			CPUUsages:    [][]float64{{30.0}},
			AmbientTemps: []float64{20.0},
		}}, 70.0, false},
		{"80", fields{PredictPCFnDummy}, args{context.Background(), 2000, NodeStatus{
			CPUUsages:    [][]float64{{35.0}},
			AmbientTemps: []float64{25.0},
		}}, 80.0, false},
		{"err", fields{PredictPCFnDummy}, args{context.Background(), 2000, NodeStatus{}}, 0.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FakePCPredictor{
				PredictFunc: tt.fields.PredictFunc,
			}
			gotWatt, err := p.Predict(tt.args.ctx, tt.args.requestCPUMilli, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("FakePCPredictor.Predict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWatt != tt.wantWatt {
				t.Errorf("FakePCPredictor.Predict() = %v, want %v", gotWatt, tt.wantWatt)
			}
		})
	}
}

func Test_newMLServerPCPredictorRequest(t *testing.T) {
	type args struct {
		cpuUsage           float64
		ambientTemp        float64
		staticPressureDiff float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"10/22/0.2",
			args{cpuUsage: 10.0, ambientTemp: 22.0, staticPressureDiff: 0.2},
			`{"inputs":[{"name":"predict-prob","shape":[1,3],"datatype":"FP32","data":[[10,22,0.2]]}]}`,
		},
		{"99.9/99.9/99.9",
			args{cpuUsage: 99.9, ambientTemp: 99.9, staticPressureDiff: 99.9},
			`{"inputs":[{"name":"predict-prob","shape":[1,3],"datatype":"FP32","data":[[99.9,99.9,99.9]]}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := newMLServerPCPredictorRequest(tt.args.cpuUsage, tt.args.ambientTemp, tt.args.staticPressureDiff)
			p, err := json.Marshal(v)
			if err != nil {
				t.Errorf("unable to encode err=%v", err)
			}
			if !reflect.DeepEqual(p, []byte(tt.want)) {
				t.Errorf("newMLServerPCPredictorRequest() = %s, want %s", p, tt.want)
			}
		})
	}
}

func Test_decodeJSON_mlServerPCPredictorResponse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  mlServerPCPredictorResponse
	}{
		{"94.76267448501928",
			`{"model_name":"model1","model_version":"v0.1.0","id":"c0d61974-f9bd-48ac-9af8-9c1e6d9a07e6","parameters":{"content_type":null,"headers":null},"outputs":[{"name":"predict","shape":[1,1],"datatype":"FP64","parameters":null,"data":[94.76267448501928]}]}`,
			mlServerPCPredictorResponse{Outputs: []struct {
				Data []float64 `json:"data"`
			}{{Data: []float64{94.76267448501928}}},
			},
		},
		{"15",
			`{"model_name":"model1","model_version":"v0.1.0","id":"c0d61974-f9bd-48ac-9af8-9c1e6d9a07e6","parameters":{"content_type":null,"headers":null},"outputs":[{"name":"predict","shape":[1,1],"datatype":"FP64","parameters":null,"data":[15]}]}`,
			mlServerPCPredictorResponse{Outputs: []struct {
				Data []float64 `json:"data"`
			}{{Data: []float64{15}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got mlServerPCPredictorResponse
			if err := json.Unmarshal([]byte(tt.input), &got); err != nil {
				t.Errorf("unable to decode err=%v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mlServerPCPredictorResponse got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMLServerPCPredictor_getURLV2Infer(t *testing.T) {
	type fields struct {
		Server  string
		Model   string
		Version string
	}
	tests := []struct {
		fields  fields
		want    string
		wantErr bool
	}{
		{fields{
			Server:  "http://localhost:8080",
			Model:   "model1",
			Version: "v0.1.0",
		}, "http://localhost:8080/v2/models/model1/versions/v0.1.0/infer", false},
		{fields{
			Server:  "https://example.com",
			Model:   "model2",
			Version: "v1.0.0",
		}, "https://example.com/v2/models/model2/versions/v1.0.0/infer", false},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			p := &MLServerPCPredictor{
				Server:  tt.fields.Server,
				Model:   tt.fields.Model,
				Version: tt.fields.Version,
			}
			got, err := p.getURLV2Infer()
			if (err != nil) != tt.wantErr {
				t.Errorf("MLServerPCPredictor.getURLV2Infer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MLServerPCPredictor.getURLV2Infer() = %v, want %v", got, tt.want)
			}
		})
	}
}
