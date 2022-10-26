package dag

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
)

// merkle dag manipulator.
// There are various ways to generate merkle dag for some content. The implementors
// of this interface should capture those settings. Note that this means that those
// settings remain the same for an instance (across all requests on that instance).
// Though some implementor can use context to change the settings per request.
type UnixfsAPI interface {
	GenerateDag(ctx context.Context, reader io.Reader) (cid.Cid, error)

	// free the DAG stored on disk, thereby freeing up space
	// to be used only if the dag is persisted on disk
	// this would also do a GC in order to make the space immediately
	// available
	RemoveDag(ctx context.Context, cid cid.Cid) error

	// get the data stored in ipfs referenced by cid.
	// If the node is online (networking enabled), this would also search
	// in other ipfs nodes (using bitswap).
	// If the dag is cleaned up from local store, it might take time for the
	// data to be pinned/available on remote nodes, which means that a "upload"
	// followed immediately by a "get" might not work.
	Get(ctx context.Context, cid cid.Cid) ([]byte, error)
}
