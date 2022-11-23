package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

var (
	emptyBytes       = []byte("")
	WEB3_JWT         = "WEB3_JWT"
	DOWNLOAD_TIMEOUT = 12 * time.Minute // download can take a lot of time if it's not locally present
	UPLOAD_TIMEOUT   = 60 * time.Second // uploads of around 6MB files happen in less than 10s typically
)

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
	nodeCreateReq := pinner.NewNodeRequest(clientCreateReq).CidVersion(1).CidComputeOnly(false)
	node := pinner.NewPinnerNode(*nodeCreateReq)

	mux.Handle("/upload", recoveryWrapper(uploadHttpHandler(node)))
	mux.Handle("/get", recoveryWrapper(downloadHttpHandler(node)))

	log.Print("Listening...")
	err := http.ListenAndServe(":"+strconv.Itoa(portNumber), mux)
	if err != nil {
		log.Println("error listening and serving on TCP network: %w", err)
	}
}

func respondError(w http.ResponseWriter, err error) {
	err_str := fmt.Sprintf("{\"error\": \"%s\"}", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte(err_str))
	if err != nil {
		log.Println("error writing data to connection: %w", err)
	}
}

func uploadHttpHandler(node pinner.PinnerNode) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		mreader, err := r.MultipartReader()
		if err != nil {
			respondError(w, err)
			return
		}

		var contents string = ""

		for {
			part, err := mreader.NextPart()
			if err == io.EOF {
				break
			}

			pcontents, err := ioutil.ReadAll(part)
			if err != nil {
				respondError(w, err)
				return
			}

			contents += string(pcontents)
		}

		ccid, err := uploadHandler(contents, node)
		if err != nil {
			respondError(w, err)
			return
		} else {
			succ_str := fmt.Sprintf("{\"cid\": \"%s\"}", ccid.String())
			_, err := w.Write([]byte(succ_str))
			if err != nil {
				log.Println("error writing data to connection: %w", err)
			}
		}
	}
	return http.HandlerFunc(fn)
}

func downloadHttpHandler(node pinner.PinnerNode) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if cidStr := r.FormValue("cid"); cidStr != "" {
			contents, err := downloadHandler(cidStr, node)
			if err != nil {
				respondError(w, err)
			} else {
				_, err := w.Write(contents)
				if err != nil {
					log.Println("error writing data to connection: %w", err)
				}
			}
		} else {
			log.Println("Please provide a cid for fetching!")
		}
	}

	return http.HandlerFunc(fn)
}

func recoveryWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func uploadHandler(contents string, node pinner.PinnerNode) (cid.Cid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), UPLOAD_TIMEOUT)
	defer cancel()

	fcid, err := node.UnixfsService().GenerateDag(ctx, bytes.NewReader([]byte(contents)))
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	log.Printf("generated dag has root cid: %s\n", fcid)

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

	_, err = carf.Seek(0, 0)
	if err != nil {
		log.Println("error writing data to connection: %w", err)
	} // reset for read
	var ccid cid.Cid
	ccid, err = node.PinService().UploadFile(ctx, carf)
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}
	log.Printf("uploaded file has root cid: %s\n", ccid)

	carf.Close()
	return ccid, nil
}

func downloadHandler(cidStr string, node pinner.PinnerNode) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DOWNLOAD_TIMEOUT)
	defer cancel()

	cid, err := cid.Parse(cidStr)
	if err != nil {
		return emptyBytes, err
	}
	return node.UnixfsService().Get(ctx, cid)
}
