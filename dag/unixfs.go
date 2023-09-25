package dag

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/covalenthq/ipfs-pinner/coreapi"
	coreiface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/coreiface/path"
	files "github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	ipldformat "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

type unixfsApi struct {
	ipfs               coreapi.CoreExtensionAPI
	offlineIpfs        coreiface.CoreAPI
	addOptions         []options.UnixfsAddOption
	httpContentFetcher *httpContentFetcher
	peeringAvailable   bool
}

var (
	emptyBytes = []byte("")
)

func NewUnixfsAPI(ipfs coreapi.CoreExtensionAPI, cidVersion int, cidGenerationOnly bool, ipfsFetchUrls []string) UnixfsAPI {
	api := unixfsApi{}
	api.addOptions = append(api.addOptions, options.Unixfs.CidVersion(cidVersion))
	api.addOptions = append(api.addOptions, options.Unixfs.HashOnly(cidGenerationOnly))
	api.addOptions = append(api.addOptions, options.Unixfs.Pin(!cidGenerationOnly))

	// default merkle dag creation options
	// we want to use the same options throughout, and provide these values explicitly
	// even if the default values by ipfs libs change in future
	api.addOptions = append(api.addOptions, options.Unixfs.Hash(mh.SHA2_256))
	api.addOptions = append(api.addOptions, options.Unixfs.Inline(false))
	api.addOptions = append(api.addOptions, options.Unixfs.InlineLimit(32))
	// for cid version, raw leaves is used by default
	api.addOptions = append(api.addOptions, options.Unixfs.RawLeaves(cidVersion == 1))
	api.addOptions = append(api.addOptions, options.Unixfs.Chunker("size-262144"))
	api.addOptions = append(api.addOptions, options.Unixfs.Layout(options.BalancedLayout))
	api.addOptions = append(api.addOptions, options.Unixfs.Nocopy(false))

	api.ipfs = ipfs

	var err error
	api.offlineIpfs, err = api.ipfs.WithOptions(options.Api.Offline(true))
	if err != nil {
		log.Fatalf("failed to start offline ipfs core")
	}

	api.peeringAvailable = len(ipfs.Config().Peering.Peers) > 0
	api.httpContentFetcher = NewHttpContentFetcher(ipfsFetchUrls)
	return &api
}

func (api *unixfsApi) GenerateDag(ctx context.Context, reader io.Reader) (cid.Cid, error) {
	node := files.NewReaderFile(reader)
	rpath, err := api.ipfs.Unixfs().Add(ctx, node, api.addOptions...)
	if err != nil {
		return cid.Undef, err
	}
	return rpath.Cid(), nil
}

func (api *unixfsApi) RemoveDag(ctx context.Context, cid cid.Cid) error {
	rp, err := api.ipfs.ResolvePath(ctx, path.New(cid.String()))
	if err != nil {
		return err
	}

	err = api.ipfs.Pin().Rm(ctx, rp, options.Pin.RmRecursive(true))
	if err != nil {
		return fmt.Errorf("error removing pin recursively: %v", err)
	}

	api.ipfs.GC().GarbageCollect(ctx)
	return nil
}

func (api *unixfsApi) Get(ctx context.Context, cid cid.Cid) ([]byte, error) {
	// if peering is available, try local+bitswap search with timeout
	// if not available, try local search with timeout
	// fallback for both cases is http gateway fetch
	cidStr := cid.String()
	log.Printf("unixfsApi.Get: getting the cid: %s\n", cidStr)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	effectiveNode := api.offlineIpfs
	if api.peeringAvailable {
		log.Println("peering available...trying local+bitswap search")
		effectiveNode = api.ipfs
	}

	node, err := effectiveNode.Unixfs().Get(timeoutCtx, path.New(cidStr))
	if ipldformat.IsNotFound(err) || errors.Is(err, context.DeadlineExceeded) {
		if api.peeringAvailable {
			log.Printf("trying out http search as ipfs p2p failed: %s\n", err)
		}

		content, err := api.httpContentFetcher.FetchCidViaHttp(ctx, cidStr)
		if err != nil {
			log.Printf("error fetching: %s", err)
			return emptyBytes, err
		}
		log.Println("got the content!")

		return content, nil
	} else if err != nil {
		log.Printf("failed to fetch via ipfs p2p: %s\n", err)
		return emptyBytes, err
	}

	//api.recurse(node)
	switch val := node.(type) {
	case files.File:
		return api.readFile(val)
	case files.Directory:
		return api.readFirstFile(val)
	default:
		return emptyBytes, fmt.Errorf("unknown node type %s fetched for %s", reflect.TypeOf(node).String(), cidStr)
	}
}

func (api *unixfsApi) readFile(fnd files.File) ([]byte, error) {
	data, err := io.ReadAll(fnd)
	if err != nil {
		return emptyBytes, fmt.Errorf("error reading data: %v", err)
	}

	return data, nil
}

// reads the first file in the directory
func (api *unixfsApi) readFirstFile(dir files.Directory) ([]byte, error) {
	it := dir.Entries()
	if it.Next() {
		node := it.Node()
		fnd, ok := node.(files.File)
		if !ok {
			return emptyBytes, fmt.Errorf("node of type: %s for %s", reflect.TypeOf(dir).String(), it.Name())
		}

		return api.readFile(fnd)
	}

	return emptyBytes, fmt.Errorf("node %s entries is empty: %v", dir.Entries().Name(), it.Err())
}

// Meant for debugging purposes, do not use in prod
//
//lint:ignore U1000 function which traverses through the node and prints debug info.
func (api *unixfsApi) recurse(dir files.Node) {
	switch val := dir.(type) {
	case files.File:
		size, err := val.Size()
		if err != nil {
			log.Println("found error")
		}
		log.Printf("file found finally (%d)\n", size)
	case files.Directory:
		it := val.Entries()
		for it.Next() {
			name := it.Name()
			log.Println("name: ", name)
			api.recurse(it.Node())
		}

		err := it.Err()
		log.Printf("some error:%s\n", err)
	default:
		log.Fatalf("oh no %s not found", val)
	}
}
