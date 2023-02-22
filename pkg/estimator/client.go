package estimator

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"

	http2curl "moul.io/http2curl/v2"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator/api"
)

func ClientOptionAddRequestHeader(k, v string) ClientOption {
	return func(c *api.Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Add(k, v)
			return nil
		})
		return nil
	}
}

func ClientOptionGetRequestAsCurl(w io.Writer) ClientOption {
	return func(c *api.Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			cmd, _ := http2curl.GetCurlCommand(req)
			w.Write([]byte(cmd.String()))
			return nil
		})
		return nil
	}
}

type Client struct {
	c       api.ClientWithResponsesInterface
	reqNS   string
	reqName string
}

func NewClient(server string, estimatorNamespace, estimatorName string, opts ...ClientOption) (*Client, error) {
	c, err := api.NewClientWithResponses(server, opts...)
	if err != nil {
		return nil, err
	}
	ec := Client{c: c, reqNS: estimatorNamespace, reqName: estimatorName}
	return &ec, nil
}

func (c *Client) EstimatePowerConsumption(ctx context.Context, cpuMilli, numWorkloads int) (pc *PowerConsumption, apiErr *Error, requestErr error) {
	body := api.PostNamespacesNsEstimatorsNameValuesPowerconsumptionJSONRequestBody{
		CpuMilli:      cpuMilli,
		NumWorkloads:  numWorkloads,
		WattIncreases: nil,
	}
	resp, err := c.c.PostNamespacesNsEstimatorsNameValuesPowerconsumptionWithResponse(ctx, c.reqNS, c.reqName, body)
	if err != nil {
		return nil, nil, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		// HACK: restore math.MaxFloat64 to math.Inf(1) (see also: server.go)
		for i := range *resp.JSON200.WattIncreases {
			if (*resp.JSON200.WattIncreases)[i] == math.MaxFloat64 {
				(*resp.JSON200.WattIncreases)[i] = math.Inf(1)
			}
		}
		return resp.JSON200, nil, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400, nil
	case http.StatusUnauthorized:
		return nil, &api.Error{Code: ErrClientUnauthorized.Error(), Message: "client unauthorized"}, nil
	case http.StatusNotFound:
		return nil, resp.JSON404, nil
	case http.StatusInternalServerError:
		return nil, resp.JSON500, nil
	default:
		return nil, nil, fmt.Errorf("%v (%w)", resp.Status(), ErrUnexpected)
	}
}
