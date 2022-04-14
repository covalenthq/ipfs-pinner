package core

import cid "github.com/ipfs/go-cid"

type CidGetter interface {
	GetCid() cid.Cid
}
