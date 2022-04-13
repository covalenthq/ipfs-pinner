package dag

import (
	"context"
	"io"

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

func (api *unixfsApi) GenerateDag(ctx context.Context, reader io.Reader) (cid.Cid, error) {
	node := files.NewReaderFile(reader)
	rpath, err := api.ipfs.Unixfs().Add(ctx, node, api.addOptions...)
	if err != nil {
		return cid.Undef, err
	}
	return rpath.Cid(), nil
}
