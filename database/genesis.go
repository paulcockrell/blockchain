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

type genesis struct {
	Balances map[common.Address]uint `json:"balances"`
}

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(genesisJSON), 0644)
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}
