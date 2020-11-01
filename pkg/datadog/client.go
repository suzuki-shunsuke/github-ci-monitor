package datadog

import (
	"context"
	"net/http"
	"net/url"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	client httpclient.Client
	apiKey string
}

func New(apiKey string) Client {
	return Client{
		client: httpclient.New("https://api.datadoghq.com"),
		apiKey: apiKey,
	}
}

type ParamCheck struct {
	Check    string   `json:"check,omitempty"`
	HostName string   `json:"host_name,omitempty"`
	Message  string   `json:"message,omitempty"`
	Status   int      `json:"status"`
	Tags     []string `json:"tags,omitempty"`
}

func (client Client) Check(ctx context.Context, params ParamCheck) (*http.Response, error) {
	return client.client.Call(ctx, httpclient.CallParams{
		Method: http.MethodPost,
		Path:   "/api/v1/check_run",
		Query: url.Values{
			"api_key": []string{client.apiKey},
		},
		RequestBody: params,
	})
}
