package pinclient

import (
	"context"
	"os"

	"github.com/covalenthq/ipfs-pinner/core"
	"github.com/ipfs/go-cid"
)

type PinServiceAPI interface {
	IsIPFSSupportedFor(ps core.PinningService) bool
	Add(ctx context.Context, cid cid.Cid, opts ...AddOption) (core.PinStatusGetter, error)
	UploadFile(ctx context.Context, file *os.File) (cid.Cid, error)
	ServiceType() core.PinningService
}
