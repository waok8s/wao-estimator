package estimator

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type ComputeLeastCostPatternsFunc func(clusterNum, podNum int, wattMatrix [][]float64) (minWatt float64, minWattPatterns [][]int, err error)

var ComputeLeastCostPatternsFn = findLeastCostPatternsExhaustive

type ComputeLeastCostsFunc func(clusterNum, podNum int, wattMatrix [][]float64) (minWatts []float64, err error)

var ComputeLeastCostsFn = findLeastCosts

func enumerateNdigitMbaseNumbers(n, m int, filter func([]int) bool, expectedReturnLength int) ([][]int, error) {
	t := time.Now()
	defer func() {
		lg.Debug().Msgf("enum elapsed=%dms", time.Since(t).Milliseconds())
	}()

	// constraints: 0 < m^n < 10^7
	// 10^6: 50ms
	// 10^7: 500ms
	// 10^8: 5000ms
	maxPowMN := 10_000_000
	total := int(math.Pow(float64(m), float64(n)))
	if total > maxPowMN {
		return nil, fmt.Errorf("combinatorial explosion m=%d n=%d m^n=%d max=%d", m, n, total, maxPowMN)
	}

	// it takes 500ms to grow an array from len=0 to len=10_000_000
	pp := make([][]int, 0, expectedReturnLength)

	// enumerate all numbers, digit=N base=M total=M^N
	for i := 0; i < total; i++ {
		v := make([]int, n)
		tmp := i
		for j := 0; j < n; j++ {
			v[j] = tmp % m
			tmp /= m
		}
		if filter(v) {
			pp = append(pp, v)
		}
	}

	lg.Debug().Msgf("n=%d m=%v total=%d filtered=%d", n, m, total, len(pp))
	return pp, nil
}

func is2DArray(vv [][]float64) (row int, col int, err error) {
	if len(vv) == 0 {
		return 0, 0, nil
	}
	row = len(vv)
	col = len(vv[0])
	for i, v := range vv {
		if len(v) != col {
			return 0, 0, fmt.Errorf("wrong 2D array idx=%d", i)
		}
	}
	return
}

func findLeastCostPatternsExhaustive(box, itemsPerBox int, costs [][]float64) (minCost float64, minCostPatterns [][]int, err error) {
	t := time.Now()
	defer func() {
		lg.Debug().Msgf("findLeastCostPatternsExhaustive total elapsed=%dms", time.Since(t).Milliseconds())
	}()

	row, col, err := is2DArray(costs)
	if err != nil {
		return 0.0, nil, nil
	}
	if row != box || col != itemsPerBox {
		return 0.0, nil, errors.New("len(costs)!=box || len(costs[*])!=itemsPerBox")
	}

	filterFn := func(v []int) bool {
		var sum int
		for _, i := range v {
			sum += i
		}
		return sum == itemsPerBox
	}

	// n=14 m=3 m^n=4782969 len(vv)=105
	// n=10 m=5 m^n=9765625 len(vv)=715
	// n=6 m=14 m^n=7529536 len(vv)=8568
	// n=5 m=25 m^n=9765625 len(vv)=20475
	// n=3 m=101 m^n=1030301 len(vv)=5151
	const expectedLen = 5000
	vv, err := enumerateNdigitMbaseNumbers(box, itemsPerBox+1, filterFn, expectedLen)
	if err != nil {
		return 0.0, nil, err
	}

	minCost = math.MaxFloat64
	for i, v := range vv {
		var sum float64
		for box, itemNum := range v {
			if itemNum == 0 {
				continue
			}
			sum += costs[box][itemNum-1]
		}
		lg.Trace().Msgf("i=%d ptn=%v cost=%v", i, v, sum)
		if sum < minCost {
			minCost = sum
			minCostPatterns = nil
		}
		if sum == minCost {
			minCostPatterns = append(minCostPatterns, v)
		}
	}

	lg.Debug().Msgf("box=%d itemsPerBox=%d minCost=%f len(minCostPatterns)=%d", box, itemsPerBox, minCost, len(minCostPatterns))
	lg.Debug().Msgf("box=%d itemsPerBox=%d costs=%v minCostPatterns=%v", box, itemsPerBox, costs, minCostPatterns)

	return minCost, minCostPatterns, nil
}

func findLeastCosts(box, itemsPerBox int, costs [][]float64) (minCosts []float64, err error) {
	t := time.Now()
	defer func() {
		lg.Info().Msgf("findLeastCosts elapsed=%dms minCosts=%v", time.Since(t).Milliseconds(), minCosts)
	}()

	row, col, err := is2DArray(costs)
	if err != nil {
		return []float64{}, nil
	}
	if row != box || col != itemsPerBox {
		return nil, errors.New("len(costs)!=box || len(costs[*])!=itemsPerBox")
	}

	minCosts = []float64{}
	for m := 1; m <= itemsPerBox; m++ {
		var curCosts [][]float64
		for i := 0; i < box; i++ {
			curCosts = append(curCosts, costs[i][:(m)])
		}
		minCost, _, err := ComputeLeastCostPatternsFn(box, m, curCosts)
		if err != nil {
			fmt.Print(err)
		}
		minCosts = append(minCosts, minCost)
	}
	return minCosts, nil
}

// toDiff returns a new matrix, make each row element
// the difference from the first element in that row.
// i.e. [[a0 a0+a1 a0+a2] [b0 b0+b1 b0+b2]] -> [[a1 a2] [b1 b2]]
//
// input:
//
//	[100 120 140]
//	[200 300 400]
//	[500 500 500]
//
// output:
//
//	[  20  40]
//	[ 100 200]
//	[   0   0]
func toDiff(vv [][]float64) ([][]float64, error) {
	row, col, err := is2DArray(vv)
	if err != nil {
		return nil, err
	}
	if row == 0 {
		return [][]float64{}, nil
	}
	if col == 0 {
		return nil, errors.New("col must be >0")
	}
	xx := make([][]float64, row)
	for i := range xx {
		xx[i] = make([]float64, col-1)
		v0 := vv[i][0]
		for j := 0; j < col-1; j++ {
			xx[i][j] = vv[i][j+1] - v0
		}
	}
	return xx, nil
}
