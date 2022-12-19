package estimator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator/api"
)

const ServerDefaultPort = 5656

type Server struct {
	Estimators *Estimators
	init       sync.Once
}

func (s *Server) initOnce() {
	s.init.Do(func() {
		if s.Estimators == nil {
			s.Estimators = &Estimators{}
		}
	})
}

var _ api.StrictServerInterface = (*Server)(nil)

func (s *Server) PostNamespacesNsEstimatorsNameValuesPowerconsumption(ctx context.Context, request api.PostNamespacesNsEstimatorsNameValuesPowerconsumptionRequestObject) (api.PostNamespacesNsEstimatorsNameValuesPowerconsumptionResponseObject, error) {
	s.initOnce()

	e, ok := s.Estimators.Get(client.ObjectKey{Namespace: request.Ns, Name: request.Name}.String())
	if !ok {
		return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption404JSONResponse{
			Code:    ErrServerEstimatorNotFound.Error(),
			Message: fmt.Sprintf("estimator %v/%v not found", request.Ns, request.Name),
		}, nil
	}
	wattIncrease, err := e.EstimatePowerConsumption(ctx, request.Body.CpuMilli, request.Body.NumWorkloads)
	if err != nil {
		switch {
		// 400
		case errors.Is(err, ErrEstimatorInvalidRequest):
			return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption400JSONResponse{
				Code:    ErrEstimatorInvalidRequest.Error(),
				Message: err.Error(),
			}, nil
		// 500
		default:
			unwrappedErr := err
			if _, ok := getErrorFromCode[err.Error()]; !ok {
				// wrapped or unexpected
				unwrappedErr = errors.Unwrap(err)
				if unwrappedErr == nil {
					unwrappedErr = ErrUnexpected
				}
			}
			return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption500JSONResponse{
				Code:    unwrappedErr.Error(),
				Message: err.Error(),
			}, nil
		}
	}
	return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption200JSONResponse{
		CpuMilli:      request.Body.CpuMilli,
		NumWorkloads:  request.Body.NumWorkloads,
		WattIncreases: &wattIncrease,
	}, nil
}

type AuthenticationFunc = openapi3filter.AuthenticationFunc

const AuthFnAPIKeyRequestHeader = "X-API-KEY"

func AuthFnAPIKey(apiKeys map[string]struct{}) AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		rh := http.CanonicalHeaderKey(AuthFnAPIKeyRequestHeader) // X-Api-Key
		h := input.RequestValidationInput.Request.Header[rh]
		if len(h) == 0 {
			return fmt.Errorf("request header %s not found", rh)
		}
		for _, k := range h {
			if _, ok := apiKeys[k]; ok {
				return nil
			}
		}
		return errors.New("authentication failed")
	}
}

func (s *Server) Handler(middlewares ...func(http.Handler) http.Handler) (http.Handler, error) {
	return s.HandlerWithAuthFn(openapi3filter.NoopAuthenticationFunc, middlewares...)
}

func (s *Server) HandlerWithAuthFn(authFn AuthenticationFunc, middlewares ...func(http.Handler) http.Handler) (http.Handler, error) {
	s.initOnce()

	spec, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	spec.Servers = nil // clear servers in the spec

	h := api.NewStrictHandler(s, nil)

	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Use(middleware.OapiRequestValidatorWithOptions(
		spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: authFn,
			},
		}))
	return api.HandlerFromMux(h, r), nil
}
