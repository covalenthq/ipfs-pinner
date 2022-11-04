package dag

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/covalenthq/ipfs-pinner/coreapi"
	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	ipldformat "github.com/ipfs/go-ipld-format"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

type unixfsApi struct {
	ipfs        coreapi.CoreExtensionAPI
	offlineIpfs coreiface.CoreAPI
	addOptions  []options.UnixfsAddOption
}

var (
	emptyBytes = []byte("")
)

func NewUnixfsAPI(ipfs coreapi.CoreExtensionAPI, cidVersion int, cidGenerationOnly bool) UnixfsAPI {
	api := unixfsApi{}
	api.addOptions = append(api.addOptions, options.Unixfs.CidVersion(cidVersion))
	api.addOptions = append(api.addOptions, options.Unixfs.HashOnly(cidGenerationOnly))
	api.addOptions = append(api.addOptions, options.Unixfs.Pin(!cidGenerationOnly))
	api.ipfs = ipfs

	var err error
	api.offlineIpfs, err = api.ipfs.WithOptions(options.Api.Offline(true))
	if err != nil {
		log.Fatalf("failed to start offline ipfs core")
	}
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
	// try fetching from local first, and if not found then go for dweb.link/ipfs instead of bitswap
	cidStr := cid.String()
	fmt.Printf("unixfsApi.Get: getting the cid: %s\n", cidStr)

	node, err := api.offlineIpfs.Unixfs().Get(ctx, path.New(cidStr))
	if ipldformat.IsNotFound(err) {
		fmt.Println("trying out dweb.link")
		resp, err := http.Get(fmt.Sprintf("https://dweb.link/ipfs/%s", cidStr))
		if err != nil {
			return emptyBytes, err
		}

		return ioutil.ReadAll(resp.Body)
	} else if err != nil {
		fmt.Println("failed to fetch from offline")
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
	data, err := ioutil.ReadAll(fnd)
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

//lint:ignore U1000 function which traverses through the node and prints debug info.
// Meant for debugging purposes, do not use in prod
func (api *unixfsApi) recurse(dir files.Node) {
	switch val := dir.(type) {
	case files.File:
		size, err := val.Size()
		if err != nil {
			fmt.Println("found error")
		}
		fmt.Printf("file found finally (%d)\n", size)
	case files.Directory:
		it := val.Entries()
		for it.Next() {
			name := it.Name()
			fmt.Println("name: ", name)
			api.recurse(it.Node())
		}

		err := it.Err()
		fmt.Printf("some error:%s\n", err)
	default:
		log.Fatalf("oh no %s not found", val)
	}
}
