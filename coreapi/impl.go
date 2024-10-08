package coreapi

import (
	"context"
	"log"

	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	icore "github.com/ipfs/kubo/core/coreapi"
	coreiface "github.com/ipfs/kubo/core/coreiface"
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

func (gci *garbageCollectorImpl) InitPeriodicGC(ctx context.Context) <-chan error {
	errc := make(chan error)
	go func() {
		errc <- corerepo.PeriodicGC(ctx, gci.node)
		close(errc)
	}()
	return errc
}
