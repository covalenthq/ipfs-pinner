package coreapi

import (
	"context"
	"log"

	coreiface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	icore "github.com/ipfs/kubo/core/coreapi"
	corerepo "github.com/ipfs/kubo/core/corerepo"
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
