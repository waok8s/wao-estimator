package estimator

import (
	"context"
	"fmt"
	"io"
	"net/http"

	http2curl "moul.io/http2curl/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator/api"
)

type PowerConsumption = api.PowerConsumption
type ClientOption = api.ClientOption
type Error = api.Error

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
	c api.ClientWithResponsesInterface
	e client.ObjectKey
}

func NewClient(server string, estimatorNamespace, estimatorName string, opts ...ClientOption) (*Client, error) {
	c, err := api.NewClientWithResponses(server, opts...)
	if err != nil {
		return nil, err
	}
	ec := Client{c: c, e: client.ObjectKey{Namespace: estimatorNamespace, Name: estimatorName}}
	return &ec, nil
}

func (c *Client) EstimatePowerConsumption(ctx context.Context, cpuMilli, numWorkloads int) (pc *PowerConsumption, apiErr *Error, requestErr error) {
	body := api.PostNamespacesNsEstimatorsNameValuesPowerconsumptionJSONRequestBody{
		CpuMilli:      cpuMilli,
		NumWorkloads:  numWorkloads,
		WattIncreases: nil,
	}
	resp, err := c.c.PostNamespacesNsEstimatorsNameValuesPowerconsumptionWithResponse(ctx, c.e.Namespace, c.e.Name, body)
	if err != nil {
		return nil, nil, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
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
