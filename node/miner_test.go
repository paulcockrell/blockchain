package node

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/paulcockrell/blockchain/database"
	"github.com/paulcockrell/blockchain/wallet"
)

func TestValidBlockHash(t *testing.T) {
	hexHash := "000000afd252fe0b0851ffc722ed879e282d642cafa65efb396436136565190f"
	var hash = database.Hash{}

	// convert to raw bytes
	hex.Decode(hash[:], []byte(hexHash))

	// validate hash
	isValid := database.IsBlockHashValid(hash)
	if !isValid {
		t.Fatalf("hash %q with 6 zeroes should be valid", hexHash)
	}
}

func TextInvalidBlockHash(t *testing.T) {
	hexHash := "000001fa04f8160395c387277f8b2f14837603383d33809a4db586086168edfa"
	var hash = database.Hash{}

	hex.Decode(hash[:], []byte(hexHash))

	isValid := database.IsBlockHashValid(hash)
	if isValid {
		t.Fatal("hash is not supposed to be valid")
	}
}

func TestMine(t *testing.T) {
	miner := database.NewAccount(wallet.PaulcAccount)
	pendingBlock := createRandomPendingBlock(miner)

	ctx := context.Background()

	minedBlock, err := Mine(ctx, pendingBlock)
	if err != nil {
		t.Fatal(err)
	}

	minedBlockHash, err := minedBlock.Hash()
	if err != nil {
		t.Fatal(err)
	}

	if !database.IsBlockHashValid(minedBlockHash) {
		t.Fatal()
	}

	if minedBlock.Header.Miner.String() != miner.String() {
		t.Fatal("mined block miner should equal miner from pending block")
	}
}

func TestMineWithTimeout(t *testing.T) {
	miner := database.NewAccount(wallet.PaulcAccount)
	pendingBlock := createRandomPendingBlock(miner)

	ctx, _ := context.WithTimeout(context.Background(), time.Microsecond*100)

	_, err := Mine(ctx, pendingBlock)
	if err == nil {
		t.Fatal(err)
	}
}

func createRandomPendingBlock(miner common.Address) PendingBlock {
	return NewPendingBlock(
		database.Hash{},
		1,
		miner,
		[]database.Tx{
			database.Tx{
				From:  database.NewAccount(wallet.PaulcAccount),
				To:    database.NewAccount(wallet.DavecAccount),
				Value: 1,
				Time:  1579451695,
				Data:  "",
			},
		},
	)
}
