package coreapi

import (
	"context"
	"log"

	"github.com/ipfs/go-ipfs/config"
	"github.com/ipfs/go-ipfs/core"
	icore "github.com/ipfs/go-ipfs/core/coreapi"
	corerepo "github.com/ipfs/go-ipfs/core/corerepo"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

type coreApiImpl struct {
	coreiface.CoreAPI
	gci    *garbageCollectorImpl
	config *config.Config
}

func NewCoreExtensionApi(ipfsNode *core.IpfsNode) CoreExtensionAPI {
	impl := &coreApiImpl{}
	impl.gci = &garbageCollectorImpl{node: ipfsNode}
	impl.CoreAPI, _ = icore.NewCoreAPI(ipfsNode)
	impl.config, _ = ipfsNode.Repo.Config()
	return impl
}

func (impl *coreApiImpl) GC() GarbageCollectAPI {
	return impl.gci
}

func (impl *coreApiImpl) Config() *config.Config {
	return impl.config
}

type garbageCollectorImpl struct {
	node *core.IpfsNode
}

func (gci *garbageCollectorImpl) GarbageCollect(ctx context.Context) {
	err := corerepo.GarbageCollect(gci.node, ctx)
	if err != nil {
		log.Println("error getting garbage collector: %w", err)
	}
}
