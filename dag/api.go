package dag

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
)

// generate the merkle dag from content provided by the reader.
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
}
