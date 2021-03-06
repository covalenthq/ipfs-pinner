package dag

import (
	"context"
	"fmt"
	"io"

	"github.com/covalenthq/ipfs-pinner/coreapi"
	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

type unixfsApi struct {
	ipfs       coreapi.CoreExtensionAPI
	addOptions []options.UnixfsAddOption
}

func NewUnixfsAPI(ipfs coreapi.CoreExtensionAPI, cidVersion int, cidGenerationOnly bool) UnixfsAPI {
	api := unixfsApi{}
	api.addOptions = append(api.addOptions, options.Unixfs.CidVersion(cidVersion))
	api.addOptions = append(api.addOptions, options.Unixfs.HashOnly(cidGenerationOnly))
	api.addOptions = append(api.addOptions, options.Unixfs.Pin(!cidGenerationOnly))
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

func (api *unixfsApi) RemoveDag(ctx context.Context, cid cid.Cid) error {
	rp, err := api.ipfs.ResolvePath(ctx, path.New(cid.String()))
	if err != nil {
		return err
	}

	err = api.ipfs.Pin().Rm(ctx, rp, options.Pin.RmRecursive(true))
	if err != nil {
		return fmt.Errorf("error removing pin recursively: %v", err)
	}

	api.ipfs.GC().GarbageCollect(ctx)
	return nil
}
