package estimator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	http2curl "moul.io/http2curl/v2"
)

type PowerConsumptionPredictor interface {
	Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error)
}

type FakePCPredictor struct {
	PredictFunc func(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error)
}

var _ PowerConsumptionPredictor = (*FakePCPredictor)(nil)

func (p *FakePCPredictor) Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error) {
	if p.PredictFunc == nil {
		return 0.0, fmt.Errorf("PredictFunc not set (%w)", ErrPCPredictor)
	}
	return p.PredictFunc(ctx, requestCPUMilli, status)
}

// PredictPCFnDummy returns ( mCPU/100 + avg(CPUUsage) + avg(AmbientTemp) )
func PredictPCFnDummy(_ context.Context, mcpu int, status NodeStatus) (float64, error) {
	aat, err := status.AverageAmbientTemp()
	if err != nil {
		return 0.0, err
	}
	acu, err := status.AverageCPUUsage()
	if err != nil {
		return 0.0, err
	}
	w := float64(mcpu)/100 + aat + acu
	return w, nil
}

type MLServerPCPredictor struct {
	// Server specifies server address
	// e.g. "http://localhost:8080"
	Server string
	// Model specifies model name
	// e.g. "model1"
	Model string
	// Version specifies version for the model
	// e.g. "v0.1.0"
	Version string
}

var _ PowerConsumptionPredictor = (*MLServerPCPredictor)(nil)

func (p *MLServerPCPredictor) Predict(ctx context.Context, requestCPUMilli int, status NodeStatus) (watt float64, err error) {

	// TODO: convert requestCPUMilli to CPU usage
	cpuUsage := 10.0
	// TODO: fetch the ambient sensor number used in the model and get the value
	ambientTemp := 22.0
	// TODO: fetch staticPressureDiff
	staticPressureDiff := 0.2

	return p.POSTPredictRequest(ctx, cpuUsage, ambientTemp, staticPressureDiff)
}

// getURLV2Infer returns the API endpoint.
// e.g. "http://localhost:8080/v2/models/model1/versions/v0.1.0/infer"
func (p *MLServerPCPredictor) getURLV2Infer() (string, error) {
	return url.JoinPath(p.Server, "v2", "models", p.Model, "versions", p.Version, "infer")
}

func (p *MLServerPCPredictor) POSTPredictRequest(ctx context.Context, cpuUsage, ambientTemp, staticPressureDiff float64) (float64, error) {

	url, err := p.getURLV2Infer()
	if err != nil {
		return 0.0, fmt.Errorf("unable to get endpoint URL: %w", err)
	}

	body, err := json.Marshal(newMLServerPCPredictorRequest(cpuUsage, ambientTemp, staticPressureDiff))
	if err != nil {
		return 0.0, fmt.Errorf("unable to marshal the request body=%+v err=%w", body, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return 0.0, fmt.Errorf("unable to create HTTP request: %w", err)
	}

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		lg.Err(err).Msgf("MLServerPCPredictor.Predict could not parse http.Request to curl command")
	} else {
		lg.Trace().Msgf("MLServerPCPredictor.Predict request=%v", curl.String())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0.0, fmt.Errorf("unable to send HTTP request: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var predictResp mlServerPCPredictorResponse
		if err := json.NewDecoder(resp.Body).Decode(&predictResp); err != nil {
			return 0.0, fmt.Errorf("could not decode resp: %w", err)
		}
		if len(predictResp.Outputs) == 0 || len(predictResp.Outputs[0].Data) == 0 {
			return 0.0, fmt.Errorf("invalid response predictResp=%+v", predictResp)
		}
		return predictResp.Outputs[0].Data[0], nil
	default:
		return 0.0, fmt.Errorf("HTTP status=%v request=%v", resp.Status, curl.String())
	}

}

// mlServerPCPredictorRequest holds a request.
//
// e.g.
//
//	{
//	  "inputs": [
//	    {
//	      "name": "predict-prob",
//	      "shape": [ 1, 3 ],
//	      "datatype": "FP32",
//	      "data": [ [10, 22, 0.2 ] ]
//	    }
//	  ]
//	}
type mlServerPCPredictorRequest struct {
	Inputs []struct {
		Name     string      `json:"name"`
		Shape    []int       `json:"shape"`
		Datatype string      `json:"datatype"`
		Data     [][]float32 `json:"data"`
	} `json:"inputs"`
}

func newMLServerPCPredictorRequest(cpuUsage, ambientTemp, staticPressureDiff float64) *mlServerPCPredictorRequest {
	const (
		name     = "predict-prob"
		datatype = "FP32"
		shapeX   = 1
		shapeY   = 3
	)
	return &mlServerPCPredictorRequest{
		Inputs: []struct {
			Name     string      `json:"name"`
			Shape    []int       `json:"shape"`
			Datatype string      `json:"datatype"`
			Data     [][]float32 `json:"data"`
		}{
			{
				Name:     name,
				Shape:    []int{shapeX, shapeY},
				Datatype: datatype,
				Data:     [][]float32{{float32(cpuUsage), float32(ambientTemp), float32(staticPressureDiff)}},
			},
		},
	}
}

// mlServerPCPredictorResponse holds a response.
// Ignore values except outputs[*].data[]
//
// e.g.
//
//	{
//	  "model_name": "model1",
//	  "model_version": "v0.1.0",
//	  "id": "0dc429d2-bd02-404b-b624-a0fa628e451e",
//	  "parameters": {
//	    "content_type": null,
//	    "headers": null
//	  },
//	  "outputs": [
//	    {
//	      "name": "predict",
//	      "shape": [ 1, 1 ],
//	      "datatype": "FP64",
//	      "parameters": null,
//	      "data": [ 94.76267448501928 ]
//	    }
//	  ]
//	}
type mlServerPCPredictorResponse struct {
	Outputs []struct {
		Data []float64 `json:"data"`
	} `json:"outputs"`
}
