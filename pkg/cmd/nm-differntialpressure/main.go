package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

var lg zerolog.Logger = log.With().Str("component", "DifferntialPressureAPINM (main)").Logger()

func main() {

	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

	server := "http://localhost:5000"
	sensor := "101037B"

	nm := &estimator.DifferentialPressureNodeMonitor{Server: server, Sensor: sensor}
	resp, err := nm.GETValueRequest(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", resp)
}
