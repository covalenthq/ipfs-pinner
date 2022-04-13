// facility to create car files out of unixfs nodes

package car

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
)

type CarExporterAPI interface {
	Export(ctx context.Context, contentRoot cid.Cid, writer io.Writer) error

	// export car file to multiple destinations.
	// failure in one writer doesn't fail write attempts to other writers
	MultiExport(ctx context.Context, contentRoot cid.Cid, writers []io.Writer) error
}
