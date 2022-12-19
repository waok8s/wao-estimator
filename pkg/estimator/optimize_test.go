package estimator

import (
	"reflect"
	"testing"
)

var (
	testCostsReq1 = [][]float64{
		{11.0, 14.0, 14.0, 16.0, 30.0, 31.0},
		{11.0, 14.0, 15.0, 15.0, 22.0, 29.0},
		{10.0, 17.0, 19.0, 25.0, 26.0, 30.0},
		{72.0, 80.0, 92.0, 99.0, 99.0, 99.0},
		{11.0, 17.0, 19.0, 25.0, 27.0, 29.0},
		{14.0, 15.0, 15.0, 26.0, 31.0, 32.0},
		{16.0, 16.0, 16.0, 16.0, 24.0, 32.0},
		{29.0, 29.0, 29.0, 29.0, 29.0, 29.0},
	}
	testCostsAns1 = []float64{10.0, 14.0, 14.0, 15.0, 22.0, 29.0}
)

func Test_findLeastCosts(t *testing.T) {
	type args struct {
		box         int
		itemsPerbox int
		costs       [][]float64
	}
	tests := []struct {
		name         string
		args         args
		wantMinCosts []float64
		wantErr      bool
	}{
		{"8,6", args{
			box:         8,
			itemsPerbox: 6,
			costs:       testCostsReq1,
		}, testCostsAns1, false},
		{"1,6", args{
			box:         1,
			itemsPerbox: 6,
			costs: [][]float64{
				{11.0, 14.0, 14.0, 16.0, 30.0, 31.0},
			},
		}, []float64{11.0, 14.0, 14.0, 16.0, 30.0, 31.0}, false},
		{"8,1", args{
			box:         8,
			itemsPerbox: 1,
			costs: [][]float64{
				{11.0},
				{11.0},
				{10.0},
				{72.0},
				{11.0},
				{14.0},
				{16.0},
				{29.0},
			},
		}, []float64{10.0}, false},
		{"1,1", args{
			box:         1,
			itemsPerbox: 1,
			costs: [][]float64{
				{11.0},
			},
		}, []float64{11.0}, false},
		{"0,0", args{
			box:         0,
			itemsPerbox: 0,
			costs:       [][]float64{},
		}, []float64{}, false},
		{"8,0", args{
			box:         8,
			itemsPerbox: 0,
			costs:       [][]float64{{}, {}, {}, {}, {}, {}, {}, {}},
		}, []float64{}, false},
		{"0,6 (array[0][6] is invalid so return error)", args{
			box:         0,
			itemsPerbox: 6,
			costs:       [][]float64{},
		}, nil, true},
		{"wrongX", args{
			box:         -1,
			itemsPerbox: 6,
			costs:       testCostsReq1,
		}, nil, true},
		{"wrongY", args{
			box:         8,
			itemsPerbox: -1,
			costs:       testCostsReq1,
		}, nil, true},
		{"wrongXY", args{
			box:         -1,
			itemsPerbox: -1,
			costs:       [][]float64{},
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMinCosts, err := findLeastCosts(tt.args.box, tt.args.itemsPerbox, tt.args.costs)
			if (err != nil) != tt.wantErr {
				t.Errorf("calcMinCosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMinCosts, tt.wantMinCosts) {
				t.Errorf("calcMinCosts() = %v, want %v", gotMinCosts, tt.wantMinCosts)
			}
		})
	}
}

var (
	testCostsReq1BeforeToDiff = [][]float64{
		{100.0, 111.0, 114.0, 114.0, 116.0, 130.0, 131.0},
		{200.0, 211.0, 214.0, 215.0, 215.0, 222.0, 229.0},
		{300.0, 310.0, 317.0, 319.0, 325.0, 326.0, 330.0},
		{400.0, 472.0, 480.0, 492.0, 499.0, 499.0, 499.0},
		{500.0, 511.0, 517.0, 519.0, 525.0, 527.0, 529.0},
		{600.0, 614.0, 615.0, 615.0, 626.0, 631.0, 632.0},
		{700.0, 716.0, 716.0, 716.0, 716.0, 724.0, 732.0},
		{800.0, 829.0, 829.0, 829.0, 829.0, 829.0, 829.0},
	}
)

func Test_toDiff(t *testing.T) {
	type args struct {
		vv [][]float64
	}
	tests := []struct {
		name    string
		args    args
		want    [][]float64
		wantErr bool
	}{
		{"testCostsReq1", args{testCostsReq1BeforeToDiff}, testCostsReq1, false},
		{"3/3", args{[][]float64{
			{100, 120, 140},
			{200, 300, 400},
			{500, 500, 500},
		}}, [][]float64{
			{20, 40},
			{100, 200},
			{0, 0},
		}, false},
		{"1/3", args{[][]float64{
			{100, 120, 140},
		}}, [][]float64{
			{20, 40},
		}, false},
		{"1/1", args{[][]float64{{123}}}, [][]float64{{}}, false},
		{"0/0", args{[][]float64{}}, [][]float64{}, false},
		{"1/0 err", args{[][]float64{{}}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toDiff(tt.args.vv)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}
