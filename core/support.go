package core

type PinningService string

const (
	Pinata      PinningService = "pinata"
	Web3Storage PinningService = "web3.storage"
	Other       PinningService = "other"
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
