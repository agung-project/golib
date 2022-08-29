package http

import (
	"github.com/go-resty/resty/v2"
)

type RestClient struct {
	HttpClient *resty.Client
}

func NewRestClient(baseUrl string) *RestClient {
	httpClient := resty.New()
	httpClient.SetHostURL(baseUrl)
	return &RestClient{HttpClient: httpClient}
}
