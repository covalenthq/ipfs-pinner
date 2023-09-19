// A car exporter based on boxo/ipld/car module
package car

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/covalenthq/ipfs-pinner/coreapi"
	"github.com/ipfs/go-cid"
	gocar "github.com/ipld/go-car"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

// implements the CarExporterAPI
type carExporter struct {
	api coreapi.CoreExtensionAPI
}

func NewCarExporter(api coreapi.CoreExtensionAPI) CarExporterAPI {
	return &carExporter{api}
}

func (exp *carExporter) Export(ctx context.Context, contentRoot cid.Cid, writer io.Writer) error {
	store := dagStore{dag: exp.api.Dag(), ctx: ctx}
	dag := gocar.Dag{Root: contentRoot, Selector: selectorparse.CommonSelector_ExploreAllRecursively}
	// TraverseLinksOnlyOnce is safe for an exhaustive selector but won't be when we allow
	// arbitrary selectors here
	scar := gocar.NewSelectiveCar(ctx, store, []gocar.Dag{dag}, gocar.TraverseLinksOnlyOnce())
	return scar.Write(writer)
}

func (exp *carExporter) MultiExport(ctx context.Context, contentRoot cid.Cid, writers []io.Writer) error {
	if len(writers) == 0 {
		return fmt.Errorf("no writers provided")
	}

	if len(writers) == 1 {
		return exp.Export(ctx, contentRoot, writers[0])
	}

	store := dagStore{dag: exp.api.Dag(), ctx: ctx}
	dag := gocar.Dag{Root: contentRoot, Selector: selectorparse.CommonSelector_ExploreAllRecursively}
	scar := gocar.NewSelectiveCar(ctx, store, []gocar.Dag{dag}, gocar.TraverseLinksOnlyOnce())

	preparedScar, err := scar.Prepare()
	if err != nil {
		return err
	}

	for _, writer := range writers {
		err := preparedScar.Dump(ctx, writer)
		if err != nil {
			log.Println("error writing car file: %w", err)
		}
	}

	return nil
}
