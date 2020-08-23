package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignCryptoParams(t *testing.T) {
	// Generate the key on the fly
	privKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Private key:")
	spew.Dump(privKey)

	// Prepare a message to digitally sign
	msg := []byte("Cats are awesome")

	// Sign it
	sig, err := Sign(msg, privKey)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the length is 65 bytes
	if len(sig) != crypto.SignatureLength {
		t.Fatal(fmt.Errorf("wrong size for signature: got %d want %d", len(sig), crypto.SignatureLength))
	}

	// Print the 3 required Ethereum signature crypto values
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := new(big.Int).SetBytes([]byte{sig[64]})

	fmt.Println("R, S, V:")
	spew.Dump(r, s, v)
}

func TestSign(t *testing.T) {
	// Generate the key on the fly
	privKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Convert the Public Key to bytes with elliptic curve settings
	pubKey := privKey.PublicKey
	pubKeyBytes := elliptic.Marshal(crypto.S256(), pubKey.X, pubKey.Y)

	// Hash the Public key to 32 bytes
	pubKeyBytesHash := crypto.Keccak256(pubKeyBytes[1:])

	// The last 20 bytes of the Public key hash will be it's public username
	account := common.BytesToAddress(pubKeyBytesHash[12:])

	msg := []byte("Cats are super amazing awww")

	// Sign a message -> generate message's signature
	sig, err := Sign(msg, privKey)
	if err != nil {
		t.Fatal(err)
	}

	// Recover a Public key from the signature
	recoveredPubKey, err := Verify(msg, sig)
	if err != nil {
		t.Fatal(err)
	}

	// Convert the Public key to username again
	recoveredPubKeyBytes := elliptic.Marshal(
		crypto.S256(),
		recoveredPubKey.X,
		recoveredPubKey.Y,
	)
	recoveredPubKeyBytesHash := crypto.Keccak256(
		recoveredPubKeyBytes[1:],
	)
	recoveredAccount := common.BytesToAddress(
		recoveredPubKeyBytesHash[12:],
	)

	// Compare the username matches meaning,
	// The signature generation and account verification
	// by extracting the Pub key from signature works
	if account.Hex() != recoveredAccount.Hex() {
		t.Fatalf(
			"msg was signed by account %s but signature recovery produced an account %s",
			account.Hex(),
			recoveredAccount.Hex(),
		)
	}
}
