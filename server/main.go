package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/ipfs/go-cid"

	pinner "github.com/covalenthq/ipfs-pinner"
	"github.com/covalenthq/ipfs-pinner/core"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
)

var WEB3_JWT = "WEB3_JWT"

func main() {
	portNumber := flag.Int("port", 3000, "port number for the server")
	token := flag.String("jwt", "", "JWT token for web3.storage")

	flag.Parse()
	setUpAndRunServer(*portNumber, *token)
}

func setUpAndRunServer(portNumber int, token string) {
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
	http.ListenAndServe(":"+strconv.Itoa(portNumber), mux)
}

func PinningHandler(node pinner.PinnerNode) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if fp := r.FormValue("filePath"); fp != "" {
			filePath := fp
			ccid, err := recoveryWrapper(filePath, node)
			if err != nil {
				err_str := fmt.Sprintf("{\"error\": \"%s\"}", err)
				w.Write([]byte(err_str))
			} else if len(ccid.String()) != 46 {
				err_str := "{\"error\": \"invalid cid generated\"}"
				w.Write([]byte(err_str))
			} else {
				succ_str := fmt.Sprintf("{\"cid\": \"%s\"}", ccid.String())
				w.Write([]byte(succ_str))
			}
		} else {
			fmt.Println("Please provide a file filePath for pinning! No file filePath found in the request.")
		}
	}
	return http.HandlerFunc(fn)
}

func recoveryWrapper(filePath string, node pinner.PinnerNode) (cid.Cid, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	return pinningCoreHandleFunc(filePath, node)
}

func pinningCoreHandleFunc(filePath string, node pinner.PinnerNode) (cid.Cid, error) {
	ctx := context.Background()

	ccid := cid.Cid{}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}
	fcid, err := node.UnixfsService().GenerateDag(ctx, file)
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	carf, err := os.CreateTemp(os.TempDir(), "*.car")
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	defer carf.Close() // should delete the file due to unlink

	err = syscall.Unlink(carf.Name())
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	log.Printf("car file location: %s\n", carf.Name())

	err = node.CarExporter().Export(ctx, fcid, carf)
	if err != nil {
		carf.Close()
		log.Printf("%v", err)
		return cid.Undef, err
	}

	carf.Seek(0, 0) // reset for read
	ccid, err = node.PinService().UploadFile(ctx, carf)
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	carf.Close()

	log.Printf("removing dag...")
	curr := time.Now().UnixMilli()
	err = node.UnixfsService().RemoveDag(ctx, ccid)
	after := time.Now().UnixMilli()
	log.Println("time taken:", after-curr)
	if err != nil {
		log.Printf("%v", err)
	}

	return ccid, nil
}
