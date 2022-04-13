package coreapi

import (
	"context"

	"github.com/covalenthq/ipfs-pinner/pinclient"
)

type PinnerNodeCreateRequest struct {
	pinServiceRequest pinclient.ClientCreateRequest
	ctx               context.Context
	cidComputeOnly    bool
	cidVersion        int
}

func NewNodeRequest() *PinnerNodeCreateRequest {
	request := new(PinnerNodeCreateRequest)
	request.cidVersion = 0
	request.cidComputeOnly = true
	return request
}

func (req *PinnerNodeCreateRequest) Context(ctx context.Context) *PinnerNodeCreateRequest {
	req.ctx = ctx
	return req
}

func (req *PinnerNodeCreateRequest) PinClientRequest(clientRequest pinclient.ClientCreateRequest) *PinnerNodeCreateRequest {
	req.pinServiceRequest = clientRequest
	return req
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
