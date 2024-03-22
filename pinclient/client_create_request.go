package pinclient

import (
	"net/http"

	"github.com/covalenthq/ipfs-pinner/core"
	"github.com/web3-storage/go-ucanto/did"
)

type ClientCreateRequest struct {
	ps                    core.PinningService
	pinningServiceBaseUrl string
	filePinBaseUrl        string
	//bearerToken           string

	W3_AgentKey            string
	W3_AgentDid            did.DID
	W3_DelegationProofPath string
	GC_Enable              bool

	httpClient *http.Client
}

func NewClientRequest(ps core.PinningService) ClientCreateRequest {
	request := ClientCreateRequest{ps: ps}
	request.pinningServiceBaseUrl = ps.GetPinningServiceBaseUrl()
	request.filePinBaseUrl = ps.GetFilePinBaseUrl()
	return request
}

func (r ClientCreateRequest) W3AgentKey(key string) ClientCreateRequest {
	r.W3_AgentKey = key
	return r
}

func (r ClientCreateRequest) W3AgentDid(did did.DID) ClientCreateRequest {
	r.W3_AgentDid = did
	return r
}

func (r ClientCreateRequest) DelegationProofPath(proofPath string) ClientCreateRequest {
	r.W3_DelegationProofPath = proofPath
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

func (r ClientCreateRequest) GcEnable(gcEnable bool) ClientCreateRequest {
	r.GC_Enable = gcEnable
	return r
}
