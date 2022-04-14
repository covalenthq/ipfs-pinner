// defines the api for the ipfs-pinner node

package coreapi

import (
	car "github.com/covalenthq/ipfs-pinner/car"
	"github.com/covalenthq/ipfs-pinner/dag"
	"github.com/covalenthq/ipfs-pinner/pinclient"
)

type PinnerNode interface {
	CarExporter() car.CarExporterAPI
	PinService() pinclient.PinServiceAPI
	UnixfsService() dag.UnixfsAPI
}
