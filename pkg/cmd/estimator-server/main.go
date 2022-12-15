package main

import (
	"net"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

// This is just an example showing how to start an estimator.Server
var lg zerolog.Logger = log.With().Str("component", "Estimator Server (main)").Logger()

func main() {

	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	apikeys := map[string]struct{}{}
	host := "localhost"
	port := estimator.ServerDefaultPort
	addr := net.JoinHostPort(host, port)

	authFn := estimator.AuthFnAPIKey(apikeys)
	if len(apikeys) == 0 {
		authFn = openapi3filter.NoopAuthenticationFunc
	}

	ns := &estimator.Nodes{}
	ns.Add("n1", estimator.NewNode("n1", nil, 30*time.Second, nil))
	es := &estimator.Estimators{}
	if ok := es.Add("default/default", estimator.NewEstimator(ns)); !ok {
		panic("es.Add not ok")
	}
	h, err := estimator.NewServer(es).HandlerWithAuthFn(authFn, middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, middleware.Heartbeat("/healthz"))
	if err != nil {
		panic(err)
	}
	panic(http.ListenAndServe(addr, h))
}
