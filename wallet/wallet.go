package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paulcockrell/blockchain/database"
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

func NewKeystoreAccount(dataDir, password string) (common.Address, error) {
	ks := keystore.NewKeyStore(
		GetKeystoreDirPath(dataDir),
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	acc, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	}

	return acc.Address, nil
}

func SignTx(tx database.Tx, privKey *ecdsa.PrivateKey) (database.SignedTx, error) {
	rawTx, err := tx.Encode()
	if err != nil {
		return database.SignedTx{}, err
	}

	sig, err := Sign(rawTx, privKey)
	if err != nil {
		return database.SignedTx{}, err
	}

	return database.NewSignedTx(tx, sig), nil
}

func SignTxWithKeystoreAccount(tx database.Tx, acc common.Address, pwd, keystoreDir string) (database.SignedTx, error) {
	ks := keystore.NewKeyStore(
		keystoreDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	ksAccount, err := ks.Find(accounts.Account{
		Address: acc,
	})
	if err != nil {
		return database.SignedTx{}, err
	}

	ksAccountJson, err := ioutil.ReadFile(ksAccount.URL.Path)
	if err != nil {
		return database.SignedTx{}, err
	}

	key, err := keystore.DecryptKey(ksAccountJson, pwd)
	if err != nil {
		return database.SignedTx{}, err
	}

	signedTx, err := SignTx(tx, key.PrivateKey)
	if err != nil {
		return database.SignedTx{}, err
	}

	return signedTx, nil
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
