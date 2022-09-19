package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	pinner "github.com/covalenthq/ipfs-pinner"
	"github.com/covalenthq/ipfs-pinner/core"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
)

var WEB3_JWT = "WEB3_JWT"
var UPLOAD_FILE = "temp.txt"

func pinningHandler(address string, node pinner.PinnerNode) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		if addr := r.FormValue("address"); addr != "" {
			address = addr
		}

		ctx := context.Background()

		file, err := os.Open(address)
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

		w.Write([]byte(ccid.String()))

		log.Printf("removing dag...")
		curr := time.Now().UnixMilli()
		err = node.UnixfsService().RemoveDag(ctx, ccid)
		after := time.Now().UnixMilli()
		log.Println("time taken:", after-curr)
		if err != nil {
			log.Fatalf("%v", err)
		}

	}

	return http.HandlerFunc(fn)

}

func main() {
	port_number := os.Args[1]

	println(port_number)

	setUpAndRunServer(port_number)
}

func setUpAndRunServer(port_number string) {
	mux := http.NewServeMux()

	token, present := os.LookupEnv(WEB3_JWT)
	if !present {
		log.Fatalf("token (%s) not found in env", WEB3_JWT)
	}

	clientCreateReq := client.NewClientRequest(core.Web3Storage).BearerToken(token)
	// check if cid compute true works with car uploads
	nodeCreateReq := pinner.NewNodeRequest(clientCreateReq).CidVersion(0).CidComputeOnly(false)
	node := pinner.NewPinnerNode(*nodeCreateReq)

	th := pinningHandler(UPLOAD_FILE, node)
	mux.Handle("/pin", th)

	log.Print("Listening...")
	http.ListenAndServe(":"+port_number, mux)
}

func assertEquals(obj1 interface{}, obj2 interface{}) {
	if obj1 != obj2 {
		log.Fatalf("fail %v and %v doesn't match", obj1, obj2)
	}
}
