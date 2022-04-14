// Pinner node which is a composer of services and APIs provided by the ipfs-pinner.
// Specifically, it aims to provide a way to upload data to pinning services using
// either IPFS remote pinning service API or direct/car file uploads.

package coreapi

import (
	car "github.com/covalenthq/ipfs-pinner/car"
	"github.com/covalenthq/ipfs-pinner/core"
	"github.com/covalenthq/ipfs-pinner/dag"
	"github.com/covalenthq/ipfs-pinner/pinclient"
	logging "github.com/ipfs/go-log/v2"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

var logger = logging.Logger("ipfs-pinner")

type pinnerNode struct {
	ipfsCore      coreiface.CoreAPI
	carExporter   car.CarExporterAPI
	unixfsService dag.UnixfsAPI
	pinApiClient  pinclient.PinServiceAPI
}

func NewPinnerNode(req PinnerNodeCreateRequest) PinnerNode {
	node := pinnerNode{}
	var err error
	if node.ipfsCore, err = core.CreateIpfsCoreApi(req.cidComputeOnly); err != nil {
		logger.Fatal("error initializing ipfs node: ", err)
	}

	node.pinApiClient = pinclient.NewClient(req.pinServiceRequest, req.cidVersion)
	node.carExporter = car.NewCarExporter(node.ipfsCore)
	node.unixfsService = dag.NewUnixfsAPI(node.ipfsCore, req.cidVersion, req.cidComputeOnly)

	return &node
}

func (node *pinnerNode) CarExporter() car.CarExporterAPI {
	return node.carExporter
}

func (node *pinnerNode) PinService() pinclient.PinServiceAPI {
	return node.pinApiClient
}

func (node *pinnerNode) UnixfsService() dag.UnixfsAPI {
	return node.unixfsService
}
