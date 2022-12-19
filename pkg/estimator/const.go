package estimator

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var lg zerolog.Logger = log.With().Str("component", "estimator").Logger()

var (
	ErrUnexpected = errors.New("ErrUnexpected")

	ErrClientUnauthorized      = errors.New("ErrClientUnauthorized")
	ErrServerEstimatorNotFound = errors.New("ErrServerEstimatorNotFound")

	ErrEstimator                 = errors.New("ErrEstimator")
	ErrEstimatorNoNodesAvailable = errors.New("ErrEstimatorNoNodesAvailable")
	ErrEstimatorInvalidRequest   = errors.New("ErrEstimatorInvalidRequest")

	ErrNodeMonitor         = errors.New("ErrNodeMonitor")
	ErrNodeMonitorNotFound = errors.New("ErrNodeMonitorNotFound")

	ErrNodeStatus = errors.New("ErrNodeStatus")

	ErrPCPredictor         = errors.New("ErrPCPredictor")
	ErrPCPredictorNotFound = errors.New("ErrPCPredictorNotFound")
)

func GetErrorFromCode(apiErr Error) error {
	err, ok := getErrorFromCode[apiErr.Code]
	if !ok {
		return ErrUnexpected
	}
	return err
}

var getErrorFromCode map[string]error = map[string]error{
	ErrUnexpected.Error(): ErrUnexpected,

	ErrClientUnauthorized.Error():      ErrClientUnauthorized,
	ErrServerEstimatorNotFound.Error(): ErrServerEstimatorNotFound,

	ErrEstimator.Error():                 ErrEstimator,
	ErrEstimatorNoNodesAvailable.Error(): ErrEstimatorNoNodesAvailable,
	ErrEstimatorInvalidRequest.Error():   ErrEstimatorInvalidRequest,

	ErrNodeMonitor.Error():         ErrNodeMonitor,
	ErrNodeMonitorNotFound.Error(): ErrNodeMonitorNotFound,

	ErrNodeStatus.Error(): ErrNodeStatus,

	ErrNodeStatus.Error(): ErrNodeStatus,

	ErrPCPredictor.Error():         ErrPCPredictor,
	ErrPCPredictorNotFound.Error(): ErrPCPredictorNotFound,
}
