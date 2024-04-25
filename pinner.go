// Pinner node which is a composer of services and APIs provided by the ipfs-pinner.
// Specifically, it aims to provide a way to upload data to pinning services using
// either IPFS remote pinning service API or direct/car file uploads.

package pinner

import (
	"context"
	"log"

	car "github.com/covalenthq/ipfs-pinner/car"
	"github.com/covalenthq/ipfs-pinner/core"
	"github.com/covalenthq/ipfs-pinner/coreapi"
	"github.com/covalenthq/ipfs-pinner/dag"
	"github.com/covalenthq/ipfs-pinner/pinclient"
	"github.com/covalenthq/ipfs-pinner/w3up"
)

type pinnerNode struct {
	ipfsCore      coreapi.CoreExtensionAPI
	carExporter   car.CarExporterAPI
	unixfsService dag.UnixfsAPI
	pinApiClient  pinclient.PinServiceAPI
}

func NewPinnerNode(ctx context.Context, req PinnerNodeCreateRequest) PinnerNode {
	node := pinnerNode{}
	ipfsNode, err := core.CreateIpfsNode(ctx, req.cidComputeOnly)
	if err != nil {
		log.Fatal("error initializing ipfs node: ", err)
	}

	node.ipfsCore = coreapi.NewCoreExtensionApi(ipfsNode)
	if req.enableGC {
		log.Print("enabling garbage collection....")
		node.ipfsCore.GC().InitPeriodicGC(ctx)
	}

	//SETUP W3UP
	log.Print("setting up w3up for uploads....")
	w3up := w3up.NewW3up(req.pinServiceRequest.W3_AgentKey, req.pinServiceRequest.W3_AgentDid, req.pinServiceRequest.W3_DelegationProofPath)
	agentDid, err := w3up.WhoAmI()
	if err != nil {
		log.Fatal("error getting agent did: ", err)
	}
	log.Printf("w3up agent did: %s\n", agentDid.String())

	spaceDid, err := w3up.SpaceAdd()
	if err != nil {
		log.Fatal("error adding space: ", err)
	}
	log.Printf("w3up space did: %s\n", spaceDid.String())
	log.Print("w3up setup complete")

	node.pinApiClient = pinclient.NewClient(req.pinServiceRequest, req.cidVersion, w3up)
	node.carExporter = car.NewCarExporter(node.ipfsCore)
	node.unixfsService = dag.NewUnixfsAPI(node.ipfsCore, req.cidVersion, req.cidComputeOnly, req.ipfsFetchUrls)

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
