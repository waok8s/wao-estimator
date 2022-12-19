package estimator

import (
	"errors"
	"testing"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator/api"
)

func TestGetErrorFromCode(t *testing.T) {
	type args struct {
		apiErr Error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"", args{api.Error{Code: ErrEstimator.Error(), Message: "hoge"}}, ErrEstimator},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetErrorFromCode(tt.args.apiErr); !errors.Is(err, tt.wantErr) {
				t.Errorf("GetErrorFromCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
