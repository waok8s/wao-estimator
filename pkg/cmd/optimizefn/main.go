package main

import (
	"fmt"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lg zerolog.Logger = log.With().Str("component", "CalcMinCosts (main)").Logger()

func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

	x := [][]float64{
		{11.0, 14.0, 14.0, 16.0, 30.0, 31.0},
		{11.0, 14.0, 15.0, 15.0, 22.0, 29.0},
		{10.0, 17.0, 19.0, 25.0, 26.0, 30.0},
		{72.0, 80.0, 92.0, 99.0, 99.0, 99.0},
		{11.0, 17.0, 19.0, 25.0, 27.0, 29.0},
		{14.0, 15.0, 15.0, 26.0, 31.0, 32.0},
		{16.0, 16.0, 16.0, 16.0, 24.0, 32.0},
		{29.0, 29.0, 29.0, 29.0, 29.0, 29.0},
	}
	clusterNum := len(x)
	podNum := 0
	if len(x) != 0 {
		podNum = len(x[0])
	}

	a, err := estimator.ComputeLeastCostsFn(clusterNum, podNum, x)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)
}
