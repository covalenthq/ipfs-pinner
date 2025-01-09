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
	"github.com/web3-storage/go-ucanto/did"
	"github.com/web3-storage/go-ucanto/principal/ed25519/signer"
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

type Config struct {
	portNumber          int
	w3AgentKey          string
	w3AgentDid          did.DID
	delegationProofPath string
	ipfsGatewayUrls     []string
	enableGC            bool
}

func NewConfig(portNumber int, w3AgentKey string, w3AgentDid did.DID, delegationProofPath string, ipfsGatewayUrls []string, enableGC bool) *Config {
	return &Config{portNumber, w3AgentKey, w3AgentDid, delegationProofPath, ipfsGatewayUrls, enableGC}
}

var (
	emptyBytes = []byte("")

	DOWNLOAD_TIMEOUT = 12 * time.Minute // download can take a lot of time if it's not locally present
	UPLOAD_TIMEOUT   = 60 * time.Second // uploads of around 6MB files happen in less than 10s typically

	maxMemory = int64(100 << 20) // 100 MB
)

func main() {
	portNumber := flag.Int("port", 3001, "port number for the server")

	w3AgentKey := flag.String("w3-agent-key", "", "w3 agent key")
	w3DelegationFile := flag.String("w3-delegation-file", "", "w3 delegation file")

	ipfsGatewayUrls := flag.String("ipfs-gateway-urls", "https://w3s.link/ipfs/%s,https://dweb.link/ipfs/%s,https://ipfs.io/ipfs/%s", "comma separated list of ipfs gateway urls")

	enableGC := flag.Bool("enable-gc", false, "enable garbage collection")

	flag.Parse()
	core.Version()

	if *w3AgentKey == "" {
		log.Fatalf("w3 agent key is required")
	}

	if *w3DelegationFile == "" {
		log.Fatalf("w3 delegation file is required")
	}

	_, err := os.ReadFile(*w3DelegationFile)
	if err != nil {
		log.Fatalf("error reading delegation proof file: %v", err)
	}

	agentSigner, err := signer.Parse(*w3AgentKey)
	if err != nil {
		log.Fatalf("error parsing agent signer: %v", err)
	}

	log.Printf("agent did: %s", agentSigner.DID().DID().String())

	config := NewConfig(*portNumber, *w3AgentKey, agentSigner.DID().DID(), *w3DelegationFile, strings.Split(*ipfsGatewayUrls, ","), *enableGC)

	setUpAndRunServer(*config)
}

func setUpAndRunServer(config Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := http.NewServeMux()
	httpState := NewState()

	clientCreateReq := client.NewClientRequest(core.Web3Storage).
		W3AgentKey(config.w3AgentKey).
		W3AgentDid(config.w3AgentDid).
		DelegationProofPath(config.delegationProofPath).
		GcEnable(config.enableGC)
	// check if cid compute true works with car uploads
	nodeCreateReq := pinner.NewNodeRequest(clientCreateReq, config.ipfsGatewayUrls, config.enableGC).CidVersion(1).CidComputeOnly(false)
	node := pinner.NewPinnerNode(ctx, *nodeCreateReq)

	mux.Handle("/upload", recoveryWrapper(uploadHttpHandler(node)))
	mux.Handle("/get", recoveryWrapper(downloadHttpHandler(node)))
	mux.Handle("/cid", recoveryWrapper(cidHttpHandler(node)))
	mux.Handle("/health", recoveryWrapper(healthHttpHandler(httpState)))
	mux.Handle("/sabotage", recoveryWrapper(sabotageHttpHandler(httpState)))
	mux.Handle("/recover", recoveryWrapper(recoverHttpHandler(httpState)))
	mux.Handle("/timeout", recoveryWrapper(timeoutHttpHandler(httpState)))

	log.Print("Listening...")
	serverStopCommandSetup()
	err := http.ListenAndServe(":"+strconv.Itoa(config.portNumber), mux)
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
		log.Println("Exit command received. Exiting now...")

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

		for _, v := range contents {
			ccid, err := uploadHandler(v, node)
			if err != nil {
				respondError(w, err)
				return
			} else {
				succStr := fmt.Sprintf("{\"cid\": \"%s\"}", ccid.String())
				_, err := w.Write([]byte(succStr))
				if err != nil {
					log.Println("error writing data to connection: %w", err)
				}
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
		contents, err := readContentFromRequest(r)
		if err != nil {
			respondError(w, err)
			return
		}

		for _, v := range contents {
			ccid, err := cidHandler(v, node)
			if err != nil {
				respondError(w, err)
				return
			} else {
				succStr := fmt.Sprintf("{\"cid\": \"%s\"}", ccid.String())
				_, err := w.Write([]byte(succStr))
				if err != nil {
					log.Println("error writing data to connection: %w", err)
				}
			}
		}
	}

	return http.HandlerFunc(fn)
}

func readContentFromRequest(r *http.Request) (map[string][]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Printf("failed to close file: %v", err)
				}
			}()

			data, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			files[fileHeader.Filename] = data
		}
	}
	return files, nil
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

func cidHandler(contents []byte, node pinner.PinnerNode) (cid.Cid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), UPLOAD_TIMEOUT)
	defer cancel()

	fcid, err := node.UnixfsService().GenerateDag(ctx, bytes.NewReader(contents))
	if err != nil {
		log.Printf("%v", err)
		return cid.Undef, err
	}

	log.Printf("cidHandler: generated dag has root cid: %s\n", fcid)
	return fcid, nil
}

func uploadHandler(contents []byte, node pinner.PinnerNode) (cid.Cid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), UPLOAD_TIMEOUT)
	defer cancel()

	fcid, err := node.UnixfsService().GenerateDag(ctx, bytes.NewReader(contents))
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

	defer func() {
		if err := carf.Close(); err != nil {
			log.Printf("failed to close car file: %v", err)
		}
	}()

	defer func() {
		err := syscall.Unlink(carf.Name())
		if err != nil {
			log.Printf("error in unlinking:%v", err)
		}
	}()

	log.Printf("car file location: %s\n", carf.Name())

	err = node.CarExporter().Export(ctx, fcid, carf)
	if err != nil {
		// For the non-deferred calls:
		if err := carf.Close(); err != nil {
			log.Printf("failed to close car file: %v", err)
		}
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

	if err := carf.Close(); err != nil {
		log.Printf("failed to close car file: %v", err)
	}
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
