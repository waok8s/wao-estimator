package estimator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"moul.io/http2curl/v2"
)

type DifferentialPressureNodeMonitor struct {
	Server string
	Sensor string
}

var _ NodeMonitor = (*DifferentialPressureNodeMonitor)(nil)

// NewDifferentialPressureNodeMonitorFromURL parses the given endpoint URL.
//
// Format: {Server}/api/sensor/{Sensor}
// Example: http://hogehoge:5000/api/sensor/101037B -> &{Server: "http://hogehoge:5000", Sensor: "101037B"}
func NewDifferentialPressureNodeMonitorFromURL(endpoint string) (*DifferentialPressureNodeMonitor, error) {
	parsedURL, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse DifferentialPressureAPI endpoint URL %w: %v", ErrNodeMonitor, err)
	}
	ss := strings.Split(parsedURL.Path, "/") // "/api/sensor/101037B" -> ["", "api", "sensor", "101037B"]
	if len(ss) < 4 {
		return nil, fmt.Errorf("could not parse DifferentialPressureAPI endpoint URL %w: url path must be in {Server}/api/sensor/{Sensor} format", ErrNodeMonitor)
	}
	return &DifferentialPressureNodeMonitor{
		Server: parsedURL.Scheme + "://" + parsedURL.Host,
		Sensor: ss[3],
	}, nil
}

func (m *DifferentialPressureNodeMonitor) FetchStatus(ctx context.Context, base *NodeStatus) error {
	if base == nil {
		base = NewNodeStatus()
	}
	v, err := m.GETValueRequest(ctx)
	if err != nil {
		return fmt.Errorf("could not read sensor value (%w): %v", ErrNodeMonitor, err)
	}
	NodeStatusSetStaticPressureDiff(base, v.Pressure)
	return nil
}

// getURLAPISensor returns the API endpoint.
// e.g. "http://localhost:5000/api/sensor/101037B"
func (m *DifferentialPressureNodeMonitor) getURLAPISensor() (string, error) {
	return url.JoinPath(m.Server, "api", "sensor", m.Sensor)
}

func (m *DifferentialPressureNodeMonitor) GETValueRequest(ctx context.Context) (sensorValue, error) {
	var v sensorValue

	url, err := m.getURLAPISensor()
	if err != nil {
		return v, fmt.Errorf("unable to get endpoint URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return v, fmt.Errorf("unable to create HTTP request: %w", err)
	}

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		lg.Err(err).Msgf("DifferentialPressureNodeMonitor.FetchStatus could not parse http.Request to curl command")
	} else {
		lg.Trace().Msgf("DifferentialPressureNodeMonitor.FetchStatus request=%v", curl.String())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return v, fmt.Errorf("unable to send HTTP request: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var apiResp differentialPressureAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return v, fmt.Errorf("could not decode resp: %w", err)
		}
		if len(apiResp.Sensors) == 0 {
			return v, fmt.Errorf("invalid response apiResp=%+v", apiResp)
		}
		return apiResp.Sensors[0], nil
	default:
		return v, fmt.Errorf("HTTP status=%v request=%v", resp.Status, curl.String())
	}
}

// differentialPressureAPIResponse holds a response.
//
// e.g.
//
//	{
//	  "code": 200,
//	  "sensor": [
//	    {
//	      "pressure": 0.01,
//	      "sensorid": "101037B",
//	      "temperature": 26.02
//	    }
//	  ]
//	}
type differentialPressureAPIResponse struct {
	StatusCode int           `json:"code"`
	Sensors    []sensorValue `json:"sensor"`
}

type sensorValue struct {
	SensorID    string  `json:"sensorid"`
	Pressure    float64 `json:"pressure"`
	Temperature float64 `json:"temperature"`
}
