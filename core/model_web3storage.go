package core

import (
	"encoding/json"
	"fmt"

	"github.com/covalenthq/ipfs-pinner/openapi"
	cid "github.com/ipfs/go-cid"
)

// PinataResponseGetter Getter for PinataResponse object
type Web3StorageResponseGetter interface {
	fmt.Stringer
	json.Marshaler
	CidGetter
}

type web3StorageResponseObject struct {
	openapi.Web3StorageCarResponse
}

func NewWeb3StorageResponseGetter(res openapi.Web3StorageCarResponse) Web3StorageResponseGetter {
	return &web3StorageResponseObject{res}
}

func (pro *web3StorageResponseObject) GetCid() cid.Cid {
	c, err := cid.Parse(*pro.Cid)
	if err != nil {
		return cid.Undef
	}

	return c
}

func (pro *web3StorageResponseObject) MarshalJSON() ([]byte, error) {
	jsonStr := fmt.Sprintf("{\"cid\": %v}", pro.GetCid())
	return []byte(jsonStr), nil
}

func (pro *web3StorageResponseObject) String() string {
	marshalled, err := json.MarshalIndent(pro, "", "\t")
	if err != nil {
		return ""
	}

	return string(marshalled)
}
