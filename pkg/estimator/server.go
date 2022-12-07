package estimator

import (
	"context"
	"errors"
	"net/http"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator/api"
)

type Server struct{ e *Estimator }

var _ api.StrictServerInterface = (*Server)(nil)

func NewServer(estimator *Estimator) *Server { return &Server{e: estimator} }

func (s *Server) PostNamespacesNsEstimatorsNameValuesPowerconsumption(ctx context.Context, request api.PostNamespacesNsEstimatorsNameValuesPowerconsumptionRequestObject) (api.PostNamespacesNsEstimatorsNameValuesPowerconsumptionResponseObject, error) {
	wattIncrease, err := s.e.EstimatePowerConsumption(ctx, request.Ns, request.Name, request.Body.CpuMilli, request.Body.NumWorkloads)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption400Response{}, nil
		case errors.Is(err, ErrEstimatorNotFound):
			return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption404Response{}, nil
		case errors.Is(err, ErrEstimatorError):
			return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption500Response{}, nil
		default:
			return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption500Response{}, nil
		}
	}
	return api.PostNamespacesNsEstimatorsNameValuesPowerconsumption200JSONResponse{
		CpuMilli:     request.Body.CpuMilli,
		NumWorkloads: request.Body.NumWorkloads,
		WattIncrease: wattIncrease,
	}, nil
}

func (s *Server) Handler(middlewares ...func(http.Handler) http.Handler) (http.Handler, error) {
	return s.HandlerWithAuthFn(openapi3filter.NoopAuthenticationFunc, middlewares...)
}

func (s *Server) HandlerWithAuthFn(authFn openapi3filter.AuthenticationFunc, middlewares ...func(http.Handler) http.Handler) (http.Handler, error) {
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
