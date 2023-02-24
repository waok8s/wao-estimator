package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

var lg zerolog.Logger = log.With().Str("component", "MLServerPCP (main)").Logger()

func main() {

	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

	server := "http://localhost:8080"
	model := "model1"
	version := "v0.1.0"
	cpuUsage := 10.0
	ambientTemp := 22.0
	staticPressureDiff := 0.2

	pcp := estimator.MLServerPCPredictor{Server: server, Model: model, Version: version}
	watt, err := pcp.POSTPredictRequest(context.Background(), cpuUsage, ambientTemp, staticPressureDiff)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println(watt)
}
