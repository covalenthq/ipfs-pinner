//nolint:deadcode,errcheck,unused
package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	pinner "github.com/covalenthq/ipfs-pinner"
	"github.com/covalenthq/ipfs-pinner/core"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
	"github.com/ipfs/go-cid"
)

var WEB3_JWT = "WEB3_JWT"
var UPLOAD_FILE = "./main.go" // uploading current file itself

func main() {
	token, present := os.LookupEnv(WEB3_JWT)
	if !present {
		log.Fatalf("token (%s) not found in env", WEB3_JWT)
	}

	clientCreateReq := client.NewClientRequest(core.Web3Storage).BearerToken(token)
	// check if cid compute true works with car uploads
	nodeCreateReq := pinner.NewNodeRequest(clientCreateReq, []string{"https://w3s.link/ipfs/%s"}).CidVersion(1).CidComputeOnly(false)
	node := pinner.NewPinnerNode(*nodeCreateReq)
	ctx := context.Background()
	//upload(ctx, node)
	core.Version()
	download(ctx, node)
}

func upload(ctx context.Context, node pinner.PinnerNode) {
	file, err := os.Open(UPLOAD_FILE)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fcid, err := node.UnixfsService().GenerateDag(ctx, file)
	if err != nil {
		log.Fatalf("%v", err)
	}

	carf, err := os.CreateTemp(os.TempDir(), "*.car")
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = syscall.Unlink(carf.Name())
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("car file location: %s\n", carf.Name())

	err = node.CarExporter().Export(ctx, fcid, carf)
	if err != nil {
		carf.Close()
		log.Fatalf("%v", err)
	}

	carf.Seek(0, 0) // reset for read
	ccid, err := node.PinService().UploadFile(ctx, carf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	carf.Close() // should delete the file due to unlink

	assertEquals(fcid, ccid)
	log.Printf("the two cids match: %s\n", ccid.String())

	log.Printf("removing dag...")
	curr := time.Now().UnixMilli()
	err = node.UnixfsService().RemoveDag(ctx, ccid)
	after := time.Now().UnixMilli()
	log.Println("time taken:", after-curr)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func download(ctx context.Context, node pinner.PinnerNode) {
	//ccid, err := cid.Parse("bafybeifzst7cbujrqemiulznrkttouzshnqkrajiib5fp5te53ojs5sl5u") // file encapsulated in folder
	//ccid, err := cid.Parse("QmeFd8e4UaAPrPnwxBWcpqY3tMpggpWB3XYqftMpyYYLWZ") // straight up file

	//ccid, err := cid.Parse("bafybeifzst7cbujrqemiulznrkttouzshnqkrajiib5fp5te53ojs5sl5u")
	//ccid, err := cid.Parse("QmZrxsDZwrCKbcJLNf1D6GaX2fobHZpesBc3DwhVBLQ33p")
	//ccid, err := cid.Parse("bafkreic7xhzqyube57gex7okhzytg7i5eq6fvz5snpte7swy547s22bs5q")
	ccid, err := cid.Parse("bafybeic7nbudupk56j2ixdrczddzmuu3qtmlyrpqjxuy4jkaeeaminxq3e") // took about 12-15 minutes

	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("searching for pid: %v\n", ccid)
	contents, err := node.UnixfsService().Get(ctx, ccid)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// print the first 20 characters to prevent polluting the terminal.
	log.Println(string(contents)[:20])
}

func assertEquals(obj1 interface{}, obj2 interface{}) {
	if obj1 != obj2 {
		log.Fatalf("fail %v and %v doesn't match", obj1, obj2)
	}
}
