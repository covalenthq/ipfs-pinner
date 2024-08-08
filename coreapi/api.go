package coreapi

import (
	"context"

	"github.com/ipfs/kubo/config"
	coreiface "github.com/ipfs/kubo/core/coreiface"
)

type CoreExtensionAPI interface {
	coreiface.CoreAPI
	GC() GarbageCollectAPI
	Config() *config.Config
}

type GarbageCollectAPI interface {
	GarbageCollect(ctx context.Context)
	InitPeriodicGC(ctx context.Context) <-chan error
}
