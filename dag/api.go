package dag

import (
	"context"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
)

// generate the merkle dag from a give node.
// There are various ways to generate merkle dag for some content. The implementors
// of this interface should capture those settings. Note that this means that those
// settings remain the same for an instance (across all requests on that instance).
// Though some implementor can use context to change the settings per request.
type UnixfsAPI interface {
	GenerateDag(ctx context.Context, node files.Node) (cid.Cid, error)
}
