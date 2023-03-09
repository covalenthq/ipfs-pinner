package coreapi

import (
	"context"

	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/config"
)

type CoreExtensionAPI interface {
	coreiface.CoreAPI
	GC() GarbageCollectAPI
	Config() *config.Config
}

type GarbageCollectAPI interface {
	GarbageCollect(ctx context.Context)
}
