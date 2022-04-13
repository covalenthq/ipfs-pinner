package car

import (
	"context"

	blocks "github.com/ipfs/go-block-format"
	coreiface "github.com/ipfs/interface-go-ipfs-core"

	"github.com/ipfs/go-cid"
)

// struct over APIDagService which implements ReadStore interface
type dagStore struct {
	dag coreiface.APIDagService
	ctx context.Context
}

func (ds dagStore) Get(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	return ds.dag.Get(ctx, c) // Node type promoted to blocks.Block
}
