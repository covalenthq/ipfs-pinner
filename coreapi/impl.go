package coreapi

import (
	"context"

	"github.com/ipfs/go-ipfs/core"
	icore "github.com/ipfs/go-ipfs/core/coreapi"
	corerepo "github.com/ipfs/go-ipfs/core/corerepo"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

type coreApiImpl struct {
	coreiface.CoreAPI
	gci *garbageCollectorImpl
}

func NewCoreExtensionApi(ipfsNode *core.IpfsNode) CoreExtensionAPI {
	impl := &coreApiImpl{}
	//impl.gci = &garbageCollectorImpl{node: }
	impl.gci = &garbageCollectorImpl{node: ipfsNode}
	impl.CoreAPI, _ = icore.NewCoreAPI(ipfsNode, options.Api.Offline(true), options.Api.FetchBlocks(false))
	return impl
}

func (impl *coreApiImpl) GC() GarbageCollectAPI {
	return impl.gci
}

type garbageCollectorImpl struct {
	node *core.IpfsNode
}

func (gci *garbageCollectorImpl) GarbageCollect(ctx context.Context) {
	corerepo.GarbageCollect(gci.node, ctx)
}
