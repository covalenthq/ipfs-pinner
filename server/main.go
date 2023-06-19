package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ipfs/go-cid"

	pinner "github.com/covalenthq/ipfs-pinner"
	"github.com/covalenthq/ipfs-pinner/core"
	client "github.com/covalenthq/ipfs-pinner/pinclient"
)

const (
	OK      = "OK"
	BAD     = "BAD"
	TIMEOUT = "TIMEOUT"
)

type State struct {
	status string
}

func NewState() *State {
	return &State{status: OK}
}

var (
	emptyBytes       = []byte("")
	WEB3_JWT         = "WEB3_JWT"
	DOWNLOAD_TIMEOUT = 12 * time.Minute // download can take a lot of time if it's not locally present
	UPLOAD_TIMEOUT   = 60 * time.Second // uploads of around 6MB files happen in less than 10s typically
)

func main() {
	portNumber := flag.Int("port", 3001, "port number for the server")
	token := flag.String("jwt", "", "JWT token for web3.storage")
	ipfsGatewayUrls := flag.String("ipfs-gateway-urls", "https://w3s.link/ipfs/%s,https://dweb.link/ipfs/%s,https://ipfs.io/ipfs/%s", "comma separated list of ipfs gateway urls")

	flag.Parse()
	core.Version()
	setUpAndRunServer(*portNumber, *token, *ipfsGatewayUrls)
}

func setUpAndRunServer(portNumber int, token, ipfsGatewayUrls string) {
	mux := http.NewServeMux()
	httpState := NewState()
	if token == "" {
		var present bool
		token, present = os.LookupEnv(WEB3_JWT)
		if !present {
			log.Fatalf("token (%s) not found in env", WEB3_JWT)
		}
	}

	var ipfsGatewayUrlArr []string
	if ipfsGatewayUrls != "" {
		ipfsGatewayUrlArr = strings.Split(ipfsGatewayUrls, ",")
		for _, ipfsUrlStr := range ipfsGatewayUrlArr {
			if !strings.Contains(ipfsUrlStr, "%s") {
				log.Fatalf("ipfs gateway url %s does not contain %%s", ipfsUrlStr)
			}

			if _, err := url.Parse(fmt.Sprintf(ipfsUrlStr, "sample_cid")); err != nil {
				log.Fatalf("ipfs gateway url %s is not a valid url", ipfsUrlStr)
			}
		}
	} else {
		log.Fatalf("ipfs gateway urls not found in params")
	}

	clientCreateReq := client.NewClientRequest(core.Web3Storage).BearerToken(token)
	// check if cid compute true works with car uploads
	nodeCreateReq := pinner.NewNodeRequest(clientCreateReq, ipfsGatewayUrlArr).CidVersion(1).CidComputeOnly(false)
	node := pinner.NewPinnerNode(*nodeCreateReq)

	mux.Handle("/upload", recoveryWrapper(uploadHttpHandler(node)))
	mux.Handle("/get", recoveryWrapper(downloadHttpHandler(node)))
	mux.Handle("/cid", recoveryWrapper(cidHttpHandler(node)))
	mux.Handle("/health", recoveryWrapper(healthHttpHandler(httpState)))
	mux.Handle("/sabotage", recoveryWrapper(sabotageHttpHandler(httpState)))
	mux.Handle("/recover", recoveryWrapper(recoverHttpHandler(httpState)))
	mux.Handle("/timeout", recoveryWrapper(timeoutHttpHandler(httpState)))

	log.Print("Listening...")
	serverStopCommandSetup()
	err := http.ListenAndServe(":"+strconv.Itoa(portNumber), mux)
	if err != nil {
		log.Println("error listening and serving on TCP network: %w", err)
	}
}

func serverStopCommandSetup() {
	// Catch the Ctrl-C and SIGTERM from kill command

	ch := make(chan os.Signal, 1)

	// Ctrl-C : sigint / os.Interrupt
	// kill command : sigterm
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)
		log.Println("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.

		os.Exit(0)

	}()
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
		contents, err := readContentFromRequest(r)
		if err != nil {
			respondError(w, err)
			return
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

func cidHttpHandler(node pinner.PinnerNode) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("not reached here")
		contents, err := readContentFromRequest(r)
		if err != nil {
			respondError(w, err)
			return
		}

		ccid, err := cidHandler(contents, node)
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

func readContentFromRequest(r *http.Request) (string, error) {
	mreader, err := r.MultipartReader()
	if err != nil {
		return "", err
	}

	var contents string = ""

	for {
		part, err := mreader.NextPart()
		if err == io.EOF {
			break
		}

		pcontents, err := io.ReadAll(part)
		if err != nil {
			return "", err
		}

		contents += string(pcontents)
	}

	return contents, nil
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

func cidHandler(contents string, node pinner.PinnerNode) (cid.Cid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), UPLOAD_TIMEOUT)
	defer cancel()

	fcid, err := node.UnixfsService().GenerateDag(ctx, bytes.NewReader([]byte(contents)))
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	log.Printf("cidHandler: generated dag has root cid: %s\n", fcid)
	return fcid, nil
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

func healthHttpHandler(s *State) http.Handler {
	// Check the health of the server and return a status code accordingly
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received /health request:", "source=", r.RemoteAddr, "status=", s.status)
		switch s.status {
		case OK:
			_, err := io.WriteString(w, "I'm healthy")
			if err != nil {
				log.Println("cannot write %w string", err)
			}
			return
		case BAD:
			http.Error(w, "Internal Error", 500)
			return
		case TIMEOUT:
			time.Sleep(30 * time.Second)
			return
		default:
			_, err := io.WriteString(w, "UNKNOWN")
			if err != nil {
				log.Println("cannot write %w string", err)
			}
			return
		}
	}
	return http.HandlerFunc(fn)
}

func sabotageHttpHandler(s *State) http.Handler {
	fn := func(w http.ResponseWriter, _ *http.Request) {
		s.status = BAD
		_, err := io.WriteString(w, "Sabotage ON")
		if err != nil {
			log.Println("cannot write %w string", err)
		}
	}
	return http.HandlerFunc(fn)
}

func recoverHttpHandler(s *State) http.Handler {
	fn := func(w http.ResponseWriter, _ *http.Request) {
		s.status = OK
		_, err := io.WriteString(w, "Recovered.")
		if err != nil {
			log.Println("cannot write %w string", err)
		}
	}
	return http.HandlerFunc(fn)
}

func timeoutHttpHandler(s *State) http.Handler {
	fn := func(w http.ResponseWriter, _ *http.Request) {
		s.status = TIMEOUT
		_, err := io.WriteString(w, "Configured to timeout.")
		if err != nil {
			log.Println("cannot write %w string", err)
		}
	}
	return http.HandlerFunc(fn)
}
