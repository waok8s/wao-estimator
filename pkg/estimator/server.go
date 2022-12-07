package estimator

import (
	"context"
	"errors"
	"net/http"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator/api"
)

type Server struct{ e *Estimator }

var _ api.StrictServerInterface = (*Server)(nil)

func NewServer(estimator *Estimator) *Server { return &Server{e: estimator} }

func (s *Server) PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumption(ctx context.Context, request api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumptionRequestObject) (api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumptionResponseObject, error) {
	wattIncrease, err := s.e.EstimatePowerConsumption(ctx, request.Namespace, request.Estimator, request.Body.CpuMilli, request.Body.NumWorkload)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			return api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumption400Response{}, nil
		case errors.Is(err, ErrEstimatorNotFound):
			return api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumption404Response{}, nil
		case errors.Is(err, ErrEstimatorError):
			return api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumption500Response{}, nil
		default:
			return api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumption500Response{}, nil
		}
	}
	return api.PostNamespacesNamespaceEstimatorsEstimatorResourcesPowerconsumption200JSONResponse{
		CpuMilli:     request.Body.CpuMilli,
		NumWorkload:  request.Body.NumWorkload,
		WattIncrease: wattIncrease,
	}, nil
}

func (s *Server) Handler(middlewares ...func(http.Handler) http.Handler) (http.Handler, error) {
	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	swagger.Servers = nil // clear servers in the spec

	h := api.NewStrictHandler(s, nil)

	r := chi.NewRouter()
	r.Use(middlewares...)
	r.Use(middleware.OapiRequestValidator(swagger))

	return api.HandlerFromMux(h, r), nil
}
