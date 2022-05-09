package coreapi

import (
	"context"

	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

type CoreExtensionAPI interface {
	coreiface.CoreAPI
	GC() GarbageCollectAPI
}

type GarbageCollectAPI interface {
	GarbageCollect(ctx context.Context)
}
