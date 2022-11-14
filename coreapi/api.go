package coreapi

import (
	"context"

	"github.com/ipfs/go-ipfs/config"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

type CoreExtensionAPI interface {
	coreiface.CoreAPI
	GC() GarbageCollectAPI
	Config() *config.Config
}

type GarbageCollectAPI interface {
	GarbageCollect(ctx context.Context)
}
