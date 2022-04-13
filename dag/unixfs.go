package dag

import (
	"context"
	"os"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

type unixfsApi struct {
	ipfs       coreiface.CoreAPI
	addOptions []options.UnixfsAddOption
}

func NewUnixfsAPI(ipfs coreiface.CoreAPI, cidVersion int, cidGenerationOnly bool) UnixfsAPI {
	api := unixfsApi{}
	api.addOptions = append(api.addOptions, options.Unixfs.CidVersion(cidVersion))
	api.addOptions = append(api.addOptions, options.Unixfs.HashOnly(cidGenerationOnly))
	api.ipfs = ipfs
	return &api
}

func (api *unixfsApi) GenerateDag(ctx context.Context, filepath string) (cid.Cid, error) {
	node, err := getNodeFor(filepath)
	if err != nil {
		return cid.Undef, err
	}
	rpath, err := api.ipfs.Unixfs().Add(ctx, node, api.addOptions...)
	if err != nil {
		return cid.Undef, err
	}
	return rpath.Cid(), nil
}

func getNodeFor(filepath string) (files.Node, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	n, err := files.NewSerialFile(filepath, true, stat)
	if err != nil {
		return nil, err
	}

	return n, nil
}
