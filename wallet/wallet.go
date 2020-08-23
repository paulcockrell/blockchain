package wallet

import "path/filepath"

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
