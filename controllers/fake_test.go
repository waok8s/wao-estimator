package controllers

import (
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

func Test_labelValueFloat(t *testing.T) {
	type args struct {
		m     map[string]string
		label string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"not_found", args{
			m:     map[string]string{},
			label: "hoge",
		}, 0, true},
		{"1.0", args{
			m:     map[string]string{"hoge": "1.0"},
			label: "hoge",
		}, 1.0, false},
		{"-1", args{
			m:     map[string]string{"fuga": "-1"},
			label: "fuga",
		}, -1.0, false},
		{"a", args{
			m:     map[string]string{"hoge": "a"},
			label: "hoge",
		}, 0, true},
		{"empty_label", args{
			m:     map[string]string{"": "1.0"},
			label: "",
		}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelValueFloat(tt.args.m, tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("labelValueFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("labelValueFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
