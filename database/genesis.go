package database

import (
	"encoding/json"
	"io/ioutil"
)

var genesisJSON = `
{
	"genesis_time": "2020-08-17T15:53:00.000000000Z",
	"chain_id": "the-blockchain-bar-ledger",
	"balances": {
		"paulc": 1000000
	}
}`

type genesis struct {
	Balances map[Account]uint `json:"balances"`
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
