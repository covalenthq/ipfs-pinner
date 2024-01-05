package w3up

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-ucanto/did"
)

type W3up struct {
	principal           string
	agentDid            did.DID
	delegationProofPath string
	w3up_dir            string
}

func NewW3up(principal string, did did.DID, delegationProofPath string) *W3up {
	tmpDir, err := os.MkdirTemp("", "w3up_config")
	if err != nil {
		log.Fatal("error creating temp dir for w3up config: ", err)
	}
	return &W3up{principal: principal, agentDid: did, delegationProofPath: delegationProofPath, w3up_dir: tmpDir}
}

func (w3up *W3up) WhoAmI() (did.DID, error) {
	command := fmt.Sprintf("W3_STORE_NAME=%s W3_PRINCIPAL=\"%s\" w3 whoami", w3up.w3up_dir, w3up.principal)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("error initializing w3up: ", string(output))
		return did.DID{}, err
	}

	pdid, err := did.Parse(strings.TrimSpace(string(output)))
	if err != nil {
		log.Fatal("error parsing did: ", err)
		return did.DID{}, err
	}
	return pdid, nil
}

func (w3up *W3up) SpaceAdd() (did.DID, error) {
	command := fmt.Sprintf("W3_STORE_NAME=%s w3 space add %s", w3up.w3up_dir, w3up.delegationProofPath)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("error adding w3up space: ", string(output))
		return did.DID{}, err
	}

	pdid, err := did.Parse(strings.TrimSpace(string(output)))
	if err != nil {
		log.Fatal("error parsing did: ", err)
		return did.DID{}, err
	}

	return pdid, nil
}

func (w3up *W3up) UploadCarFile(carFile *os.File) (cid.Cid, error) {
	command := fmt.Sprintf("W3_STORE_NAME=%s w3 up %s --json --no-wrap --car", w3up.w3up_dir, carFile.Name())
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("error uploading car file: ", string(output))
		return cid.Undef, err
	}

	jsons := strings.TrimSpace(string(output))

	log.Printf("w3 up output: %s\n", jsons)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsons), &result)
	if err != nil {
		log.Fatal("error parsing json: ", err)
		return cid.Undef, err
	}

	root, ok := result["root"].(map[string]interface{})
	if !ok {
		log.Fatal("error type asserting root")
		return cid.Undef, fmt.Errorf("error type asserting root")
	}

	rcid, ok := root["/"].(string)
	if !ok {
		log.Fatal("error type asserting cid")
		return cid.Undef, fmt.Errorf("error type asserting cid")
	}

	return cid.Parse(rcid)
}
