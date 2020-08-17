package main

import (
	"fmt"
	"os"
	"time"

	"github.com/paulcockrell/blockchain/database"
)

func main() {
	cwd, _ := os.Getwd()
	state, err := database.NewStateFromDisk(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("paulc", "paulc", 3, ""),
			database.NewTx("paulc", "paulc", 700, "reward"),
		},
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persist()

	block1 := database.NewBlock(
		block0hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("paulc", "davec", 2000, ""),
			database.NewTx("paulc", "paulc", 100, "reward"),
			database.NewTx("davec", "paulc", 1, ""),
			database.NewTx("davec", "kimc", 1000, ""),
			database.NewTx("davec", "paulc", 50, ""),
			database.NewTx("paulc", "paulc", 600, "reward"),
		},
	)

	state.AddBlock(block1)
	state.Persist()
}
