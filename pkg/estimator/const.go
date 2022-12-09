package estimator

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lg zerolog.Logger = log.With().Str("component", "estimator").Logger()

var (
	ErrInvalidRequest      = errors.New("ErrInvalidRequest")
	ErrEstimator           = errors.New("ErrEstimator")
	ErrEstimatorNotFound   = errors.New("ErrEstimatorNotFound")
	ErrNodeMonitor         = errors.New("ErrNodeMonitor")
	ErrNodeMonitorNotFound = errors.New("ErrNodeMonitorNotFound")
	ErrNodeStatus          = errors.New("ErrNodeStatus")
	ErrNodeStatusNotFound  = errors.New("ErrNodeStatusNotFound")
	ErrPCPredictor         = errors.New("ErrPCPredictor")
	ErrPCPredictorNotFound = errors.New("ErrPCPredictorNotFound")
)
