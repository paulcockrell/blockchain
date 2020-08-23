package database

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
)

var genesisJSON = `
{
	"genesis_time": "2020-08-17T15:53:00.000000000Z",
	"chain_id": "the-blockchain-bar-ledger",
	"balances": {
		"0xb61E2B65e6066b0575EdD91f992B8ee8Dbd96481": 1000000
	}
}`

type Genesis struct {
	Balances map[common.Address]uint `json:"balances"`
}

func writeGenesisToDisk(path string, genesis []byte) error {
	return ioutil.WriteFile(path, genesis, 0644)
}

func loadGenesis(path string) (Genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}

	var loadedGenesis Genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return Genesis{}, err
	}

	return loadedGenesis, nil
}
