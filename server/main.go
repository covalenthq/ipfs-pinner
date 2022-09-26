package main

import (
	"log"
	"net/http"
	"os"

	main2 "github.com/covalenthq/ipfs-pinner/binary"

	pinner "github.com/covalenthq/ipfs-pinner"
	"github.com/covalenthq/ipfs-pinner/core"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
)

var WEB3_JWT = "WEB3_JWT"
var UPLOAD_FILE = "temp.txt"

// run this: go run main.go "3000"
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

	th := main2.PinningHandler(UPLOAD_FILE, node)
	mux.Handle("/pin", th)

	log.Print("Listening...")
	http.ListenAndServe(":"+port_number, mux)
}
