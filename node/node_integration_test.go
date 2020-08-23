package node

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/paulcockrell/blockchain/database"
	"github.com/paulcockrell/blockchain/fs"
	"github.com/paulcockrell/blockchain/wallet"
)

func TestNode_Run(t *testing.T) {
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	n := New(datadir, "127.0.0.1", 8085, database.NewAccount(wallet.PaulcAccount), PeerNode{})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err = n.Run(ctx)
	if err.Error() != "http: Server closed" {
		t.Fatal("node server was suppose to close after 5s")
	}
}

func TestNode_Mining(t *testing.T) {
	andrej := database.NewAccount(wallet.PaulcAccount)
	babayaga := database.NewAccount(wallet.DavecAccount)

	// Remove the test directory if it already exists
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	// Required for AddPendingTX() to describe
	// from what node the TX came from (local node in this case)
	nInfo := NewPeerNode(
		"127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true,
	)

	// Construct a new Node instance and configure
	// Paulc as a miner
	n := New(datadir, nInfo.IP, nInfo.Port, andrej, nInfo)

	// Allow the mining to run for 30 mins, in the worst case
	ctx, closeNode := context.WithTimeout(
		context.Background(),
		time.Minute*30,
	)

	// Schedule a new TX in 3 seconds from now, in a separate thread
	// because the n.Run() few lines below is a blocking call
	go func() {
		time.Sleep(time.Second * miningIntervalSeconds / 3)
		tx := database.NewTx(andrej, babayaga, 1, "")

		_ = n.AddPendingTX(tx, nInfo)
	}()

	// Schedule a new TX in 12 seconds from now simulating
	// that it came in - while the first TX is being mined
	go func() {
		time.Sleep(time.Second*miningIntervalSeconds + 2)
		tx := database.NewTx(andrej, babayaga, 2, "")

		_ = n.AddPendingTX(tx, nInfo)
	}()

	go func() {
		// Periodically check if we mined the 2 blocks
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-ticker.C:
				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()

	// Run the node, mining and everything in a blocking call (hence the go-routines before)
	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("2 pending TX not mined into 2 under 30m")
	}
}

func TestNode_MiningStopsOnNewSyncedBlock(t *testing.T) {
	andrej := database.NewAccount(wallet.PaulcAccount)
	babayaga := database.NewAccount(wallet.DavecAccount)

	// Remove the test directory if it already exists
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	// Required for AddPendingTX() to describe
	// from what node the TX came from (local node in this case)
	nInfo := NewPeerNode(
		"127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true,
	)

	n := New(datadir, nInfo.IP, nInfo.Port, babayaga, nInfo)

	// Allow the test to run for 30 mins, in the worst case
	ctx, closeNode := context.WithTimeout(context.Background(), time.Minute*30)

	tx1 := database.NewTx(andrej, babayaga, 1, "")
	tx2 := database.NewTx(andrej, babayaga, 2, "")
	tx2Hash, _ := tx2.Hash()

	// Pre-mine a valid block without running the `n.Run()`
	// with Paulc as a miner who will receive the block reward,
	// to simulate the block came on the fly from another peer
	validPreMinedPb := NewPendingBlock(database.Hash{}, 0, andrej, []database.Tx{tx1})
	validSyncedBlock, err := Mine(ctx, validPreMinedPb)
	if err != nil {
		t.Fatal(err)
	}

	// Add 2 new TXs into the Davec's node
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds - 2))

		err := n.AddPendingTX(tx1, nInfo)
		if err != nil {
			t.Fatal(err)
		}

		err = n.AddPendingTX(tx2, nInfo)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Once the Davec is mining the block, simulate that
	// Paulc mined the block with TX1 in it faster
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should be mining")
		}

		_, err := n.state.AddBlock(validSyncedBlock)
		if err != nil {
			t.Fatal(err)
		}
		// Mock the Paulc's block came from a network
		n.newSyncedBlocks <- validSyncedBlock

		time.Sleep(time.Second * 2)
		if n.isMining {
			t.Fatal("synced block should have canceled mining")
		}

		// Mined TX1 by Paulc should be removed from the Mempool
		_, onlyTX2IsPending := n.pendingTXs[tx2Hash.Hex()]

		if len(n.pendingTXs) != 1 && !onlyTX2IsPending {
			t.Fatal("synced block should have canceled mining of already mined TX")
		}

		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should be mining again the 1 TX not included in synced block")
		}
	}()

	go func() {
		// Regularly check whenever both TXs are now mined
		ticker := time.NewTicker(time.Second * 10)

		for {
			select {
			case <-ticker.C:
				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 2)

		// Take a snapshot of the DB balances
		// before the mining is finished and the 2 blocks
		// are created.
		startingPaulcBalance := n.state.Balances[andrej]
		startingDavecBalance := n.state.Balances[babayaga]

		// Wait until the 30 mins timeout is reached or
		// the 2 blocks got already mined and the closeNode() was triggered
		<-ctx.Done()

		endPaulcBalance := n.state.Balances[andrej]
		endDavecBalance := n.state.Balances[babayaga]

		// In TX1 Paulc transferred 1 TBB token to Davec
		// In TX2 Paulc transferred 2 TBB tokens to Davec
		expectedEndPaulcBalance := startingPaulcBalance - tx1.Value - tx2.Value + database.BlockReward
		expectedEndDavecBalance := startingDavecBalance + tx1.Value + tx2.Value + database.BlockReward

		if endPaulcBalance != expectedEndPaulcBalance {
			t.Fatalf("Paulc expected end balance is %d not %d", expectedEndPaulcBalance, endPaulcBalance)
		}

		if endDavecBalance != expectedEndDavecBalance {
			t.Fatalf("Davec expected end balance is %d not %d", expectedEndDavecBalance, endDavecBalance)
		}

		t.Logf("Starting Paulc balance: %d", startingPaulcBalance)
		t.Logf("Starting Davec balance: %d", startingDavecBalance)
		t.Logf("Ending Paulc balance: %d", endPaulcBalance)
		t.Logf("Ending Davec balance: %d", endDavecBalance)
	}()

	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("was suppose to mine 2 pending TX into 2 valid blocks under 30m")
	}

	if len(n.pendingTXs) != 0 {
		t.Fatal("no pending TXs should be left to mine")
	}
}

func getTestDataDirPath() string {
	return filepath.Join(os.TempDir(), ".tbb_test")
}
