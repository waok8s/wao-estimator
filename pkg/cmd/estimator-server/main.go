package main

import (
	"fmt"
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

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

	apikeys := map[string]struct{}{}
	host := "localhost"
	port := fmt.Sprint(estimator.ServerDefaultPort)
	addr := net.JoinHostPort(host, port)

	authFn := estimator.AuthFnAPIKey(apikeys)
	if len(apikeys) == 0 {
		authFn = openapi3filter.NoopAuthenticationFunc
	}

	ns := &estimator.Nodes{}
	ns.Add("n1", estimator.NewNode("n1", nil, 30*time.Second, nil))
	ns.Add("n2", estimator.NewNode("n2", nil, 30*time.Second, nil))
	ns.Add("n3", estimator.NewNode("n3", nil, 30*time.Second, nil))
	// ns.Add("n4", estimator.NewNode("n4", nil, 30*time.Second, nil))
	// ns.Add("n5", estimator.NewNode("n5", nil, 30*time.Second, nil))
	// ns.Add("n6", estimator.NewNode("n6", nil, 30*time.Second, nil))
	// ns.Add("n7", estimator.NewNode("n7", nil, 30*time.Second, nil))
	// ns.Add("n8", estimator.NewNode("n8", nil, 30*time.Second, nil))
	es := &estimator.Estimators{}
	if ok := es.Add("default/default", &estimator.Estimator{Nodes: ns}); !ok {
		panic("es.Add not ok")
	}
	sv := &estimator.Server{Estimators: es}
	h, err := sv.HandlerWithAuthFn(authFn, middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, middleware.Heartbeat("/healthz"))
	if err != nil {
		panic(err)
	}
	panic(http.ListenAndServe(addr, h))
}
