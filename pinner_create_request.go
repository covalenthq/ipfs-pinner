package pinner

import (
	"github.com/covalenthq/ipfs-pinner/pinclient"
)

type PinnerNodeCreateRequest struct {
	pinServiceRequest pinclient.ClientCreateRequest
	cidComputeOnly    bool
	cidVersion        int
}

func NewNodeRequest(clientRequest pinclient.ClientCreateRequest) *PinnerNodeCreateRequest {
	request := new(PinnerNodeCreateRequest)
	request.cidVersion = 0
	request.cidComputeOnly = true
	request.pinServiceRequest = clientRequest
	return request
}

func (req *PinnerNodeCreateRequest) CidVersion(version int) *PinnerNodeCreateRequest {
	req.cidVersion = version
	return req
}

// If only cid computation is required from the node (and no persistence of the ipfs dags to filesystem)
func (req *PinnerNodeCreateRequest) CidComputeOnly(cidComputeOnly bool) *PinnerNodeCreateRequest {
	req.cidComputeOnly = cidComputeOnly
	return req
}
