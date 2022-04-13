package pinclient

import (
	"net/http"

	"github.com/covalenthq/ipfs-pinner/core"
)

type ClientCreateRequest struct {
	ps                    core.PinningService
	pinningServiceBaseUrl string
	filePinBaseUrl        string
	bearerToken           string
	httpClient            *http.Client
}

func NewClientRequest(ps core.PinningService) ClientCreateRequest {
	request := ClientCreateRequest{ps: ps}
	request.pinningServiceBaseUrl = ps.GetPinningServiceBaseUrl()
	request.filePinBaseUrl = ps.GetFilePinBaseUrl()
	return request
}

func (r ClientCreateRequest) BearerToken(token string) ClientCreateRequest {
	r.bearerToken = token
	return r
}

func (r ClientCreateRequest) PinningServiceBaseUrl(url string) ClientCreateRequest {
	r.pinningServiceBaseUrl = url
	return r
}

func (r ClientCreateRequest) FilePinBaseUrl(url string) ClientCreateRequest {
	r.filePinBaseUrl = url
	return r
}

func (r ClientCreateRequest) HttpClient(client http.Client) ClientCreateRequest {
	r.httpClient = &client
	return r
}
