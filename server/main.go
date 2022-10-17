package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/ipfs/go-cid"

	pinner "github.com/covalenthq/ipfs-pinner"
	"github.com/covalenthq/ipfs-pinner/core"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
)

var WEB3_JWT = "WEB3_JWT"

// run this: go run main.go "3000"
func main() {

	port_number := flag.String("port", "3000", "port number for the server")
	token := flag.String("jwt", "", "jwt token for web3.storage")
	// port_number := os.Args[1]

	flag.Parse()

	setUpAndRunServer(*port_number, *token)
}

func setUpAndRunServer(port_number string, token string) {
	mux := http.NewServeMux()
	if token == "" {
		var present bool
		token, present = os.LookupEnv(WEB3_JWT)
		if !present {
			log.Fatalf("token (%s) not found in env", WEB3_JWT)
		}
	}

	clientCreateReq := client.NewClientRequest(core.Web3Storage).BearerToken(token)
	// check if cid compute true works with car uploads
	nodeCreateReq := pinner.NewNodeRequest(clientCreateReq).CidVersion(0).CidComputeOnly(false)
	node := pinner.NewPinnerNode(*nodeCreateReq)

	th := PinningHandler(node)
	mux.Handle("/pin", th)

	log.Print("Listening...")
	http.ListenAndServe(":"+port_number, mux)
}

func PinningHandler(node pinner.PinnerNode) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if addr := r.FormValue("address"); addr != "" {
			address := addr
			ccid := recovery_wrapper(address, node)
			if ccid.String() == "b"{
				w.Write([]byte("no cid generated"))
			} else {
			w.Write([]byte(ccid.String()))
			}
		} else {
			fmt.Println("Please provide a file address for pinning! No file address found in the request.")
		}
	}
	return http.HandlerFunc(fn)
}

func recovery_wrapper(address string, node pinner.PinnerNode) cid.Cid {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	return myHandleFunc(address, node)
}

func myHandleFunc(address string, node pinner.PinnerNode) cid.Cid {
	ctx := context.Background()

	ccid := cid.Cid{}

	file, err := os.Open(address)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	fcid, err := node.UnixfsService().GenerateDag(ctx, file)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	carf, err := os.CreateTemp(os.TempDir(), "*.car")
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	err = syscall.Unlink(carf.Name())
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	log.Printf("car file location: %s\n", carf.Name())

	err = node.CarExporter().Export(ctx, fcid, carf)
	if err != nil {
		carf.Close()
		log.Printf("%v", err)
		panic(err)
	}

	carf.Seek(0, 0) // reset for read
	ccid, err = node.PinService().UploadFile(ctx, carf)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
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
		log.Printf("%v", err)
		panic(err)
	}

	return ccid
}

func assertEquals(obj1 interface{}, obj2 interface{}) {
	if obj1 != obj2 {
		log.Printf("fail %v and %v doesn't match", obj1, obj2)
	}
}
