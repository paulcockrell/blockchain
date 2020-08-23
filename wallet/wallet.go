package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
)

const keystoreDirName = "keystore"

// Datastore filepath ~/.tbb_paulc
const PaulcAccount = "0xb61E2B65e6066b0575EdD91f992B8ee8Dbd96481"

// Datastore filepath ~/.tbb_davec
const DavecAccount = "0x86d4ba480A6b65C21A1d652C51E15576b269A5E7"

// Datastore filepath ~/.tbb_kimc
const KimcAccount = "0xEC205dcb742008680CA0a53a0748Fd5C14E3f769"

func GetKeystoreDirPath(dataDir string) string {
	return filepath.Join(dataDir, keystoreDirName)
}

func Sign(msg []byte, privKey *ecdsa.PrivateKey) (sig []byte, err error) {
	msgHash := sha256.Sum256(msg)

	return crypto.Sign(msgHash[:], privKey)
}

func Verify(msg, sig []byte) (*ecdsa.PublicKey, error) {
	msgHash := sha256.Sum256(msg)

	recoveredPubKey, err := crypto.SigToPub(msgHash[:], sig)
	if err != nil {
		return nil, fmt.Errorf("unable to verify message signature %s", err.Error())
	}

	return recoveredPubKey, nil
}
