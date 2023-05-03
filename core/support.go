package core

import (
	"fmt"
	"os"
	"runtime"
)

type PinningService string

const (
	Pinata      PinningService = "pinata"
	Web3Storage PinningService = "web3.storage"
	Other       PinningService = "other"
	// IpfsPinnerVersionMajor is Major version component of the current release
	IpfsPinnerVersionMajor = 0
	// IpfsPinnerVersionMinor is Minor version component of the current release
	IpfsPinnerVersionMinor = 1
	// IpfsPinnerVersionPatch is Patch version component of the current release
	IpfsPinnerVersionPatch = 12
	clientIdentifier       = "ipfs-pinner" // Client identifier to advertise over the network
)

func (f PinningService) GetPinningServiceBaseUrl() string {
	switch f {
	case Pinata:
		return "https://api.pinata.cloud/psa"
	case Web3Storage:
		return "https://api.web3.storage"
	}

	return ""
}

func (f PinningService) GetFilePinBaseUrl() string {
	switch f {
	case Pinata:
		return "https://api.pinata.cloud"
	case Web3Storage:
		return "https://api.web3.storage"
	}
	panic("unsupported file pin support")
}

func (f PinningService) String() string {
	return string(f)
}

// IpfsPinnerVersion holds the textual version string.
var IpfsPinnerVersion = func() string {
	return fmt.Sprintf("%d.%d.%d", IpfsPinnerVersionMajor, IpfsPinnerVersionMinor, IpfsPinnerVersionPatch)
}()

// Version Provides version info on bsp agent binary
func Version() {
	fmt.Println(clientIdentifier)
	fmt.Println("ipfs-pinner Version:", IpfsPinnerVersion)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
}
