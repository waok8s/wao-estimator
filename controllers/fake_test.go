package controllers

import (
	"reflect"
	"testing"
)

func Test_labelValueString(t *testing.T) {
	type args struct {
		m     map[string]string
		label string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"not_found", args{
			m:     map[string]string{},
			label: "hoge",
		}, "", true},
		{"empty", args{
			m:     map[string]string{"hoge": ""},
			label: "hoge",
		}, "", false},
		{"fuga", args{
			m:     map[string]string{"hoge": "fuga"},
			label: "hoge",
		}, "fuga", false},
		{"num", args{
			m:     map[string]string{"hoge": "1.5"},
			label: "hoge",
		}, "1.5", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelValueString(tt.args.m, tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("labelValueString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("labelValueString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_labelValueInt(t *testing.T) {
	type args struct {
		m     map[string]string
		label string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"not_found", args{
			m:     map[string]string{},
			label: "hoge",
		}, 0, true},
		{"1", args{
			m:     map[string]string{"hoge": "1"},
			label: "hoge",
		}, 1, false},
		{"-1", args{
			m:     map[string]string{"fuga": "-1"},
			label: "fuga",
		}, -1, false},
		{"a", args{
			m:     map[string]string{"hoge": "a"},
			label: "hoge",
		}, 0, true},
		{"empty_label", args{
			m:     map[string]string{"": "1"},
			label: "",
		}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelValueInt(tt.args.m, tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("labelValueInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("labelValueInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_labelValueJSON(t *testing.T) {
	type args struct {
		m     map[string]string
		label string
	}
	tests := []struct {
		name     string
		args     args
		wantType any
		wantVal  any
		wantErr  bool
	}{
		{"not_found", args{
			m:     map[string]string{},
			label: "hoge",
		}, []float64{}, []float64{}, true},
		{"VF64/empty", args{
			m: map[string]string{
				"a": "bd",
			},
			label: "a",
		}, []float64{}, []float64{}, false},
		{"VF64/normal", args{
			m: map[string]string{
				"a": "b1.0p1.5d",
			},
			label: "a",
		}, []float64{}, []float64{1.0, 1.5}, false},
		{"VF64/wrong1", args{
			m: map[string]string{
				"a": "aaaaaaa",
			},
			label: "a",
		}, []float64{}, nil, true},
		{"VF64/wrong2", args{
			m: map[string]string{
				"a": "b__b1.0p1.5d__p__b2.0p2.5d__d",
			},
			label: "a",
		}, []float64{}, nil, true},
		{"VVF64/empty", args{
			m: map[string]string{
				"a": "bd",
			},
			label: "a",
		}, [][]float64{}, [][]float64{}, false},
		{"VVF64/normal", args{
			m: map[string]string{
				"a": "b__b1.0p1.5d__p__b2.0p2.5d__d",
			},
			label: "a",
		}, [][]float64{}, [][]float64{{1.0, 1.5}, {2.0, 2.5}}, false},
		{"VVF64/wrong", args{
			m: map[string]string{
				"a": "b__b1.0p1.5d__p__b2.0p2.5d__d",
			},
			label: "a",
		}, []float64{}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			var err error

			switch tt.wantType.(type) {
			case []float64:
				got, err = labelValueJSON[[]float64](tt.args.m, tt.args.label)
			case [][]float64:
				got, err = labelValueJSON[[][]float64](tt.args.m, tt.args.label)
			default:
				t.Errorf("got type=%T", tt.wantType)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("labelValueJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// NOTE: it's hard to compare the returned value in error cases (e.g. any([]float64(nil),) vs. any() )
			if err == nil {
				if !reflect.DeepEqual(got, tt.wantVal) {
					t.Errorf("labelValueJSON() = %v (%T), want %v (%T)", got, got, tt.wantVal, tt.wantVal)
				}
			}
		})
	}
}
