package main

import (
	"context"
	"log"
	"os"

	"github.com/covalenthq/ipfs-pinner/core"
	"github.com/covalenthq/ipfs-pinner/coreapi"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
)

var WEB3_JWT = "WEB3_JWT"
var UPLOAD_FILE = "/Users/sudeep/Downloads/data.out"

func main() {

	token, present := os.LookupEnv(WEB3_JWT)
	if !present {
		log.Fatalf("token (%s) not found in env", WEB3_JWT)
	}
	ctx := context.Background()
	clientCreateReq := client.NewClientRequest(core.Web3Storage).BearerToken(token)
	// check if cid compute true works with car uploads
	nodeCreateReq := coreapi.NewNodeRequest(clientCreateReq).CidVersion(0).CidComputeOnly(false)
	node := coreapi.NewPinnerNode(*nodeCreateReq)

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

	carFileName := carf.Name()
	log.Printf("car file location: %s\n", carf.Name())

	err = node.CarExporter().Export(ctx, fcid, carf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	carf.Close()

	carf, err = os.Open(carFileName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	ccid, err := node.PinService().UploadFile(ctx, carf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	assertEquals(fcid, ccid)
	log.Printf("the two cids match: %s\n", ccid.String())
}

func assertEquals(obj1 interface{}, obj2 interface{}) {
	if obj1 != obj2 {
		log.Fatalf("fail %v and %v doesn't match", obj1, obj2)
	}
}
