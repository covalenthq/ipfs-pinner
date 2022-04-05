package ipfs_pin_lib

type PinningService string

const (
	Pinata      PinningService = "pinata"
	Web3Storage                = "web3.storage"
	Other                      = "other"
)

func (f PinningService) filePinningSupported() bool {
	return f != Other
}

func (f PinningService) getPinningServiceBaseUrl() string {
	switch f {
	case Pinata:
		return "https://api.pinata.cloud/psa"
	case Web3Storage:
		return "https://api.web3.storage"
	}

	return ""
}

func (f PinningService) getFilePinBaseUrl() string {
	switch f {
	case Pinata:
		return "https://api.pinata.cloud"
	case Web3Storage:
		return "https://api.web3.storage"
	}
	panic("unsupported file pin support")
}
