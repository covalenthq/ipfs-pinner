// unixFsOptions := []options.UnixfsAddOption{options.Unixfs.CidVersion(1)}
// 	if hashOnly {
// 		unixFsOptions = append(unixFsOptions, options.Unixfs.HashOnly(true))
// 	}
// 	rpath, err := api.Unixfs().Add(ctx, node, unixFsOptions...) //, options.Unixfs.HashOnly(true))
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(rpath.Cid())

package dag

import (
	"context"

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

func (api *unixfsApi) GenerateDag(ctx context.Context, node files.Node) (cid.Cid, error) {
	rpath, err := api.ipfs.Unixfs().Add(ctx, node, api.addOptions...)
	if err != nil {
		return cid.Undef, err
	}
	return rpath.Cid(), nil
}
